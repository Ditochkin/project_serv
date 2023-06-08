package API

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"net/http"
	"time"

	"db_lab7/db"
	"db_lab7/types"

	"github.com/dgrijalva/jwt-go"
)

const (
	salt       = "asjhdjahsdjahsdas"
	signingKey = "%*FG67G%f786^G%&()(&J*H)(_I*K{76534d5D"
	tokenTTL   = 24 * 3600 * time.Second
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int64  `json:"user_id"`
	Role   string `json:"role"`
}

func (a *API) ParseToken(accessToken string) (int64, string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, "", errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, claims.Role, nil
}

func (a *API) getUserByUserNameAndPassword(username, password string) (int64, string, error) {
	rows, err := a.store.Query(db.GetUserQuery, username, generatePasswordHash(password))
	if err != nil {
		return 0, "", err
	}
	defer rows.Close()
	var user types.User
	isThereAnyRow := rows.Next()
	if !isThereAnyRow {
		rows.Close()
		return 0, "", errors.New("login or password is incorrect")
	}
	err = rows.Scan(&user.Id, &user.Name, &user.Username, &user.Password, &user.Role)
	return user.Id, user.Role, err
}

func (a *API) generateTokensByCred(username, password string) (string, error) {
	userID, role, err := a.getUserByUserNameAndPassword(username, password)
	if err != nil {
		return "", err
	}
	return a.generateTokens(userID, role)
}

func (a *API) generateTokens(userID int64, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userID,
		role,
	})

	ttk, err := token.SignedString([]byte(signingKey))
	return ttk, err
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func setTokenCookies(writer http.ResponseWriter, token string) {
	http.SetCookie(writer, &http.Cookie{
		Name:    "session_token",
		Value:   token,
		Expires: time.Now().Add(tokenTTL),
	})
}

func (a *API) GetIDAndRoleFromToken(writer http.ResponseWriter, request *http.Request) (int64, string, error) {
	ckc, err := request.Cookie("session_token")
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		fmt.Println("Im here! 1")
		return 0, "", err
	}
	if err == nil {
		userID, role, err := a.ParseToken(ckc.Value)
		if err == nil {
			fmt.Println(userID, role)
			return userID, role, nil
		}
	}
	fmt.Println("Im here! 2")

	return 0, "", errors.New("")
}
