package apiserver

import (
	"testing"

	apitesting "k8s.io/apimachinery/pkg/api/testing"
)

func TestRoundTripTypes(t *testing.T) {
	apitesting.RoundTripTestForScheme(t, Scheme, nil)
}
