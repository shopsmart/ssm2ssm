package service

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// SSMClient abstracts out the single function used within the ssmiface
type SSMClient interface {
	GetParametersByPathPages(getParametersByPathInput *ssm.GetParametersByPathInput, fn func(resp *ssm.GetParametersByPathOutput, lastPage bool) bool) error
	PutParameter(input *ssm.PutParameterInput) (*ssm.PutParameterOutput, error)
}

// Service will abstract out calls to AWS
type Service interface {
	GetParameters(searchPath string) (map[string]string, error)
	PutParameters(path string, params map[string]string, overwrite bool) error
}

// Impl will be the implementation of the Service interface
type Impl struct {
	SSMClient SSMClient
}

// New creates a new service
func New() (Service, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return NewFromClient(ssm.New(sess))
}

// NewFromClient creates a new service from an implementation of the SSMClient
func NewFromClient(ssmClient SSMClient) (Service, error) {
	return &Impl{
		SSMClient: ssmClient,
	}, nil
}

// GetParameters fetches all of the parameters under a path into a map
func (svc *Impl) GetParameters(searchPath string) (map[string]string, error) {
	params := map[string]string{}
	getParametersByPathInput := &ssm.GetParametersByPathInput{
		Path:           aws.String(searchPath),
		WithDecryption: aws.Bool(true),
	}

	err := svc.SSMClient.GetParametersByPathPages(getParametersByPathInput, func(resp *ssm.GetParametersByPathOutput, lastPage bool) bool {
		for _, param := range resp.Parameters {
			name := strings.TrimPrefix(*param.Name, fmt.Sprintf("%s/", searchPath))
			params[name] = *param.Value
		}
		return true
	})

	return params, err
}

// PutParameters saves all of the parameters from a map to a path
func (svc *Impl) PutParameters(path string, params map[string]string, overwrite bool) error {
	var err error

	existingParams := map[string]string{}
	if !overwrite {
		getParametersByPathInput := &ssm.GetParametersByPathInput{
			Path:           aws.String(path),
			WithDecryption: aws.Bool(true),
		}

		err = svc.SSMClient.GetParametersByPathPages(getParametersByPathInput, func(resp *ssm.GetParametersByPathOutput, lastPage bool) bool {
			for _, param := range resp.Parameters {
				name := strings.TrimPrefix(*param.Name, fmt.Sprintf("%s/", path))
				existingParams[name] = *param.Value
			}
			return true
		})
	}

	for key, param := range params {
		paramPath := fmt.Sprintf("%s/%s", path, key)

		if !overwrite {
			if _, ok := existingParams[key]; ok {
				return fmt.Errorf("cowardly refusing to overwrite: %s", paramPath)
			}
		}

		putParameterInput := &ssm.PutParameterInput{
			Name:      aws.String(paramPath),
			DataType:  aws.String("text"),
			Type:      aws.String(ssm.ParameterTypeSecureString),
			Overwrite: aws.Bool(overwrite),
			Value:     aws.String(param),
		}

		_, err = svc.SSMClient.PutParameter(putParameterInput)
	}

	return err
}
