package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"

	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/security"
)

func GenerateLinkPayload(payload *exports.LinkPayload) (string, error) {
	payloadStr, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	key, err := security.ReadPrivateKey()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	rng := rand.Reader
	cipherText, err := rsa.EncryptPKCS1v15(rng, &key.PublicKey, payloadStr)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	sig, err := security.Sign(cipherText)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	un := append(cipherText, sig...)
	linkPayload := security.Base64Encode(un)
	return linkPayload, nil
}

func ProcessInviteLink(str string) (*exports.LinkPayload, error) {
	fromBase64, err := security.Base64UrlDecode(str)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	key, err := security.ReadPrivateKey()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	signLen := key.PublicKey.Size()
	fmt.Println(signLen)
	hashBytes := fromBase64
	signature := hashBytes[len(hashBytes)-signLen:]
	originalMsgCipherText := hashBytes[:len(hashBytes)-signLen]
	err = security.VerifyHash(hashBytes, signature)
	if err != nil {
		fmt.Println(err.Error())
		//return nil, err
	}
	fmt.Println("valid")
	fmt.Println(originalMsgCipherText)
	originalMsg, err := rsa.DecryptPKCS1v15(nil, key, originalMsgCipherText)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	payload := &exports.LinkPayload{}
	err = json.Unmarshal(originalMsg, payload)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	fmt.Println(payload)
	return payload, nil
}
