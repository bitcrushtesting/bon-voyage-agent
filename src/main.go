package main

import (
	"bon-voyage-agent/connection"
	"bon-voyage-agent/utils"
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

const version = "1.0.0"

var Commit string

func main() {
	fmt.Println("Agent version", version)

	pluginFolder := flag.String("plugin", "./plugins", "Directory containing plugins")

	// Display help information
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n", "bon-voyage-agent")
		flag.PrintDefaults()
	}
	flag.Parse()

	config := utils.NewConfig()
	err := config.LoadConfig()
	if err != nil {
		fmt.Println("Could not load config:", err)
		os.Exit(1)
	}

	plugins, err := utils.LoadPlugins(*pluginFolder)
	if err != nil {
		fmt.Println("Could not load plugins:", err)
		os.Exit(1)
	}

	router := connection.NewRouter()

	for _, p := range plugins {
		router.HandleFunc(p.Init(nil), p.Call)
	}

	c := connection.Connection{
		Host: config.Server.Host,
		Port: config.Server.Port,
		Name: config.Agent.Name,
		Uuid: config.Agent.Id,
		Key:  config.Server.Key,
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	dataSocket, err := c.ConnectDataSocket()
	if err != nil {
		fmt.Println("Could not connect data socket:", err)
		os.Exit(1)
	}
	configSocket, err := c.ConnectConfigSocket()
	if err != nil {
		fmt.Println("Could not connect config socket:", err)
		os.Exit(1)
	}

	done := make(chan struct{})
	fmt.Println("------------------------------------")
	go func() { // Service the DATA socket
		defer close(done)
		for {
			_, message, err := dataSocket.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					return
				}
				fmt.Println("read:", err)
				return
			}
			fmt.Printf("Received data: %s\n", message)
		}
	}()

	go func() { // Service the CONFIG socket
		defer close(done)
		for {
			_, message, err := configSocket.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					return
				}
				fmt.Println("read:", err)
				return
			}

			reply, err := router.ParseMessage(message)
			if err != nil {
				fmt.Println("Parsing error:", err)
			}
			fmt.Println("Reply:", string(reply))
			configSocket.WriteMessage(websocket.TextMessage, reply)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text()
			err := dataSocket.WriteMessage(websocket.TextMessage, []byte(text))
			if err != nil {
				fmt.Println("write:", err)
				return
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading from stdin:", err)
		}
	}()

	select {
	case <-done:
		fmt.Println("WebSocket connection closed")
	case <-interrupt:
		fmt.Println("\nInterrupt signal received, closing connection...")
		err := dataSocket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			fmt.Println("write close:", err)
		}
		err = configSocket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			fmt.Println("write close:", err)
			return
		}
		select {
		case <-done:
		case <-time.After(time.Second):
		}
	}
	fmt.Println("------------------------------------")
}
