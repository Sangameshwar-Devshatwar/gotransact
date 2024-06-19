package utils

import (
	"gotransact/apps/Accounts/models"
	transmodel "gotransact/apps/transaction/models"
	"gotransact/apps/transaction/structutils"
	"gotransact/apps/transaction/validators"
	"gotransact/logger"
	"gotransact/pkg/db"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func PostPayment(payreq structutils.PaymentRequest, user models.User) (int, string, map[string]interface{}) {
	logger.InfoLogger.WithFields(logrus.Fields{}).Info("attempted postpayment method with email", user.Email, " and company ", user.Company)
	if err := validators.GetValidator().Struct(payreq); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := make(map[string]string)
			for _, fieldErr := range validationErrors {
				fieldName := fieldErr.Field()
				// Using the field name to get the custom error message
				if customMessage, found := validators.CustomErrorMessages[fieldName]; found {
					errors[fieldName] = customMessage
				} else {
					errors[fieldName] = err.Error() // Default error message
				}
			}
			return http.StatusBadRequest, "Validation error", map[string]interface{}{}
		}

		return http.StatusBadRequest, "error while validating", map[string]interface{}{}
	}

	var paymentgateway transmodel.PaymentGateway
	if err := db.DB.Where("slug=?", "card").First(&paymentgateway).Error; err != nil {
		logger.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error searching for slug")
		return http.StatusInternalServerError, "no user found in database", map[string]interface{}{}
	}
	//var TransactionReq transmodel.TransactionRequest
	transactiondetails := transmodel.TransactionRequest{
		UserID:                 user.ID,
		Status:                 transmodel.StatusProcessing,
		Description:            payreq.Description,
		Amount:                 payreq.Amount,
		PaymentGatewayMethodID: paymentgateway.ID,
	}
	if err := db.DB.Create(&transactiondetails).Error; err != nil {
		logger.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error creating transaction details in database")
		return http.StatusInternalServerError, "failed to create transactiondetails in database", map[string]interface{}{}
	}
	logger.InfoLogger.WithFields(logrus.Fields{
		"data": "transaction details",
	}).Info("transaction details added to database successfully")
	transactionHistory := transmodel.TransactionHistory{
		TransactionID: transactiondetails.ID,
		Status:        transactiondetails.Status,
		Description:   transactiondetails.Description,
		Amount:        transactiondetails.Amount,
	}
	if err := db.DB.Create(&transactionHistory).Error; err != nil {
		logger.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error adding transaction to transaction history")
		return http.StatusInternalServerError, "failed to create transactiondetails in database", map[string]interface{}{}
	}
	logger.InfoLogger.WithFields(logrus.Fields{}).Info("transaction details added to transaction history successfully")
	go SendMail(user, transactiondetails)
	return http.StatusOK, "mail sent for confirmation of transaction to user", map[string]interface{}{}
}
func ConfirmationPay(id string, status string) (int, string, map[string]interface{}) {
	logger.InfoLogger.WithFields(logrus.Fields{}).Info("attempted confirm payement method", "transaction id:", id)
	transid, err := uuid.Parse(id)
	if err != nil {
		return http.StatusBadRequest, "Invalid transaction ID", map[string]interface{}{}
	}

	if status == "true" {
		var transreq transmodel.TransactionRequest
		if err := db.DB.Where("internal_id=?", transid).First(&transreq).Error; err != nil {
			return http.StatusBadRequest, "transaction request not found", map[string]interface{}{}
		}

		// if err:=db.DB.Update("status",models.StatusSuccess)
		if err := db.DB.Model(&transreq).Where("internal_id = ?", transid).Update("status", transmodel.StatusSuccess).Error; err != nil {
			logger.ErrorLogger.WithFields(logrus.Fields{"error": err.Error()}).Error("error whike updating status as success")
			return http.StatusInternalServerError, "failed to update status", map[string]interface{}{}
		}
		transhistory := transmodel.TransactionHistory{
			TransactionID: transreq.ID,
			Status:        transmodel.StatusSuccess,
			Description:   transreq.Description,
			Amount:        transreq.Amount,
		}
		if err := db.DB.Create(&transhistory).Error; err != nil {
			logger.ErrorLogger.WithFields(logrus.Fields{"error": err.Error()}).Error("error while creating transaction history in database")
			return http.StatusInternalServerError, "error while creating record in database", map[string]interface{}{}
		}
		//return http.StatusOK, "payment done", map[string]interface{}{}

	} else if status == "false" {
		var transreq transmodel.TransactionRequest
		if err := db.DB.Where("internal_id=?", transid).First(&transreq).Error; err != nil {
			return http.StatusInternalServerError, "can't find in database", map[string]interface{}{}
		}

		// if err:=db.DB.Update("status",models.StatusSuccess)
		if err := db.DB.Model(&transreq).Where("internal_id = ?", transid).Update("status", transmodel.StatusFailed).Error; err != nil {
			return http.StatusInternalServerError, "failed to update status", map[string]interface{}{}
		}
		transhistory := transmodel.TransactionHistory{
			TransactionID: transreq.ID,
			Status:        transmodel.StatusFailed,
			Description:   transreq.Description,
			Amount:        transreq.Amount,
		}
		if err := db.DB.Create(&transhistory).Error; err != nil {
			logger.ErrorLogger.WithFields(logrus.Fields{"error": err.Error()}).Error("error while creating transaction history in database")
			return http.StatusInternalServerError, "error while creating record in database", map[string]interface{}{}
		}
		logger.InfoLogger.WithFields(logrus.Fields{}).Info("transaction Cancelled")
		return http.StatusOK, "Transaction Cancelled", map[string]interface{}{}
	}
	logger.InfoLogger.WithFields(logrus.Fields{}).Info("transaction successfull")
	return http.StatusOK, "Transaction successfull", map[string]interface{}{}
}
