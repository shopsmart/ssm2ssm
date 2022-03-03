package service_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/shopsmart/ssm2ssm/internal/testutils"
)

func TestService(t *testing.T) {
	testutils.Setup(t)
	defer testutils.Teardown()

	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Suite")
}
