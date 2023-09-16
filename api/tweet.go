package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/jonathangloria/mini-twitter-clone/db/sqlc"
)

type createTweetRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Body     string `json:"body" binding:"required,max=200"`
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

	user, err := server.store.GetUserByUsername(ctx, req.Username)
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

	ctx.JSON(http.StatusOK, rsp)
}

type getFeedRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type getFeedQuery struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
}

func (server *Server) getFeed(ctx *gin.Context) {
	var req getFeedRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var query getFeedQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
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
