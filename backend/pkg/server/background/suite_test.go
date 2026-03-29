package background

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBackground(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Background Suite")
}
