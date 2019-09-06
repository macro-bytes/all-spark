// +build !cli

package main

import (
	"cloud"
	"context"
	"errors"
	"template"

	"github.com/aws/aws-lambda-go/lambda"
)

type lambdaAction struct {
	command string
	params  template.AwsTemplate
}

type lambdaResponse struct {
	message string
}

func HandleRequest(ctx context.Context, action lambdaAction) lambdaResponse {
	client := &cloud.AwsEnvironment{}

	var err error = nil
	switch action.command {
	case CreateCluster:
		_, err = client.CreateClusterHelper(action.params)
		break
	case DestroyCluster:
		err = client.DestroyClusterHelper(action.params)
		break
	default:
		err = errors.New("invalid command " + action.command)
		break
	}

	if err != nil {
		return lambdaResponse{"error: " + err.Error()}
	}
	return lambdaResponse{"success"}
}

func main() {
	lambda.Start(HandleRequest)
}
