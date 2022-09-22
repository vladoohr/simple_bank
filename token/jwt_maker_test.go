package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
	"github.com/vladoohr/simple_bank/util"
)

var (
	secretKey = util.RandomString(32)
)

func TestJWTToken(t *testing.T) {
	jwtMaker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, jwtMaker)

	username := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	tokenString, payload, err := jwtMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)
	require.NotEmpty(t, payload)

	payload, err = jwtMaker.VerifyToken(tokenString)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, payload.ExpireAt, expiredAt, time.Second)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	jwtMaker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, jwtMaker)

	username := util.RandomOwner()
	duration := -time.Minute

	tokenString, payload, err := jwtMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)
	require.NotEmpty(t, payload)

	payload, err = jwtMaker.VerifyToken(tokenString)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	username := util.RandomOwner()
	duration := -time.Minute

	payload, err := NewPayload(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	token := jwt.NewWithClaims(jwt.SigningMethodNone, payload)

	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

	jwtMaker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, jwtMaker)

	payload, err = jwtMaker.VerifyToken(tokenString)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
