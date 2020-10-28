package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/dstotijn/tonny"
)

var addr = ":8080"

func main() {
	teeLn, err := tonny.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		ConnState: func(conn net.Conn, state http.ConnState) {
			if state != http.StateIdle && state != http.StateClosed {
				return
			}
			teeConn, ok := conn.(tonny.TeeConn)
			if !ok {
				return
			}

			io.Copy(os.Stdout, teeConn.ReadBuffer)
			io.Copy(os.Stdout, teeConn.WriteBuffer)
			fmt.Printf("\n")
		},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "OK")
		}),
		// Disable HTTP/2, because ConnState will not work as expected when it's
		// enabled.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	fmt.Printf("Listening on %v ...\n", addr)
	srv.Serve(teeLn)
}
