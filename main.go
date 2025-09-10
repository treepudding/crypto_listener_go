package main

import (
	"context"
	"log"
	"trading/infra/ws"

	githubjson "github.com/cloudwego/hertz/pkg/common/json"
	"github.com/gorilla/websocket"
)

//	func main() {
//		// 创建一个 Hertz 服务器实例
//		h := server.Default()
//
//		RegisterRoutes(h)
//		// 启动服务器，默认监听 0.0.0.0:8888
//		h.Spin()
//	}

func main() {
	ctx := context.Background()
	m, err := ws.Init(ctx, ws.Config{
		EnablePublic:     true,
		EnablePrivate:    true,
		SimulatedTrading: true,
		APIKey:           "565f6b64-78b6-43e7-89b9-9d767ccf609d", // 从安全来源注入
		SecretKey:        "3BD0978D9FD4854D4A89D1C620A01067",     // 从安全来源注入
		Passphrase:       "J1J1wanokx!",                          // 从安全来源注入
	})
	if err != nil {
		panic(err)
	}

	// 私有：登录成功后订阅账户和持仓
	if m.PrivateConn() != nil {
		subscribeReq := map[string]interface{}{
			"id":   "12312",
			"op":   "subscribe",
			"args": []interface{}{map[string]string{"channel": "balance_and_position"}},
		}
		subMsg, _ := githubjson.Marshal(subscribeReq)
		if err := m.PrivateConn().WriteMessage(websocket.TextMessage, subMsg); err != nil {
			log.Fatal("Subscribe write error:", err)
		}
		go func() {
			for {
				_, message, err := m.PrivateConn().ReadMessage()
				if err != nil {
					log.Println("Read error:", err)
					return
				}
				log.Printf("Received: %s", message)
			}
		}()
	}

	// 公共：示例订阅
	if m.PublicConn() != nil {
		_ = m.SubscribePublic([]map[string]string{
			{"channel": "tickers", "instId": "BTC-USDT"},
			{"channel": "tickers", "instId": "ETH-USDT"},
		})
		go func() {
			for {
				_, message, err := m.PublicConn().ReadMessage()
				if err != nil {
					log.Println("Read error:", err)
					return
				}
				log.Printf("Received: %s", message)
			}
		}()
	}

	select {}
}

func demo() {

}
