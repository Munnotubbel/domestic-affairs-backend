package main

import (
	qrcode "github.com/skip2/go-qrcode"
)

// CreateQrCode generates a qrcode
func CreateQrCode(url string) ([]byte, error) {

	var err error
	var png []byte
	png, err = qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}

	return png, nil
}

// WriteQrCode write a file driectly
func WriteQrCode(url string, path string) error {

	err := qrcode.WriteFile(url, qrcode.Medium, 256, path)
	if err != nil {
		return err
	}
	return nil
}
