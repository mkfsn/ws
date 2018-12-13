package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	headers, cookies, url := parseArgs()
	requestHeader := http.Header{}
	for k, v := range headers.Map() {
		requestHeader.Set(k, v)
	}
	for _, v := range []string(cookies) {
		requestHeader.Add("Cookie", v)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		client, err := NewClient(url, requestHeader)
		if err != nil {
			panic(err)
		}
		defer client.Close()
		if err != nil {
			log.Println("write close:", err)
			return
		}

		for {
			select {
			case <-client.done:
				return
			case message := <-client.Receive():
				log.Println(string(message))
			}
		}
	}()

	<-interrupt
	log.Println("interrupt")
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: ./ws [-cookie value] [-header value] host\n\nOptions:\n")
	flag.PrintDefaults()
}

func parseArgs() (arrayFlags, arrayFlags, string) {
	var headers arrayFlags
	var cookies arrayFlags

	flag.Var(&headers, "header", "HTTP Header: -header KEY=VALUE, e.g. -header accept=text/html")
	flag.Var(&cookies, "cookie", "HTTP Cookie: -cookie KEY=VALUE, e.g. -cookie FOO=bar")
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		log.Fatal(fmt.Errorf("Invalid arguments"))
	}

	return headers, cookies, args[0]
}
