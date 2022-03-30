package utils

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
)

// AvailablePort return any available ports
func AvailablePort() int {
	ln, err := net.Listen("tcp", ":")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	return ln.Addr().(*net.TCPAddr).Port
}

// URLToListenAddr parse URL and return correspend listen address
//
// http://127.0.0.1:8080/xx returns (127.0.0.1:80, 80, http://127.0.0.1:8080/xx)
func URLToListenAddr(addr string) (string, string, string, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return "", "", "", err
	}

	hostname := u.Hostname()
	port := u.Port()

	switch port {
	case "":
		switch u.Scheme {
		case "http":
			port = "80"
		case "https":
			port = "443"
		default:
			return "", "", "", fmt.Errorf("unsupported scheme: %s", u.Scheme)
		}
	case "0":
		port = strconv.Itoa(AvailablePort())
	}

	u.Host = hostname + ":" + port

	return hostname + ":" + port, port, u.String(), nil
}
