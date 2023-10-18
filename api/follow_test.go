package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/jonathangloria/mini-twitter-clone/db/mock"
	db "github.com/jonathangloria/mini-twitter-clone/db/sqlc"
	"github.com/jonathangloria/mini-twitter-clone/token"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestFollowUserAPI(t *testing.T) {
	user, _ := randomUser(t)
	user2, _ := randomUser(t)
	follow := randomFollowing(user, user2)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"followed_user": user2.Username,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).Return(user, nil)
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(user2.Username)).
					Times(1).Return(user2, nil)
				arg := db.CreateFollowingParams{
					UserID:     user2.ID,
					FollowerID: user.ID,
				}
				store.EXPECT().CreateFollowing(gomock.Any(), gomock.Eq(arg)).
					Times(1).Return(follow, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchFollow(t, recorder.Body, follow)
			},
		},
		{
			name: "FollowTheSameAccountTwice",
			body: gin.H{
				"followed_user": user2.Username,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).Return(user, nil)
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(user2.Username)).
					Times(1).Return(user2, nil)
				arg := db.CreateFollowingParams{
					UserID:     user2.ID,
					FollowerID: user.ID,
				}
				store.EXPECT().CreateFollowing(gomock.Any(), gomock.Eq(arg)).
					Times(1).Return(db.Follow{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "UserNotFound",
			body: gin.H{
				"followed_user": "usernotfound",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).Return(user, nil)
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq("usernotfound")).
					Times(1).Return(db.User{}, sql.ErrNoRows)
				store.EXPECT().CreateFollowing(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"followed_user": user2.Username,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).Return(user, nil)
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Eq(user2.Username)).
					Times(1).Return(user2, nil)
				store.EXPECT().CreateFollowing(gomock.Any(), gomock.Any()).
					Times(1).Return(db.Follow{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidUserID",
			body: gin.H{
				"followed_user": "*()&>>1w_",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateFollowing(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "FollowOwnAccount",
			body: gin.H{
				"followed_user": user.Username,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateFollowing(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/followers"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomFollowing(user db.User, followedUser db.User) db.Follow {
	follow := db.Follow{
		UserID:     followedUser.ID,
		FollowerID: user.ID,
	}
	return follow
}

func requireBodyMatchFollow(t *testing.T, body *bytes.Buffer, follow db.Follow) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotFollow followerResponse
	fmt.Println(gotFollow)
	err = json.Unmarshal(data, &gotFollow)

	require.NoError(t, err)
	require.Equal(t, follow.UserID, gotFollow.FollowedID)
	require.Equal(t, follow.FollowerID, gotFollow.UserID)
}
