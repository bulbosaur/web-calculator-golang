package agent

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bulbosaur/web-calculator-golang/internal/handlers"
	"github.com/spf13/viper"
)

func RunAgent() {
	host := viper.GetString("server.AG_HOST")
	port := viper.GetString("server.AG_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	http.HandleFunc("POST /internal/task", handlers.TaskHandler)

	log.Printf("Agent starting on %s", addr)
	err := http.ListenAndServe(addr, nil)

	if err != nil {
		log.Fatal("Agent server error:", err)
	}
}
