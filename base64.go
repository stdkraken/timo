package main

import (
	"encoding/base64"
	"strings"
)

func IsBase64PNG(URL string) bool {
	return strings.HasPrefix(URL, "data:image/png;base64,")
}

func DecodeBase64(Code string) []byte {
	dec, err := base64.StdEncoding.DecodeString(strings.Split(Code, ",")[1])
	if err != nil {
		panic(err)
	}
	return dec
}
