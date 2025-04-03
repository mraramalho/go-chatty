package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	secretKey  string
	clients    map[net.Conn]*Client
	mutex      sync.Mutex
	broadcast  chan Message
	shutdownWg sync.WaitGroup
}

type Client struct {
	conn     net.Conn
	nickname string
}

type Message struct {
	sender  *Client
	message string
}

func NewServer() *Server {
	return &Server{
		clients:   make(map[net.Conn]*Client),
		mutex:     sync.Mutex{},
		broadcast: make(chan Message),
	}
}

func (s *Server) ServeAndListen() {

	serverIP := os.Getenv("SERVER_IP")
	if serverIP == "" {
		serverIP = "0.0.0.0"
	}
	fmt.Println("Server IP:", serverIP)

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "3333"
	}

	s.secretKey = getSecretKey()

	listenAddr := net.JoinHostPort(serverIP, serverPort)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("Server listening on", listenAddr)

	go s.gracefulShutdown(listener)
	go s.handleBroadcast()

	for {
		conn, err := listener.Accept()
		if err != nil {
			// Check if the error is because the listener was closed
			if strings.Contains(err.Error(), "use of closed network connection") {
				fmt.Println("Listener closed, stopping accept loop")
				break
			}
			fmt.Println("Error accepting connection:", err)
			continue
		}
		s.shutdownWg.Add(1)
		go s.handleAuthentication(conn)
	}
}

// getSecretKey gets the secret key from the user to be used in the server
// as a password to access the server and encript the messages.
func getSecretKey() string {
	fmt.Println("Define a server secret key:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func (s *Server) gracefulShutdown(listener net.Listener) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	fmt.Println("Shutting down server...")

	listener.Close()

	// Give clients a chance to disconnect gracefully
	fmt.Println("Closing all client connections...")
	s.mutex.Lock()
	for conn, client := range s.clients {
		fmt.Printf("Closing connection for client %s\n", client.nickname)
		conn.Close()
	}
	s.mutex.Unlock()

	// Wait for all client handlers to finish
	fmt.Println("Waiting for all client handlers to finish...")
	s.shutdownWg.Wait()

	// Close the broadcast channel
	close(s.broadcast)

	fmt.Println("All connections closed. Server shut down gracefully.")
	os.Exit(0)
}

func (s *Server) handleAuthentication(conn net.Conn) {
	defer func() {
		conn.Close()
		s.shutdownWg.Done()
	}()

	// Set a deadline for authentication
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	// Send authentication prompt
	fmt.Fprintf(conn, "Enter the secret key to access the server:\n")

	reader := bufio.NewReader(conn)
	clientSecretKey, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading secret key:", err)
		return
	}

	if strings.TrimSpace(clientSecretKey) != s.secretKey {
		_, err := fmt.Fprintf(conn, "Invalid secret key. Connection closed.\n")
		if err != nil {
			fmt.Println("Error sending invalid secret key message:", err)
		}
		fmt.Println("Client provided invalid secret key. Connection closed.")
		return
	}

	// Reset deadline after successful authentication
	conn.SetDeadline(time.Time{})

	// Ask for nickname
	fmt.Fprintf(conn, "Welcome to the chat! Type your nickname:\n")
	nickname, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading nickname:", err)
		return
	}

	client := &Client{
		conn:     conn,
		nickname: strings.TrimSpace(nickname),
	}

	s.mutex.Lock()
	s.clients[conn] = client
	s.mutex.Unlock()

	// Notify all clients about the new user
	s.broadcast <- Message{
		sender:  client,
		message: "has joined the chat!",
	}

	s.handleClient(conn)

}

func (s *Server) handleClient(conn net.Conn) {
	defer func() {
		s.mutex.Lock()
		client, exists := s.clients[conn]
		if exists {
			// Notify others that this client has disconnected
			s.broadcast <- Message{
				sender:  client,
				message: "has left the chat.",
			}
			fmt.Printf("Client %s disconnected\n", client.nickname)
			delete(s.clients, conn)
		}
		s.mutex.Unlock()
		conn.Close()
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()

		s.mutex.Lock()
		client, exists := s.clients[conn]
		s.mutex.Unlock()

		if !exists {
			return
		}

		s.broadcast <- Message{
			sender:  client,
			message: message,
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Erro na conexão:", err)
	}

}

func (s *Server) handleBroadcast() {
	for message := range s.broadcast {
		s.mutex.Lock()
		for _, client := range s.clients {
			if client.conn != message.sender.conn {
				_, err := fmt.Fprintf(client.conn, "[%s]: %s\n", message.sender.nickname, message.message)
				if err != nil {
					fmt.Printf("Error sending message to %s: %v\n", client.nickname, err)
				}
			}
		}
		s.mutex.Unlock()
	}
}
