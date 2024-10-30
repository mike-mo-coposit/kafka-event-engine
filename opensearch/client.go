package opensearch

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/elastic/go-elasticsearch/v8"
)

// Client is a global instance of the elasticsearch client
var Client *elasticsearch.Client

// InitClient initializes the OpenSearch client with environment-based configuration
func InitClient() {
	// Get environment (staging, uat, production) from the environment variable
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		log.Fatal("ENVIRONMENT variable is not set")
	}

	// Fetch OpenSearch configuration from AWS Parameter Store
	endpoint, username, password, err := fetchOpenSearchConfig(environment)
	if err != nil {
		log.Fatalf("Failed to fetch OpenSearch configuration: %s", err)
	}

	// Set up OpenSearch client with fetched configuration
	cfg := elasticsearch.Config{
		Addresses: []string{endpoint},
		Username:  username,
		Password:  password,
	}

	Client, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating OpenSearch client: %s", err)
	}
}

// fetchOpenSearchConfig retrieves OpenSearch configuration for the specified environment
func fetchOpenSearchConfig(environment string) (endpoint, username, password string, err error) {
	region := os.Getenv("AWS_OPENSEARCH_REGION")
	if region == "" {
		// region = "your-region" // Fallback region if not set as an environment variable
		log.Fatalf("Environment AWS_OPENSEARH_REGION is NOT SET")
	}

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if accessKey == "" || secretKey == "" {
		log.Fatalf("AWS credentials are not set")
	}

	// Create session with explicit credentials
	ssmClient := ssm.New(session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})))

	// Define the base path for the parameters
	path := fmt.Sprintf("/ammoze/querycentre/%s/", environment)

	// Fetch parameters by path
	input := &ssm.GetParametersByPathInput{
		Path:           aws.String(path),
		WithDecryption: aws.Bool(true), // Enables decryption for secure strings
		Recursive:      aws.Bool(true),
	}

	result, err := ssmClient.GetParametersByPath(input)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get parameters: %w", err)
	}

	// Map each parameter to the correct variable
	for _, param := range result.Parameters {
		switch *param.Name {
		case path + "OpenSearchEndpoint":
			endpoint = *param.Value
		case path + "OpenSearchUsername":
			username = *param.Value
		case path + "OpenSearchPassword":
			password = *param.Value
		}
	}

	if endpoint == "" || username == "" || password == "" {
		return "", "", "", fmt.Errorf("incomplete OpenSearch configuration")
	}

	return endpoint, username, password, nil
}
