package lib

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"encoding/hex"
	"io"
	"math/rand"

	"sync"

	"github.com/astaxie/beego"
)

var tokenSecret cipher.Block
var tokenLock sync.Once

func lazyInit() {
	tokenLock.Do(func() {
		if secret, secretErr := hex.DecodeString(beego.AppConfig.String("TokenSecret")); secretErr != nil {
			panic(secretErr)
		} else {
			if tokenSecret, secretErr = aes.NewCipher(secret); secretErr != nil {
				panic(secretErr)
			}
		}
	})
}

func GenerateAccessToken(user uint64) string {
	lazyInit()

	var tokenBuf = bytes.NewBuffer(nil)
	binary.Write(tokenBuf, binary.BigEndian, rand.Int63())
	binary.Write(tokenBuf, binary.BigEndian, user)
	var encBuf = make([]byte, tokenSecret.BlockSize())
	tokenSecret.Encrypt(encBuf, tokenBuf.Bytes())
	tokenBuf.Write(encBuf)
	return hex.EncodeToString(tokenBuf.Bytes())
}

func ValidateAccessToken(token string) (uid uint64, valid bool) {
	lazyInit()

	var tokenData, tokenErr = hex.DecodeString(token)
	if tokenErr != nil {
		return
	}

	if len(tokenData) != 32 {
		return
	}
	var tokenBuf = bytes.NewBuffer(tokenData)
	var rnd int64
	binary.Read(tokenBuf, binary.BigEndian, &rnd)
	binary.Read(tokenBuf, binary.BigEndian, &uid)
	var encBuf = bytes.NewBuffer(nil)
	io.CopyN(encBuf, tokenBuf, int64(tokenSecret.BlockSize()))
	var decBuf = make([]byte, tokenSecret.BlockSize())
	tokenSecret.Decrypt(decBuf, encBuf.Bytes())
	valid = bytes.Compare(decBuf, tokenData[0:16]) == 0
	return
}
