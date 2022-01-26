// @Version 1.0.0
// @Title API-Gateway
// @Description Backend for the API-Gateway Software-Architectures Project
// @Server http://www.fake.com Server-1
// @Security AuthorizationHeader read write
// @SecurityScheme AuthorizationHeader http bearer Input your token

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/team-api-gateway/api-gw-backend/internal/crawler"
	"github.com/team-api-gateway/api-gw-backend/internal/database"
	"github.com/team-api-gateway/api-gw-backend/internal/handler"
	"github.com/team-api-gateway/api-gw-backend/internal/sync"
)

func main() {
	godotenv.Load()

	db, err := database.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}

	crawl := crawler.New([]string{
		"https://gist.githubusercontent.com/neinkob15/51f88128f392921fa07628d0a17f8487/raw/04814a678d7bece550fb23ce9670e69f24986c02/gistfile1.txt",
		"https://raw.githubusercontent.com/OAI/OpenAPI-Specification/main/examples/v3.0/callback-example.json",
		"https://raw.githubusercontent.com/OAI/OpenAPI-Specification/main/examples/v3.0/link-example.yaml",
	})
	specs, err := crawl.LoadSpecs()
	if err != nil {
		fmt.Println(err)
		return
	}

	cron := gocron.NewScheduler(time.Local)
	cron.Every(10).Minutes().Do(func() {
		err = sync.Sync(db, specs)
		if err != nil {
			fmt.Println(err)
		}
	})
	cron.StartAsync()

	router := handler.NewRouter()
	h := handler.NewHandler(db)

	router.Mount("/", h.Routes())
	fmt.Println("APP IS WORKING!!!")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println(err)
		return
	}
}
