package utility

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func RunIntegrationTests() bool {
	return os.Getenv("RUN_INTEGRATION_TESTS") == "true"
}

func GetTestProvider() {

}

func GetTestConfig() string {
	runIntegrationTests := RunIntegrationTests()

	client_id := "mock-client-id"
	client_secret := "mock-client"

	if runIntegrationTests {
		client_id = os.Getenv("EDGIO_CLIENT_ID")
		client_secret = os.Getenv("EDGIO_CLIENT_SECRET")
	}

	return fmt.Sprintf(`
	provider "edgio" {
		client_id     = "%s"
		client_secret = "%s"
	}`, client_id, client_secret)
}

func RandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteRune(letters[rand.Intn(len(letters))])
	}
	return sb.String()
}
