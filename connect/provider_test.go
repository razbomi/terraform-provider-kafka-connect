package connect

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testProvider *schema.Provider
var testProviders map[string]terraform.ResourceProvider

var urlEnvVars = []string{
	"KAFKA_CONNECT_URL",
}

var certEnvVars = []string{
	"KAFKA_CLIENT_CERT",
}

var keyEnvVars = []string{
	"KAFKA_CLIENT_KEY",
}

func init() {
	testProvider = Provider().(*schema.Provider)
	testProviders = map[string]terraform.ResourceProvider{
		"kafka-connect": testProvider,
	}

}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	client := testProvider.Meta()
	log.Printf("[INFO] Checking KafkaConnect client")
	if client == nil {
		//t.Fatal("No client")
	}

	if v := multiEnvSearch(urlEnvVars); v != "https://kafka-connect.dev.moneynp.xinja.com.au" {
		t.Fatalf("One of %s must be set to https://kafka-connect.dev.moneynp.xinja.com.au for acceptance tests", strings.Join(urlEnvVars, ", "))
	}

	if v := multiEnvSearch(certEnvVars); v != "/Users/stuart/Documents/2018/Xinja/Development/kafka/keys/stulox.dev.pem" {
		t.Fatalf("One of %s must be set to cert for acceptance tests", strings.Join(certEnvVars, ", "))
	}

	if v := multiEnvSearch(keyEnvVars); v != "/Users/stuart/Documents/2018/Xinja/Development/kafka/keys/stulox.key" {
		t.Fatalf("One of %s must be set to key for acceptance tests", strings.Join(keyEnvVars, ", "))
	}
}

func multiEnvSearch(ks []string) string {
	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			log.Printf("[INFO] %s", v)
			return v
		}
	}
	log.Printf("[INFO] Nothing to Return")
	return ""
}
