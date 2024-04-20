package redis

import (
	"ForumWeb/pkg/ws"
	"github.com/gorilla/websocket"
)

func StartRedisPubSub() {
	pubsub := rdb.Subscribe("postVotes")
	defer pubsub.Close()
	for {
		msg, err := pubsub.ReceiveMessage()
		if err != nil {
			continue
		}
		clients := ws.Clients
		// 向所有连接的 WebSocket 客户端发送消息
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}
