package utility

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
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

func ToPtr[T any](value T) *T {
	return &value
}

func ToPtrString(value types.String) *string {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	valueString := value.ValueString()

	return &valueString
}

func ToPtrInt64(value types.Int64) *int64 {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	valueInt64 := value.ValueInt64()

	return &valueInt64
}

func ToPtrBool(value types.Bool) *bool {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	valueBool := value.ValueBool()

	return &valueBool
}
