package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mraramalho/go-chatty/peer"
	"github.com/mraramalho/go-chatty/server"
)

func main() {
	// Define command line flags
	instanceType := flag.String(
		"instanceType",
		"server",
		"Choose 'server' or 'client' to run a Server or a Client instance. Default is server.",
	)

	serverIP := flag.String(
		"ip",
		"",
		"Server IP address. Overrides SERVER_IP environment variable.",
	)

	serverPort := flag.String(
		"port",
		"",
		"Server port. Overrides SERVER_PORT environment variable.",
	)

	// Parse command line arguments
	flag.Parse()

	// Set environment variables if provided via flags
	if *serverIP != "" {
		os.Setenv("SERVER_IP", *serverIP)
	}

	if *serverPort != "" {
		os.Setenv("SERVER_PORT", *serverPort)
	}

	// Print application banner
	printBanner()

	// Run the appropriate mode
	switch strings.ToLower(*instanceType) {
	case "server":
		fmt.Println("Starting in SERVER mode...")
		server := server.NewServer()
		server.ServeAndListen()

	case "client":
		fmt.Println("Starting in CLIENT mode...")
		client := peer.NewClient()
		client.Connect()

	default:
		fmt.Printf("Invalid mode: %s\n", *instanceType)
		fmt.Println("Use 'server' to initialize a server or 'client' to connect to a running chat.")
		flag.Usage()
		os.Exit(1)
	}
}

// printBanner prints a welcome message
func printBanner() {
	banner := `
    ____            ____ _           _   _         
   / ___| ___      / ___| |__   __ _| |_| |_ _   _ 
  | |  _ / _ \____| |   | '_ \ / _' | __| __| | | |
  | |_| | (_) |___| |___| | | | (_| | |_| |_| |_| |
   \____|\___/     \____|_| |_|\__,_|\__|\__|\__, |
                                             |___/ 
   A simple TCP chat application
   ---------------------------
`
	fmt.Println(banner)
}
