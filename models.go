package main
import SexDB "github.com/Plankiton/SexPistol/Cartridge"

type Category struct {
	SexDB.Model
	Name       string `json:"name,omitempty" gorm:"column:descricao"`
	Meals      []Meal `json:"meals" gorm:"foreignKey:CatID"`
}

type Meal struct {
	SexDB.Model
	Name       string  `json:"Name,omitempty" gorm:"column:descricao"`
	Desc   string  `json:"Desc,omitempty" gorm:"column:descricao_detalhada"`
	Price      float64 `json:"Price,omitempty" gorm:"column:preco_venda"`

	CatID      uint    `json:"-" gorm:"column:id_categoria"`
}

func(*Category) TableName() string{
	return "categoria"
}
func(*Meal) TableName() string{
	return "produtos"
}
