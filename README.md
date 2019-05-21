# `terraform-plugin-kafka-connect`

A [Terraform][1] plugin for managing [Apache Kafka Connect][2].

## Installation

Download and extract the [latest
release](https://github.com/razbomi/terraform-provider-kafka-connect/releases/latest) to
your [terraform plugin directory][third-party-plugins] (typically `~/.terraform.d/plugins/`)

## Example

Configure the provider directly by setting the following Env Variables:
*  `KAFKA_CONNECT_URL` 
*  `CLIENT_CERT` 
*  `CLIENT_KEY`
```hcl
provider "kafka-connect" {
  url = "https://kafka-connect.dev.moneynp.xinja.com.au"

}

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
	  "topics"                 = "platform-event.v1"
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
```

## Testing

Execute the following to run terraform provider acceptance tests.
```sh
# Sets the Kafka url
export URL=https://localhost:8083

# Params required to do mutual TLS
export CERT=client.pem
export KEY=client.key

# Runs terraform acceptance testing
make testacc URL=URL CERT=CERT KEY=KEY
```
[1]: https://www.terraform.io
[2]: https://kafka.apache.org/documentation/#connect
[third-party-plugins]: https://www.terraform.io/docs/configuration/providers.html#third-party-plugins
