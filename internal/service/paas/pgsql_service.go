package paas

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/service/paas/services"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// ResourcePostgreSQLService manages a PaaS PostgreSQL service (aws_paas_pgsql).
//
// It is a PostgreSQL-only view of aws_paas_service: the schema is derived from
// aws_paas_service, so that new common arguments appear in both resources, and the
// CRUD functions of aws_paas_service are reused. Users and databases are always
// managed by the aws_paas_pgsql_user and aws_paas_pgsql_database resources.
func ResourcePostgreSQLService() *schema.Resource {
	base := ResourceService()

	serviceSchema := make(map[string]*schema.Schema, len(base.Schema))
	for k, v := range base.Schema {
		serviceSchema[k] = v
	}

	// Blocks of other service types are not supported.
	for _, serviceType := range services.ManagedServiceTypes() {
		delete(serviceSchema, serviceType)
	}

	// The pgsql block is required and has no inline user and database blocks.
	pgsql := services.PostgreSQL.ResourceSchema()
	pgsql.Optional = false
	pgsql.Required = true
	pgsql.ExactlyOneOf = nil
	pgsqlElem := pgsql.Elem.(*schema.Resource)
	delete(pgsqlElem.Schema, "user")
	delete(pgsqlElem.Schema, "database")
	serviceSchema[services.ServiceTypePostgreSQL] = pgsql

	// Users and databases are always managed externally. The shared CRUD functions
	// read these flags, so they are kept as computed attributes set to true.
	for _, k := range []string{"users_managed_externally", "databases_managed_externally"} {
		serviceSchema[k] = &schema.Schema{
			Type:     schema.TypeBool,
			Computed: true,
		}
	}

	// These ClickHouse and Kafka arguments are read back by the shared Read function,
	// so they are kept as read-only attributes.
	for _, k := range []string{"additional_roles", "coordinator"} {
		serviceSchema[k] = computedOnlySchema(base.Schema[k])
	}

	return &schema.Resource{
		CreateContext: resourcePostgreSQLServiceCreate,
		ReadContext:   resourcePostgreSQLServiceRead,
		UpdateContext: resourcePostgreSQLServiceUpdate,
		DeleteContext: resourceServiceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: base.Timeouts,

		// The aws_paas_service validations of other service types read their blocks,
		// which are absent in this schema, so only the applicable validations are used.
		CustomizeDiff: validatePostgreSQLServiceResourceConfiguration,

		Schema: serviceSchema,
	}
}

func validatePostgreSQLServiceResourceConfiguration(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	if err := validatePostgreSQLReplicationMode(d); err != nil {
		return err
	}

	return validateAutoDiskExpansionConfiguration(ctx, d, meta)
}

// validatePostgreSQLReplicationMode rejects a high availability service without
// replication_mode before the cloud rejects the CreateService request.
func validatePostgreSQLReplicationMode(d *schema.ResourceDiff) error {
	const replicationModeKey = "pgsql.0.replication_mode"
	if !d.NewValueKnown("high_availability") || !d.NewValueKnown(replicationModeKey) {
		return nil
	}

	if d.Get("high_availability").(bool) && d.Get(replicationModeKey).(string) == "" {
		return fmt.Errorf("%s must be set if high_availability is true", replicationModeKey)
	}

	return nil
}

func resourcePostgreSQLServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	setPostgreSQLServiceExternalManagement(d)

	return resourceServiceCreate(ctx, d, meta)
}

func resourcePostgreSQLServiceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// On import the service type is not known yet. The shared Read function expects the
	// block of the read service type in the schema, so other types are rejected first.
	importing := d.Get("service_type").(string) == ""
	if importing {
		conn := meta.(*conns.AWSClient).PaaSConn

		service, err := findServiceByID(ctx, conn, d.Id())

		if !d.IsNewResource() && tfresource.NotFound(err) {
			log.Printf("[WARN] PaaS Service (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}

		if err != nil {
			return diag.Errorf("error reading PaaS Service (%s): %s", d.Id(), err)
		}

		if serviceType := aws.StringValue(service.ServiceType); serviceType != services.ServiceTypePostgreSQL {
			return diag.Errorf("PaaS Service (%s) has type %q, aws_paas_pgsql supports only %q services, use aws_paas_service instead",
				d.Id(), serviceType, services.ServiceTypePostgreSQL)
		}
	}

	setPostgreSQLServiceExternalManagement(d)

	diags := resourceServiceRead(ctx, d, meta)
	if diags.HasError() || d.Id() == "" {
		return diags
	}

	// DescribeService does not return arbitratorRequired, so after import the attribute is
	// missing and its ForceNew default would replace the service. The value is restored
	// from the instance roles when it can be determined, and is left unset otherwise.
	if importing {
		if arbitratorRequired, ok := postgreSQLServiceArbitratorRequired(d); ok {
			if err := d.Set("arbitrator_required", arbitratorRequired); err != nil {
				return diag.Errorf("error setting arbitrator_required: %s", err)
			}
		}
	}

	return diags
}

// postgreSQLServiceArbitratorRequired reports whether the service has an arbitrator, and
// whether this can be determined. An arbitrator can only be used with high availability,
// and it is listed among the service instances with the arbitrator role.
func postgreSQLServiceArbitratorRequired(d *schema.ResourceData) (bool, bool) {
	if !d.Get("high_availability").(bool) {
		return false, true
	}

	instances := d.Get("instances").(*schema.Set).List()
	for _, instance := range instances {
		if instance, ok := instance.(map[string]interface{}); ok && instance["role"] == postgreSQLArbitratorRole {
			return true, true
		}
	}

	// Without instances, for example while the service is being created, the value is unknown.
	return false, len(instances) > 0
}

const postgreSQLArbitratorRole = "arbitrator"

func resourcePostgreSQLServiceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	setPostgreSQLServiceExternalManagement(d)

	return resourceServiceUpdate(ctx, d, meta)
}

// setPostgreSQLServiceExternalManagement makes the shared CRUD functions treat users and
// databases as managed by the aws_paas_pgsql_user and aws_paas_pgsql_database resources.
func setPostgreSQLServiceExternalManagement(d *schema.ResourceData) {
	d.Set("users_managed_externally", true)
	d.Set("databases_managed_externally", true)
}

// computedOnlySchema returns a read-only copy of the schema with read-only nested attributes.
func computedOnlySchema(s *schema.Schema) *schema.Schema {
	computed := &schema.Schema{
		Type:     s.Type,
		Computed: true,
	}

	switch elem := s.Elem.(type) {
	case *schema.Resource:
		nested := make(map[string]*schema.Schema, len(elem.Schema))
		for k, v := range elem.Schema {
			nested[k] = computedOnlySchema(v)
		}
		computed.Elem = &schema.Resource{Schema: nested}
	case *schema.Schema:
		computed.Elem = &schema.Schema{Type: elem.Type}
	}

	return computed
}
