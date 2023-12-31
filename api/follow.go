package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/jonathangloria/mini-twitter-clone/db/sqlc"
	"github.com/jonathangloria/mini-twitter-clone/token"
	"github.com/lib/pq"
)

type createFollowerRequest struct {
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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if req.FollowedUser == authPayload.Username {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("cannot follow your own account")))
		return
	}

	follower, err := server.store.GetUserByUsername(ctx, authPayload.Username)
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
		Username:     authPayload.Username,
		FollowedID:   following.UserID,
		FollowedUser: req.FollowedUser,
	}

	ctx.JSON(http.StatusCreated, rsp)
}

type getFollowsRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) listFollower(ctx *gin.Context) {
	var req getFollowsRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetUser(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	followers, err := server.store.ListFollower(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, followers)
}

func (server *Server) listFollowing(ctx *gin.Context) {
	var req getFollowsRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetUser(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	following, err := server.store.ListFollowing(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, following)
}
