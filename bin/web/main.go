package main

import (
	"github.com/MicBun/go-100-coverage-docker-crud/database"
	"github.com/MicBun/go-100-coverage-docker-crud/docs"
	"github.com/MicBun/go-100-coverage-docker-crud/service"
	"github.com/MicBun/go-100-coverage-docker-crud/web"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default env")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Unable to connect to db %v", err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Unable to migrate db %v", err)
	}

	description := "This is a sample Go User API server. \n\n" +
		"To get Bearer Token, first you need to login. \n\n" +
		"Login by POST /login username: admin, password: admin to get Admin Bearer Token. \n\n" +
		"Login by POST /login username: usera@email.com password: password123 to get User Bearer Token. \n\n" +
		"Then you can use the Bearer Token to access the other endpoints. \n\n" +
		"Admin can access CRUD endpoints, while User can only access /users/get GET endpoint. \n\n" +
		"Checkout my Github: https://github.com/MicBun\n\n" +
		"Checkout my Linkedin: https://www.linkedin.com/in/MicBun\n\n"

	docs.SwaggerInfo.Title = "Go User API"
	docs.SwaggerInfo.Description = description

	c := service.New(db)
	service.SeedData(c)
	web.RegisterAPIRoutes(c)
	c.Web.Run()
}
