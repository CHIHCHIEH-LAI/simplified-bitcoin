package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	nodeType := flag.String("type", "", "Type of process: 'node' or 'wallet'")
	port := flag.String("port", "", "Port to run the node(only for node process)")
	flag.Parse()

	if *nodeType == "" {
		log.Println("Please specify the type of process to run")
		flag.Usage()
		os.Exit(1)
	}

	switch *nodeType {
	case "node":
		if *port == "" {
			log.Println("Please specify the port to run the node")
			flag.Usage()
			os.Exit(1)
		}
	case "wallet":
	default:
		log.Println("Invalid process type. Use 'node' or 'wallet'")
		flag.Usage()
		os.Exit(1)
	}
}
