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
	if err != nil {
		Sex.Die(err)
	}

	db.AddModels(
		new(Category),
		new(Meal),
	)

	pistol := Sex.NewPistol().
	Add("/", func (r Sex.Request) Sex.Json {
		cats := [] Category {}
		if db.Find(&cats).Error != nil {
			return nil
		}

		page, err := StrConv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			page = 1
		}

		limit, err := StrConv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			limit = 10
		}

		for c, cat := range cats {
			db.Joins("join categoria cat on cat.ID = id_categoria").
			Offset((page - 1) * limit).
			Limit(limit).
			Find(&cat.Meals)
			cat.Name = Cap(cat.Name)
			for m, meal := range cat.Meals {
				cat.Meals[m].Name = Cap(meal.Name)
				cat.Meals[m].Desc = Cap(meal.Desc)
			}

			cats[c] = cat
		}

		return cats
	})

	Sex.Err(Http.ListenAndServe(":8000", Cors(pistol)))
	// Sex.Err(pistol.Run(Cors))
}

func Cap(t string) string {
	if len(t) <= 1 {
		return Str.ToUpper(t)
	}

	return Str.ToUpper(t[:1])+Str.ToLower(t[1:])
}

