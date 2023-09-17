package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/jonathangloria/mini-twitter-clone/db/sqlc"
	"github.com/jonathangloria/mini-twitter-clone/token"
)

type createTweetRequest struct {
	Body string `json:"body" binding:"required,max=200"`
}

type tweetResponse struct {
	TweetID   int64     `json:"tweet_id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

func newTweetResponse(tweet db.Tweet, user db.User) tweetResponse {
	return tweetResponse{
		TweetID:   tweet.ID,
		UserID:    user.ID,
		Username:  user.Username,
		Body:      tweet.Body,
		CreatedAt: tweet.CreatedAt,
	}
}

func (server *Server) createTweet(ctx *gin.Context) {
	var req createTweetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	user, err := server.store.GetUserByUsername(ctx, authPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateTweetParams{
		UserID: user.ID,
		Body:   req.Body,
	}

	tweet, err := server.store.CreateTweet(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newTweetResponse(tweet, user)

	ctx.JSON(http.StatusCreated, rsp)
}

type getTweetRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getTweet(ctx *gin.Context) {
	var req getTweetRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tweet, err := server.store.GetTweet(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, tweet)
}

func (server *Server) deleteTweet(ctx *gin.Context) {
	var req getTweetRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tweet, err := server.store.GetTweet(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, tweet.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if user.Username != authPayload.Username {
		err := errors.New("the account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = server.store.DeleteTweet(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusNoContent, err)
}

type getFeedQuery struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
}

func (server *Server) getFeed(ctx *gin.Context) {
	var req getTweetRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var query getFeedQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if user.Username != authPayload.Username {
		err := errors.New("the account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	feed, err := server.store.GetFeed(ctx, db.GetFeedParams{
		FollowerID: req.ID,
		Offset:     (query.PageID - 1) * 10,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, feed)
}
