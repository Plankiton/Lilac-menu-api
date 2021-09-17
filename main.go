package main

import (
	Http "net/http"
	StrConv "strconv"
	Str "strings"

	Sex "github.com/Plankiton/SexPistol"
	SexDB "github.com/Plankiton/SexPistol/Cartridge"
	mysql "gorm.io/driver/mysql"
)

func main() {
	driver, uri := mysql.Open, Sex.GetEnv("PREAMAR_DATABASE_URL", "test.db")
	if Sex.GetEnv("SEX_DEBUG", "false") != "false" {
		driver, uri = SexDB.Sqlite, "test.db"
	}

	db, err := SexDB.Open(uri, driver)
	db.SetLogLevel("info")
	if err != nil {
		Sex.Die(err)
	}

	db.AddModels(
		new(Category),
		new(Meal),
	)

	pistol := Sex.NewPistol().
		Add("/cats", func(r Sex.Request) Sex.Json {
			cats := []Category{}
			if err := db.Find(&cats).Error; err != nil {
				return nil
			}

			for i, cat := range cats {
				cats[i].Name = Cap(cat.Name)
				cats[i].Index = i
			}

			return cats
		}).
		Add(`/cat/{id}/meals`, func(r Sex.Request) Sex.Json {
			page, err := StrConv.Atoi(r.URL.Query().Get("page"))
			if err != nil {
				page = 1
			}

			limit, err := StrConv.Atoi(r.URL.Query().Get("limit"))
			if err != nil {
				limit = 10
			}

			catId, err := StrConv.Atoi(r.PathVars["id"])
			if err != nil {
				catId = 0
			}

			cat := Category{}
			if db.First(&cat, "id = ?", catId).Error != nil {
				return nil
			}

			db.
				Where("id_categoria = ?", cat.ID).
				Joins("join categoria cat on cat.ID = id_categoria").
				Offset((page - 1) * limit).
				Limit(limit).
				Find(&cat.Meals)

			db.Model(Meal{}).
				Where("id_categoria = ?", cat.ID).
				Joins("join categoria cat on cat.ID = id_categoria").
				Count(&cat.Len)

			cat.Name = Cap(cat.Name)
			for m, meal := range cat.Meals {
				cat.Meals[m].Name = Cap(meal.Name)
				cat.Meals[m].Desc = Cap(meal.Desc)
			}

			return cat
		}).
		Add(`/meals`, func(r Sex.Request) Sex.Json {
			page, err := StrConv.Atoi(r.URL.Query().Get("page"))
			if err != nil {
				page = 1
			}

			limit, err := StrConv.Atoi(r.URL.Query().Get("limit"))
			if err != nil {
				limit = 10
			}

			if query := r.URL.Query().Get("query"); query != "" {
				query = "%" + query + "%"

				meals := []Meal{}
				if err := db.
					Where("produtos.descricao like ? or produtos.descricao_detalhada like ?", query, query).
					Find(&meals).
					Error; err != nil {
					return Sex.Bullet{
						Message: err.Error(),
					}
				}

				for i, meal := range meals {
					if err := db.
						Joins("join produtos meal on meal.id_categoria = categoria.id").
						First(&meal.Cat, "categoria.id = ?", meal.CatID).
						Error; err != nil {
						return Sex.Bullet{
							Message: err.Error(),
						}
					}

					meals[i] = meal
				}

				return meals
			}

			cats := []*Category{}
			if err := db.
				Find(&cats).
				Error; err != nil {
				return Sex.Bullet{
					Message: err.Error(),
				}
			}

			for i, cat := range cats {
				db.
					Joins("join categoria cat on cat.ID = id_categoria and cat.ID = ?", cat.ID).
					Offset((page - 1) * limit).
					Limit(limit).
					Find(&cat.Meals)
				cat.Name = Cap(cat.Name)

				db.Model(Meal{}).
					Where("id_categoria = ?", cat.ID).
					Joins("join categoria cat on cat.ID = id_categoria").
					Count(&cat.Len)

				for m, meal := range cat.Meals {
					cat.Meals[m].Name = Cap(meal.Name)
					cat.Meals[m].Desc = Cap(meal.Desc)
				}

				cats[i].Index = i
			}

			return cats
		})

	// go Sex.Err(pistol.Run(Cors))
	Sex.Err(Http.ListenAndServe(":8000", Cors(pistol)))
}

func Cap(t string) string {
	if len(t) <= 1 {
		return Str.ToUpper(t)
	}

	return Str.ToUpper(t[:1]) + Str.ToLower(t[1:])
}
