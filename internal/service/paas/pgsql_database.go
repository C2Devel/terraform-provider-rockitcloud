package paas

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/paas"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/service/paas/services"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func ResourcePostgreSQLDatabase() *schema.Resource {
	databaseSchema := map[string]*schema.Schema{
		"backup_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"database_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
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
		"users": {
			Type:     schema.TypeSet,
			Optional: true,
			MaxItems: 1000,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}

	// owner, extensions, encoding, locale, backup_id and backup_db_name share validation
	// with the database block of aws_paas_service.
	for k, v := range services.PostgreSQL.DatabaseParametersSchema() {
		databaseSchema[k] = v
	}

	// Encoding and locale cannot be changed in an existing database,
	// and a database can be restored from a backup only on creation.
	for _, k := range []string{"backup_db_name", "backup_id", "encoding", "locale"} {
		databaseSchema[k].ForceNew = true
	}
	for _, k := range []string{"encoding", "locale"} {
		databaseSchema[k].Computed = true
	}
	// A case-only difference in encoding (e.g. utf8 and UTF8) must not recreate the database.
	databaseSchema["encoding"].DiffSuppressFunc = func(_, old, new string, _ *schema.ResourceData) bool {
		return strings.EqualFold(old, new)
	}

	return &schema.Resource{
		CreateContext: resourcePostgreSQLDatabaseCreate,
		ReadContext:   resourcePostgreSQLDatabaseRead,
		UpdateContext: resourcePostgreSQLDatabaseUpdate,
		DeleteContext: resourcePostgreSQLDatabaseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		// Each operation switches the service to UPDATING for about a minute, and operations
		// on users and databases of the same service are serialized.
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: validatePostgreSQLDatabaseOwnerNotInUsers,

		Schema: databaseSchema,
	}
}

// validatePostgreSQLDatabaseOwnerNotInUsers rejects the owner listed in users at plan time,
// because the API returns DatabaseOwnerCannotBeDatabaseUser for such a request.
func validatePostgreSQLDatabaseOwnerNotInUsers(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	if !d.NewValueKnown("owner") || !d.NewValueKnown("users") {
		return nil
	}

	owner := d.Get("owner").(string)
	if users, ok := d.Get("users").(*schema.Set); ok && users.Contains(owner) {
		return fmt.Errorf("users must not contain the database owner %q: the owner already has access to the database", owner)
	}

	return nil
}

func resourcePostgreSQLDatabaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn
	serviceID := d.Get("service_id").(string)
	name := d.Get("name").(string)

	service, err := findServiceByID(ctx, conn, serviceID)

	if err != nil {
		return diag.Errorf("error reading PaaS Service (%s): %s", serviceID, err)
	}

	if serviceType := aws.StringValue(service.ServiceType); serviceType != services.ServiceTypePostgreSQL {
		return diag.Errorf("PaaS Service (%s) has type %q, aws_paas_pgsql_database supports only %q services",
			serviceID, serviceType, services.ServiceTypePostgreSQL)
	}

	input := &paas.CreateDatabaseInput{
		ServiceId: aws.String(serviceID),
		Database: services.PostgreSQL.ExpandDatabase(map[string]interface{}{
			"backup_db_name": d.Get("backup_db_name").(string),
			"backup_enabled": d.Get("backup_enabled").(bool),
			"backup_id":      d.Get("backup_id").(string),
			"encoding":       d.Get("encoding").(string),
			"extensions":     d.Get("extensions").(*schema.Set),
			"locale":         d.Get("locale").(string),
			"name":           name,
			"owner":          d.Get("owner").(string),
			"user":           expandPostgreSQLDatabaseUserNames(d.Get("users").(*schema.Set)),
		}),
	}

	lockPaaSService(serviceID)
	defer unlockPaaSService(serviceID)

	log.Printf("[DEBUG] Creating PaaS PostgreSQL Database (%s) in Service (%s)", name, serviceID)
	outputRaw, err := tfresource.RetryWhenAWSErrCodeEqualsContext(ctx, d.Timeout(schema.TimeoutCreate), func() (interface{}, error) {
		return conn.CreateDatabaseWithContext(ctx, input)
	}, OperationNotAllowedCode)

	if err != nil {
		return diag.Errorf("error creating PaaS PostgreSQL Database (%s) in Service (%s): %s", name, serviceID, err)
	}

	output := outputRaw.(*paas.CreateDatabaseOutput)
	if output == nil || output.Database == nil || aws.StringValue(output.Database.Id) == "" {
		return diag.Errorf("error creating PaaS PostgreSQL Database (%s) in Service (%s): API returned no database ID", name, serviceID)
	}
	d.SetId(PostgreSQLDatabaseCreateResourceID(serviceID, aws.StringValue(output.Database.Id)))

	if _, err := waitServiceUpdated(ctx, conn, serviceID, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for PaaS Service (%s) to update: %s", serviceID, err)
	}

	return resourcePostgreSQLDatabaseRead(ctx, d, meta)
}

func resourcePostgreSQLDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn

	serviceID, id, err := PostgreSQLDatabaseParseResourceID(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	database, err := findPostgreSQLDatabaseByID(ctx, conn, serviceID, id)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] PaaS PostgreSQL Database (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("error reading PaaS PostgreSQL Database (%s): %s", d.Id(), err)
	}

	d.Set("backup_enabled", database.BackupEnabled)
	d.Set("database_id", database.Id)
	d.Set("name", database.Name)
	d.Set("service_id", serviceID)

	parameters := database.Parameters
	if v, ok := parameters["owner"].(string); ok {
		d.Set("owner", v)
	}

	// Encoding and locale are set only when the API returns them,
	// so that the values chosen by the cloud are kept as computed.
	for _, key := range []string{"encoding", "locale"} {
		if v, ok := parameters[key].(string); ok && v != "" {
			d.Set(key, v)
		}
	}

	// Backup parameters are used only on creation and are kept from the configuration
	// unless the API returns them.
	for _, keys := range [][2]string{{"backupId", "backup_id"}, {"backupDbName", "backup_db_name"}} {
		for _, key := range keys {
			if v, ok := parameters[key].(string); ok && v != "" {
				d.Set(keys[1], v)
			}
		}
	}

	extensions, _ := parameters["extensions"].([]interface{})
	if err := d.Set("extensions", extensions); err != nil {
		return diag.Errorf("error setting PaaS PostgreSQL Database (%s) extensions: %s", d.Id(), err)
	}

	users := make([]interface{}, 0, len(database.Users))
	for _, user := range database.Users {
		if user != nil && aws.StringValue(user.Name) != "" {
			users = append(users, aws.StringValue(user.Name))
		}
	}
	if err := d.Set("users", users); err != nil {
		return diag.Errorf("error setting PaaS PostgreSQL Database (%s) users: %s", d.Id(), err)
	}

	return nil
}

func resourcePostgreSQLDatabaseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn

	serviceID, id, err := PostgreSQLDatabaseParseResourceID(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges("backup_enabled", "extensions", "owner", "users") {
		oldExtensions, _ := d.GetChange("extensions")
		oldUsers, _ := d.GetChange("users")

		// ModifyDatabase requires all current database parameters, not only the changed ones.
		input := &paas.ModifyDatabaseInput{
			DatabaseId: aws.String(id),
			ServiceId:  aws.String(serviceID),
			Database: &paas.DatabaseModifyRequest{
				BackupEnabled: aws.Bool(d.Get("backup_enabled").(bool)),
				Parameters: expandPostgreSQLDatabaseModifyParameters(
					d.Get("owner").(string),
					d.Get("encoding").(string),
					d.Get("locale").(string),
					d.Get("extensions").(*schema.Set),
					oldExtensions.(*schema.Set),
				),
				Users: expandPostgreSQLDatabaseModifyUsers(
					d.Get("users").(*schema.Set),
					oldUsers.(*schema.Set),
				),
			},
		}

		lockPaaSService(serviceID)
		defer unlockPaaSService(serviceID)

		log.Printf("[DEBUG] Modifying PaaS PostgreSQL Database (%s)", d.Id())
		_, err := tfresource.RetryWhenAWSErrCodeEqualsContext(ctx, d.Timeout(schema.TimeoutUpdate), func() (interface{}, error) {
			return conn.ModifyDatabaseWithContext(ctx, input)
		}, OperationNotAllowedCode)

		if err != nil {
			return diag.Errorf("error modifying PaaS PostgreSQL Database (%s): %s", d.Id(), err)
		}

		if _, err := waitServiceUpdated(ctx, conn, serviceID, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.Errorf("error waiting for PaaS Service (%s) to update: %s", serviceID, err)
		}
	}

	return resourcePostgreSQLDatabaseRead(ctx, d, meta)
}

func resourcePostgreSQLDatabaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn

	serviceID, id, err := PostgreSQLDatabaseParseResourceID(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	input := &paas.DeleteDatabaseInput{
		DatabaseId: aws.String(id),
		ServiceId:  aws.String(serviceID),
	}

	lockPaaSService(serviceID)
	defer unlockPaaSService(serviceID)

	log.Printf("[DEBUG] Deleting PaaS PostgreSQL Database (%s)", d.Id())
	_, err = tfresource.RetryWhenAWSErrCodeEqualsContext(ctx, d.Timeout(schema.TimeoutDelete), func() (interface{}, error) {
		return conn.DeleteDatabaseWithContext(ctx, input)
	}, OperationNotAllowedCode)

	if tfawserr.ErrCodeEquals(err, ServiceNotFoundCode, ResourceDoesNotExistCode) {
		log.Printf("[WARN] PaaS PostgreSQL Database (%s) not found, removing from state", d.Id())
		return nil
	}

	if err != nil {
		return diag.Errorf("error deleting PaaS PostgreSQL Database (%s): %s", d.Id(), err)
	}

	if _, err := waitServiceUpdated(ctx, conn, serviceID, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for PaaS Service (%s) to update: %s", serviceID, err)
	}

	return nil
}

// expandPostgreSQLDatabaseUserNames converts the set of user names to the list of
// database user blocks accepted by ExpandDatabase.
func expandPostgreSQLDatabaseUserNames(users *schema.Set) []interface{} {
	if users == nil || users.Len() == 0 {
		return nil
	}

	tfList := make([]interface{}, 0, users.Len())
	for _, name := range users.List() {
		tfList = append(tfList, map[string]interface{}{"name": name.(string)})
	}

	return tfList
}

// expandPostgreSQLDatabaseModifyParameters builds the full set of database parameters for
// ModifyDatabase. An empty extensions array is sent only when previously installed
// extensions are being removed.
func expandPostgreSQLDatabaseModifyParameters(owner, encoding, locale string, extensions, oldExtensions *schema.Set) map[string]interface{} {
	parameters := map[string]interface{}{
		"owner": owner,
	}

	if encoding != "" {
		parameters["encoding"] = encoding
	}

	if locale != "" {
		parameters["locale"] = locale
	}

	switch {
	case extensions != nil && extensions.Len() > 0:
		parameters["extensions"] = extensions.List()
	case oldExtensions != nil && oldExtensions.Len() > 0:
		parameters["extensions"] = []interface{}{}
	}

	return parameters
}

// expandPostgreSQLDatabaseModifyUsers builds the full list of database users for ModifyDatabase.
// An empty list is sent only when access of all previous users is being revoked.
func expandPostgreSQLDatabaseModifyUsers(users, oldUsers *schema.Set) []*paas.UserCreateRequest {
	if users != nil && users.Len() > 0 {
		apiObjects := make([]*paas.UserCreateRequest, 0, users.Len())
		for _, name := range users.List() {
			apiObjects = append(apiObjects, &paas.UserCreateRequest{Name: aws.String(name.(string))})
		}
		return apiObjects
	}

	if oldUsers != nil && oldUsers.Len() > 0 {
		return []*paas.UserCreateRequest{}
	}

	return nil
}

const postgreSQLDatabaseResourceIDSeparator = "/"

func PostgreSQLDatabaseCreateResourceID(serviceID, databaseID string) string {
	parts := []string{serviceID, databaseID}
	id := strings.Join(parts, postgreSQLDatabaseResourceIDSeparator)

	return id
}

func PostgreSQLDatabaseParseResourceID(id string) (string, string, error) {
	parts := strings.Split(id, postgreSQLDatabaseResourceIDSeparator)

	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unexpected format for ID (%[1]s), expected service-id%[2]sdatabase-id", id, postgreSQLDatabaseResourceIDSeparator)
}
