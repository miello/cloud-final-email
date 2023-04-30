package ses

import (
	"bytes"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"
	"gopkg.in/gomail.v2"
)

const (

	// Specify a configuration set. To use a configuration
	// set, comment the next line and line 92.
	//ConfigurationSet = "ConfigSet"

	// The subject line for the email.
	Subject = "Your Jerm Resume is ready! (Ref ID: %s)"

	// The HTML body for the email.

	//The email body for recipients with non-HTML email clients.
	TextBody = "This email was sent with Amazon SES using the AWS SDK for Go."

	// The character encoding for the email.
	CharSet = "UTF-8"
)

func SendEmail(svc *ses.SES, sender string, targetEmail string, file []byte, fileId string) {
	HtmlBody := fmt.Sprintf(`<h1>Your Jerm Resume is ready (Ref id: %s)</h1>`, fileId)
	filename := fmt.Sprintf("jermed-%s.pdf", fileId)
	SubjectFormated := fmt.Sprintf(Subject, fileId)

	var emailRaw bytes.Buffer

	goMail := gomail.NewMessage()
	goMail.SetHeader("From", fmt.Sprintf("Jerm Resume <%s>", sender))
	goMail.SetHeader("To", targetEmail)
	goMail.SetHeader("Subject", SubjectFormated)
	goMail.SetBody("text/html", HtmlBody)
	goMail.Attach(filename, gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(file)
		return err
	}))

	goMail.WriteTo(&emailRaw)

	input := &ses.SendRawEmailInput{
		Destinations: []*string{
			aws.String(targetEmail),
		},
		Source: aws.String(sender),
		RawMessage: &ses.RawMessage{
			Data: emailRaw.Bytes(),
		},
	}

	result, err := svc.SendRawEmail(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}

		return
	}

	fmt.Println("Email Sent to address: " + targetEmail)
	fmt.Println(result)
}
