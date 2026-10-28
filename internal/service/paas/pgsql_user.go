package paas

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/paas"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/service/paas/services"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func ResourcePostgreSQLUser() *schema.Resource {
	userSchema := map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"service_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"user_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	// password and roles share validation with the user block of aws_paas_service.
	for k, v := range services.PostgreSQL.UserParametersSchema() {
		userSchema[k] = v
	}

	return &schema.Resource{
		CreateContext: resourcePostgreSQLUserCreate,
		ReadContext:   resourcePostgreSQLUserRead,
		UpdateContext: resourcePostgreSQLUserUpdate,
		DeleteContext: resourcePostgreSQLUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		// Each operation switches the service to UPDATING for about a minute,
		// and operations on the same service are serialized.
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: userSchema,
	}
}

func resourcePostgreSQLUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn
	serviceID := d.Get("service_id").(string)
	name := d.Get("name").(string)

	service, err := findServiceByID(ctx, conn, serviceID)

	if err != nil {
		return diag.Errorf("error reading PaaS Service (%s): %s", serviceID, err)
	}

	if serviceType := aws.StringValue(service.ServiceType); serviceType != services.ServiceTypePostgreSQL {
		return diag.Errorf("PaaS Service (%s) has type %q, aws_paas_pgsql_user supports only %q services",
			serviceID, serviceType, services.ServiceTypePostgreSQL)
	}

	input := &paas.CreateUserInput{
		ServiceId: aws.String(serviceID),
		User: services.PostgreSQL.ExpandUser(map[string]interface{}{
			"name":     name,
			"password": d.Get("password").(string),
			"roles":    d.Get("roles").(*schema.Set),
		}, false),
	}

	lockPaaSService(serviceID)
	defer unlockPaaSService(serviceID)

	log.Printf("[DEBUG] Creating PaaS PostgreSQL User (%s) in Service (%s)", name, serviceID)
	outputRaw, err := tfresource.RetryWhenAWSErrCodeEqualsContext(ctx, d.Timeout(schema.TimeoutCreate), func() (interface{}, error) {
		return conn.CreateUserWithContext(ctx, input, request.WithLogLevel(aws.LogOff))
	}, OperationNotAllowedCode)

	if err != nil {
		return diag.Errorf("error creating PaaS PostgreSQL User (%s) in Service (%s): %s", name, serviceID, err)
	}

	output := outputRaw.(*paas.CreateUserOutput)
	if output == nil || output.User == nil || aws.StringValue(output.User.Id) == "" {
		return diag.Errorf("error creating PaaS PostgreSQL User (%s) in Service (%s): API returned no user ID", name, serviceID)
	}
	d.SetId(PostgreSQLUserCreateResourceID(serviceID, aws.StringValue(output.User.Id)))

	if _, err := waitServiceUpdated(ctx, conn, serviceID, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for PaaS Service (%s) to update: %s", serviceID, err)
	}

	return resourcePostgreSQLUserRead(ctx, d, meta)
}

func resourcePostgreSQLUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn

	serviceID, id, err := PostgreSQLUserParseResourceID(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	user, err := findPostgreSQLUserByID(ctx, conn, serviceID, id)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] PaaS PostgreSQL User (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("error reading PaaS PostgreSQL User (%s): %s", d.Id(), err)
	}

	d.Set("name", user.Name)
	d.Set("service_id", serviceID)
	d.Set("user_id", user.Id)

	// DescribeUser returns the password in plain text, so changes made outside of
	// Terraform are detected and the password is restored on import.
	if v, ok := user.Parameters["password"].(string); ok && v != "" {
		d.Set("password", v)
	}

	roles, _ := user.Parameters["roles"].([]interface{})
	if err := d.Set("roles", roles); err != nil {
		return diag.Errorf("error setting PaaS PostgreSQL User (%s) roles: %s", d.Id(), err)
	}

	return nil
}

func resourcePostgreSQLUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn

	serviceID, id, err := PostgreSQLUserParseResourceID(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges("password", "roles") {
		oldRoles, _ := d.GetChange("roles")

		// ModifyUser requires all current user parameters, not only the changed ones.
		input := &paas.ModifyUserInput{
			ServiceId: aws.String(serviceID),
			UserId:    aws.String(id),
			User: &paas.UserModifyRequest{
				Parameters: expandPostgreSQLUserModifyParameters(
					d.Get("password").(string),
					d.Get("roles").(*schema.Set),
					oldRoles.(*schema.Set),
				),
			},
		}

		lockPaaSService(serviceID)
		defer unlockPaaSService(serviceID)

		log.Printf("[DEBUG] Modifying PaaS PostgreSQL User (%s)", d.Id())
		_, err := tfresource.RetryWhenAWSErrCodeEqualsContext(ctx, d.Timeout(schema.TimeoutUpdate), func() (interface{}, error) {
			return conn.ModifyUserWithContext(ctx, input, request.WithLogLevel(aws.LogOff))
		}, OperationNotAllowedCode)

		if err != nil {
			return diag.Errorf("error modifying PaaS PostgreSQL User (%s): %s", d.Id(), err)
		}

		if _, err := waitServiceUpdated(ctx, conn, serviceID, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.Errorf("error waiting for PaaS Service (%s) to update: %s", serviceID, err)
		}
	}

	return resourcePostgreSQLUserRead(ctx, d, meta)
}

func resourcePostgreSQLUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn

	serviceID, id, err := PostgreSQLUserParseResourceID(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	input := &paas.DeleteUserInput{
		ServiceId: aws.String(serviceID),
		UserId:    aws.String(id),
	}

	lockPaaSService(serviceID)
	defer unlockPaaSService(serviceID)

	// The API rejects deleting a database owner with CannotDeleteUser, so the
	// error is returned as is and Terraform keeps the user in state.
	log.Printf("[DEBUG] Deleting PaaS PostgreSQL User (%s)", d.Id())
	_, err = tfresource.RetryWhenAWSErrCodeEqualsContext(ctx, d.Timeout(schema.TimeoutDelete), func() (interface{}, error) {
		return conn.DeleteUserWithContext(ctx, input)
	}, OperationNotAllowedCode)

	if tfawserr.ErrCodeEquals(err, ServiceNotFoundCode, ResourceDoesNotExistCode) {
		log.Printf("[WARN] PaaS PostgreSQL User (%s) not found, removing from state", d.Id())
		return nil
	}

	if err != nil {
		return diag.Errorf("error deleting PaaS PostgreSQL User (%s): %s", d.Id(), err)
	}

	if _, err := waitServiceUpdated(ctx, conn, serviceID, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for PaaS Service (%s) to update: %s", serviceID, err)
	}

	return nil
}

// expandPostgreSQLUserModifyParameters builds the full set of user parameters for ModifyUser.
// Older PaaS environments do not accept roles, so an empty roles array is sent only
// when previously assigned roles are being removed (see clearPostgreSQLUserRoles).
func expandPostgreSQLUserModifyParameters(password string, roles, oldRoles *schema.Set) map[string]interface{} {
	parameters := map[string]interface{}{
		"password": password,
	}

	switch {
	case roles != nil && roles.Len() > 0:
		parameters["roles"] = roles.List()
	case oldRoles != nil && oldRoles.Len() > 0:
		parameters["roles"] = []interface{}{}
	}

	return parameters
}

// lockPaaSService serializes user and database operations within one PaaS service,
// because each operation switches the service to UPDATING until it is applied.
func lockPaaSService(serviceID string) {
	conns.GlobalMutexKV.Lock(paasServiceMutexKey(serviceID))
}

func unlockPaaSService(serviceID string) {
	conns.GlobalMutexKV.Unlock(paasServiceMutexKey(serviceID))
}

func paasServiceMutexKey(serviceID string) string {
	return "paas-service-" + serviceID
}

const postgreSQLUserResourceIDSeparator = "/"

func PostgreSQLUserCreateResourceID(serviceID, userID string) string {
	parts := []string{serviceID, userID}
	id := strings.Join(parts, postgreSQLUserResourceIDSeparator)

	return id
}

func PostgreSQLUserParseResourceID(id string) (string, string, error) {
	parts := strings.Split(id, postgreSQLUserResourceIDSeparator)

	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unexpected format for ID (%[1]s), expected service-id%[2]suser-id", id, postgreSQLUserResourceIDSeparator)
}
