package connect

import (
	"fmt"
	"testing"

	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	kc "github.com/razbomi/go-kafka-connect/lib/connectors"
)

func TestAccConnectorConfigUpdate(t *testing.T) {
	r.Test(t, r.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: testResourceConnector_initialConfig,
				Check:  testResourceConnector_initialCheck,
			},
			{
				Config: testResourceConnector_updateConfig,
				Check:  testResourceConnector_updateCheck,
			},
		},
	})
}

func testResourceConnector_initialCheck(s *terraform.State) error {
	resourceState := s.Modules[0].Resources["kafka-connect_connector.test"]
	if resourceState == nil {
		return fmt.Errorf("resource not found in state")
	}

	instanceState := resourceState.Primary
	if instanceState == nil {
		return fmt.Errorf("resource has no primary instance")
	}

	name := instanceState.ID

	if name != instanceState.Attributes["name"] {
		return fmt.Errorf("id doesn't match name")
	}

	client := testProvider.Meta().(kc.HighLevelClient)

	c, err := client.GetConnector(kc.ConnectorRequest{Name: "acc-test"})
	if err != nil {
		return err
	}

	flushSize := c.Config["flush.size"]
	expected := "10"
	if flushSize != expected {
		return fmt.Errorf("flush.size should be %s, got %s connector not updated. \n %v", expected, flushSize, c.Config)
	}

	return nil
}

func testResourceConnector_updateCheck(s *terraform.State) error {
	client := testProvider.Meta().(kc.HighLevelClient)

	c, err := client.GetConnector(kc.ConnectorRequest{Name: "acc-test"})
	if err != nil {
		return err
	}

	flushSize := c.Config["flush.size"]
	expected := "3"
	if flushSize != expected {
		return fmt.Errorf("flush.size should be %s, got %s connector not updated. \n %v", expected, flushSize, c.Config)
	}

	return nil
}

const testResourceConnector_initialConfig = `
resource "kafka-connect_connector" "test" {
	name = "acc-test"
	
	config = {
	  "name"                   = "acc-test"
	  "connector.class"        = "io.confluent.connect.s3.S3SinkConnector"
	  "partition.duration.ms"  = "3600000"
	  "s3.region"              = "ap-southeast-2"
	  "flush.size"             = "10"
	  "tasks.max"              = "1"
	  "timezone"               = "UTC"
	  "topics"                 = "mobile-platform.sap-transaction-events.transactions.v4"
	  "locale"                 = "US"
	  "format.class"           = "io.confluent.connect.s3.format.avro.AvroFormat"
	  "partitioner.class"      = "io.confluent.connect.storage.partitioner.TimeBasedPartitioner"
	  "schema.generator.class" = "io.confluent.connect.storage.hive.schema.DefaultSchemaGenerator"
	  "storage.class"          = "io.confluent.connect.s3.storage.S3Storage"
	  "s3.bucket.name"         = "xinja-data-nonprod-platform-events"
	  "path.format"            = "'year'=YYYY/'month'=MM/'day'=dd/'hour'=HH/"
	  "timestamp.extractor"    = "Record"
	}
  }
`

const testResourceConnector_updateConfig = `
resource "kafka-connect_connector" "test" {
	name = "acc-test"
	
	config = {
	  "name"                   = "acc-test"
	  "connector.class"        = "io.confluent.connect.s3.S3SinkConnector"
	  "partition.duration.ms"  = "3600000"
	  "s3.region"              = "ap-southeast-2"
	  "flush.size"             = "3"
	  "tasks.max"              = "1"
	  "timezone"               = "UTC"
	  "topics"                 = "mobile-platform.sap-transaction-events.transactions.v4"
	  "locale"                 = "US"
	  "format.class"           = "io.confluent.connect.s3.format.avro.AvroFormat"
	  "partitioner.class"      = "io.confluent.connect.storage.partitioner.TimeBasedPartitioner"
	  "schema.generator.class" = "io.confluent.connect.storage.hive.schema.DefaultSchemaGenerator"
	  "storage.class"          = "io.confluent.connect.s3.storage.S3Storage"
	  "s3.bucket.name"         = "xinja-data-nonprod-platform-events"
	  "path.format"            = "'year'=YYYY/'month'=MM/'day'=dd/'hour'=HH/"
	  "timestamp.extractor"    = "Record"
	}
  }
`
