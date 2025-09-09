package main

import (
	"github.com/Archer1A/go-kits/utils/mask"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Resp struct {
	Code int
	Data interface{}
}

type Role struct {
	Name string `json:"name"`
}

type User struct {
	Username mask.UsernameMask `json:"username"`
	Phone    mask.PhoneMask    `json:"phone"`
}

func ListUsers(ctx *gin.Context) {
	users := []User{
		{
			Username: "张大来",
			Phone:    "+8615088888888",
		},
	}
	ctx.JSON(http.StatusOK, users)
}

func main() {
	engine := gin.New()

	engine.GET("/users", ListUsers)
	if err := engine.Run(":8080"); err != nil {
		panic(err)
	}

}
