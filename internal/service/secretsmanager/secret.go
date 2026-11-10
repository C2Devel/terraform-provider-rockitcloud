package secretsmanager

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func ResourceSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecretCreate,
		Read:   resourceSecretRead,
		Update: resourceSecretUpdate,
		Delete: resourceSecretDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"secret_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name_prefix"},
				ValidateFunc:  validSecretName,
			},
			"name_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
				ValidateFunc:  validSecretNamePrefix,
			},
			"recovery_window_in_days": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
				ValidateFunc: validation.Any(
					validation.IntBetween(7, 30),
					validation.IntInSlice([]int{0}),
				),
			},
		},
	}
}

func resourceSecretCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SecretsManagerConn

	secretName := create.Name(d.Get("name").(string), d.Get("name_prefix").(string))
	input := &secretsmanager.CreateSecretInput{
		Description: aws.String(d.Get("description").(string)),
		Name:        aws.String(secretName),
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		input.KmsKeyId = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Creating Secrets Manager Secret: %s", input)

	// Retry for secret recreation after deletion
	var output *secretsmanager.CreateSecretOutput
	err := resource.Retry(PropagationTimeout, func() *resource.RetryError {
		var err error
		output, err = conn.CreateSecret(input)
		// Temporarily retry on these errors to support immediate secret recreation:
		// InvalidRequestException: You can’t perform this operation on the secret because it was deleted.
		// InvalidRequestException: You can't create this secret because a secret with this name is already scheduled for deletion.
		if tfawserr.ErrMessageContains(err, secretsmanager.ErrCodeInvalidRequestException, "scheduled for deletion") || tfawserr.ErrMessageContains(err, secretsmanager.ErrCodeInvalidRequestException, "was deleted") {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if tfresource.TimedOut(err) {
		output, err = conn.CreateSecret(input)
	}
	if err != nil {
		return fmt.Errorf("error creating Secrets Manager Secret: %w", err)
	}

	// sem returns SecretId (UUID) instead of ARN; use it as the resource ID.
	d.SetId(aws.StringValue(output.SecretId))
	d.Set("secret_id", output.SecretId)

	return resourceSecretRead(d, meta)
}

func resourceSecretRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SecretsManagerConn

	outputRaw, err := tfresource.RetryWhenNewResourceNotFound(PropagationTimeout, func() (interface{}, error) {
		return FindSecretByID(conn, d.Id())
	}, d.IsNewResource())

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, secretsmanager.ErrCodeResourceNotFoundException) {
		log.Printf("[WARN] Secrets Manager Secret (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("reading Secrets Manager Secret (%s): %w", d.Id(), err)
	}

	output := outputRaw.(*secretsmanager.DescribeSecretOutput)

	d.Set("description", output.Description)
	d.Set("kms_key_id", output.KmsKeyId)
	d.Set("name", output.Name)
	d.Set("name_prefix", create.NamePrefixFromName(aws.StringValue(output.Name)))
	d.Set("secret_id", output.SecretId)

	return nil
}

func resourceSecretUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SecretsManagerConn

	if d.HasChanges("description", "kms_key_id") {
		input := &secretsmanager.UpdateSecretInput{
			Description: aws.String(d.Get("description").(string)),
			SecretId:    aws.String(d.Id()),
		}

		if v, ok := d.GetOk("kms_key_id"); ok {
			input.KmsKeyId = aws.String(v.(string))
		}

		log.Printf("[DEBUG] Updating Secrets Manager Secret: %s", input)
		_, err := conn.UpdateSecret(input)

		if err != nil {
			return fmt.Errorf("updating Secrets Manager Secret (%s): %w", d.Id(), err)
		}
	}

	return resourceSecretRead(d, meta)
}

func resourceSecretDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SecretsManagerConn

	input := &secretsmanager.DeleteSecretInput{
		SecretId: aws.String(d.Id()),
	}

	recoveryWindowInDays := d.Get("recovery_window_in_days").(int)
	if recoveryWindowInDays == 0 {
		input.ForceDeleteWithoutRecovery = aws.Bool(true)
	} else {
		input.RecoveryWindowInDays = aws.Int64(int64(recoveryWindowInDays))
	}

	log.Printf("[DEBUG] Deleting Secrets Manager Secret: %s", d.Id())
	_, err := conn.DeleteSecret(input)

	if tfawserr.ErrCodeEquals(err, secretsmanager.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("deleting Secrets Manager Secret (%s): %w", d.Id(), err)
	}

	_, err = tfresource.RetryUntilNotFound(PropagationTimeout, func() (interface{}, error) {
		return FindSecretByID(conn, d.Id())
	})

	if err != nil {
		return fmt.Errorf("waiting for Secrets Manager Secret (%s) delete: %w", d.Id(), err)
	}

	return nil
}

func findSecret(conn *secretsmanager.SecretsManager, input *secretsmanager.DescribeSecretInput) (*secretsmanager.DescribeSecretOutput, error) {
	output, err := conn.DescribeSecret(input)

	if tfawserr.ErrCodeEquals(err, secretsmanager.ErrCodeResourceNotFoundException) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}

func FindSecretByID(conn *secretsmanager.SecretsManager, id string) (*secretsmanager.DescribeSecretOutput, error) {
	input := &secretsmanager.DescribeSecretInput{
		SecretId: aws.String(id),
	}

	output, err := findSecret(conn, input)

	if err != nil {
		return nil, err
	}

	if output.DeletedDate != nil {
		return nil, &resource.NotFoundError{LastRequest: input}
	}

	return output, nil
}
