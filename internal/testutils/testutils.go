package testutils

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/golang/mock/gomock"
)

//go:generate mockgen -source=../../pkg/service/service.go -destination=gen_service_mocks.go -package testutils SSMClient,Service

var (
	// MockController is a gomock controller
	MockController *gomock.Controller

	// ParametersMap is a map of key-value pairs to represent parameters
	ParametersMap = map[string]string{
		"foo": "<value-to-be-quoted>",
		"bar": "contains a single quote y'all",
		"baz": "uses multiple* special &characters, y'all",
	}
)

// Setup initializes the mock controller
func Setup(t *testing.T) {
	MockController = gomock.NewController(t)
}

// Teardown cleans up after gomock
func Teardown() {
	MockController.Finish()
}

// ParametersSlice converts the ParametersMap into a slice of ssm Parameters
func ParametersSlice(prefix string) []*ssm.Parameter {
	params := []*ssm.Parameter{}
	for key, value := range ParametersMap {
		params = append(params, &ssm.Parameter{
			Name:  aws.String(fmt.Sprintf("%s/%s", prefix, key)),
			Value: aws.String(value),
		})
	}

	return params
}

// MockGetParametersByPath mocks out the GetParametersByPath function on the mock client
// If params is nil, the ParametersSlice will be called passing in path to create the params
// slice
func MockGetParametersByPath(client *MockSSMClient, path string, params []*ssm.Parameter) {
	if params == nil {
		params = ParametersSlice(path)
	}

	in := ssm.GetParametersByPathInput{
		Path:           aws.String(path),
		WithDecryption: aws.Bool(true),
	}

	var _ = in

	client.
		EXPECT().
		GetParametersByPathPages(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ *ssm.GetParametersByPathInput, fn func(*ssm.GetParametersByPathOutput, bool) bool) error {
			resp := ssm.GetParametersByPathOutput{
				NextToken:  nil,
				Parameters: params,
			}
			_ = fn(&resp, true)
			return nil
		}).
		Times(1)
}

// ExpectAllParametersToBePut mocks out the PutParameters function on the mock client for all
// parameters within ParametersMap
func ExpectAllParametersToBePut(client *MockSSMClient, path string, overwrite bool) {
	for key, param := range ParametersMap {
		client.
			EXPECT().
			PutParameter(&ssm.PutParameterInput{
				Name:      aws.String(fmt.Sprintf("%s/%s", path, key)),
				DataType:  aws.String("text"),
				Type:      aws.String(ssm.ParameterTypeSecureString),
				Overwrite: aws.Bool(overwrite),
				Value:     aws.String(param),
			}).
			Times(1)
	}
}
