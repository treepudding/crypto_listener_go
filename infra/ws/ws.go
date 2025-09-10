package ws

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Infra 层：只负责连接与心跳，不涉及业务

type ConnType int

const (
	ConnPublic ConnType = iota
	ConnPrivate
)

type Config struct {
	// 公共与私有连接是否启用
	EnablePublic  bool
	EnablePrivate bool
	// 模拟盘
	SimulatedTrading bool
	// 鉴权凭证（仅私有连接需要）
	APIKey     string
	SecretKey  string
	Passphrase string
}

type Manager struct {
	publicConn  *websocket.Conn
	privateConn *websocket.Conn
	mu          sync.RWMutex
}

var (
	managerInstance *Manager
	once            sync.Once
	initErr         error
)

// Init 在包外部调用，一次性初始化公共/私有连接
func Init(ctx context.Context, cfg Config) (*Manager, error) {
	once.Do(func() {
		m := &Manager{}

		if cfg.EnablePublic {
			c, err := dial(ctx, PublicWsURL, defaultHeaders(cfg.SimulatedTrading))
			if err != nil {
				initErr = err
				return
			}
			m.publicConn = c
			startHeartbeat(ctx, c, "Public")
		}

		if cfg.EnablePrivate {
			c, err := dial(ctx, PrivateWsURL, defaultHeaders(cfg.SimulatedTrading))
			if err != nil {
				initErr = err
				return
			}
			m.privateConn = c
			startHeartbeat(ctx, c, "Private")

			// 登录仅针对私有连接
			if err := m.login(ctx, cfg); err != nil {
				initErr = err
				return
			}
		}

		managerInstance = m
	})

	if initErr != nil {
		return nil, initErr
	}
	return managerInstance, nil
}

func dial(ctx context.Context, rawURL string, header http.Header) (*websocket.Conn, error) {
	u, _ := url.Parse(rawURL)
	log.Printf("Connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.DialContext(ctx, u.String(), header)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func startHeartbeat(ctx context.Context, conn *websocket.Conn, tag string) {
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
					log.Printf("%s Connection: Ping failed: %v", tag, err)
					return
				}
				log.Printf("%s Connection: Sent ping heartbeat.", tag)
			}
		}
	}()
}

func (m *Manager) login(ctx context.Context, cfg Config) error {
	if m.privateConn == nil {
		return nil
	}
	if cfg.APIKey == "" || cfg.SecretKey == "" || cfg.Passphrase == "" {
		return nil
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	method := "GET"
	requestPath := "/users/self/verify"
	body := ""
	signature := GenerateSignature(cfg.SecretKey, timestamp, method, requestPath, body)

	loginParam := map[string]string{
		"apiKey":     cfg.APIKey,
		"passphrase": cfg.Passphrase,
		"timestamp":  timestamp,
		"sign":       signature,
	}
	loginReq := map[string]interface{}{
		"op":   "login",
		"args": []interface{}{loginParam},
	}
	payload, _ := json.Marshal(loginReq)
	if err := m.privateConn.WriteMessage(websocket.TextMessage, payload); err != nil {
		return err
	}
	_, message, err := m.privateConn.ReadMessage()
	if err != nil {
		return err
	}
	log.Printf("Login Response: %s", message)
	return nil
}

// --- 对外暴露的简洁API ---

func (m *Manager) PublicConn() *websocket.Conn {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.publicConn
}

func (m *Manager) PrivateConn() *websocket.Conn {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.privateConn
}

// SubscribePublic 仅演示；保持原公共订阅逻辑
func (m *Manager) SubscribePublic(args []map[string]string) error {
	c := m.PublicConn()
	if c == nil {
		return nil
	}
	req := map[string]interface{}{
		"op":   "subscribe",
		"args": []interface{}{},
	}
	for _, a := range args {
		req["args"] = append(req["args"].([]interface{}), a)
	}
	b, _ := json.Marshal(req)
	log.Println("Public Connection: Subscribe Request:", string(b))
	return c.WriteMessage(websocket.TextMessage, b)
}
