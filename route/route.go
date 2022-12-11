package route

import (
	// "net/http"

	ctr "github.com/ayo-ajayi/rest_api_template/controllers"
	"github.com/ayo-ajayi/rest_api_template/db"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func Router() *gin.Engine {
	database := db.DBinit()
	ctrl := ctr.NewController(database)

	router := gin.Default()
	CORSMiddleware := cors.Default
	
	router.Use(CORSMiddleware())
	router.GET("/favicon.ico", func(ctx *gin.Context) {
		ctx.File("./favicon.png")
	})
	choice := router.Group("/")
	{
		
		choiceByID := router.Group("/:id")
		{	choiceByID.Use(ctrl.CheckID)
			choiceByID.GET("/", ctrl.GetChoiceByIDCtr)
			choiceByID.PUT("/", ctrl.UpdateChoiceCtr)
			choiceByID.DELETE("/", ctrl.DeleteChoiceCtr)
		}
		choice.GET("/", ctrl.GetChoiceCtr)
		choice.POST("/", ctrl.PostChoiceCtr)
	}
	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{
			"ERROR": "NOT FOUND",
		})
	})
	return router
}

