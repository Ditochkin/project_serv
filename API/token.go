package API

import (
	"crypto/sha1"
	"db_lab7/db"
	"db_lab7/types"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	salt            = "asjhdjahsdjahsdas"
	signingKey      = "%*FG67G%f786^G%&()(&J*H)(_I*K{76534d5D"
	tokenTTL        = 5 * time.Second
	refreshTokenTTL = 5 * 24 * time.Hour
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

func (a *API) ParseRefreshToken(refreshToken string) (*types.User, error) {
	rows, err := a.store.Query(db.GetUserByRefreshTokenQuery, refreshToken)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var user types.User
	isThereAnyRow := rows.Next()
	if !isThereAnyRow {
		return nil, errors.New("no such refresh token")
	}
	err = rows.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.RefreshToken, &user.RefreshTokenEAT, &user.Role)
	if err != nil {
		return nil, err
	}

	if time.Now().After(time.Unix(user.RefreshTokenEAT.Int64, 0)) {
		return nil, errors.New("refresh token expired")
	}
	return &user, nil
}

func (a *API) getUserByUserNameAndPassword(username, password string) (int64, string, error) {
	fmt.Println(username, generatePasswordHash(password))
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
	err = rows.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.RefreshToken, &user.RefreshTokenEAT, &user.Role)
	fmt.Println(user.Id, user.Username, user.Password, user.Email)
	return user.Id, user.Role, err
}

func (a *API) generateTokensByCred(username, password string) (string, string, error) {
	userID, role, err := a.getUserByUserNameAndPassword(username, password)
	if err != nil {
		return "", "", err
	}
	return a.generateTokens(userID, role)
}

func (a *API) generateTokens(userID int64, role string) (string, string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userID,
		role,
	})

	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}
	_, err = a.store.Exec(db.UpdateRefreshQuery, refreshToken, time.Now().Add(refreshTokenTTL).Unix(), userID)
	if err != nil {
		return "", "", err
	}
	ttk, err := token.SignedString([]byte(signingKey))
	return ttk, refreshToken, err
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func newRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
