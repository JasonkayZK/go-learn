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
	var pet models.Pet
	err := c.Bind(pet)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	insert, err := p.DB.Table("pets").Insert(&pet)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	if insert != 1 {
		c.JSON(404, gin.H{"error": "Fail to insert, maybe exist?"})
	}

	c.JSON(200, gin.H{"status": 200, "message": "pet item created."})
}

func (p *PetController) Show(c *gin.Context) {
	var pet models.Pet
	has, err := p.DB.Table("pets").Where("id = ?", c.Param("id")).Get(&pet)
	if err != nil {
		c.JSON(404, gin.H {"status": 404, "message": "pet select error"})
		return
	}
	if !has {
		c.JSON(404, gin.H {"status": 404, "message": "pet not found"})
		return
	}

	if pet.Id == 0 {
		c.JSON(404, gin.H {"status": 404, "message": "pet not found"})
		return
	}

	c.JSON(200, gin.H {"status": 200, "data": pet})
}

func (p *PetController) Update(c *gin.Context) {

}

func (p *PetController) UploadImage(c *gin.Context) {

}
