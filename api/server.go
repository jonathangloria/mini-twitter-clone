package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/jonathangloria/mini-twitter-clone/db/sqlc"
	"github.com/jonathangloria/mini-twitter-clone/token"
	"github.com/jonathangloria/mini-twitter-clone/util"
)

type Server struct {
	config     util.Config
	store      *db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store *db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users/login", server.loginUser)
	router.POST("/users", server.createUser)
	router.GET("/tweets/:id", server.getTweet)
	router.GET("/users/:id/tweets", server.listTweet)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.GET("/users/:id", server.getUser)
	authRoutes.POST("/followers", server.followUser)
	authRoutes.POST("/tweets", server.createTweet)
	authRoutes.GET("/users/:id/feed", server.getFeed)
	authRoutes.DELETE("/tweets/:id", server.deleteTweet)
	authRoutes.PATCH("/tweets/:id", server.updateTweet)
	// authRoutes.GET("/users/:id/followers", server.listFollower)
	// authRoutes.GET("/users/:id/following", server.listFollowing)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
