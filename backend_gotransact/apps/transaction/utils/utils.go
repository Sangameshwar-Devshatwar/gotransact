package utils

import (
	"bytes"
	"fmt"
	"gotransact/apps/Accounts/models"
	mod "gotransact/apps/transaction/models"
	"log"
	"text/template"
	"time"

	"github.com/google/uuid"
	"gopkg.in/mail.v2"
)

type EmailData struct {
	Username         string
	TransactionID    uuid.UUID
	Amount           string
	DateTime         time.Time
	ConfirmationLink string
	FailPaymentLink  string
	AdminEmail       string
	CompanyName      string
	CompanyContact   string
	CompanyWebsite   string
}

func SendMail(user models.User, transdetails mod.TransactionRequest) {
	fmt.Println("====================startof the mail=============")
	tmpl, err := template.ParseFiles("/home/trellis/Desktop/folder-1/golang-training/backend_gotransact/apps/transaction/utils/mail_template.html")
	if err != nil {
		log.Fatal("Error loading email template:", err)
	}
	//iid, err := uuid.Parse(transdetails.InternalID)

	data := EmailData{
		Username:         user.FirstName,
		TransactionID:    transdetails.InternalID,
		Amount:           transdetails.Amount,
		DateTime:         transdetails.CreatedAt,
		ConfirmationLink: "http:" + "//localhost:8080/api/confirm_payment?transactionid=" + transdetails.InternalID.String() + "&status=true", // Adjust the confirmation link as needed
		FailPaymentLink:  "http://" + "localhost:8080/api/confirm_payment?transactionid=" + transdetails.InternalID.String() + "&status=false",
		AdminEmail:       "Admin123@gmail.com", // Adjust the admin email
		CompanyName:      "trellisMagic",
		CompanyContact:   "Tower 4,infocity ,Gandhinagar,Gujrat",
		CompanyWebsite:   "http://trellisMagic.com",
	}

	// Execute the template with the data
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		log.Fatal("Error executing email template:", err)
	}

	// Create the email
	m := mail.NewMessage()
	m.SetHeader("From", "sangatrellis123@gmail.com")
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Payment Confirmation")
	m.SetBody("text/html", body.String())

	// Configure the SMTP server
	d := mail.NewDialer("smtp.gmail.com", 587, "sangatrellis123@gmail.com", "mhch pnah ljze lsyw")

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Error sending email:", err)
	}
	fmt.Println("====================end of the mail=============")

}
