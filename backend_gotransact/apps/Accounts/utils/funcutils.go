package utils

import (
	"gotransact/apps/Accounts/models"
	validators "gotransact/apps/Accounts/validator"
	astructutils "gotransact/apps/Astructutils"
	"gotransact/logger"

	//"gotransact/apps/logger"
	"gotransact/pkg/db"
	"time"

	"net/http"

	//"github.com/sirupsen/logrus"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// var (
// 	ctx = context.Background()
// 	// rdb = redis.NewClient(&redis.Options{
// 	// 	Addr: "localhost:6379", // Redis server address
// 	// })
// )

func Signup(requestInputs astructutils.SignupUser) (int, string, map[string]interface{}) {
	logger.InfoLogger.WithFields(logrus.Fields{}).Info("attempted signup method with email", requestInputs.Email, "and company", requestInputs.CompanyName)
	if err := validators.GetValidator().Struct(requestInputs); err != nil {
		return http.StatusBadRequest, "Password should contain atleast one upper case character,one lower case character,one number and one special character", map[string]interface{}{}
	}
	var count int64
	if err := db.DB.Model(&models.User{}).Where("email = ?", requestInputs.Email).Count(&count).Error; err != nil {
		return http.StatusInternalServerError, "Database error", map[string]interface{}{}
	}
	if count > 0 {
		return http.StatusBadRequest, "email already exists", map[string]interface{}{}
	}
	if err := db.DB.Model(&models.Company{}).Where("name = ?", requestInputs.CompanyName).Count(&count).Error; err != nil {
		return http.StatusInternalServerError, "Database error", map[string]interface{}{}
	}
	if count > 0 {
		return http.StatusBadRequest, "company already exists", map[string]interface{}{}
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestInputs.Password), bcrypt.DefaultCost)
	if err != nil {
		return http.StatusInternalServerError, "failed to hash the password", map[string]interface{}{}
	}

	//values stored into original model to add user to database
	user := models.User{
		FirstName: requestInputs.FirstName,
		LastName:  requestInputs.LastName,
		Email:     requestInputs.Email,
		Password:  string(hashedPassword),
		Company: models.Company{
			Name: requestInputs.CompanyName,
		},
	}
	//user created by using user passesd values binded into struct
	if err := db.DB.Create(&user).Error; err != nil {
		logger.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error while creating database")
		return http.StatusInternalServerError, "failed to create user database", map[string]interface{}{}
	}
	//sent mail using goroutines
	logger.InfoLogger.WithFields(logrus.Fields{}).Info("User created successfully with email", requestInputs.Email, "and company", requestInputs.CompanyName)
	go SendMail(user.Email)
	return http.StatusOK, "User Created Successfully", map[string]interface{}{}
}

// login function
func Login(inputuser astructutils.LoginInput) (int, string, map[string]interface{}) {
	logger.InfoLogger.WithFields(logrus.Fields{}).Info("attempted login method using mail", inputuser.Email)

	if err := validators.GetValidator().Struct(inputuser); err != nil {
		return http.StatusBadRequest, "error while validating fields", map[string]interface{}{}
	}
	var user models.User
	if err := db.DB.Where("email=?", inputuser.Email).First(&user).Error; err != nil {
		logger.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error finding mail in database")
		return http.StatusInternalServerError, "no user found in database", map[string]interface{}{}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputuser.Password)); err != nil {
		return http.StatusInternalServerError, "invalid password", map[string]interface{}{}
	}

	token, err := GeneratePasetoToken(user)
	if err != nil {
		return http.StatusInternalServerError, "error generating token", map[string]interface{}{}
	}
	logger.InfoLogger.WithFields(logrus.Fields{
		
	}).Info("Logged in successfully using mail",inputuser.Email)
	return http.StatusOK, "logged in successfully", map[string]interface{}{"data": token}
}

// logout function to implement logout handler
func Logout(authHeader string) (int, string, map[string]interface{}) {
	logger.InfoLogger.WithFields(logrus.Fields{
		"message": "in logout func",
	}).Info("attempted logout method")

	if authHeader == "" {
		return http.StatusUnauthorized, "authorization header missing", map[string]interface{}{}
	}

	//tokenStr := authHeader[len("Bearer "):]

	_, err := VerifyPasetoToken(authHeader)
	if err != nil {
		return http.StatusUnauthorized, "invalid token", map[string]interface{}{}
	}

	// Blacklist the token by storing it in Redis with an expiration time
	err = rdb.Set(ctx, authHeader, "Blacklisted", 24*time.Hour).Err() // adjust expiration time as needed
	if err != nil {
		return http.StatusInternalServerError, "failed to blacklist token", map[string]interface{}{}
	}
	logger.InfoLogger.WithFields(logrus.Fields{"logout": "success"}).Info("Logged out successfully")
	return http.StatusOK, "logged out successfully", map[string]interface{}{}
}
