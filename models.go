package main

import SexDB "github.com/Plankiton/SexPistol/Cartridge"

type Category struct {
	SexDB.Model
	Index int    `json:"index" gorm:"-"`
	ID    uint   `json:"db_id"`
	Len   int64  `json:"meal_count"`
	Name  string `json:"name,omitempty" gorm:"column:descricao"`
	Meals []Meal `json:"meals" gorm:"foreignKey:CatID"`
}

type Meal struct {
	SexDB.Model
	ID    uint    `json:"db_id"`
	Name  string  `json:"Name,omitempty" gorm:"column:descricao"`
	Desc  string  `json:"Desc,omitempty" gorm:"column:descricao_detalhada"`
	Price float64 `json:"Price,omitempty" gorm:"column:preco_venda"`

	CatID uint     `json:"-" gorm:"column:id_categoria"`
	Cat   Category `json:"cat" gorm:"-"`
}

func (*Category) TableName() string {
	return "categoria"
}
func (*Meal) TableName() string {
	return "produtos"
}
