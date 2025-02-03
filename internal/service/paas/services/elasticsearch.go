package services

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/experimental/nullable"
)

type elasticSearchManager struct {
	service
}

var ElasticSearch = elasticSearchManager{
	service{
		name:               ServiceTypeElasticSearch,
		class:              []string{ServiceClassSearch},
		defaultClass:       ServiceClassSearch,
		allowArbitrator:    true,
		allowBackup:        false,
		dataVolumeRequired: true,
		usersEnabled:       false,
		databasesEnabled:   false,
		loggingEnabled:     true,
		monitoringEnabled:  true,
	},
}

func (s elasticSearchManager) serviceParametersSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allow_anonymous": {
			Type:         nullable.TypeNullableBool,
			Optional:     true,
			ForceNew:     true,
			Computed:     true,
			ValidateFunc: nullable.ValidateTypeStringNullableBool,
		},
		"anonymous_role": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Computed:     true,
			ValidateFunc: validation.StringInSlice([]string{"viewer", "editor"}, false),
		},
		"kibana": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},
		"options": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"password": {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			ForceNew:     true,
			ValidateFunc: validation.StringDoesNotContainAny("^-!:;%'`\"\\"),
		},
		"version": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
	}
}

func (s elasticSearchManager) serviceParametersDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allow_anonymous": {
			Type:     nullable.TypeNullableBool,
			Optional: true,
		},
		"anonymous_role": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"kibana": {
			Type:     schema.TypeBool,
			Optional: true,
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
		"version": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (s elasticSearchManager) expandServiceParameters(tfMap map[string]interface{}) ServiceParameters {
	if tfMap == nil {
		return nil
	}

	serviceParameters := ServiceParameters{}

	if v, null, _ := nullable.Bool(tfMap["allow_anonymous"].(string)).Value(); !null {
		serviceParameters["allow_anonymous"] = v
	}

	if v, ok := tfMap["anonymous_role"].(string); ok && v != "" {
		serviceParameters["anonymous_role"] = v
	}

	if v, ok := tfMap["kibana"].(bool); ok {
		serviceParameters["kibana"] = v
	}

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

func (s elasticSearchManager) flattenServiceParameters(serviceParameters ServiceParameters) map[string]interface{} {
	if serviceParameters == nil {
		return map[string]interface{}{}
	}

	tfMap := map[string]interface{}{}

	if v, ok := serviceParameters["allowAnonymous"].(bool); ok {
		tfMap["allow_anonymous"] = strconv.FormatBool(v)
	}

	if v, ok := serviceParameters["anonymousRole"].(string); ok {
		tfMap["anonymous_role"] = v
	}

	if v, ok := serviceParameters["kibana"].(bool); ok {
		tfMap["kibana"] = v
	}

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
