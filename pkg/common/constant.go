package common

import "time"

const (
	ResourceName		string = "test.com/device"
	DevicePath 			string = "/etc/demo"
	DeviceSocket 		string = "demo.sock"
	ConnectTimeout			   = time.Second * 5
)
