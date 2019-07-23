package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// port
const PORT = ":13556"

type message struct {
	UserName string `json:"username"`
	Text     string `json:"text"`
}

type file struct {
	FileName string `json:"filename"`
	Data     []byte `json:""`
}

// Setup TLS connection and listen
func listener(config *tls.Config, wg *sync.WaitGroup, alivech chan bool) {
	defer wg.Done()
	peer, err := tls.Listen("tcp", PORT, config)
	if err != nil {
		log.Println("TLS listening error!")
		return
	}
	defer peer.Close()

	for {
		fmt.Println()
		log.Println("Start listening on port:", PORT)
		conn, err := peer.Accept()
		if err != nil {
			continue
		}
		
		var msg = message{}
		for {
			err := json.NewDecoder(conn).Decode(&msg)
			// Connection not alive
			if err != nil {
				alivech <- false
				log.Println("Disconnected!")
				break
			}
			if msg.Text != "" {
				log.Printf("%s: %s", msg.UserName, msg.Text)
				fmt.Println(">>")
			}
		}
	}
}

// Send the message to peer
func sender(remoteAddr string, config *tls.Config, msgch chan message, alivech chan bool) {
CONNECTION:
	for {
		conn, err := tls.Dial("tcp", remoteAddr+PORT, config)
		if err != nil {
			log.Println("Connecting...")	
			time.Sleep(2 * time.Second)
			continue
		}
		log.Println("Connected to another peer!")
		fmt.Println(">>")

		for {
			select {
				case msg := <-msgch: 
					err := json.NewEncoder(conn).Encode(msg)
					if err != nil {
						continue CONNECTION
					}
				case alive := <-alivech:
					if !alive {
						continue CONNECTION
					}
			}
		}
	}
}

// Read message from stdin 
func reader(msgch chan message, username string) {
	var in = bufio.NewReader(os.Stdin)
	var msg = message{UserName: username}
	for {
		fmt.Println(">>")
		if text, _ := in.ReadString('\n'); text != "\n" {
			msg.Text = strings.TrimSpace(text)
			msgch <- msg
		} 
	}
}


func main() {
	certSender, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
	if err != nil {
		log.Fatalln("client: load keys error", err)
	}
	configSender := tls.Config{
		Certificates: []tls.Certificate{certSender},
		InsecureSkipVerify: false,
	}

	certListener, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
	if err != nil {
		log.Fatalln("server: load keys error", err)
	}
	configListener := tls.Config{
		Certificates: []tls.Certificate{certListener},
		InsecureSkipVerify: false,
	}

	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("\tpeer <username> <target_ip>")
		return
	}

	username := os.Args[1]
	remoteAddr := os.Args[2]
	var wg sync.WaitGroup
	msgch := make(chan message)
	alivech := make(chan bool)

	fmt.Println("Welcome to talk free,", username, "\tRemote target:", remoteAddr)

	go listener(&configListener, &wg, alivech)
	go sender(remoteAddr, &configSender, msgch, alivech)
	go reader(msgch, username)
	wg.Add(3)

	wg.Wait()
}
