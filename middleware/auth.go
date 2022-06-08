package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Authentication(c *gin.Context) bool {
	_, err := c.Request.Cookie("seasalt")
	log.Println("outb")
	if err != nil {
		log.Println("outback")
		log.Println(err)
		c.Redirect(http.StatusPermanentRedirect, "/sellersignup")
		return false
	}
	return true
}
