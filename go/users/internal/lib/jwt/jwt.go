package jwt

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/justcgh9/discord-clone-users/internal/models"
)

type Claims struct {
	UserID string
	Email  string
	AppID  string
}

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrExpired = errors.New("token expired")
	ErrInvalidClaims = errors.New("invalid claims")
	ErrInvalidExp = errors.New("missing or invalid exp")

)

func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string, secret string) (*Claims, error) {
	
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {		
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {		
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, ErrExpired
			}
		} else {
			return nil, ErrInvalidExp
		}

		uid, _ := claims["uid"].(string)
		email, _ := claims["email"].(string)
		appID, _ := claims["app_id"].(string)

		return &Claims{
			UserID: uid,
			Email:  email,
			AppID:  appID,
		}, nil
	}

	return nil, ErrInvalidClaims
}

func ExtractAppID(tokenString string) (int, error) {
	parser := jwt.NewParser()
	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return 0, err
	}

	claims := token.Claims.(jwt.MapClaims)
	appIdStr, ok := claims["app_id"].(string)
	if !ok {
		return 0, errors.New("invalid or missing app_id in token")
	}

	appId, err := strconv.Atoi(appIdStr)
	if err != nil {
		return  0, errors.New("invalid app_id")
	}

	return appId, nil
}
