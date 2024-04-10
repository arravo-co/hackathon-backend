package security

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/dvsekhvalnov/jose2go/base64url"
)

func GenerateKeys() {
	var privateKey *rsa.PrivateKey
	privateKeyFile, err := os.ReadFile("private_key.pem")
	if err != nil {
		//panic(err)
		privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			panic(err)
		}
		privatePEM := pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		}
		privateKeyFile, err := os.Create("private_key.pem")
		if err != nil {
			panic(err)
		}
		err = pem.Encode(privateKeyFile, &privatePEM)
		if err != nil {
			panic(err)
		}
		err = privateKeyFile.Close()
		if err != nil {
			panic(err)
		}
	} else {
		privateKeyBlock, _ := pem.Decode(privateKeyFile)
		privateKey, err = x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
		if err != nil {
			fmt.Println("Error parsing private key:", err)
			return
		}
	}

	// Extract the public key from the private key
	publicKey := &privateKey.PublicKey

	// Encode the public key to the PEM format
	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(publicKey),
	}
	publicKeyFile, err := os.Create("public_key.pem")
	if err != nil {
		fmt.Println("Error creating public key file:", err)
		os.Exit(1)
	}
	pem.Encode(publicKeyFile, publicKeyPEM)
	publicKeyFile.Close()

	fmt.Println("RSA key pair generated successfully!")
}

func ReadPrivateKey() (*rsa.PrivateKey, error) {
	// Load the private key from the file
	privateKeyFile, err := os.ReadFile("private_key.pem")
	if err != nil {
		fmt.Println("Error loading private key file:", err)
		os.Exit(1)
	}
	privateKeyBlock, _ := pem.Decode(privateKeyFile)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		fmt.Println("Error parsing private key:", err)
		return nil, err
	}
	return privateKey, nil
}

func Base64Encode(str []byte) string {
	return base64url.Encode(str)
}

func Base64UrlDecode(str string) ([]byte, error) {
	return base64url.Decode(str)
	/*
		data, err := base64.StdEncoding.DecodeString(str)
		if err != nil {
			fmt.Println("error:", err)
			return nil, err
		}
		return data, nil
	*/
}

func VerifyHash(byt []byte, signature []byte) error {
	key, err := ReadPrivateKey()
	if err != nil {
		fmt.Println(err)
		return err
	}
	hash := sha256.Sum256(byt)
	err = rsa.VerifyPKCS1v15(&key.PublicKey, crypto.SHA256, hash[:], signature)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func Sign(hash []byte) ([]byte, error) {
	key, err := ReadPrivateKey()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	hashed := sha256.Sum256(hash)
	//rng := rand.Reader
	sigBytes, err := rsa.SignPKCS1v15(nil, key, crypto.SHA256, hashed[:])
	fmt.Println("kkkkbjjgyfytfftttttttttttttttttttttttttttttttttdrtdrdrtseee")
	return sigBytes, err
}
