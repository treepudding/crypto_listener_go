package ws

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
)

func GenerateSignature(secretKey, timestamp, method, requestPath, body string) string {
	message := timestamp + method + requestPath + body
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func defaultHeaders(simulatedTrading bool) http.Header {
	h := http.Header{}
	if simulatedTrading {
		h.Set("x-simulated-trading", "1")
	}
	return h
}
