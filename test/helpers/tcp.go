package helpers

import (
	"net"
)

// GetFreePortTCP asks the kernel for a free open port that is ready to use.
func GetFreePortTCP() (port int, err error) {
	a, err := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
	if err != nil {
		return
	}

	l, err := net.ListenTCP("tcp", a)
	if err != nil {
		return
	}
	defer l.Close()

	port = l.Addr().(*net.TCPAddr).Port
	return
}
