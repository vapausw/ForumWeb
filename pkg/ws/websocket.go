package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 注意安全性设置，这里应根据实际情况配置
	},
}

var Clients = make(map[*websocket.Conn]bool) // 连接的客户端集合

func WebsocketHandler(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set websocket upgrade: " + err.Error()})
		return
	}
	defer ws.Close()

	Clients[ws] = true

	for {
		// 读取消息
		mt, message, err := ws.ReadMessage()
		if err != nil {
			delete(Clients, ws)
			break
		}
		// 可以在这里处理收到的消息
		println(string(message))

		// 发送消息
		err = ws.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}
