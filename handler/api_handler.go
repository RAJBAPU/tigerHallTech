package handler

import (
	models "simpl_pr/model"

	"simpl_pr/service"

	driver "simpl_pr/service/driver"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func SignupPage(dependency *service.User) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var newUser service.SignUpInput
		response := ResponseWithModel(401, "invalid request body", nil)

		err := ctx.BindJSON(&newUser)
		if err != nil {
			ctx.JSON(200, response)
			return

		}
		user := driver.NewUser(
			dependency.Configs,
			dependency.User,
		)
		errCode, err := user.SignUpUser(newUser)
		if err != nil {
			response = ResponseWithModel(460, "internal server error", nil)
			ctx.JSON(200, response)
			return
		}

		if errCode == 452 || errCode == 453 { // 452: invalid email or password, 453: password mismatch
			response = ResponseWithModel(errCode, "Please Retry", nil)
			ctx.JSON(200, response)
			return
		}
		response = ResponseWithModel(200, "success", nil)
		ctx.JSON(200, response)

	}
}

func VerifyEmail(dependency *service.User) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		verificationCode := ctx.Param("verificationCode")
		response := ResponseWithModel(401, "invalid request body", nil)

		user := driver.NewUser(
			dependency.Configs,
			dependency.User,
		)
		err := user.VerifyEmail(verificationCode)
		if err != nil {
			response = ResponseWithModel(460, "internal server error", nil)
			ctx.JSON(200, response)
			return
		}
		response = ResponseWithModel(200, "success", nil)
		ctx.JSON(200, response)

	}
}

func SignInUser(dependency *service.User) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var newUser service.SignInInput
		response := ResponseWithModel(401, "invalid request body", nil)
		err := ctx.BindJSON(&newUser)
		if err != nil {
			ctx.JSON(200, response)
			return

		}
		user := driver.NewUser(
			dependency.Configs,
			dependency.User,
		)
		responseBody, err := user.SignInUser(newUser, ctx)
		if err != nil {
			response = ResponseWithModel(460, "internal server error", nil)
			ctx.JSON(200, response)
			return
		}

		response = ResponseWithModel(200, "success", responseBody)
		ctx.JSON(200, response)

	}

}
func GetUserDetails(dependency *service.User) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUser, exists := ctx.Get("currentUser")
		response := ResponseWithModel(401, "invalid request body", nil)
		if !exists {
			ctx.JSON(200, response)
			return
		}

		currentUuser := currentUser.(*models.TgUser)
		user := driver.NewUser(
			dependency.Configs,
			dependency.User,
		)

		userDetails, err := user.GetUserDetails(currentUuser)
		if err != nil {
			response = ResponseWithModel(460, "internal server error", nil)
			ctx.JSON(200, response)
			return
		}
		response = ResponseWithModel(200, "success", userDetails)
		ctx.JSON(200, response)

	}
}

func PostTigerDetails(dependency *service.Tiger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var newTiger service.TigerDetails
		response := ResponseWithModel(401, "invalid request body", nil)

		err := ctx.BindJSON(&newTiger)
		if err != nil {
			ctx.JSON(200, response)
			return

		}
		tiger := driver.NewTiger(
			dependency.Configs,
			dependency.User,
			dependency.TigerDetails,
			dependency.TigerImages,
		)

		err = tiger.PostTigerDetails(newTiger)
		if err != nil {
			response = ResponseWithModel(460, "internal server error", nil)
			ctx.JSON(200, response)
			return
		}
		response = ResponseWithModel(200, "success", nil)
		ctx.JSON(200, response)
	}

}

func GetAllTigers(dependency *service.Tiger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response := ResponseWithModel(401, "invalid request body", nil)
		page := cast.ToInt(ctx.DefaultQuery("page", "1"))
		pageSize := cast.ToInt(ctx.DefaultQuery("pageSize", "10"))

		tiger := driver.NewTiger(
			dependency.Configs,
			dependency.User,
			dependency.TigerDetails,
			dependency.TigerImages,
		)
		tigerDetails, err := tiger.GetAllTigers(page, pageSize)
		if err != nil {
			response = ResponseWithModel(460, "internal server error", nil)
			ctx.JSON(200, response)
			return
		}
		response = ResponseWithModel(200, "success", tigerDetails)
		ctx.JSON(200, response)
	}
}

func PostSightingDetails(dependency *service.Tiger) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var tigerSighting service.TigerDetails
		response := ResponseWithModel(401, "invalid request body", nil)

		err := ctx.BindJSON(&tigerSighting)
		if err != nil {
			ctx.JSON(200, response)
			return

		}
		currentUser, exists := ctx.Get("currentUser")
		if !exists {
			ctx.JSON(200, response)
			return
		}
		user := currentUser.(*models.TgUser)

		tiger := driver.NewTiger(
			dependency.Configs,
			dependency.User,
			dependency.TigerDetails,
			dependency.TigerImages,
		)

		errorCode, err := tiger.PostSightingDetails(tigerSighting, user)
		if errorCode == 452 { // 452: Tiger sighted in 5Km radius
			response = ResponseWithModel(errorCode, "Tiger was already spotted in range", nil)
			ctx.JSON(200, response)
			return
		}

		if err != nil {
			response = ResponseWithModel(460, "internal server error", nil)
			ctx.JSON(200, response)
			return
		}
		response = ResponseWithModel(200, "success", nil)
		ctx.JSON(200, response)
	}

}

func GetSightingDetails(dependency *service.Tiger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tigerId := ctx.Query("tigerId")
		page := cast.ToInt(ctx.DefaultQuery("page", "1"))
		pageSize := cast.ToInt(ctx.DefaultQuery("pageSize", "10"))
		response := ResponseWithModel(401, "invalid request body", nil)
		if tigerId == "" {
			ctx.JSON(200, response)
			return
		}

		tiger := driver.NewTiger(
			dependency.Configs,
			dependency.User,
			dependency.TigerDetails,
			dependency.TigerImages,
		)
		responseBody, err := tiger.GetSightingDetails(cast.ToInt(tigerId), page, pageSize)
		if err != nil {
			response = ResponseWithModel(460, "internal server error", nil)
			ctx.JSON(200, response)
			return
		}
		response = ResponseWithModel(200, "success", responseBody)
		ctx.JSON(200, response)
	}
}
func ResponseWithModel(code int, msg string, model interface{}) Response {
	return Response{
		Code:  code,
		Msg:   msg,
		Model: model,
	}
}

type Response struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Model interface{} `json:"model"`
}
