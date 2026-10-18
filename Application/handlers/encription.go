package handlers

import (
	Provider "GameApplication/entity/entityProvider"
	"context"
	"fmt"
	"math/rand"
	"time"

	"gitlab.com/johnbass/phc"
	"golang.org/x/crypto/argon2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func getRandomSalt() []byte {
	var res []byte
	for i := 0; i <= 8; i++ {
		res = append(res, letterRunes[rand.Intn(len(letterRunes))])
	}
	return res
}

func getSaltFromPass(s string) []byte {

	digest, err := phc.Parse(s)
	if err != nil {
		panic(err)
	}

	fmt.Println("Salt:", digest.Salt)
	fmt.Println("Hash:", digest.Hash)

	return []byte(digest.Salt.String())
}

func encription(ctx context.Context, p Provider.Provider, s string, password string) string {
	var salt []byte
	if password == "" {
		salt = getRandomSalt()
	} else {
		salt = getSaltFromPass(password)
	}
	key := argon2.IDKey([]byte(s), salt, 1, 2*1024*1024, 4, 32)
	return string(key)
}
