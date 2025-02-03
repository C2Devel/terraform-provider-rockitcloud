package services

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type rabbitMQManager struct {
	service
}

var RabbitMQ = rabbitMQManager{
	service{
		name:               ServiceTypeRabbitMQ,
		class:              []string{ServiceClassMessageBroker},
		defaultClass:       ServiceClassMessageBroker,
		allowArbitrator:    false,
		allowBackup:        false,
		dataVolumeRequired: true,
		usersEnabled:       false,
		databasesEnabled:   false,
		loggingEnabled:     true,
		monitoringEnabled:  true,
	},
}

func (s rabbitMQManager) serviceParametersSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"options": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"password": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(8, 32),
				validation.StringDoesNotContainAny(`-|[]\'";`),
			),
		},
		"version": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
	}
}

func (s rabbitMQManager) serviceParametersDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"options": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"password": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"version": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (s rabbitMQManager) expandServiceParameters(tfMap map[string]interface{}) ServiceParameters {
	if tfMap == nil {
		return nil
	}

	serviceParameters := ServiceParameters{}

	if v, ok := tfMap["password"].(string); ok && v != "" {
		serviceParameters["password"] = v
	}

	if v, ok := tfMap["options"].(map[string]interface{}); ok && len(v) > 0 {
		serviceParameters["options"] = v
	}

	if v, ok := tfMap["version"].(string); ok {
		serviceParameters["version"] = v
	}

	return serviceParameters
}

func (s rabbitMQManager) flattenServiceParameters(serviceParameters ServiceParameters) map[string]interface{} {
	if serviceParameters == nil {
		return map[string]interface{}{}
	}

	tfMap := map[string]interface{}{}

	if v, ok := serviceParameters["password"].(string); ok {
		tfMap["password"] = v
	}

	if v, ok := serviceParameters["options"].(map[string]interface{}); ok {
		tfMap["options"] = v
	}

	if v, ok := serviceParameters["version"].(string); ok {
		tfMap["version"] = v
	}

	return tfMap
}
