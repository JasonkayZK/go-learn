package dao

import (
	"fmt"
	"go-mysql-server-demo/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type PetDAO struct {
	DB *xorm.Engine
}

func (p *PetDAO) CreatePet(pet *models.Pet) error {
	insert, err := p.DB.Table("pets").Insert(pet)
	if err != nil {
		return err
	}
	if insert != 1 {
		return fmt.Errorf("error, fail to insert, maybe exist")
	}

	return nil
}

func (p *PetDAO) FindPetById(id int) (*models.Pet, error) {
	var pet models.Pet
	has, err := p.DB.Table("pets").Where("id = ?", id).Get(&pet)
	if err != nil {
		return nil, err
	}
	if !has || pet.Id == 0 {
		return nil, fmt.Errorf("pet not found")
	}

	return &pet, nil
}

func (p *PetDAO) Update(petId, petAge int, petName string) error {
	res, err := p.DB.Exec("UPDATE `pets` SET `name` = ?, `age` = ? WHERE `id` = ?", petName, petAge, petId)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected != 1 {
		return fmt.Errorf("fail to update, maybe record not exist")
	}

	return nil
}

func (p *PetDAO) DeleteById(petId int) error {
	res, err := p.DB.Exec("DELETE FROM `pets` WHERE `id` = ?", petId)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected != 1 {
		return fmt.Errorf("fail to delete, maybe record not exist")
	}

	return nil
}
