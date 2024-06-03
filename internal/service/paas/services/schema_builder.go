package services

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/experimental/nullable"
)

// ResourceSchema returns a full schema of service parameters for resource.
//
// It includes:
//   - common `class` parameter that holds chosen class of service
//   - service-specific parameters including `user`, `database`, `logging`, and `monitoring` parameter blocks
//     if service can hold them
func (s service) ResourceSchema() *schema.Schema {
	serviceSchema := map[string]*schema.Schema{
		"class": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      s.defaultClass,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(s.class, false),
		},
	}

	for k, v := range s.toInterface().serviceParametersSchema() {
		serviceSchema[k] = v
	}

	if s.usersEnabled {
		serviceSchema["user"] = s.userSchema(s.toInterface().userParametersSchema())
	}

	if s.databasesEnabled {
		serviceSchema["database"] = s.databaseSchema(s.toInterface().databaseParametersSchema())
	}

	if s.loggingEnabled {
		serviceSchema["logging"] = &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"log_to": {
						Type:     schema.TypeString,
						Required: true,
					},
					"logging_tags": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validation.StringLenBetween(1, 256),
						},
					},
				},
			},
		}
	}

	if s.monitoringEnabled {
		serviceSchema["monitoring"] = &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"monitor_by": {
						Type:     schema.TypeString,
						Required: true,
					},
					"monitoring_labels": {
						Type:             schema.TypeMap,
						Optional:         true,
						ValidateDiagFunc: validation.MapKeyLenBetween(1, 64),
						Elem:             &schema.Schema{Type: schema.TypeString},
					},
				},
			},
		}
	}

	var requiredWith, conflictsWith []string
	if s.dataVolumeRequired {
		requiredWith = append(requiredWith, "data_volume")
	}

	if !s.allowArbitrator {
		conflictsWith = append(conflictsWith, "arbitrator_required")
	}

	if !s.allowBackup {
		conflictsWith = append(conflictsWith, "backup_settings")
	}

	return &schema.Schema{
		Type:          schema.TypeList,
		MaxItems:      1,
		Optional:      true,
		ForceNew:      true,
		ExactlyOneOf:  ManagedServiceTypes(),
		ConflictsWith: conflictsWith,
		RequiredWith:  requiredWith,
		Elem: &schema.Resource{
			Schema: serviceSchema,
		},
	}
}

// userSchema returns a schema of user parameters for resource.
//
// It includes:
//   - common parameter name that holds username
//   - service-specific user parameters, if service has them
func (s service) userSchema(userParametersSchema map[string]*schema.Schema) *schema.Schema {
	userSchema := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
	}

	for k, v := range userParametersSchema {
		userSchema[k] = v
	}

	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1000,
		Elem: &schema.Resource{
			Schema: userSchema,
		},
	}
}

// databaseSchema returns a schema of database parameters for resource.
//
// It includes:
//   - common `backup_enabled` and `name` parameters
//   - service-specific database user parameters if service has them
//   - service-specific database parameters, if service has them
func (s service) databaseSchema(databaseParametersSchema map[string]*schema.Schema) *schema.Schema {
	databaseSchema := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
	}

	if s.allowBackup {
		databaseSchema["backup_enabled"] = &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		}
	}

	if s.usersEnabled {
		databaseSchema["user"] = s.userSchema(s.toInterface().databaseUserParametersSchema())
	}

	for k, v := range databaseParametersSchema {
		databaseSchema[k] = v
	}

	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1000,
		Elem: &schema.Resource{
			Schema: databaseSchema,
		},
	}
}

// serviceParametersSchema represents resource schema for service-specific parameters except for databases and users.
//
// If PaaS service has specific parameters, it should override this method.
func (s service) serviceParametersSchema() map[string]*schema.Schema {
	return nil
}

// userParametersSchema represents resource schema for service-specific user parameters.
//
// If PaaS service has specific user parameters, it should override this method.
func (s service) userParametersSchema() map[string]*schema.Schema {
	return nil
}

// databaseParametersSchema represents resource schema for service-specific database parameters.
//
// If PaaS service has specific database parameters, it should override this method.
func (s service) databaseParametersSchema() map[string]*schema.Schema {
	return nil
}

