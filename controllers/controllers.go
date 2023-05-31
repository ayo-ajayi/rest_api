package controllers

import (
	//"encoding/json"

	"database/sql"
	"net/http"

	"github.com/ayo-ajayi/rest_api_template/model"
	"github.com/gin-gonic/gin"
)

type DBInterface interface {
	GetChoice() ([]model.Choice, error)
	CheckID(id string) (*model.Choice, error)
	DeleteChoice(id string) error
	PostChoice(newChoice *model.Choice) error
	UpdateChoice(updateChoice model.Choice) error
}

type Controller struct {
	DB DBInterface
}

func NewController(db DBInterface) *Controller {
	return &Controller{DB: db}
}
func (ctr *Controller) GetChoiceCtr(c *gin.Context) {
	choice, err := ctr.DB.GetChoice()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.JSON(200, choice)
}

func (ctr *Controller) CheckID(c *gin.Context) {
	id := c.Param("id")
	a, err := ctr.DB.CheckID(id)
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
	res := c.MustGet("values").(model.Choice)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.JSON(200, &res)
}

func (ctr *Controller) DeleteChoiceCtr(c *gin.Context) {
	res := c.MustGet("values").(model.Choice)
	if err := ctr.DB.DeleteChoice(res.ID); err != nil {
		c.JSON(500, err.Error())
		return
	}
	// b, _ :=json.Marshal(&res)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.JSON(200, gin.H{
		"message": "deleted successfully",
		"id":      res.ID,
	})
}

func (ctr *Controller) PostChoiceCtr(c *gin.Context) {
	newChoice := model.Choice{}
	if err := c.BindJSON(&newChoice); err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	if err := ctr.DB.PostChoice(&newChoice); err != nil {
		c.JSON(404, err)
		return
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.JSON(200, gin.H{
		"res": newChoice,
	})
}

func (ctr *Controller) UpdateChoiceCtr(c *gin.Context) {
	id := c.MustGet("values").(model.Choice).ID
	updateChoice := model.Choice{}
	if err := c.BindJSON(&updateChoice); err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	updateChoice.ID = id

	if err := ctr.DB.UpdateChoice(updateChoice); err != nil {
		c.JSON(400, err.Error())
		return
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.JSON(200, gin.H{"message": "success",
		"record": updateChoice})

}
