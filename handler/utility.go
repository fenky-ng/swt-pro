package handler

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"strings"
	"time"

	"github.com/fenky-ng/swt-pro/constant"
	"github.com/fenky-ng/swt-pro/generated"
	"github.com/fenky-ng/swt-pro/model"
	"github.com/fenky-ng/swt-pro/repository"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
)

var (
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
)

func getPrivateKey(filepath string) *rsa.PrivateKey {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(bytes)
	if err != nil {
		panic(err)
	}

	return key
}

func getPublicKey(filepath string) *rsa.PublicKey {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(bytes)
	if err != nil {
		panic(err)
	}

	return key
}

func getSignKey(customPrivateKeyPath string) *rsa.PrivateKey {
	if signKey == nil {
		privateKeyPath := constant.PrivateKeyPath
		if customPrivateKeyPath != "" {
			privateKeyPath = customPrivateKeyPath
		}
		signKey = getPrivateKey(privateKeyPath)
	}
	return signKey
}

func getVerifyKey(customPublicKeyPath string) *rsa.PublicKey {
	if verifyKey == nil {
		publicKeyPath := constant.PublicKeyPath
		if customPublicKeyPath != "" {
			publicKeyPath = customPublicKeyPath
		}
		verifyKey = getPublicKey(publicKeyPath)
	}
	return verifyKey
}

func hashAndSalt(input string) (salt string, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.MinCost)
	return string(hash), err
}

func comparePasswords(hashedPassword string, plainPassword string) bool {
	hashedPasswordBytes := []byte(hashedPassword)
	plainPasswordBytes := []byte(plainPassword)

	err := bcrypt.CompareHashAndPassword(hashedPasswordBytes, plainPasswordBytes)
	if err != nil {
		return false
	}

	return true
}

func generateJwtToken(user repository.User) (signedToken string, err error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))

	t.Claims = model.SessionClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    constant.ApplicationName,
			ExpiresAt: time.Now().Add(constant.LoginExpirationDuration).Unix(),
		},
		UserID:      user.ID,
		PhoneNumber: user.PhoneNumber,
	}

	return t.SignedString(getSignKey(""))
}

func getSessionClaims(ctx echo.Context) (sc model.SessionClaims, err error) {
	tokenString := ctx.Request().Header.Get("Authorization")
	tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")
	if tokenString == "" {
		err = errors.New("JWT token not found")
		return sc, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &model.SessionClaims{}, func(*jwt.Token) (interface{}, error) {
		return getVerifyKey(""), nil
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), jwt.ErrTokenExpired.Error()) {
			err = errors.New("Session is expired")
		} else {
			log.Errorf("ParseWithClaims error: %s", err.Error())
			err = errors.New("There was an error when parsing JWT")
		}
		return sc, err
	}

	if token.Claims.(*model.SessionClaims) == nil {
		err = errors.New("No session")
		return sc, err
	}

	sc = *token.Claims.(*model.SessionClaims)

	return sc, nil
}

func generateResponseHeader(errorCode int, errorMessages []string, successful bool) generated.ResponseHeader {
	var res generated.ResponseHeader
	if errorCode != 0 {
		res.ErrorCode = &errorCode
	}
	if len(errorMessages) != 0 {
		res.ErrorMessages = &errorMessages
	}
	if successful != false {
		res.Successful = &successful
	}
	return res
}
