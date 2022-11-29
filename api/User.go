package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type User struct {
}

func (u User) GetUserInfo(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func PostUserAdd(c *gin.Context) {
	c.String(http.StatusOK, "postUserAdd")
}
