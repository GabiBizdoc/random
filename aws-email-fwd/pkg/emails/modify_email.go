package emails

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"mime/multipart"
	"net/mail"
	"net/textproto"
	"os"
)

func MyReceiverEmail() string {
	return os.Getenv("RECEIVER_EMAIL")
}

func MySenderEmail() string {
	return os.Getenv("SENDER_EMAIL")
}

// Deprecated: getPreview Text area content. A shot summary about the email.
// No longer used because `message/rfc822` does this anyway
func getPreview(original *mail.Message) []byte {
	sb := bytes.Buffer{}
	appendData := func(key string, value string) {
		sb.WriteString(fmt.Sprintf("%s: %s\n", key, value))
	}

	appendData("Subject", original.Header.Get("Subject"))
	appendData("From", original.Header.Get("From"))
	appendData("To", original.Header.Get("To"))
	appendData("Return-Path", original.Header.Get("Return-Path"))
	appendData("CC", original.Header.Get("CC"))

	// TODO: link to the original file
	//appendData("File !Ref", "S3 file")

	return sb.Bytes()
}

// WrapEmail Wraps the email into a forward email and returns the raw body.
func WrapEmail(original *mail.Message) []byte {
	msg := bytes.Buffer{}
	multipartWriter := multipart.NewWriter(&msg)

	contentType := "multipart/mixed; boundary=" + multipartWriter.Boundary()
	msg.WriteString("Content-Type: " + contentType + "\r\n")
	msg.WriteString("MIME-Version: 1.0\n")
	msg.WriteString("From: " + MySenderEmail() + "\n")
	msg.WriteString("To: " + MyReceiverEmail() + "\n")
	msg.WriteString("Subject: [Fwd]: " + original.Header.Get("Subject") + "\n\n")

	originalEmailPart, _ := multipartWriter.CreatePart(textproto.MIMEHeader{
		"Content-Type": {"message/rfc822"},
	})

	originalEmailPart.Write(FormatEmail(original))
	multipartWriter.Close()
	return msg.Bytes()
}

// ModifyEmail This function overrides the headers in the original email
// due to SES's requirement for a verified sender and return path.
// We change the 'To' address and adjust the title to compensate for any lost information.
func ModifyEmail(original *mail.Message) (*string, []*string) {
	const fromSeparator = "• from •"
	const toSeparator = "• to •"

	fakeFormat := func(name string, email string, separator string) string {
		return fmt.Sprintf("%s %s <%s>", name, separator, email)
	}

	originalFrom, err := mail.ParseAddress(original.Header.Get("From"))
	if err != nil {
		originalFrom = &mail.Address{Name: "Unknown From", Address: MySenderEmail()}
	}
	newFrom := mail.Address{
		Name:    fakeFormat(originalFrom.Name, originalFrom.Address, fromSeparator),
		Address: MySenderEmail(),
	}

	originalTo, err := mail.ParseAddress(original.Header.Get("To"))
	if err != nil {
		originalTo = &mail.Address{Name: "Unknown To", Address: MySenderEmail()}
	}
	newTo := mail.Address{Name: fakeFormat(originalTo.Name, originalTo.Address, toSeparator), Address: MyReceiverEmail()}

	original.Header["From"] = []string{newFrom.String()}
	original.Header["To"] = []string{newTo.String()}
	original.Header["Return-Path"] = []string{newFrom.String()}
	original.Header["Subject"] = []string{
		fmt.Sprintf("[PROXY]: %s %s %s %s", original.Header.Get("Subject"), fromSeparator, originalFrom.Name, originalFrom.Address),
	}

	source := aws.String(newFrom.String())
	destinations := []*string{aws.String(newTo.String())}

	return source, destinations
}
