package internal

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
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

func handle(conn net.Conn, forward string) {
	req, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		conn.Close()
		log.Printf("failed to read request: %v", err)
		return
	}

	if req.Method != http.MethodConnect {
		conn.Close()
		log.Printf("unsupported method: %s", req.Method)
		return
	}

	client, err := net.DialTimeout("tcp", "localhost"+forward, 10*time.Second)
	// client, err := net.DialTimeout("tcp", req.Host, 10*time.Second)
	if err != nil {
		conn.Close()
		log.Printf("failed to dial: %v", err)
		conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
		return
	}

	conn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))

	go transfer(client, conn)
	go transfer(conn, client)
}

func transfer(w io.WriteCloser, r io.ReadCloser) {
	defer w.Close()
	defer r.Close()
	io.Copy(w, r)
}
