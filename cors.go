package main

import (
	"net/http"

	"github.com/Plankiton/SexPistol"
	"github.com/rs/cors"
)

func Cors(p *Sex.Pistol) http.Handler {
	return cors.Default().Handler(cors.AllowAll().Handler(p))
}
