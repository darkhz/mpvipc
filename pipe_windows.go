// +build windows

package mpvipc

import (
	"net"

	"gopkg.in/natefinch/npipe.v2"
)

func dial(path string) (net.Conn, error) {
	// npipe's Dial normally blocks indefinitely until a
	// connection is found. Use DialTimeout.
	return npipe.DialTimeout(path, 1000)
}
