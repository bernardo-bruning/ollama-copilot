package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func Proxy(port string, forward string) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("failed to accept: %v", err)
		}

		go handle(conn, forward)
	}
}

var hosts = []string{
	"api.githubcopilot.com",
	"api.github.com",
	"copilot-proxy.githubusercontent.com",
	"proxy.individual.githubcopilot.com",
	"proxy.business.githubcopilot.com",
	"copilot-telemetry.githubusercontent.com",
}

func handle(conn net.Conn, forward string) {
	req, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		conn.Close()
		log.Printf("failed to read request: %v", err)
		return
	}

	address := fmt.Sprintf("%s:%s", req.URL.Hostname(), req.URL.Port())

	for _, host := range hosts {
		if strings.Contains(req.URL.Hostname(), host) {
			// This is a host we know and want to forward back to ourselves
			address = "localhost" + forward
			break
		}
	}

	if req.Method != http.MethodConnect {
		conn.Close()
		log.Printf("unsupported method: %s", req.Method)
		return
	}

	client, err := net.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		conn.Close()
		log.Printf("failed to dial: %v", err)
		_, err = conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
		if err != nil {
			log.Printf("failed to write response: %v", err)
		}
		return
	}

	_, err = conn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	if err != nil {
		conn.Close()
		log.Printf("failed to write response: %v", err)
		return
	}

	go transfer(client, conn)
	go transfer(conn, client)
}

func transfer(w io.WriteCloser, r io.ReadCloser) {
	defer w.Close()
	defer r.Close()
	_, err := io.Copy(w, r)
	if errors.Is(err, net.ErrClosed) {
		return
	}

	if err == net.ErrClosed {
		return
	}

	if err != nil {
		log.Printf("failed to transfer: %v", err)
	}
}
