package socket

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
)

// Hub menyimpan semua koneksi aktif
type Hub struct {
	Clients    map[string]*websocket.Conn
	ClientsMux sync.Mutex
}

// Inisialisasi Hub Global
var GlobalHub = &Hub{
	Clients: make(map[string]*websocket.Conn),
}

// Register untuk menambah client
func (h *Hub) Register(userId string, conn *websocket.Conn) {
	h.ClientsMux.Lock()
	defer h.ClientsMux.Unlock()
	h.Clients[userId] = conn
}

// Unregister untuk menghapus client
func (h *Hub) Unregister(userId string) {
	h.ClientsMux.Lock()
	defer h.ClientsMux.Unlock()
	delete(h.Clients, userId)
}

// Emit untuk mengirim pesan ke user tertentu
func (h *Hub) Emit(userId string, payload interface{}) {
	h.ClientsMux.Lock()
	defer h.ClientsMux.Unlock()
	if conn, ok := h.Clients[userId]; ok {
		conn.WriteJSON(payload)
	}
}
