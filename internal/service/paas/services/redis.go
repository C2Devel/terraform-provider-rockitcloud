package services

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type redisManager struct {
	service
}

var Redis = redisManager{
	service{
		name:               ServiceTypeRedis,
		class:              []string{ServiceClassDatabase, ServiceClassCacher},
		defaultClass:       ServiceClassCacher,
		allowArbitrator:    false,
		allowBackup:        false,
		dataVolumeRequired: true,
		usersEnabled:       false,
		databasesEnabled:   false,
		loggingEnabled:     true,
		monitoringEnabled:  true,
	},
}

func (s redisManager) serviceParametersSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cluster_type": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"native", "sentinel"}, false),
		},
		"databases": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 2147483647),
		},
		"maxmemory_policy": {
			Type:     schema.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"noeviction",
				"allkeys-lru",
				"allkeys-lfu",
				"volatile-lru",
				"volatile-lfu",
				"allkeys-random",
				"volatile-random",
				"volatile-ttl",
			}, false),
		},
		"options": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"password": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
			ForceNew:  true,
			ValidateFunc: validation.All(
				validation.StringIsNotEmpty,
				validation.StringDoesNotContainAny("`'\"\\"),
			),
		},
		"persistence_aof": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"persistence_rdb": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"timeout": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 2147483647),
		},
		"tcp_backlog": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 4096),
		},
		"tcp_keepalive": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(0),
		},
		"version": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
	}
}

func (s redisManager) serviceParametersDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cluster_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"databases": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"maxmemory_policy": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"options": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"password": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"persistence_aof": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"persistence_rdb": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"timeout": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"tcp_backlog": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"tcp_keepalive": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"version": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (s redisManager) expandServiceParameters(tfMap map[string]interface{}) ServiceParameters {
	if tfMap == nil {
		return nil
	}

	serviceParameters := ServiceParameters{}

	if v, ok := tfMap["cluster_type"].(string); ok && v != "" {
		serviceParameters["cluster_type"] = v
	}

	if v, ok := tfMap["databases"].(int); ok && v != 0 {
		serviceParameters["databases"] = int64(v)
	}

	if v, ok := tfMap["maxmemory_policy"].(string); ok && v != "" {
		serviceParameters["maxmemory-policy"] = v
	}

	if v, ok := tfMap["password"].(string); ok && v != "" {
		serviceParameters["password"] = v
	}

	if v, ok := tfMap["persistence_aof"].(bool); ok {
		serviceParameters["persistence_aof"] = v
	}

	if v, ok := tfMap["persistence_rdb"].(bool); ok {
		serviceParameters["persistence_rdb"] = v
	}

	if v, ok := tfMap["options"].(map[string]interface{}); ok && len(v) > 0 {
		serviceParameters["options"] = v
	}

	if v, ok := tfMap["timeout"].(int); ok {
		serviceParameters["timeout"] = int64(v)
	}

	if v, ok := tfMap["tcp_backlog"].(int); ok && v != 0 {
		serviceParameters["tcp-backlog"] = int64(v)
	}

	if v, ok := tfMap["tcp_keepalive"].(int); ok {
		serviceParameters["tcp-keepalive"] = int64(v)
	}

	if v, ok := tfMap["version"].(string); ok {
		serviceParameters["version"] = v
	}

	return serviceParameters
}

func (s redisManager) flattenServiceParameters(serviceParameters ServiceParameters) map[string]interface{} {
	if serviceParameters == nil {
		return map[string]interface{}{}
	}

	tfMap := map[string]interface{}{}

	if v, ok := serviceParameters["clusterType"].(string); ok {
		tfMap["cluster_type"] = v
	}

	if v, ok := serviceParameters["databases"].(int64); ok {
		tfMap["databases"] = v
	}

	if v, ok := serviceParameters["maxmemory-policy"].(string); ok {
		tfMap["maxmemory_policy"] = v
	}

	if v, ok := serviceParameters["password"].(string); ok {
		tfMap["password"] = v
	}

	if v, ok := serviceParameters["persistenceAof"].(bool); ok {
		tfMap["persistence_aof"] = v
	}

	if v, ok := serviceParameters["persistenceRdb"].(bool); ok {
		tfMap["persistence_rdb"] = v
	}

	if v, ok := serviceParameters["options"].(map[string]interface{}); ok {
		tfMap["options"] = v
	}

	if v, ok := serviceParameters["timeout"].(int64); ok {
		tfMap["timeout"] = v
	}

	if v, ok := serviceParameters["tcp-backlog"].(int64); ok {
		tfMap["tcp_backlog"] = v
	}

	if v, ok := serviceParameters["tcp-keepalive"].(int64); ok {
		tfMap["tcp_keepalive"] = v
	}

	if v, ok := serviceParameters["version"].(string); ok {
		tfMap["version"] = v
	}

	return tfMap
}
