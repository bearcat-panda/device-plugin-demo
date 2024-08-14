package device_plugin

import (
	"context"
	"device-plugin-demo/pkg/common"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	"log"
	"net"
	"os"
	"path"
	"syscall"
	"time"
)

type DevicePluginDemo struct {
	server *grpc.Server
	stop 	chan struct{} // this channel signals to stop the device plugin
	dd     *DeviceDemo
}

func NewDevicePluginDemo() *DevicePluginDemo {
	return &DevicePluginDemo{
		server: grpc.NewServer(grpc.EmptyServerOption{}),
		stop: 	make(chan struct{}),
		dd:		NewDeviceDemo(common.DevicePath),
	}
}

// Run start gRPC server and watcher
func (c *DevicePluginDemo) Run() error {
	err := c.dd.List()
	if err != nil {
		log.Fatalf("list device error: %v", err)
	}

	go func() {
		if err = c.dd.watch(); err != nil {
			log.Println("watch devices error")
		}
	}()

	pluginapi.RegisterDevicePluginServer(c.server, c)
	// delete old unix socket before start
	socket := path.Join(pluginapi.DevicePluginPath, common.DeviceSocket)
	err = syscall.Unlink(socket)
	if err != nil && !os.IsNotExist(err) {
		return errors.WithMessagef(err, "delete socket %s failed", socket)
	}

	sock, err := net.Listen("unix", socket)
	if err != nil {
		return errors.WithMessagef(err, "listen unix %s failed", sock)
	}

	go c.server.Serve(sock)

	// Wait for server to start by launching a blocking connection
	conn, err := connect(common.DeviceSocket, 5*time.Second)
	if err != nil {
		return err
	}
	conn.Close()

	return nil
}

// dial establishes the gRPC communication with the registered device plugin.
func connect(sockerPath string, timeout time.Duration) (*grpc.ClientConn, error)  {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	c, err := grpc.DialContext(ctx, sockerPath,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			if deadline, ok := ctx.Deadline(); ok {
				return net.DialTimeout("unix", addr, time.Until(deadline))
			}
			return net.DialTimeout("unix", addr, common.ConnectTimeout)
		}),
	)
	if err != nil {
		return nil, err
	}

	return c, nil
}




