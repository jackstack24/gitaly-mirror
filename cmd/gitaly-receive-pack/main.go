package main

import (
	"log"
	"os"
	"strings"

	"gitlab.com/gitlab-org/gitaly/client"
)

func main() {
	if !(len(os.Args) >= 3 && strings.HasPrefix(os.Args[2], "git-receive-pack")) {
		log.Fatalf("Not a valid command")
	}

	addr := os.Getenv("GITALY_SOCKET")
	if len(addr) == 0 {
		log.Fatalf("GITALY_SOCKET not set")
	}

	cli, err := client.NewClient(addr)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer cli.Close()

	code, err := cli.ReceivePack(os.Stdin, os.Stdout, os.Stderr, os.Getenv("GL_REPOSITORY"), os.Getenv("GL_ID"))
	if err != nil {
		log.Printf("Error: %v", err)
		os.Exit(int(code))
	}

	os.Exit(int(code))
}
