package testing

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
)

func LoadEnvFromFile(filename string) {
	err := godotenv.Load(filename)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func BuildValidConfig(url string) []byte {
	return []byte(fmt.Sprintf(`version: 1.0
apiUrl: %v
apiUser: admin
apiPassword: admin
`, url))
}
