package services

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type mongoDBManager struct {
	service
}

var MongoDB = mongoDBManager{
	service{
		name:               ServiceTypeMongoDB,
		class:              []string{ServiceClassDatabase},
		defaultClass:       ServiceClassDatabase,
		allowArbitrator:    true,
		allowBackup:        false,
		dataVolumeRequired: true,
		usersEnabled:       true,
		databasesEnabled:   true,
		loggingEnabled:     true,
		monitoringEnabled:  true,
	},
}

func (s mongoDBManager) serviceParametersSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"journal_commit_interval": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 500),
		},
		"maxconns": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(10, 51200),
		},
		"options": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"profile": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"off", "slowOp", "all"}, false),
		},
		"slowms": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 36000000),
		},
		"storage_engine_cache_size": s.floatByteParameterSchema(),
		"quiet": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"verbositylevel": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"v", "vv", "vvv", "vvvv", "vvvvv"}, false),
		},
		"version": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"3.6.23",
				"4.0.28",
				"4.2.23",
				"4.4.17",
				"5.0.13",
			}, false),
		},
	}
}

func (s mongoDBManager) serviceParametersDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"journal_commit_interval": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"maxconns": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"options": {
			Type:     schema.TypeMap,
			Computed: true,
		},
		"profile": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"slowms": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"storage_engine_cache_size": s.floatByteParameterDataSourceSchema(),
		"quiet": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"verbositylevel": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"version": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (s mongoDBManager) userParametersSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"password": {
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringDoesNotContainAny("`'\"\\"),
		},
	}
}

func (s mongoDBManager) userParametersDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"password": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (s mongoDBManager) databaseUserParametersSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"roles": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					"read",
					"readWrite",
					"dbAdmin",
					"dbOwner",
				}, false),
			},
		},
	}
}

func (s mongoDBManager) databaseUserParametersDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"roles": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func (s mongoDBManager) expandServiceParameters(tfMap map[string]interface{}) ServiceParameters {
	if tfMap == nil {
		return nil
	}

	serviceParameters := ServiceParameters{}

	if v, ok := tfMap["journal_commit_interval"].(int); ok && v != 0 {
		serviceParameters["journal_commit_interval"] = int64(v)
	}

	if v, ok := tfMap["maxconns"].(int); ok && v != 0 {
		serviceParameters["maxconns"] = int64(v)
	}

	if v, ok := tfMap["options"].(map[string]interface{}); ok && len(v) > 0 {
		serviceParameters["options"] = v
	}

	if v, ok := tfMap["profile"].(string); ok && v != "" {
		serviceParameters["profile"] = v
	}

	if v, ok := tfMap["slowms"].(int); ok {
		serviceParameters["slowms"] = int64(v)
	}

	if v, ok := tfMap["storage_engine_cache_size"].([]interface{}); ok && len(v) > 0 {
		var vMap = v[0].(map[string]interface{})
		serviceParameters["storage_engine_cache_size"] = map[string]interface{}{
			"value":     vMap["value"].(float64),
			"dimension": vMap["dimension"].(string),
		}
	}

	// incorrect naming of parameter in api
	if v, ok := tfMap["quiet"].(bool); ok {
		serviceParameters["verbose"] = v
	}

	if v, ok := tfMap["verbositylevel"].(string); ok && v != "" {
		serviceParameters["verbositylevel"] = v
	}

	if v, ok := tfMap["version"].(string); ok && v != "" {
		serviceParameters["version"] = v
	}

	return serviceParameters
}

func (s mongoDBManager) expandUserParameters(tfMap map[string]interface{}) UserParameters {
	if tfMap == nil {
		return nil
	}

	userParameters := UserParameters{}

	if v, ok := tfMap["password"].(string); ok && v != "" {
		userParameters["password"] = v
	}

	return userParameters
}

func (s mongoDBManager) expandDatabaseUserParameters(tfMap map[string]interface{}) DatabaseUserParameters {
	if tfMap == nil {
		return nil
	}

	databaseUserParameters := DatabaseUserParameters{}

	if v, ok := tfMap["roles"].(*schema.Set); ok && v.Len() > 0 {
		databaseUserParameters["roles"] = v.List()
	}

	return databaseUserParameters
}

func (s mongoDBManager) flattenServiceParameters(serviceParameters ServiceParameters) map[string]interface{} {
	if serviceParameters == nil {
		return map[string]interface{}{}
	}

	tfMap := map[string]interface{}{}

	if v, ok := serviceParameters["journalCommitInterval"].(int64); ok {
		tfMap["journal_commit_interval"] = v
	}

	if v, ok := serviceParameters["maxconns"].(int64); ok {
		tfMap["maxconns"] = v
	}

	if v, ok := serviceParameters["options"].(map[string]interface{}); ok {
		tfMap["options"] = v
	}

	if v, ok := serviceParameters["profile"].(string); ok {
		tfMap["profile"] = v
	}

	if v, ok := serviceParameters["slowms"].(int64); ok {
		tfMap["slowms"] = v
	}

	if vMap, okMap := serviceParameters["storageEngineCacheSize"].(map[string]interface{}); okMap {
		// todo: refactor
		if v, ok := vMap["value"]; ok && vMap["dimension"].(string) == GiB {
			if int64Val, int64Ok := v.(int64); int64Ok {
				tfMap["storage_engine_cache_size"] = float64(int64Val)
			} else if floatVal, floatOk := vMap["value"].(float64); floatOk {
				tfMap["storage_engine_cache_size"] = floatVal
			} else {
				log.Printf("[ERROR] Unknown type '%T' for storageEngineCacheSize %s", v, v)
			}
		}
		//tfMap["storage_engine_cache_size"] = flattenFloatValueWithDimension(vMap)
	}

	// incorrect naming of parameter in api
	if v, ok := serviceParameters["verbose"].(bool); ok {
		tfMap["quiet"] = v
	}

	if v, ok := serviceParameters["verbositylevel"].(string); ok {
		tfMap["verbositylevel"] = v
	}

	if v, ok := serviceParameters["version"].(string); ok {
		tfMap["version"] = v
	}

	return tfMap
}

func (s mongoDBManager) flattenUserParameters(userParameters UserParameters) map[string]interface{} {
	if userParameters == nil {
		return map[string]interface{}{}
	}

	tfMap := map[string]interface{}{}

	if v, ok := userParameters["password"].(string); ok {
		tfMap["password"] = v
	}

	return tfMap
}

func (s mongoDBManager) flattenDatabaseUserParameters(databaseUserParameters DatabaseUserParameters) map[string]interface{} {
	if databaseUserParameters == nil {
		return map[string]interface{}{}
	}

	tfMap := map[string]interface{}{}

	if v, ok := databaseUserParameters["roles"].([]interface{}); ok {
		tfMap["roles"] = v
	}

	return tfMap
}
