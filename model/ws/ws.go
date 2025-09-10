package ws

type LoginRequest struct {
	Op   string        `json:"op"`
	Args []LoginParams `json:"args"`
}

type LoginParams struct {
	APIKey     string `json:"apiKey"`
	Passphrase string `json:"passphrase"`
	Timestamp  string `json:"timestamp"`
	Sign       string `json:"sign"`
	Chanel     string `json:"channel"`
	InstType   string `json:"instType"`
}
