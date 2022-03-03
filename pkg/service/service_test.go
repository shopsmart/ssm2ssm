package service_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/golang/mock/gomock"

	"github.com/shopsmart/ssm2ssm/internal/testutils"
	"github.com/shopsmart/ssm2ssm/pkg/service"
)

const (
	inputPath  = "/test/ssm2ssm"
	outputPath = "/test/ssm2ssm-new"
)

var _ = Describe("Service", func() {

	var (
		client *testutils.MockSSMClient
		svc    service.Service
		err    error
	)

	BeforeEach(func() {
		client = testutils.NewMockSSMClient(testutils.MockController)
		svc, err = service.NewFromClient(client)
		Expect(err).Should(BeNil())
	})

	It("Should get all parameters from path", func() {
		testutils.MockGetParametersByPath(client, inputPath, nil)

		params, err := svc.GetParameters(inputPath)
		Expect(err).Should(BeNil())
		Expect(params).Should(Equal(testutils.ParametersMap))
	})

	It("Should copy all params to the new path", func() {
		testutils.MockGetParametersByPath(client, inputPath, []*ssm.Parameter{})
		testutils.ExpectAllParametersToBePut(client, outputPath, false)

		err := svc.PutParameters(outputPath, testutils.ParametersMap, false)
		Expect(err).Should(BeNil())
	})

	It("Should cowardly refuse to overwrite existing params", func() {
		testutils.MockGetParametersByPath(client, outputPath, nil)

		err := svc.PutParameters(outputPath, testutils.ParametersMap, false)
		Expect(err).ShouldNot(BeNil())
	})

	It("Should overwrite existing params", func() {
		client.
			EXPECT().
			GetParametersByPathPages(gomock.Any(), gomock.Any()).
			Times(0)

		testutils.ExpectAllParametersToBePut(client, outputPath, true)

		err := svc.PutParameters(outputPath, testutils.ParametersMap, true)
		Expect(err).Should(BeNil())
	})

})
