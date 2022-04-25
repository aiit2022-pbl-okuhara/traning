package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/aiit2022-pbl-okuhara/traning/micchie/sample-app/config"
	"github.com/aiit2022-pbl-okuhara/traning/micchie/sample-app/infrastructure/db"
	"github.com/aiit2022-pbl-okuhara/traning/micchie/sample-app/infrastructure/server"
)

func main() {
	conn, err := db.NewDB()
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}
	defer func(conn *sql.DB) {
		if err := conn.Close(); err != nil {
			log.Fatalf("Failed to close the database connection: %v", err)
		}
	}(conn)

	srv := server.NewServer(
		sampleHandler(conn),
	)

	log.Printf("Serving on localhost:%v\n", config.Config.ServerPort)
	log.Fatal(srv.ListenAndServe())
}

func sampleHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprint(w, "Hello World")
		// TODO: select from DB
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