// databaseUserParametersSchema represents resource schema for service-specific database user parameters.
//
// If PaaS service has specific database user parameters, it should override this method.
func (s service) databaseUserParametersSchema() map[string]*schema.Schema {
	return nil
}

// todo: doc
// todo: parametrize default value
func (s service) intByteParameterSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value": {
					Type:     nullable.TypeNullableInt,
					Required: true,
				},
				"dimension": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice(
						[]string{B, KiB, MiB, GiB, TiB},
						false,
					),
				},
			},
		},
	}
}

// todo: doc
// todo: parametrize forceNew
func (s service) floatByteParameterSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value": {
					Type:     schema.TypeFloat,
					Required: true,
				},
				"dimension": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice(
						[]string{B, KiB, MiB, GiB, TiB},
						false,
					),
				},
			},
		},
	}
}

// DataSourceSchema returns a full schema of service parameters for datasource.
//
// It includes the same blocks as ResourceSchema.
func (s service) DataSourceSchema() *schema.Schema {
	serviceSchema := map[string]*schema.Schema{
		"class": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	for k, v := range s.toInterface().serviceParametersDataSourceSchema() {
		serviceSchema[k] = v
	}

	if s.usersEnabled {
		serviceSchema["user"] = s.userDataSourceSchema(s.toInterface().userParametersDataSourceSchema())
	}

	if s.databasesEnabled {
		serviceSchema["database"] = s.databaseDataSourceSchema(s.toInterface().databaseParametersSchema())
	}

	if s.loggingEnabled {
		serviceSchema["logging"] = &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"log_to": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"logging_tags": {
						Type:     schema.TypeSet,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
				},
			},
		}
	}

	if s.monitoringEnabled {
		serviceSchema["monitoring"] = &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"monitor_by": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"monitoring_labels": {
						Type:     schema.TypeMap,
						Computed: true,
					},
				},
			},
		}
	}

	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: serviceSchema,
		},
	}
}

// userDataSourceSchema returns a schema of user parameters for datasource.
//
// It includes the same blocks as userSchema.
func (s service) userDataSourceSchema(userParametersSchema map[string]*schema.Schema) *schema.Schema {
	userSchema := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	for k, v := range userParametersSchema {
		userSchema[k] = v
	}

	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: userSchema,
		},
	}
}

// databaseSchema returns a schema of database parameters for resource.
//
// It includes the same blocks as databaseSchema.
func (s service) databaseDataSourceSchema(databaseParametersSchema map[string]*schema.Schema) *schema.Schema {
	databaseSchema := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	if s.allowBackup {
		databaseSchema["backup_enabled"] = &schema.Schema{
			Type:     schema.TypeBool,
			Computed: true,
		}
	}

	if s.usersEnabled {
		databaseSchema["user"] = s.userDataSourceSchema(s.toInterface().databaseUserParametersDataSourceSchema())
	}

	for k, v := range databaseParametersSchema {
		databaseSchema[k] = v
	}

	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: databaseSchema,
		},
	}
}

// serviceParametersDataSourceSchema represents datasource schema for service-specific parameters
// except for databases and users.
//
// If PaaS service has specific parameters, it should override this method.
func (s service) serviceParametersDataSourceSchema() map[string]*schema.Schema {
	return nil
}

// userParametersDataSourceSchema represents datasource schema for service-specific user parameters.
//
// If PaaS service has specific user parameters, it should override this method.
func (s service) userParametersDataSourceSchema() map[string]*schema.Schema {
	return nil
}

// databaseParametersDataSourceSchema represents datasource schema for service-specific database parameters.
//
// If PaaS service has specific database parameters, it should override this method.
func (s service) databaseParametersDataSourceSchema() map[string]*schema.Schema {
	return nil
}

// databaseUserParametersDataSourceSchema represents datasource schema for service-specific database user parameters.
//
// If PaaS service has specific database user parameters, it should override this method.
func (s service) databaseUserParametersDataSourceSchema() map[string]*schema.Schema {
	return nil
}

// todo: doc
func (s service) intByteParameterDataSourceSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value": {
					Type:     nullable.TypeNullableInt,
					Computed: true,
				},
				"dimension": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

// todo: doc
func (s service) floatByteParameterDataSourceSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value": {
					Type:     schema.TypeFloat,
					Computed: true,
				},
				"dimension": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}
