package kong

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kevholditch/gokong"
)

func resourceKongSni() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongSniCreate,
		Read:   resourceKongSniRead,
		Delete: resourceKongSniDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"certificate_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceKongSniCreate(d *schema.ResourceData, meta interface{}) error {

	sniRequest := createKongSniRequestFromResourceData(d)

	sni, err := meta.(*gokong.KongAdminClient).Snis().Create(sniRequest)

	if err != nil {
		return fmt.Errorf("failed to create kong sni: %v error: %v", sniRequest, err)
	}

	d.SetId(sni.Name)

	return resourceKongSniRead(d, meta)
}

func resourceKongSniRead(d *schema.ResourceData, meta interface{}) error {

	sni, err := meta.(*gokong.KongAdminClient).Snis().GetByName(d.Id())

	if err != nil {
		return fmt.Errorf("could not find kong sni: %v", err)
	}

	d.Set("name", sni.Name)
	d.Set("certificate_id", sni.SslCertificateId)

	return nil
}

func resourceKongSniDelete(d *schema.ResourceData, meta interface{}) error {

	err := meta.(*gokong.KongAdminClient).Snis().DeleteByName(d.Id())

	if err != nil {
		return fmt.Errorf("could not delete kong sni: %v", err)
	}

	return nil
}

func createKongSniRequestFromResourceData(d *schema.ResourceData) *gokong.SnisRequest {

	sniRequest := &gokong.SnisRequest{}

	sniRequest.Name = readStringFromResource(d, "name")
	sniRequest.SslCertificateId = readStringFromResource(d, "certificate_id")

	return sniRequest
}
