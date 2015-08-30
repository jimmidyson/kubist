package kubist

import (
	"fmt"

	"github.com/fabric8io/kubist/api"
)

type KubernetesError struct {
	*api.Status
	StatusCode int
	StatusText string
	Response   string
}

func (e *KubernetesError) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.StatusText)
}
