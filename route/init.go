package route

import (
	"github.com/gin-gonic/gin"
)

func InitRoute() *gin.Engine {
	r := gin.Default()
	//TODO:register your route here
	//for example:
	r.POST("/submit", HandleSubmit)
	r.GET("/leaderboard", HandleGetBoard)
	r.GET("/history/:name", HandleUserHistory)
	v := r.Group("/vote", CheckUserAgent)
	{
		v.POST("", HandleVote)
	}

	return r
}
