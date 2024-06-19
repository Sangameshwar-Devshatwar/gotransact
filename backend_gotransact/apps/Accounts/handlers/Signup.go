package handlers

import (
	"gotransact/logger"
	"gotransact/responses"
	"net/http"

	utils "gotransact/apps/Accounts/utils"
	astructutils "gotransact/apps/Astructutils"

	//"gotransact/apps/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// type SignupUser struct {
// 	FirstName   string `json:"firstname" binding:"required,alpha"`
// 	LastName    string `json:"lastname" binding:"required,alpha"`
// 	Email       string `json:"email" binding:"required,email"`
// 	CompanyName string `json:"companyname" binding:"required"`
// 	Password    string `json:"password" binding:"required" validate:"password_complexity"`
// }

func SignupHandler(c *gin.Context) {
	logger.InfoLogger.WithFields(logrus.Fields{
		"method":c.Request.Method,
		"url":c.Request.URL.String(),
	}).Info("attempted signup")
	var req astructutils.SignupUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}
	
	statuscode, message, data := utils.Signup(req)
	c.JSON(statuscode, responses.UserResponse{Status: statuscode,
		Message: message,
		Data:    data})
}

// if err := validators.GetValidator().Struct(req); err != nil {
// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 	return
// }
// fmt.Println(req.Password)
// var existingUser models.User
// if err := db.DB.Where("email=?", req.Email).First(&existingUser).Error; err == nil {
// 	c.JSON(http.StatusConflict, responses.UserResponse{Status: http.StatusConflict, Message: "Email already registrered", Data: map[string]interface{}{"data": "use another mail for register"}})
// 	return
// }
// var existingCompany models.Company
// if err := db.DB.Where("name=?", req.CompanyName).First(&existingCompany).Error; err == nil {
// 	c.JSON(http.StatusConflict, responses.UserResponse{Status: http.StatusConflict, Message: "company name already registrered", Data: map[string]interface{}{"data": "use another company name to register"}})
// 	return
// }
// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
// if err != nil {
// 	c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "failed to hash the password", Data: map[string]interface{}{"data": err.Error()}})
// 	return
// }
// fmt.Println(string(hashedPassword))
// //values stored into original model to add user to database
// user := models.User{
// 	FirstName: req.FirstName,
// 	LastName:  req.LastName,
// 	Email:     req.Email,
// 	Password:  string(hashedPassword),
// 	Company: models.Company{
// 		Name: req.CompanyName,
// 	},
// }
// //user created by using user passesd values binded into struct
// if err := db.DB.Create(&user).Error; err != nil {
// 	c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "failed to create user database", Data: map[string]interface{}{"data": err.Error()}})
// 	return
// }
// //sent mail using goroutines
// go utils.SendMail(req.Email)
