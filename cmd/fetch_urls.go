package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/florinutz/go-tests-kvstore"
)

func init() {
	if len(os.Args) == 1 || os.Args[1] == "help" {
		fmt.Println("each argument is an url to fetch. All of their contents will be outputted to stdout, so you'd " +
			"better save it straight to a file")
		return
	}

	if len(os.Args) < 2 {
		log.Fatal("you have to provide at least 1 url for fetching")
	}
}

func main() {
	kvs, _ := kvstore.FetchUrls(os.Args[1:], 20*time.Second, nil, nil)

	// writes the output to stdout
	err := kvstore.Write(kvs, os.Stdout, kvstore.DefaultEncoder, kvstore.DefaultContentEncoder)
	if err != nil {
		log.Fatal(err)
	}
}

func gotUrl(url string, contents []byte, err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "* failed fetching url %s\nerr: %s\n", url, err)
		return
	}
	fmt.Printf("* fetched url %s\n", url)
}
