// +build windows

package mpvipc

import (
	"net"

	"gopkg.in/natefinch/npipe.v2"
)

func dial(path string) (net.Conn, error) {
	return npipe.DialTimeout(path, 1000)
}
