package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/google/go-safeweb/safehttp"

	"github.com/aiit2022-pbl-okuhara/traning/micchie/sample-web-app/secure"
	"github.com/aiit2022-pbl-okuhara/traning/micchie/sample-web-app/server"
	"github.com/aiit2022-pbl-okuhara/traning/micchie/sample-web-app/storage"
)

var (
	port = flag.Int("port", 8080, "Port for the HTTP server")
	dev  = flag.Bool("dev", false, "Run in development mode")
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	flag.Parse()
	safehttp.UseLocalDev()
	if *dev {
		safehttp.UseLocalDev()
	}
	db := storage.NewDB()

	addr := net.JoinHostPort("localhost", strconv.Itoa(*port))
	mux := secure.NewMuxConfig(db, addr).Mux()
	server.Load(db, mux)

	log.Printf("Listening on %q", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
