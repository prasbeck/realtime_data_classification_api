package websocket

import (
    "log"
    "net/http"
    "sync"

    "github.com/gorilla/websocket"
)

type Client struct {
    Conn *websocket.Conn
    Send chan []byte
}

type Manager struct {
    Clients    map[*Client]bool
    Broadcast  chan []byte
    Register   chan *Client
    Unregister chan *Client
    Mutex      sync.Mutex
}

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func NewManager() *Manager {
    return &Manager{
        Clients:    make(map[*Client]bool),
        Broadcast:  make(chan []byte),
        Register:   make(chan *Client),
        Unregister: make(chan *Client),
    }
}

func (m *Manager) Start() {
    for {
        select {
        case client := <-m.Register:
            m.Clients[client] = true
        case client := <-m.Unregister:
            if _, ok := m.Clients[client]; ok {
                delete(m.Clients, client)
                close(client.Send)
            }
        case message := <-m.Broadcast:
            for client := range m.Clients {
                select {
                case client.Send <- message:
                default:
                    close(client.Send)
                    delete(m.Clients, client)
                }
            }
        }
    }
}

func (m *Manager) HandleConnections(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    client := &Client{Conn: conn, Send: make(chan []byte)}
    m.Register <- client

    go m.HandleMessages(client)
}

func (m *Manager) HandleMessages(client *Client) {
    defer func() {
        m.Unregister <- client
        client.Conn.Close()
    }()

    for {
        _, message, err := client.Conn.ReadMessage()
        if err != nil {
            log.Println(err)
            break
        }
        m.Broadcast <- message
    }
}

func (m *Manager) WriteMessages(client *Client) {
    for message := range client.Send {
        err := client.Conn.WriteMessage(websocket.TextMessage, message)
        if err != nil {
            log.Println(err)
            client.Conn.Close()
            break
        }
    }
}
