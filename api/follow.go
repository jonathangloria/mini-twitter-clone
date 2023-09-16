package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/jonathangloria/mini-twitter-clone/db/sqlc"
	"github.com/lib/pq"
)

type createFollowerRequest struct {
	Username     string `json:"username" binding:"required,alphanum"`
	FollowedUser string `json:"followed_user" binding:"required,alphanum"`
}

type followerResponse struct {
	UserID       int64  `json:"user_id"`
	Username     string `json:"username"`
	FollowedID   int64  `json:"followed_id"`
	FollowedUser string `json:"followed_user"`
}

func (server *Server) followUser(ctx *gin.Context) {
	var req createFollowerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.FollowedUser == req.Username {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("cannot follow your own account")))
		return
	}

	follower, err := server.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByUsername(ctx, req.FollowedUser)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateFollowingParams{
		UserID:     user.ID,
		FollowerID: follower.ID,
	}

	following, err := server.store.CreateFollowing(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := followerResponse{
		UserID:       following.FollowerID,
		Username:     req.Username,
		FollowedID:   following.UserID,
		FollowedUser: req.FollowedUser,
	}

	ctx.JSON(http.StatusOK, rsp)
}
