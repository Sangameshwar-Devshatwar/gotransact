package handlers

import (
	"gotransact/apps/transaction/utils"
	"gotransact/logger"
	"gotransact/responses"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// On link click of the attached payment link, redirect to transaction status page displaying
// transaction ID and payment status.
func ConfirmationPayment(c *gin.Context) {
	logger.InfoLogger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"url":    c.Request.URL.String(),
	}).Info("attempted confirmpayment")
	id := c.Query("transactionid")
	status := c.Query("status")
	statuses, message, data := utils.ConfirmationPay(id, status)
	c.JSON(statuses, responses.UserResponse{
		Status:  statuses,
		Message: message,
		Data:    data,
	})
	// transid, err := uuid.Parse(id)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "error while parsing internal_id"})
	// 	return
	// }

	// if status == "true" {
	// 	var transreq models.TransactionRequest
	// 	if err := db.DB.Where("internal_id=?", transid).First(&transreq).Error; err != nil {
	// 		c.JSON(http.StatusInternalServerError, responses.UserResponse{
	// 			Status:  http.StatusInternalServerError,
	// 			Message: "can't find in database",
	// 			Data: map[string]interface{}{
	// 				"error": err.Error(),
	// 			},
	// 		})
	// 	}

	// 	// if err:=db.DB.Update("status",models.StatusSuccess)
	// 	if err := db.DB.Model(&transreq).Where("internal_id = ?", transid).Update("status", models.StatusSuccess).Error; err != nil {
	// 		c.JSON(http.StatusInternalServerError, responses.UserResponse{
	// 			Status:  http.StatusInternalServerError,
	// 			Message: "failed to update status",
	// 			Data: map[string]interface{}{
	// 				"data": err.Error(),
	// 			},
	// 		})
	// 		return
	// 	}
	// 	transhistory := models.TransactionHistory{
	// 		TransactionID: transreq.ID,
	// 		Status:        models.StatusSuccess,
	// 		Description:   transreq.Description,
	// 		Amount:        transreq.Amount,
	// 	}
	// 	if err := db.DB.Create(&transhistory).Error; err != nil {
	// 		c.JSON(http.StatusInternalServerError, responses.UserResponse{
	// 			Status:  http.StatusInternalServerError,
	// 			Message: "error while creating record in database",
	// 			Data: map[string]interface{}{
	// 				"data": err.Error(),
	// 			},
	// 		})
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, responses.UserResponse{
	// 		Status:  http.StatusOK,
	// 		Message: "payment done",
	// 		Data: map[string]interface{}{
	// 			"data": "success",
	// 		},
	// 	})
	// } else if status == "false" {
	// 	var transreq models.TransactionRequest
	// 	if err := db.DB.Where("internal_id=?", transid).First(&transreq).Error; err != nil {
	// 		c.JSON(http.StatusInternalServerError, responses.UserResponse{
	// 			Status:  http.StatusInternalServerError,
	// 			Message: "can't find in database",
	// 			Data: map[string]interface{}{
	// 				"error": err.Error(),
	// 			},
	// 		})
	// 	}

	// 	// if err:=db.DB.Update("status",models.StatusSuccess)
	// 	if err := db.DB.Model(&transreq).Where("internal_id = ?", transid).Update("status", models.StatusFailed).Error; err != nil {
	// 		c.JSON(http.StatusInternalServerError, responses.UserResponse{
	// 			Status:  http.StatusInternalServerError,
	// 			Message: "failed to update status",
	// 			Data: map[string]interface{}{
	// 				"data": err.Error(),
	// 			},
	// 		})
	// 		return
	// 	}
	// 	transhistory := models.TransactionHistory{
	// 		TransactionID: transreq.ID,
	// 		Status:        models.StatusFailed,
	// 		Description:   transreq.Description,
	// 		Amount:        transreq.Amount,
	// 	}
	// 	if err := db.DB.Create(&transhistory).Error; err != nil {
	// 		c.JSON(http.StatusInternalServerError, responses.UserResponse{
	// 			Status:  http.StatusInternalServerError,
	// 			Message: "error while creating record in database",
	// 			Data: map[string]interface{}{
	// 				"data": err.Error(),
	// 			},
	// 		})
	// 		return
	// 	}

	// }

}
