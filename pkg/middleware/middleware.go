package middleware

import (
	"github.com/gin-gonic/gin"

	"Yamy-Gin/pkg/auth"

)

func AuthMiddleware() gin.HandlerFunc{
	return func(c *gin.Context){
		cookie, err := c.Cookie("access_token")
		if err != nil{
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}

		user, errJWTtoken := auth.ValidateToken(cookie)

		if errJWTtoken != nil{
			c.AbortWithStatusJSON(401, gin.H{"invalid JWT-token": errJWTtoken.Error()})
			return
		}
		
		c.Set("user_id", user.ID)
		c.Next()
	}
}