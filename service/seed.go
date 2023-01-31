package service

import (
	"log"
)

func SeedData(c *Container) {
	_, err := c.Admin.RegisterUser("admin", "admin", "Admin")
	if err != nil {
		log.Fatalf("Unable to seed data %v", err)
		return
	}
	_, err = c.Admin.RegisterUser("usera@email.com", "password123", "usera")
	if err != nil {
		log.Fatalf("Unable to seed data %v", err)
		return
	}
	_, err = c.Admin.RegisterUser("userb@email.com", "password123", "userb")
	if err != nil {
		log.Fatalf("Unable to seed data %v", err)
		return
	}
}
