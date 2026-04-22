package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func SignUser(id int, uuid, username, secret string) string {
	payload := fmt.Sprintf("%d|%s|%s", id, uuid, username)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}
