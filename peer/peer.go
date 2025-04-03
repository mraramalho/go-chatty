package peer

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type Client struct {
	ChatWg sync.WaitGroup
}

func NewClient() *Client {
	return &Client{
		ChatWg: sync.WaitGroup{},
	}
}

func (c *Client) Connect() {

	serverIP := os.Getenv("SERVER_IP")
	if serverIP == "" {
		serverIP = "0.0.0.0"
	}
	fmt.Println("Server IP:", serverIP)

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "3333"
	}

	dialAddress := net.JoinHostPort(serverIP, serverPort)

	conn, err := net.Dial("tcp", dialAddress)
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected to server!")

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Create channels for communication
	msgChan := make(chan string)
	errChan := make(chan error)
	done := make(chan struct{})

	// Read server prompt for secret key
	serverReader := bufio.NewReader(conn)
	prompt, err := serverReader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading from server: %v\n", err)
		return
	}
	fmt.Print(prompt) // Display server prompt

	// Read secret key from user
	stdinReader := bufio.NewReader(os.Stdin)
	secretKey, err := stdinReader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading secret key: %v\n", err)
		return
	}

	// Send secret key to server
	_, err = fmt.Fprint(conn, secretKey)
	if err != nil {
		fmt.Printf("Error sending secret key: %v\n", err)
		return
	}

	// Read server response after authentication
	response, err := serverReader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading server response: %v\n", err)
		return
	}
	fmt.Print(response) // Display server response

	// Check if authentication failed
	if strings.Contains(strings.ToLower(response), "invalid") {
		fmt.Println("Authentication failed. Exiting.")
		return
	}

	// Read nickname from user
	nickname, err := stdinReader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading nickname: %v\n", err)
		return
	}

	// Send nickname to server
	_, err = fmt.Fprint(conn, nickname)
	if err != nil {
		fmt.Printf("Error sending nickname: %v\n", err)
		return
	}

	nickname = strings.TrimSpace(nickname)
	fmt.Printf("Connected to chat as %s\n", nickname)
	fmt.Println("Type your messages. Press Ctrl+C to exit.")

	// Goroutine to read messages from server
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			msgChan <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			errChan <- fmt.Errorf("server connection error: %v", err)
		}
		close(done) // Signal that the connection is closed
	}()

	// Goroutine to read messages from user
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("[me]: ")

		for scanner.Scan() {
			message := scanner.Text()

			if message != "" {
				// sends user message to server
				_, err := fmt.Fprintln(conn, message)
				if err != nil {
					errChan <- fmt.Errorf("error sending message: %v", err)
					return
				}
			}
			fmt.Print("[me]: ")
		}

		if err := scanner.Err(); err != nil {
			errChan <- fmt.Errorf("error reading input: %v", err)
		}

	}()

	// Main event loop
	for {
		select {
		case <-done:
			fmt.Println("\nServer closed the connection.")
			return
		case msg := <-msgChan:
			// Limpa a linha atual (remove o prompt)
			fmt.Print("\r\033[K")
			// Exibe a mensagem recebida
			fmt.Println(msg)
			// Exibe o prompt novamente
			fmt.Print("[me]: ")
		case err := <-errChan:
			fmt.Printf("\nError: %v\n", err)
			return
		case <-sigChan:
			fmt.Println("\nReceived interrupt signal. Closing connection...")
			// Try to gracefully close the connection
			conn.Close()
			return
		}
	}

}
