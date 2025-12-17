package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// a small util to generate an api secret
func main() {
	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	out := base64.StdEncoding.EncodeToString(buf)
	fmt.Println(out)
}
