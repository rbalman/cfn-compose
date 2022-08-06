package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

func getAWSSession() (*session.Session, error) {
	region := os.Getenv("AWS_REGION")
	return session.NewSessionWithOptions(session.Options{
		// Profile: profile,
		Config: aws.Config{
			Region: &region,
			Retryer: client.DefaultRetryer{ //https://github.com/aws/aws-sdk-go/tree/main/example/aws/request/customRetryer
				NumMaxRetries: 0,
			},
		},
		SharedConfigState: session.SharedConfigEnable,
	})
}

func getCallerIdentity(sess *session.Session) (*sts.GetCallerIdentityOutput, error) {
	svc := sts.New(sess)
	input := &sts.GetCallerIdentityInput{}

	return svc.GetCallerIdentity(input)
}

func printCallerIdentity(identity *sts.GetCallerIdentityOutput) {
	fmt.Printf("Account: %s\n", *identity.Account)
	fmt.Printf("Region: %s\n", os.Getenv("AWS_REGION"))
	fmt.Printf("User: %s\n", *identity.UserId)
}
