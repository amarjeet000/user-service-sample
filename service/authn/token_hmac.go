package authn

import (
	"errors"
	"log"
	"time"
	"user-service/errorx"
	"user-service/timesource"

	"github.com/golang-jwt/jwt/v5"
)

/*
GenerateHMACSignedToken generates a jwt token, using HMAC signing mechanism.
The token is valid for 30 min for the purpose of this sample service to offer enough duration for easy testing.
However, in production, it is advisable to have lower validity period, such as 10 mins.
*/
func GenerateHMACSignedToken(id, issuer, secret string) (string, error) {
	if id == "" || issuer == "" || secret == "" {
		return "", errors.New("missing id, issuer, or secret")
	}
	// The user claims just contains an id.
	// Therefore, there is no need to encrypt the claims.
	userClaims := ClientClaims{
		UserID: id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":        id,
		"aud":        clientId,
		"iss":        issuer,
		"iat":        timesource.CurrentTime().Unix(),
		"exp":        timesource.CurrentTime().Add(time.Minute * 30).Unix(),
		"userClaims": userClaims,
	}, nil)

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Println("ERROR: error during signing", err)
		return "", err
	}

	return signedToken, nil
}

func ValidateHMACSignedToken(token, issuer, secret string) (string, error) {
	if token == "" || issuer == "" || secret == "" {
		return "", errors.New("missing token, issuer, or secret")
	}
	t, err := jwt.Parse(
		token,
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				log.Println("DEBUG: invalid signing method")
				return nil, errorx.Error{Code: errorx.InvalidToken}
			}
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
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
