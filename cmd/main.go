package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/team-api-gateway/api-gw-backend/internal/crawler"
	"github.com/team-api-gateway/api-gw-backend/internal/database"
)

func main() {
	godotenv.Load()

	db, err := database.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}

	crawl := crawler.New([]string{
		"https://raw.githubusercontent.com/OAI/OpenAPI-Specification/main/examples/v3.0/petstore.yaml",
		"https://raw.githubusercontent.com/OAI/OpenAPI-Specification/main/examples/v2.0/yaml/petstore-simple.yaml",
		"https://raw.githubusercontent.com/OAI/OpenAPI-Specification/main/examples/v3.0/callback-example.json",
		"https://raw.githubusercontent.com/OAI/OpenAPI-Specification/main/examples/v3.0/link-example.yaml",
	})
	specs, err := crawl.LoadSpecs()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(specs)
	for _, spec := range specs {
		err := db.InsertAPI(spec)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
