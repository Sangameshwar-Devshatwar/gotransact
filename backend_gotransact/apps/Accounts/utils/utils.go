package utils

import (
	"context"
	"fmt"
	"gotransact/apps/Accounts/models"
	"gotransact/logger"
	"log"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	gomail "gopkg.in/mail.v2"
)

var (
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
	})
)

func SendMail(mail string) {
	logger.InfoLogger.WithFields(logrus.Fields{}).Info("attempted Sendmail() email to", mail)
	abc := gomail.NewMessage()

	abc.SetHeader("From", "sangatrellis123@gmail.com")
	abc.SetHeader("To", mail)
	abc.SetHeader("Subject", "confirmation")
	abc.SetBody("text/plain", "user has been created successfully ,this is a confirmation mail")

	a := gomail.NewDialer("smtp.gmail.com", 587, "sangatrellis123@gmail.com", "mhch pnah ljze lsyw")
	if err := a.DialAndSend(abc); err != nil {
		logger.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error sending mail")
		log.Fatal(err.Error())
	}
	logger.InfoLogger.WithFields(logrus.Fields{}).Info("mail sent successfully to ", mail)
}

// implementing creation of token
var secretKey = paseto.NewV4AsymmetricSecretKey() // don't share this!!!
var publicKey = secretKey.Public()

func GeneratePasetoToken(user models.User) (string, error) {
	now := time.Now()
	exp := now.Add(24 * time.Hour)
	token := paseto.NewToken()
	token.SetIssuedAt(now)
	token.SetExpiration(exp)

	token.Set("User", user)
	signed := token.V4Sign(secretKey, nil)
	return signed, nil
}

func VerifyPasetoToken(signed string) (any, error) {
	val, err := rdb.Get(ctx, signed).Result()
	if err == nil && val == "Blacklisted" {
		return nil, fmt.Errorf("token has been revoked")
	}

	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())
	token, err := parser.ParseV4Public(publicKey, signed, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token %w", err)
	}
	var user models.User
	err = token.Get("User", &user)
	if err != nil {
		return nil, fmt.Errorf("subject claim not found in token")
	}
	return user, nil
}
