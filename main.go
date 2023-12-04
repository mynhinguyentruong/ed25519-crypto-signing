package main

import (
	"crypto/ed25519"
	"errors"
	"log"
	"os"
  "net/http"
  "encoding/base64"

	"golang.org/x/crypto/ssh"
	// "github.com/joho/godotenv"
  "github.com/gin-gonic/gin"

)

func mustNot(err error) {
    if err != nil {
        panic(err)
    }
}

func main() {
  // 1. Get private key

  myKey, err := LoadPrivateKeyFromPEMFile()

  // myKey, err := LoadPrivateKeyFromENV()
  mustNot(err)

  // 2. Get public key
  ed25519Pubkey, _ := GetPublicKeyFromPrivateKey(myKey) 

  // 3. Create a server
  r := gin.Default()
  r.GET("/ping", func(c *gin.Context) {
    signature := c.GetHeader("Magiclip-Signature") 
    if signature == "" {
      c.JSON(http.StatusBadRequest, gin.H{
        "error": "missing Magiclip-Signature Header",
      })
      return 
    }
    
    message := c.GetHeader("Magiclip-Message")
    if message == "" {
      c.JSON(http.StatusBadRequest, gin.H{
        "error": "missing Magiclip-Message Header",
      })
      return 
    }

    if isLegit := verifySignatureComingFromMagiclip(signature, message, ed25519Pubkey); isLegit {
      c.JSON(http.StatusOK, gin.H{
        "message": "pong",
      })

      return 
    } 
    
    c.JSON(http.StatusUnauthorized, gin.H{
        "message": "You're not from Magiclip, you're a scammer!",
    })

  })

  if err := r.Run(":4242"); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
  // listen and serve on 0.0.0.0:4242


}

func verifySignatureComingFromMagiclip(signature string, message string, ed25519Pubkey ed25519.PublicKey ) bool {
  // decode base64 into string then []byte
  sig, err := base64.StdEncoding.DecodeString(signature)
  if err != nil {
    log.Println("Error while decoding signature")
  }

  msg := []byte(message)

  if val := VerifySignature(ed25519Pubkey, sig, msg); !val {
    return false
  }

  return true

}

func LoadPrivateKeyFromENV() (ed25519.PrivateKey, error) {

  pemString := os.Getenv("PEM_PRIVATE_KEY")

  privateKey, err := ssh.ParseRawPrivateKey([]byte(pemString))
  if err != nil {
    return nil, errors.New("failed to parse raw private_key from pem string")
  }
  ed25519PrivateKey, ok := privateKey.(*ed25519.PrivateKey)
  if !ok {
    return nil, errors.New("failed to assert type into *ed25519.PrivateKey")
  }
  myKey := *ed25519PrivateKey

  return myKey, nil

}

func LoadPrivateKeyFromPEMFile() (ed25519.PrivateKey, error) {
  bytes, err := os.ReadFile("id_ed25519")
  mustNot(err)

  privateKey, err := ssh.ParseRawPrivateKey(bytes)
  mustNot(err)

  ed25519PrivateKey, ok := privateKey.(*ed25519.PrivateKey)
  if !ok {
    return nil, errors.New("failed to assert type into *ed25519.PrivateKey")
  }

  return *ed25519PrivateKey, nil
}

func GetPublicKeyFromPrivateKey(private_key ed25519.PrivateKey) (ed25519.PublicKey, error) {
  cryptoPublicKey := private_key.Public()

  // type assertion
  pubkey, ok := cryptoPublicKey.(ed25519.PublicKey)
  if !ok {
    return nil, errors.New("failed to assert into ed25519.PublicKey type")
  }

  return pubkey, nil
}

func VerifySignature(public_key ed25519.PublicKey, signature []byte, message []byte) bool {
  value := ed25519.Verify(public_key, message, signature)

  return value
}

