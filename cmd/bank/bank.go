package main

import (
	"github.com/lozovoya/gohomework14_1/pkg/card"
	"gohomework14/cmd/bank/app"
	"log"
	"net"
	"net/http"
	"os"
)

const defaultPort = "9999"
const defaultHost = "0.0.0.0"

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = defaultHost
	}

	log.Println(host)
	log.Println(port)

	if err := execute(net.JoinHostPort(host, port)); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(addr string) (err error) {
	cardSvc := card.NewService()
	mux := http.NewServeMux()
	application := app.NewServer(cardSvc, mux)

	cardSvc.AddHolder("Ivan Ivanov")
	cardSvc.AddHolder("Vasily Petrov")
	cardSvc.AddHolder("Petr Sidorov")
	cardSvc.AddCard("visa", 0, "plastic")
	cardSvc.AddCard("visa", 1, "plastic")
	cardSvc.AddCard("visa", 2, "plastic")

	application.Init()
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return server.ListenAndServe()
}
