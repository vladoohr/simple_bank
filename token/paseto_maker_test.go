package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vladoohr/simple_bank/util"
)

func TestPasetoToken(t *testing.T) {
	pasetoMaker, err := NewPasetoMaker([]byte(util.RandomString(32)))
	require.NoError(t, err)
	require.NotEmpty(t, pasetoMaker)

	username := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	tokenString, payload, err := pasetoMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)
	require.NotEmpty(t, payload)

	payload, err = pasetoMaker.VerifyToken(tokenString)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, payload.ExpireAt, expiredAt, time.Second)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	pasetoMaker, err := NewPasetoMaker([]byte(util.RandomString(32)))
	require.NoError(t, err)
	require.NotEmpty(t, pasetoMaker)

	username := util.RandomOwner()
	duration := -time.Minute

	tokenString, payload, err := pasetoMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)
	require.NotEmpty(t, payload)

	payload, err = pasetoMaker.VerifyToken(tokenString)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
}
