package authn

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"
	"user-service/config"
	"user-service/errorx"
	"user-service/timesource"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// Hardcoded file names
	privateKeyFile = "private.pem"
	publicKeyFile  = "public.pem"
)

func ReadPrivatekey(filePath string) (*rsa.PrivateKey, error) {
	keyBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("ERROR: error reading private key file", err)
		return nil, err
	}
	keyBlock, _ := pem.Decode(keyBytes)
	if keyBlock == nil || keyBlock.Type != "PRIVATE KEY" {
		log.Println("ERROR: error decoding private key block")
		return nil, errors.New("error decoding private key block")
	}
	pKey, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		log.Println("ERROR: error parsing private key", err)
		return nil, err
	}
	if _, ok := pKey.(*rsa.PrivateKey); !ok {
		log.Println("ERROR: invalid key type found")
		return nil, errors.New("invalid key type found")
	}
	return pKey.(*rsa.PrivateKey), nil
}

func ReadPublickey(filePath string) (*rsa.PublicKey, error) {
	keyBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("ERROR: error reading public key file", err)
		return nil, err
	}
	keyBlock, _ := pem.Decode(keyBytes)
	if keyBlock == nil || keyBlock.Type != "PUBLIC KEY" {
		log.Println("ERROR: error decoding public key block")
		return nil, errors.New("error decoding public key block")
	}
	pKey, err := x509.ParsePKIXPublicKey(keyBlock.Bytes)
	if err != nil {
		log.Println("ERROR: error parsing private key", err)
		return nil, err
	}
	if _, ok := pKey.(*rsa.PublicKey); !ok {
		log.Println("ERROR: invalid key type found")
		return nil, errors.New("invalid key type found")
	}
	return pKey.(*rsa.PublicKey), nil
}

/*
GenerateRSASignedToken generates a jwt token, using RSA signing mechanism.
The token is valid for 30 min for the purpose of this sample service to offer enough duration for easy testing.
However, in production, it is advisable to have lower validity period, such as 10 mins.
*/
func GenerateRSASignedToken(cfg *config.Config, id, issuer string) (string, error) {
	if id == "" || issuer == "" {
		return "", errors.New("missing id or issuer")
	}
	// The user claims just contains an id.
	// Therefore, there is no need to encrypt the claims.
	userClaims := ClientClaims{
		UserID: id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":        id,
		"aud":        clientId,
		"iss":        issuer,
		"iat":        timesource.CurrentTime().Unix(),
		"exp":        timesource.CurrentTime().Add(time.Minute * 30).Unix(),
		"userClaims": userClaims,
	}, nil)

	pKeyFilePath := filepath.Join(cfg.KeyDir, privateKeyFile)
	privateKey, err := ReadPrivatekey(pKeyFilePath)
	if err != nil {
		log.Println("ERROR: error reading private key for signing", err)
		return "", err
	}
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		log.Println("ERROR: error during signing", err)
		return "", err
	}

	return signedToken, nil
}

func ValidateRSASignedToken(cfg *config.Config, token, issuer string) (string, error) {
	if token == "" || issuer == "" {
		return "", errors.New("missing token or issuer")
	}
	pKeyFilePath := filepath.Join(cfg.KeyDir, publicKeyFile)
	pubKey, err := ReadPublickey(pKeyFilePath)
	// Note that in production scenario, reading key file from disk for each validation
	// is not desired. There should be a caching mechanism, and a way to keep the cache in sync with
	// the state of the file on disk.
	// In the interest of time, I have not implemented such mechanism for this sample service.
	if err != nil {
		log.Println("ERROR: error reading public key for validation", err)
		return "", err
	}
	t, err := jwt.Parse(
		token,
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodRSA)
			if !ok {
				log.Println("DEBUG: invalid signing method")
				return nil, errorx.Error{Code: errorx.InvalidToken}
			}
			return pubKey, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}),
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
		jwt.WithAudience(clientId),
		jwt.WithIssuer(issuer),
		jwt.WithTimeFunc(func() time.Time { return timesource.CurrentTime() }),
	)

	if err != nil {
		log.Println("DEBUG: error during token parsing and validation", err)
		return "", errorx.Error{Code: errorx.InvalidToken}
	}
	if !t.Valid {
		log.Println("DEBUG: token not valid")
		return "", errorx.Error{Code: errorx.InvalidToken}
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("DEBUG: invalid claims")
		return "", errorx.Error{Code: errorx.InvalidToken}
	}

	id, ok := claims["userClaims"].(map[string]interface{})["user_id"].(string)
	if !ok {
		log.Println("DEBUG: invalid user_id in claims")
		return "", errorx.Error{Code: errorx.InvalidToken}
	}
	if claims["sub"] != id {
		log.Println("DEBUG: invalid sub in token claims")
		return "", errorx.Error{Code: errorx.InvalidToken}
	}

	return id, nil
}
