//go:build kind

package kind_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestKindE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kind E2E Suite")
}
