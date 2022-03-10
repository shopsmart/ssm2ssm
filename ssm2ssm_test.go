package ssm2ssm_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/golang/mock/gomock"

	"github.com/shopsmart/ssm2ssm"
	"github.com/shopsmart/ssm2ssm/internal/testutils"
	"github.com/shopsmart/ssm2ssm/pkg/service"
)

const (
	InputPath  = "/aws/service/global-infrastructure/regions"
	OutputPath = "/test/ssm2ssm"
)

var _ = Describe("Ssm2ssm", func() {

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

	It("Should copy the parameters from input to output", func() {
		testutils.MockGetParametersByPath(client, InputPath, nil)
		testutils.MockGetParametersByPath(client, OutputPath, []*ssm.Parameter{})
		testutils.ExpectAllParametersToBePut(client, OutputPath, false)

		err = ssm2ssm.Copy(svc, InputPath, OutputPath, false)
		Expect(err).Should(BeNil())
	})

	It("Should skip over parameters that already exist", func() {
		testutils.MockGetParametersByPath(client, InputPath, nil)
		testutils.MockGetParametersByPath(client, OutputPath, nil)

		client.
			EXPECT().
			PutParameter(gomock.Any()).
			Times(0)

		err = ssm2ssm.Copy(svc, InputPath, OutputPath, false)
		Expect(err).Should(BeNil())
	})

	It("Should overwrite the parameters currently under output", func() {
		testutils.MockGetParametersByPath(client, InputPath, nil)
		testutils.ExpectAllParametersToBePut(client, OutputPath, true)

		err = ssm2ssm.Copy(svc, InputPath, OutputPath, true)
		Expect(err).Should(BeNil())
	})

})
