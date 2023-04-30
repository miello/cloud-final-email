package ses

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"
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

	Template = `From: Jerm Resume <%s>
To: %s
Subject: %s
MIME-version: 1.0
Content-type: multipart/mixed; boundary="NextPart"

--NextPart
Content-Type: text/html charset=UTF-8
Content-Transfer-Encoding: quoted-printable

%s

--NextPart
Content-Type: application/pdf;
Content-Transfer-Encoding: base64
Content-Disposition: attachment; filename="%s"

%s

--NextPart--`
)

func SendEmail(svc *ses.SES, sender string, targetEmail string, file string, fileId string) {
	HtmlBody := fmt.Sprintf(`<h1>Your Jerm Resume is ready (Ref id: %s)</h1>`, fileId)
	filename := fmt.Sprintf("jermed-%s.pdf", fileId)
	SubjectFormated := fmt.Sprintf(Subject, fileId)

	// TextBody := "Your Jerm Resume is ready! Please see in the attachment."
	input := &ses.SendRawEmailInput{
		Destinations: []*string{
			aws.String(targetEmail),
		},
		Source: aws.String(sender),
		RawMessage: &ses.RawMessage{
			Data: []byte(fmt.Sprintf(Template, sender, targetEmail, SubjectFormated, HtmlBody, filename, file)),
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
