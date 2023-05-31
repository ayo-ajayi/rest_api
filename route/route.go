package route

import (
	"net/http"

	ctr "github.com/ayo-ajayi/rest_api_template/controllers"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)




func SetupRouter(ctrl *ctr.Controller) (*gin.Engine, error) {
	router := gin.Default()
	router.GET("/favicon.ico", func(ctx *gin.Context) {
		ctx.File("./favicon.png")
	})
	CORSMiddleware := cors.Default()

	router.Use(CORSMiddleware)

	choiceByID := router.Group("/:id")
	{
		choiceByID.Use(ctrl.CheckID)
		choiceByID.GET("/", ctrl.GetChoiceByIDCtr)
		choiceByID.PUT("/", ctrl.UpdateChoiceCtr)
		choiceByID.DELETE("/", ctrl.DeleteChoiceCtr)
	}
	choiceAll := router.Group("/")
	{
		choiceAll.GET("", ctrl.GetChoiceCtr)
		choiceAll.POST("", ctrl.PostChoiceCtr)
	}

	router.NoRoute(func(ctx *gin.Context) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "NOT FOUND"})
	})
	return router, nil
}
