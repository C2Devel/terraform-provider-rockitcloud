package paas

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/paas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/service/paas/services"
)

func validatePostgreSQLServiceConfiguration(d *schema.ResourceDiff) error {
	const userKey = "pgsql.0.user"
	if _, ok := d.GetOk("pgsql"); !ok || !d.NewValueKnown(userKey) || !d.NewValueKnown("pgsql.0.version") {
		return nil
	}
	major, _ := strconv.Atoi(strings.Split(d.Get("pgsql.0.version").(string), ".")[0])
	for index, user := range d.Get(userKey).([]interface{}) {
		roles := user.(map[string]interface{})["roles"].(*schema.Set)
		for role, minimum := range map[string]int{"pg_create_subscription": 16, "pg_maintain": 17} {
			if roles.Contains(role) && major < minimum {
				return fmt.Errorf("%s.%d.roles: %s requires PostgreSQL %d or later", userKey, index, role, minimum)
			}
		}
	}
	return nil
}

// serviceUsersManagedExternally reports whether users of the service are managed
// by standalone aws_paas_pgsql_user resources instead of the user blocks.
func serviceUsersManagedExternally(d *schema.ResourceData, serviceType string) bool {
	return serviceType == services.ServiceTypePostgreSQL && d.Get("users_managed_externally").(bool)
}

// serviceDatabasesManagedExternally reports whether databases of the service are managed
// by standalone aws_paas_pgsql_database resources instead of the database blocks.
func serviceDatabasesManagedExternally(d *schema.ResourceData, serviceType string) bool {
	return serviceType == services.ServiceTypePostgreSQL && d.Get("databases_managed_externally").(bool)
}

// validateManagedExternally rejects configurations where the same users or databases would be
// managed both by the service and by standalone resources.
func validateManagedExternally(d *schema.ResourceDiff) error {
	for _, v := range []struct {
		flag, block, resource string
	}{
		{"users_managed_externally", "user", "aws_paas_pgsql_user"},
		{"databases_managed_externally", "database", "aws_paas_pgsql_database"},
	} {
		if !d.Get(v.flag).(bool) {
			continue
		}

		if _, ok := d.GetOk(services.ServiceTypePostgreSQL); !ok {
			return fmt.Errorf("%s is supported only for %s services", v.flag, services.ServiceTypePostgreSQL)
		}

		key := services.ServiceTypePostgreSQL + ".0." + v.block
		if d.NewValueKnown(key) && len(d.Get(key).([]interface{})) > 0 {
			return fmt.Errorf("%s must not be set when %s is true, use %s resources instead", key, v.flag, v.resource)
		}
	}

	return nil
}

// Preserve compatibility with old environments by sending an empty roles array
// only for users whose previously assigned roles are being removed.
func clearPostgreSQLUserRoles(users []*paas.UserCreateRequest, previous []interface{}) {
	assigned := make(map[string]bool)
	for _, user := range previous {
		user := user.(map[string]interface{})
		if roles, ok := user["roles"].(*schema.Set); ok && roles.Len() > 0 {
			assigned[user["name"].(string)] = true
		}
	}
	for _, user := range users {
		if _, hasRoles := user.Parameters["roles"]; !hasRoles && assigned[aws.StringValue(user.Name)] {
			user.Parameters["roles"] = []interface{}{}
		}
	}
}
