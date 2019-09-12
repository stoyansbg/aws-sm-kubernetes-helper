/*
	Input comes from environment variables:

		SECRET_NAME
		SECRET_DEST_PATH		- optional, default is /var/run/secrets/aws-sm/.secret
		AWS_REGION
		AWS_ACCESS_KEY_ID
		AWS_SECRET_ACCESS_KEY
*/

package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

var (
	secretName string
	region     string

	secret string

	secretDestPath string

	err error
)

func getSecret() string {

	//Create a Secrets Manager client
	svc := secretsmanager.New(session.New(&aws.Config{Region: aws.String(region)}))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	// Simplified error checking. Check sample code from AWS if you run into issues.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html

	result, err := svc.GetSecretValue(input)
	if err != nil {
		log.Fatal(err.Error())
	}

	if result.SecretString == nil {
		log.Fatal("the secret is empty")
	}

	return *result.SecretString

}

func saveSecret(secret string, dest string) {
	if err := ioutil.WriteFile(dest, []byte(secret), 0600); err != nil {
		log.Fatalf("failed to save secret to %s: %s", dest, err)
	}
	return
}

func main() {

	log.Println("Starting the AWS Secrets Manager Kubernetes Helper...")

	secretName = os.Getenv("SECRET_NAME")
	if secretName == "" {
		log.Fatal("SECRET_NAME must be set and not empty")
	}

	region = os.Getenv("AWS_REGION")
	if region == "" {
		log.Fatal("AWS_REGION must be set and not empty")
	}

	secretDestPath = os.Getenv("SECRET_DEST_PATH")
	if secretDestPath == "" {
		secretDestPath = "/var/run/secrets/aws-sm/.secret"
	}

	secret = getSecret()
	log.Println("pulled secret from AWS Secret manager")
	saveSecret(secret, secretDestPath)
	log.Printf("saved secret to %s", secretDestPath)

}
