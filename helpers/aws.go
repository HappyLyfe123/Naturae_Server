package helpers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"os"
)

var currSession *s3.S3

//ConnectToS3 :
func ConnectToS3() {
	token := ""
	cred := credentials.NewStaticCredentials(os.Getenv("AWS_KEY_ID"), os.Getenv("AWS_ACCESS_KEY"), token)
	_, err := cred.Get()
	if err != nil {
		log.Fatalf("Connecting to AWS S3 error %v\n", err)
		// handle error
	}
	cfg := aws.NewConfig().WithRegion("us-west-2").WithCredentials(cred)
	currSession = s3.New(session.New(), cfg)
	log.Printf("Successfully connected to S3 %s\n", currSession.ServiceID)
}

//GetS3Session : return the current s3 session
func GetS3Session() *s3.S3 {
	return currSession
}
