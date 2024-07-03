package connection

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

type Connection struct {
	Host  string
	Port  string
	Name  string
	Uuid  string
	Key   string
	Token string
}

func (c Connection) serverUrl(path string) *url.URL {
	return &url.URL{
		Scheme: "ws",
		Host:   c.Host + ":" + c.Port,
		Path:   path,
	}
}

func (c Connection) header() http.Header {
	h := http.Header{}
	h.Add("X-Device-UUID", c.Uuid)
	h.Add("X-Device-Name", c.Name)
	if c.Key != "" {
		h.Add("X-Device-Api-Key", c.Key)
	}
	return h
}

func (c Connection) ConnectDataSocket() (*websocket.Conn, error) {

	d := websocket.DefaultDialer
	d.Subprotocols = []string{"serial-tunnel-v1"}

	socket, _, err := d.Dial(c.serverUrl("/api/device/data").String(), c.header())
	if err != nil {
		return nil, fmt.Errorf("dial: %v", err)
	}

	return socket, nil

	// done := make(chan struct{})

	// go func() {
	// 	defer close(done)
	// 	for {
	// 		_, message, err := socket.ReadMessage()
	// 		if err != nil {
	// 			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
	// 				return
	// 			}
	// 			fmt.Println("read:", err)
	// 			return
	// 		}
	// 		fmt.Printf("Received: %s\n", message)
	// 	}
	// }()

	// go func() {
	// 	scanner := bufio.NewScanner(os.Stdin)
	// 	for scanner.Scan() {
	// 		text := scanner.Text()
	// 		err := socket.WriteMessage(websocket.TextMessage, []byte(text))
	// 		if err != nil {
	// 			fmt.Println("write:", err)
	// 			return
	// 		}
	// 	}
	// 	if err := scanner.Err(); err != nil {
	// 		fmt.Println("Error reading from stdin:", err)
	// 	}
	// }()

	// select {
	// case <-done:
	// 	fmt.Println("WebSocket connection closed")
	// case <-interrupt:
	// 	fmt.Println("\nInterrupt signal received, closing connection...")
	// 	err := socket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	// 	if err != nil {
	// 		fmt.Println("write close:", err)
	// 		return err
	// 	}
	// 	select {
	// 	case <-done:
	// 	case <-time.After(time.Second):
	// 	}
	// }
}

func (c Connection) ConnectConfigSocket() (*websocket.Conn, error) {

	d := websocket.DefaultDialer
	d.Subprotocols = []string{"config-tunnel-v1"}

	socket, _, err := d.Dial(c.serverUrl("/api/device/config").String(), c.header())
	if err != nil {
		return nil, fmt.Errorf("dial: %v", err)
	}

	return socket, nil
}
