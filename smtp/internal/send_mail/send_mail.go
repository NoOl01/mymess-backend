package send_mail

import (
	"fmt"
	"net/smtp"
	"smtp/internal/storage"
)

func SendOtp(email string, otp int32) error {
	from := storage.Env.Email
	password := storage.Env.EmailPassword
	smtpHost := storage.Env.SmthHost
	smtpPort := storage.Env.SmtpPort

	auth := smtp.PlainAuth("", from, password, smtpHost)

	message := []byte(fmt.Sprintf("Ваш код: %d", otp))

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{email}, message)
}
