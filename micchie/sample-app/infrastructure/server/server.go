package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aiit2022-pbl-okuhara/traning/micchie/sample-app/config"
)

func NewServer(handler http.Handler) *http.Server {
	c := config.Config
	return &http.Server{
		Handler:      handler,
		Addr:         fmt.Sprintf("0.0.0.0:%s", c.ServerPort),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}
