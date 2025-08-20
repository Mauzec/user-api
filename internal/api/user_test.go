package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/mauzec/user-api/db/mock"
	db "github.com/mauzec/user-api/db/sqlc"
	"github.com/mauzec/user-api/internal/token"
	"github.com/mauzec/user-api/internal/util"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func randomUser() db.User {
	return db.User{
		ID:       util.RandomInt(1, 10000),
		Username: util.RandomString(8),
		FullName: util.RandomString(10),
		Sex:      "M",
		Age:      int32(util.RandomInt(18, 60)),
		Phone:    util.RandomPhone(),
		Email:    util.RandomEmail(),
	}
}

func assertBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	assert.NoError(t, err)
	var got userResponse
	err = json.Unmarshal(data, &got)
	assert.NoError(t, err)
	want := newUserResponse(user)
	assert.Equal(t, want, got)
}

func assertBodyNoUser(t *testing.T, body *bytes.Buffer) {
	data, err := io.ReadAll(body)
	assert.NoError(t, err)

	var user userResponse
	err = json.Unmarshal(data, &user)
	assert.NoError(t, err)
	assert.Zero(t, user)
}

func TestGetUserAPI(t *testing.T) {
	user1 := randomUser()

	testCases := []struct {
		name          string
		username      string
		setupAuth     func(t *testing.T, req *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			username: user1.Username,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, req, tokenMaker,
					authTypeBearer,
					user1.Username,
					time.Minute*15,
				)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user1.Username)).
					Times(1).
					Return(user1, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				assertBodyMatchUser(t, recorder.Body, user1)
			},
		},
		{
			name:     "NotFound",
			username: user1.Username,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, req, tokenMaker,
					authTypeBearer,
					user1.Username,
					time.Minute*15,
				)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user1.Username)).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
				assertBodyNoUser(t, recorder.Body)
			},
		},
		{
			name:     "InternalError",
			username: user1.Username,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthHeader(t, req, tokenMaker,
					authTypeBearer,
					user1.Username,
					time.Minute*15,
				)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user1.Username)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
				assertBodyNoUser(t, recorder.Body)
			},
		},
		{
			name:      "NoAuth",
			username:  user1.Username,
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user1.Username)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/users/%s", tc.username)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			assert.NoError(t, err)

			tc.setupAuth(t, req, server.tokenMaker)
			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}

}

// not implemented; like testgetuserapi
func TestUpdateUserAPI(t *testing.T) {

}
