package handlers

import (
	modelss "gotransact/apps/Accounts/models"
	"gotransact/apps/transaction/structutils"
	"gotransact/apps/transaction/utils"
	"gotransact/logger"
	"gotransact/responses"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func PostPayment(c *gin.Context) {
	logger.InfoLogger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"url":    c.Request.URL.String(),
	}).Info("attempted postpayment")
	var payreq structutils.PaymentRequest
	if err := c.ShouldBindJSON(&payreq); err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error while binding",
			Data: map[string]interface{}{
				"data": err.Error(),
			},
		})
		return
	}
	tokenuser, exist := c.Get("User")
	if !exist {
		c.JSON(http.StatusNotFound, gin.H{"error": "didnt get user from token "})
		return
	}
	user, err := tokenuser.(modelss.User)
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{"error": "didn't get user from tokenuser "})
		return
	}
	status, message, data := utils.PostPayment(payreq, user)
	c.JSON(status, responses.UserResponse{Status: http.StatusOK, Message: message, Data: data})
}

// if err := validators.GetValidator().Struct(payreq); err != nil {
// 	if validationErrors, ok := err.(validator.ValidationErrors); ok {
// 		errors := make(map[string]string)
// 		for _, fieldErr := range validationErrors {
// 			fieldName := fieldErr.Field()
// 			// Using the field name to get the custom error message
// 			if customMessage, found := validators.CustomErrorMessages[fieldName]; found {
// 				errors[fieldName] = customMessage
// 			} else {
// 				errors[fieldName] = err.Error() // Default error message
// 			}
// 		}

// 		c.JSON(http.StatusBadRequest, responses.UserResponse{
// 			Status:  http.StatusBadRequest,
// 			Message: "Validation error",
// 			Data: map[string]interface{}{
// 				"validation_errors": errors,
// 			},
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusBadRequest, responses.UserResponse{
// 		Status:  http.StatusBadRequest,
// 		Message: "error while validating",
// 		Data: map[string]interface{}{
// 			"data": err.Error(),
// 		},
// 	})
// 	return
// }
// tokenuser, exist := c.Get("User")
// if !exist {
// 	c.JSON(http.StatusNotFound, gin.H{"error": "didnt get user from token "})
// 	return
// }
// user, err := tokenuser.(modelss.User)
// if !err {
// 	c.JSON(http.StatusBadRequest, gin.H{"error": "didn't get user from tokenuser "})
// 	return
// }
// var paymentgateway transmodel.PaymentGateway
// if err := db.DB.Where("slug=?", "card").First(&paymentgateway).Error; err != nil {
// 	fmt.Println("problem is here")
// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

// 	return
// }
// //var TransactionReq transmodel.TransactionRequest
// transactiondetails := transmodel.TransactionRequest{
// 	UserID:                 user.ID,
// 	Status:                 transmodel.StatusProcessing,
// 	Description:            payreq.Description,
// 	Amount:                 payreq.Amount,
// 	PaymentGatewayMethodID: paymentgateway.ID,
// }
// if err := db.DB.Create(&transactiondetails).Error; err != nil {
// 	c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "failed to create transactiondetails in database", Data: map[string]interface{}{"data": err.Error()}})
// 	return
// }
// transactionHistory := transmodel.TransactionHistory{
// 	TransactionID: transactiondetails.ID,
// 	Status:        transactiondetails.Status,
// 	Description:   transactiondetails.Description,
// 	Amount:        transactiondetails.Amount,
// }
// if err := db.DB.Create(&transactionHistory).Error; err != nil {
// 	c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "failed to create transactiondetails in database", Data: map[string]interface{}{"data": err.Error()}})
// 	return
// }
// go utils.SendMail(user, transactiondetails)
