package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"homework.31/pkg/app"
)

func main() {
	var port string
	flag.StringVar(&port, "port", "default port", "set port")
	flag.Parse()
	fmt.Println(port)
	Localhost := "localhost:" + port
	app.Run(Localhost)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

}
