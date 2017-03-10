package main

import (
	"crypto/aes"
	log "github.com/thinkboy/log4go"
	myaes "goim/libs/crypto/aes"
	//	"goim/libs/define"
	//	"strconv"
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

	key := []byte("1234567887654321")

	cipherText := []byte(token)
	cipher, err := aes.NewCipher(key)

	if err != nil {
		log.Error("failed to create cipher", err)
	}

	plainText, err := myaes.ECBDecrypt(cipher, cipherText)
	if err != nil {
		log.Error("decrypt failed", err)
	}
	originText := unpadding(plainText)
	// originText is the token
	userId = queryUser(originText)
	// debug only
	log.Debug("authenticaed userId %d", userId)
	roomId = 1
	/*
		if userId, err = strconv.ParseInt(token, 10, 64); err != nil {
			userId = 0
			roomId = define.NoRoom
		} else {
			roomId = 1 // only for debug
		}
	*/
	return
}

func unpadding(src []byte) (dst []byte) {
	var padLen, srcLen int
	srcLen = len(src)

	for padLen = 0; padLen < srcLen && src[srcLen-padLen-1] == 0; padLen++ {
	}
	return src[:srcLen-padLen]
}
