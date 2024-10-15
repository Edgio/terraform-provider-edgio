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

	// Create a new random source
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteRune(letters[r.Intn(len(letters))])
	}
	return sb.String()
}

// The following are functions that return pointers to primitive types as this
// seems like a major flaw of go's type system.

// PtrString returns a pointer to a string value.
func PtrString(s string) *string {
	return &s
}

// PtrInt returns a pointer to an int value.
func PtrInt(i int) *int {
	return &i
}

// PtrFloat64 returns a pointer to a float64 value.
func PtrFloat64(f float64) *float64 {
	return &f
}

// PtrBool returns a pointer to a bool value.
func PtrBool(b bool) *bool {
	return &b
}
