package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/security"
)

func EncryptAndSignPayload(payloadStr []byte) (string, error) {
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
	base64Payload := security.Base64Encode(un)
	return base64Payload, nil
}

func GenerateEmailVerificationLink(payload *exports.EmailVerificationLinkPayload) (string, error) {
	url := ""
	if payload.ServerUrl != "" {
		url = payload.ServerUrl
	} else {
		url = config.GetServerURL()
	}
	payloadStr, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	linkPayload, err := EncryptAndSignPayload(payloadStr)
	return strings.Join([]string{strings.Join([]string{url, "api/v1/auth/verification/email/completion"}, "/"),
		strings.Join([]string{"token", linkPayload}, "=")}, "?"), err
}

func GenerateEmailVerificationLinkPayload(payload *exports.EmailVerificationLinkPayload) (string, error) {
	payloadStr, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	linkPayload, err := EncryptAndSignPayload(payloadStr)
	return linkPayload, err
}

func GeneratePasswordRecoveryLinkPayload(payload *exports.PaswordRecoveryPayload) (string, error) {
	payloadStr, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	linkPayload, err := EncryptAndSignPayload(payloadStr)
	return linkPayload, err
}

func ProcessPasswordRecoveryLink(str string) (*exports.PaswordRecoveryPayload, error) {

	originalMsg, err := UnencryptAndVerifyLink(str)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	payload := &exports.PaswordRecoveryPayload{}
	err = json.Unmarshal(originalMsg, payload)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return payload, nil
}

func GenerateTeamInviteLinkPayload(payload *exports.TeamInviteLinkPayload) (string, error) {
	payloadStr, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	linkPayload, err := EncryptAndSignPayload(payloadStr)
	return linkPayload, err
}

func ProcessTeamInviteLink(str string) (*exports.TeamInviteLinkPayload, error) {
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
	originalMsg, err := rsa.DecryptPKCS1v15(nil, key, originalMsgCipherText)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	payload := &exports.TeamInviteLinkPayload{}
	err = json.Unmarshal(originalMsg, payload)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	fmt.Println(payload)
	fmt.Println("here is the payload for the message")
	return payload, nil
}

func ProcessEmailVerificationLink(str string) (*exports.EmailVerificationLinkPayload, error) {

	originalMsg, err := UnencryptAndVerifyLink(str)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	payload := &exports.EmailVerificationLinkPayload{}
	err = json.Unmarshal(originalMsg, payload)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return payload, nil
}

func UnencryptAndVerifyLink(str string) ([]byte, error) {
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
	hashBytes := fromBase64
	signature := hashBytes[len(hashBytes)-signLen:]
	originalMsgCipherText := hashBytes[:len(hashBytes)-signLen]
	err = security.VerifyHash(originalMsgCipherText, signature)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	fmt.Println(key)
	originalMsg, err := rsa.DecryptPKCS1v15(nil, key, originalMsgCipherText)
	return originalMsg, err
}

func FindProjectRoot(marker string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		// Check if the marker file/folder exists in the current directory
		markerPath := filepath.Join(currentDir, marker)
		if _, err := os.Stat(markerPath); err == nil {
			return currentDir, nil // Found the marker, return the directory
		}

		// Move to the parent directory
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// Reached the root of the filesystem without finding the marker
			break
		}
		currentDir = parentDir
	}

	return "", fmt.Errorf("marker %s not found", marker)
}

func GenerateParticipantID(emails []string) (string, error) {
	slices.Sort[[]string](emails)
	joined := strings.Join(emails, ":")
	h := sha256.New()
	_, err := h.Write([]byte(joined))
	if err != nil {
		return "", err
	}
	hashByte := h.Sum(nil)
	hashedString := fmt.Sprintf("%x", hashByte)
	slicesOfHash := strings.Split(hashedString, "")
	prefixSlices := slicesOfHash[0:5]
	postFix := slicesOfHash[len(slicesOfHash)-5:]
	sub := strings.Join([]string{"PARTICIPANT_ID_", strings.Join(append(prefixSlices, postFix...), "")}, "")
	return sub, nil
}
