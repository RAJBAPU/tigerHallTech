package middleware

import (
	"fmt"
	"net/http"
	"strings"

	models "simpl_pr/model"
	util "simpl_pr/utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func DeserializeUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token string
		cookie, err := ctx.Cookie("token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			token = fields[1]
		} else if err == nil {
			token = cookie
		}

		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		//config, _ := initializers.LoadConfig(".")
		configs := models.GetAllConfigs()
		tokeSecret := configs["TokenSecret"]
		sub, err := util.ValidateToken(token, tokeSecret)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		user, err := models.GetYpUserById(cast.ToInt(fmt.Sprint(sub)))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
			return
		}

		ctx.Set("currentUser", user)
		fmt.Println(ctx.Get("currentUser"))
		ctx.Next()
	}
}
