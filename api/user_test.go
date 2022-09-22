package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/vladoohr/simple_bank/db/mock"
	db "github.com/vladoohr/simple_bank/db/sqlc"
	"github.com/vladoohr/simple_bank/util"
)

// eqUserParams holds user params need to compate the passwords and the CreateUserParams objects
type eqUserParams struct {
	arg      db.CreateUserParams
	password string
}

// EqUser returns a matcher that matches on equality for db.User object.
func eqUser(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqUserParams{arg, password}
}

// Matches checks that the hashed password and raw pass
func (e eqUserParams) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

// String describes what the matcher matches.
func (e eqUserParams) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func TestCreateUser(t *testing.T) {
	user, password := randomUser(t)

	// create new Controller
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	// create mock store
	store := mockdb.NewMockStore(ctrl)

	// build stubs
	createUserParams := db.CreateUserParams{
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
	}
	store.EXPECT().CreateUser(gomock.Any(), eqUser(createUserParams, password)).Return(user, nil).Times(1)

	// start test server and send request
	server := newTestServer(t, store)
	recorder := httptest.NewRecorder()

	body := gin.H{
		"username": user.Username,
		"password": password,
		"fullname": user.FullName,
		"email":    user.Email,
	}
	data, err := json.Marshal(body)
	require.NoError(t, err)

	url := "/users"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, req)

	// check response
	require.Equal(t, http.StatusCreated, recorder.Code)

}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPasssword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPasssword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	return
}
