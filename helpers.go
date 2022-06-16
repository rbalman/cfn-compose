package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

func getAWSSession(profile string, region string) (*session.Session, error) {
	return session.NewSessionWithOptions(session.Options{
		Profile: profile,
		Config: aws.Config{
			Region: &region,
		},
		SharedConfigState: session.SharedConfigEnable,
	})
}

func getCallerIdentity(sess *session.Session) (*sts.GetCallerIdentityOutput, error) {
	svc := sts.New(sess)
	input := &sts.GetCallerIdentityInput{}

	return svc.GetCallerIdentity(input)
}
