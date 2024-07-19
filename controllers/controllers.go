package controllers

import (
	//"encoding/json"

	"database/sql"
	"net/http"

	"github.com/ayo-ajayi/rest_api_template/db"
	"github.com/ayo-ajayi/rest_api_template/model"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	DB db.DBRepo
}

func NewController(repo db.DBRepo) *Controller {
	return &Controller{DB: repo}
}
func (ctr *Controller) GetChoiceCtr(c *gin.Context) {
	choice, err := ctr.DB.GetChoice(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"choices": choice,
	})
}

func (ctr *Controller) CheckID(c *gin.Context) {
	id := c.Param("id")
	a, err := ctr.DB.CheckID(c, id)
	switch {
	case err == sql.ErrNoRows:
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	case err != nil:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Set("values", *a)
	c.Next()
}

func (ctr *Controller) GetChoiceByIDCtr(c *gin.Context) {
	res, exists := c.Get("values")
	if !exists {
		c.JSON(500, gin.H{"error": "no value found"})
		return
	}
	choice := res.(model.Choice)
	c.JSON(200,
		gin.H{
			"message": "success",
			"choice":  choice})
}

func (ctr *Controller) DeleteChoiceCtr(c *gin.Context) {
	res, exists := c.Get("values")
	if !exists {
		c.JSON(500, gin.H{"error": "no value found"})
		return
	}
	choice := res.(model.Choice)
	if err := ctr.DB.DeleteChoice(c, choice.ID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"message": "deleted successfully",
		"id":      choice.ID,
	})
}

func (ctr *Controller) PostChoiceCtr(c *gin.Context) {
	newChoice := model.Choice{}
	if err := c.ShouldBindJSON(&newChoice); err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	if err := ctr.DB.PostChoice(c, &newChoice); err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"choice":  newChoice,
	})
}

func (ctr *Controller) UpdateChoiceCtr(c *gin.Context) {
	res, exists := c.Get("values")
	if !exists {
		c.JSON(500, gin.H{"error": "no value found"})
		return
	}
	id := res.(model.Choice).ID
	updateChoice := model.Choice{}
	if err := c.ShouldBindJSON(&updateChoice); err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	updateChoice.ID = id

	if err := ctr.DB.UpdateChoice(c, updateChoice); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "success",
		"choice": updateChoice,
	})

}
