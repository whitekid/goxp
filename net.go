package goxp

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
)

// AvailablePort return any available TCP ports
func AvailablePort() (int, error) {
	r, err := AvailablePorts(1)
	if err != nil {
		return 0, err
	}

	return r[0], nil
}

// AvailablePorts return any available TCP ports
func AvailablePorts(count int) ([]int, error) {
	ports := make([]int, count)
	for i := 0; i < len(ports); i++ {
		ln, err := net.Listen("tcp", ":")
		if err != nil {
			return nil, err
		}
		defer ln.Close()

		ports[i] = ln.Addr().(*net.TCPAddr).Port
	}

	return ports, nil
}

// AvailableUdpPort return any available UDP ports
func AvailableUdpPort() (int, error) {
	r, err := AvailableUdpPorts(1)
	if err != nil {
		return 0, err
	}

	return r[0], nil
}

// AvailableUdpPorts return any available UDP ports
func AvailableUdpPorts(count int) ([]int, error) {
	ports := make([]int, count)
	for i := 0; i < len(ports); i++ {
		ln, err := net.ListenUDP("udp", &net.UDPAddr{Port: 0, IP: net.ParseIP(":")})
		if err != nil {
			return nil, err
		}
		defer ln.Close()
		ports[i] = ln.LocalAddr().(*net.UDPAddr).Port
	}

	return ports, nil
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
		p, err := AvailablePort()
		if err != nil {
			return "", "", "", err
		}
		port = strconv.Itoa(p)
	}

	u.Host = hostname + ":" + port

	return hostname + ":" + port, port, u.String(), nil
}
