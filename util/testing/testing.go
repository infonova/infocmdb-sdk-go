package testing

import (
	"fmt"
	"os"
)

var Url = "http://localhost"

func init() {
	if url := os.Getenv("WORKFLOW_TEST_URL"); url != "" {
		Url = url
	}
}

func BuildValidConfig(url string) []byte {
	return []byte(fmt.Sprintf(`version: 1.0
apiUrl: %v
apiUser: admin
apiPassword: admin
`, url))
}
