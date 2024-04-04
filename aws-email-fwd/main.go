package main

import (
	"bytes"
	"github.com/GabiBizdoc/random/aws-email-fwd/pkg/emails"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ses"
	"io"
	"log"
	"net/mail"
	"os"
	"strings"
)

func readS3File(svcS3 *s3.S3, bucket *string, key *string) ([]byte, error) {
	config := &s3.GetObjectInput{
		Bucket: bucket,
		Key:    key,
	}
	rawObject, err := svcS3.GetObject(config)
	if err != nil {
		return nil, err
	}
	defer rawObject.Body.Close()

	body, err := io.ReadAll(rawObject.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func modifyAndSend(svcSES *ses.SES, rawEmail []byte) error {
	originalEmail, err := mail.ReadMessage(bytes.NewReader(rawEmail))
	if err != nil {
		return err
	}

	// log.Println("Mail Received From: ", originalEmail.Header.Get("From"))
	// for k, v := range originalEmail.Header {
	// 	log.Println(k, v)
	// }

	rawEmailData := emails.WrapEmail(originalEmail)
	source := aws.String(emails.MySenderEmail())
	destinations := []*string{aws.String(emails.MyReceiverEmail())}

	//Source, Destinations := emails.ModifyEmail(originalEmail)
	//rawEmailData := emails.FormatEmail(originalEmail)

	email := &ses.SendRawEmailInput{
		Source:       source,
		Destinations: destinations,
		RawMessage: &ses.RawMessage{
			Data: rawEmailData,
		},
	}

	result, err := svcSES.SendRawEmail(email)
	if err != nil {
		return err
	}
	if result != nil {
		println("Email sent. Message ID: ", *result.MessageId)
	}
	return nil
}

func readEmailFromS3(svcS3 *s3.S3, record *events.S3EventRecord) ([]byte, error) {
	bucketName := aws.String(record.S3.Bucket.Name)
	objectKey := aws.String(record.S3.Object.Key)
	rawEmail, err := readS3File(svcS3, bucketName, objectKey)
	if err != nil {
		return nil, err
	}
	return rawEmail, nil
}

func handler(s3event events.S3Event) error {
	region := os.Getenv("REGION")
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	svcS3 := s3.New(sess)
	svcSES := ses.New(sess)

	log.Printf("Send Email Triggerd by: %v\n", s3event.Records)

	sb := strings.Builder{}
	for _, record := range s3event.Records {
		rawEmail, err := readEmailFromS3(svcS3, &record)
		if err != nil {
			sb.WriteString(err.Error())
			sb.WriteString("  |  ")
			continue
		}
		err = modifyAndSend(svcSES, rawEmail)
		if err != nil {
			sb.WriteString(err.Error())
			sb.WriteString("  |  ")
			continue
		}
	}
	if sb.Len() > 0 {
		panic(sb.String())
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
