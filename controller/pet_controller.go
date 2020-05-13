package controller

import (
	"fmt"
	"go-restful-xorm/models"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

type PetController struct {
	DB *xorm.Engine
}

func (p *PetController) Index(c *gin.Context) {
	var pets []models.Pet
	err := p.DB.Table("pets").Find(&pets)
	if err != nil {
		fmt.Printf("%v", err)
	}

	if len(pets) <= 0 {
		c.JSON(404, gin.H{"status": 404, "message": "not found."})
		return
	}
	c.JSON(200, gin.H{"status": 200, "data": pets})
}

func (p *PetController) Create(c *gin.Context) {

}

func (p *PetController) Show(c *gin.Context) {

}

func (p *PetController) Update(c *gin.Context) {

}

func (p *PetController) UploadImage(c *gin.Context) {

}
