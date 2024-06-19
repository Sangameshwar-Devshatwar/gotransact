package handlers

import (
	astructutils "gotransact/apps/Astructutils"
	"gotransact/logger"
	"gotransact/responses"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	utils "gotransact/apps/Accounts/utils"
)

// type LoginInput struct {
// 	Email    string `json:"email" binding:"required,email"`
// 	Password string `json:"password" binding:"required" validate:"password_complexity"`
// }

func LoginHandler(c *gin.Context) {
	logger.InfoLogger.WithFields(logrus.Fields{
		"method":c.Request.Method,
		"url":c.Request.URL.String(),
	}).Info("attempted login")
	var inputuser astructutils.LoginInput
	if err := c.ShouldBindJSON(&inputuser); err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data: map[string]interface{}{
				"data": err.Error()},
		})
		return
	}

	status, message, data := utils.Login(inputuser)

	c.JSON(status, responses.UserResponse{
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	})
}

// 	if err := validators.GetValidator().Struct(inputuser); err != nil {
// 		c.JSON(http.StatusBadRequest, responses.UserResponse{
// 			Status:  http.StatusBadRequest,
// 			Message: "error while validating fields",
// 			Data: map[string]interface{}{
// 				"data": err.Error()},
// 		})
// 		return
// 	}

// 	var user models.User
// 	// if err := db.DB.Where("email=? AND password=?", inputuser.Email, string(hashedPassword)).First(&user).Error; err != nil {
// 	// 	c.JSON(http.StatusInternalServerError, responses.UserResponse{
// 	// 		Status:  http.StatusInternalServerError,
// 	// 		Message: "no user found in database",
// 	// 		Data: map[string]interface{}{
// 	// 			"data": err.Error(),
// 	// 		},
// 	// 	})
// 	// 	return
// 	// }
// 	if err := db.DB.Where("email=?", inputuser.Email).First(&user).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, responses.UserResponse{
// 			Status:  http.StatusInternalServerError,
// 			Message: "no user found in database",
// 			Data: map[string]interface{}{
// 				"data": err.Error(),
// 			},
// 		})
// 		return
// 	}
// 	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(inputuser.Password), bcrypt.DefaultCost)
// 	// if err != nil {
// 	// 	c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "failed to hash the password", Data: map[string]interface{}{"data": err.Error()}})
// 	// 	return
// 	// }
// 	//fmt.Println(inputuser.Password)

// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputuser.Password)); err != nil {
// 		c.JSON(http.StatusInternalServerError, responses.UserResponse{
// 			Status:  http.StatusInternalServerError,
// 			Message: "invalid password",
// 			Data: map[string]interface{}{
// 				"data": err.Error(),
// 			},
// 		})
// 		return
// 	}
// 	// if string(hashedPassword) == user.Password {

// 	// } else {
// 	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "password is not correct "})
// 	// 	fmt.Println(string(hashedPassword))
// 	// 	fmt.Println(user.Password)
// 	// 	return
// 	// }

// 	token, err := utils.GeneratePasetoToken(user)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, responses.UserResponse{
// 			Status:  http.StatusInternalServerError,
// 			Message: "error generating token",
// 			Data: map[string]interface{}{
// 				"data": err.Error(),
// 			},
// 		})
// 		return
// 	}
