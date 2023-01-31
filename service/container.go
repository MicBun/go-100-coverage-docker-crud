package service

import (
	"github.com/MicBun/go-microservice-kubernetes/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Container struct {
	Web   *gin.Engine
	DB    *gorm.DB
	Admin user.AuthInterface
}

func New(mainDB *gorm.DB) *Container {
	ginEngine := gin.Default()

	admin := user.AdminAuth(mainDB)

	return &Container{
		Web:   ginEngine,
		DB:    mainDB,
		Admin: admin,
	}
}
