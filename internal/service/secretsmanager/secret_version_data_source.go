package secretsmanager

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func DataSourceSecretVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSecretVersionRead,

		Schema: map[string]*schema.Schema{
			"secret_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"secret_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"secret_binary": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"version_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"version_stage": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "CLOUDCURRENT",
			},
			"version_stages": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceSecretVersionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SecretsManagerConn
	secretID := d.Get("secret_id").(string)
	var version string

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretID),
	}

	if v, ok := d.GetOk("version_id"); ok {
		versionID := v.(string)
		input.VersionId = aws.String(versionID)
		version = versionID
	} else {
		versionStage := d.Get("version_stage").(string)
		input.VersionStage = aws.String(versionStage)
		version = versionStage
	}

	log.Printf("[DEBUG] Reading Secrets Manager Secret Version: %s", input)

	var output *secretsmanager.GetSecretValueOutput

	err := resource.Retry(PropagationTimeout, func() *resource.RetryError {
		var err error

		output, err = conn.GetSecretValue(input)

		if tfawserr.ErrCodeEquals(err, secretsmanager.ErrCodeResourceNotFoundException) {
			return resource.RetryableError(err)
		}

		if tfawserr.ErrMessageContains(err, secretsmanager.ErrCodeInvalidRequestException, "You can’t perform this operation on the secret because it was deleted") {
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		output, err = conn.GetSecretValue(input)
	}

	if err != nil {
		if tfawserr.ErrCodeEquals(err, secretsmanager.ErrCodeResourceNotFoundException) {
			return fmt.Errorf("Secrets Manager Secret %q Version %q not found", secretID, version)
		}
		if tfawserr.ErrMessageContains(err, secretsmanager.ErrCodeInvalidRequestException, "You can’t perform this operation on the secret because it was deleted") {
			return fmt.Errorf("Secrets Manager Secret %q Version %q not found", secretID, version)
		}
		return fmt.Errorf("error reading Secrets Manager Secret Version: %w", err)
	}

	d.SetId(fmt.Sprintf("%s|%s", secretID, version))
	d.Set("secret_id", secretID)
	d.Set("secret_string", output.SecretString)
	d.Set("version_id", output.VersionId)
	d.Set("secret_binary", string(output.SecretBinary))

	if err := d.Set("version_stages", flex.FlattenStringList(output.VersionStages)); err != nil {
		return fmt.Errorf("error setting version_stages: %w", err)
	}

	return nil
}
