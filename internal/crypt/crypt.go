package crypt

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	gocrypt "github.com/amoghe/go-crypt"
)

type HashOptions struct {
	Salt      string
	Algorithm string
}

func Hash(plaintext string, opt HashOptions) (string, error) {
	switch opt.Algorithm {
	// default to sha512 if no algorithm is specified
	case "":
		fallthrough
	case "sha512":
		if opt.Salt == "" {
			// generate 16 bytes of random data for the salt
			saltBytes := make([]byte, 16)

			_, err := rand.Read(saltBytes)
			if err != nil {
				return "", err
			}

			opt.Salt = base64.StdEncoding.EncodeToString(saltBytes)
		}

		pw, err := gocrypt.Crypt(plaintext, "$6$"+opt.Salt)
		if err != nil {
			return "", err
		}

		return pw, nil
	default:
		return "", fmt.Errorf("unsupported algorithm: %s", opt.Algorithm)
	}
}
