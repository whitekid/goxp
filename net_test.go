package goxp

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAvailablePort(t *testing.T) {
	port, err := AvailablePort()
	require.NoError(t, err)

	require.Greater(t, port, 1024)
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	require.NoError(t, err)
	defer ln.Close()
}

func TestAvailableUdpPort(t *testing.T) {
	port, err := AvailableUdpPort()
	require.NoError(t, err)

	require.Greater(t, port, 1024)
	ln, err := net.ListenUDP("udp", &net.UDPAddr{Port: port, IP: net.ParseIP("0.0.0.0")})
	require.NoError(t, err)
	defer ln.Close()

}

func TestURLToListenAddr(t *testing.T) {
	tests := [...]struct {
		name     string
		wantErr  bool
		wantHost string
		wantFull string
	}{
		{"http://127.0.0.1", false, "127.0.0.1:80", "http://127.0.0.1:80"},
		{"http://127.0.0.1:8080/xx", false, "127.0.0.1:8080", "http://127.0.0.1:8080/xx"},
		{"http://:8080/xx", false, ":8080", "http://:8080/xx"},
		{"http://127.0.0.1:0/xx", false, "127.0.0.1:*", "http://127.0.0.1:*/xx"},
		{"http://:0", false, ":*", "http://:*"},
		{"tcp://:0", false, ":*", "tcp://:*"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			host, port, full, err := URLToListenAddr(tt.name)
			if (err != nil) != tt.wantErr {
				require.Failf(t, "URLToListenAddr failed", "error = %v, wantErr = %t", err, tt.wantErr)
			}

			if strings.ContainsRune(tt.wantHost, '*') {
				tt.wantHost = strings.ReplaceAll(tt.wantHost, "*", port)
			}

			if strings.ContainsRune(tt.wantFull, '*') {
				tt.wantFull = strings.ReplaceAll(tt.wantFull, "*", port)
			}

			require.Equal(t, tt.wantHost, host)

			if v, err := strconv.Atoi(port); err != nil && v > 1024 {
				if _, err := net.Listen("tcp", host); err != nil {
					require.Failf(t, "listen failed", "addr = %s", host)
				}
			}

			require.Equal(t, tt.wantFull, full)
		})
	}
}
