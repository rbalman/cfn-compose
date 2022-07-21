package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"sync"
	"time"
	"math/rand"
)

func main() {
	workerCount := 20000
	var wg sync.WaitGroup
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		fmt.Printf("Error while creating session: %s\n", err)
	}

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(i int) {
			seconds := rand.Intn(100)
			fmt.Printf("Worker: %d started and sleeping for: %d seconds\n", i, seconds)
			time.Sleep(time.Second * time.Duration(seconds))
			_, err := DescribeStacks(sess, "balman-sqs-consumer")
			if err != nil {
				fmt.Printf("Worker: %d, Error while describing stack: %s\n", i, err)
			}	
			wg.Done()
		}(i)
	}

	wg.Wait()
}


func DescribeStacks(sess *session.Session, stackName string) (*cloudformation.DescribeStacksOutput, error) {
	input := &cloudformation.DescribeStacksInput{
		StackName:    aws.String(stackName),
	}

	svc := cloudformation.New(sess)
	return svc.DescribeStacks(input)
}