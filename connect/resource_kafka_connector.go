package connect

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	kc "github.com/razbomi/go-kafka-connect/lib/connectors"
	"log"
)

func kafkaConnectorResource() *schema.Resource {
	return &schema.Resource{
		Create: connectorCreate,
		Read:   connectorRead,
		Update: connectorUpdate,
		Delete: connectorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the connector",
			},
			"config": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    false,
				Description: "A map of string k/v attributes",
			},
		},
	}
}

func connectorCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(kc.HighLevelClient)
	name := nameFromRD(d)
	config := configFromRD(d)

	req := kc.CreateConnectorRequest{
		ConnectorRequest: kc.ConnectorRequest{
			Name: name,
		},
		Config: config,
	}

	connectorResponse, err := c.CreateConnector(req, true)
	fmt.Printf("[INFO] Created the connector %v\n", connectorResponse)

	if err == nil {
		d.SetId(name)
	}

	return err
}

func nameFromRD(d *schema.ResourceData) string {
	return d.Get("name").(string)
}

func connectorDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(kc.HighLevelClient)

	name := nameFromRD(d)
	req := kc.ConnectorRequest{
		Name: name,
	}

	fmt.Printf("[INFO] Deleing the connector %s\n", name)
	_, err := c.DeleteConnector(req, true)

	if err == nil {
		d.SetId("")
	}

	return err
}

func connectorUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(kc.HighLevelClient)

	name := nameFromRD(d)
	config := configFromRD(d)

	req := kc.CreateConnectorRequest{
		ConnectorRequest: kc.ConnectorRequest{
			Name: name,
		},
		Config: config,
	}

	log.Printf("[INFO] Looking for %s", name)
	conn, err := c.UpdateConnector(req, true)

	if err == nil {
		log.Printf("[INFO] Config updated %v", conn.Config)
		d.Set("config", conn.Config)
	}

	return err
}

func connectorRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(kc.HighLevelClient)

	name := d.Get("name").(string)
	req := kc.ConnectorRequest{
		Name: name,
	}
	log.Printf("[INFO] Looking for %s", name)
	conn, err := c.GetConnector(req)

	if err == nil {
		log.Printf("[INFO] found the config %v", conn.Config)
		d.Set("config", conn.Config)
	}

	return err
}

func configFromRD(d *schema.ResourceData) map[string]interface{} {
	return d.Get("config").(map[string]interface{})
}
