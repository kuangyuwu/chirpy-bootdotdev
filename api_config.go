package main

import "github.com/kuangyuwu/chirpy-bootdev/internal/database"

type apiConfig struct {
	fileserverHits int
	db             *database.DB
}
