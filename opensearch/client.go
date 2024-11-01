package opensearch

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/elastic/go-elasticsearch/v8"
)

type OpenSearchConfig struct {
	ENVIRONMENT string
	REGION      string
	AWS         AWSCredentialConfig
}

type AWSCredentialConfig struct {
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
}

// Client is a global instance of the elasticsearch client
var Client *elasticsearch.Client

// InitClient initializes the OpenSearch client with environment-based configuration
func InitClient(config OpenSearchConfig) (*elasticsearch.Client, error) {
	// Fetch OpenSearch configuration from AWS Parameter Store
	endpoint, username, password, err := fetchOpenSearchConfig(config)
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

	return Client, nil
}

// fetchOpenSearchConfig retrieves OpenSearch configuration for the specified environment
func fetchOpenSearchConfig(config OpenSearchConfig) (endpoint, username, password string, err error) {
	// Create session with explicit credentials
	ssmClient := ssm.New(session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(config.REGION),
		Credentials: credentials.NewStaticCredentials(config.AWS.AWS_ACCESS_KEY_ID, config.AWS.AWS_SECRET_ACCESS_KEY, ""),
	})))

	// Define the base path for the parameters
	path := fmt.Sprintf("/ammoze/querycentre/%s/", config.ENVIRONMENT)

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
