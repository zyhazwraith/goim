package main

import (
	"crypto/aes"
	"goim/libs/crypto/aes"
	"goim/libs/define"
	"strconv"

	log "github.com/thinkboy/log4go"
)

// TODO Add AES-ECB padding and mysql connection part
//
// developer could implement "Auth" interface for decide how get userId, or roomId
type Auther interface {
	Auth(token string) (userId int64, roomId int32)
}

type DefaultAuther struct {
}

func NewDefaultAuther() *DefaultAuther {
	return &DefaultAuther{}
}

// TODO token is the ciphertext, decrypt cipher and then unpadding it,
//	then we have the desired string
func (a *DefaultAuther) Auth(token string) (userId int64, roomId int32) {
	var err error
	cipherText := token

	cipher, err := aes.NewCipher(key)

	if err != nil {
		log.Error("", err)
	}

	plainText := aes.ECBDecrypt(cipher, cipherText)
	if userId, err = strconv.ParseInt(token, 10, 64); err != nil {
		userId = 0
		roomId = define.NoRoom
	} else {
		roomId = 1 // only for debug
	}
	return
}

func unpadding(src []byte) (dst []byte) {
	var len, i int
}
