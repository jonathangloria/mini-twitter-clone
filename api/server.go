package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/jonathangloria/mini-twitter-clone/db/sqlc"
	"github.com/jonathangloria/mini-twitter-clone/util"
)

type Server struct {
	config util.Config
	store  *db.Store
	router *gin.Engine
}

func NewServer(config util.Config, store *db.Store) *Server {
	server := &Server{
		config: config,
		store:  store,
	}
	router := gin.Default()

	router.POST("/users", server.createUser)
	// router.POST("/follow", server.followUser)
	// router.POST("/tweets", server.postTweet)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
