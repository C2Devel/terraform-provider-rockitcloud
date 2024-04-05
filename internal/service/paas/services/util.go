package services

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	ServiceTypeElasticSearch = "elasticsearch"
	ServiceTypeMemcached     = "memcached"
	ServiceTypeMongoDB       = "mongodb"
	ServiceTypeMySQL         = "mysql"
	ServiceTypePostgreSQL    = "pgsql"
	ServiceTypeRabbitMQ      = "rabbitmq"
	ServiceTypeRedis         = "redis"
)

func ServiceTypeValues() []string {
	return []string{
		ServiceTypeElasticSearch,
		ServiceTypeMemcached,
		ServiceTypeMongoDB,
		ServiceTypeMySQL,
		ServiceTypePostgreSQL,
		ServiceTypeRabbitMQ,
		ServiceTypeRedis,
	}
}

const (
	ServiceClassCacher        = "cacher"
	ServiceClassDatabase      = "database"
	ServiceClassMessageBroker = "message_broker"
	ServiceClassSearch        = "search"
)

func ServiceClassValues() []string {
	return []string{
		ServiceClassCacher,
		ServiceClassDatabase,
		ServiceClassMessageBroker,
		ServiceClassSearch,
	}
}

// TODO: move to some common package
const (
	Byte = 1 << (10 * iota)
	Kilobyte
	Megabyte
	Gigabyte
	Terabyte
)

const (
	B   = "B"
	KiB = "KiB"
	MiB = "MiB"
	GiB = "GiB"
	TiB = "TiB"
)

// DimensionToBytes converts a dimension string to its corresponding value in bytes.
// It returns an error if the dimension is not recognized.
func DimensionToBytes(dimension string) (int64, error) {
	switch dimension {
	case B:
		return Byte, nil
	case KiB:
		return Kilobyte, nil
	case MiB:
		return Megabyte, nil
	case GiB:
		return Gigabyte, nil
	case TiB:
		return Terabyte, nil
	default:
		return 0, fmt.Errorf("unsupported dimension: %s", dimension)
	}
}

// parseBytes parses the given value and dimension, returning the value in bytes.
func parseBytes(value int64, dimension string) (int64, error) {
	bytes, err := DimensionToBytes(dimension)

	if err != nil {
		log.Printf("[ERROR] Error parsing value `%d %s` to bytes: %s", value, dimension, err)
		return value, err
	}

	return value * bytes, nil
}

// Map with ServiceManager objects for each supported PaaS service.
var managers = map[string]ServiceManager{
	ElasticSearch.ServiceType(): ElasticSearch,
	Memcached.ServiceType():     Memcached,
	MongoDB.ServiceType():       MongoDB,
	MySQL.ServiceType():         MySQL,
	PostgreSQL.ServiceType():    PostgreSQL,
	Redis.ServiceType():         Redis,
	RabbitMQ.ServiceType():      RabbitMQ,
}

func ManagedServiceTypes() []string {
	keys := make([]string, 0, len(managers))

	for k := range managers {
		keys = append(keys, k)
	}
	return keys
}

func Manager(serviceType string) ServiceManager {
	if v, ok := managers[serviceType]; ok {
		return v
	}

	log.Printf("[ERROR] Unknown service type: %s", serviceType)
	return nil
}

func getIntValueWithDimensionResourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"value": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"dimension": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  KiB,
				ValidateFunc: validation.StringInSlice(
					[]string{KiB, MiB, GiB, TiB},
					false,
				),
			},
		},
	}
}

func getFloatValueWithDimensionResourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"value": {
				Type:         schema.TypeFloat,
				Optional:     true,
				ValidateFunc: validation.FloatAtLeast(0.25),
			},
			"dimension": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  KiB,
				ValidateFunc: validation.StringInSlice(
					[]string{KiB, MiB, GiB, TiB},
					false,
				),
			},
		},
	}
}

func flattenIntValueWithDimension(valMap map[string]interface{}) []map[string]interface{} {
	var value = map[string]interface{}{
		"value":     valMap["value"].(int64),
		"dimension": valMap["dimension"].(string),
	}
	return []map[string]interface{}{value}
}

func flattenFloatValueWithDimension(valMap map[string]interface{}) []map[string]interface{} {
	var value = map[string]interface{}{
		"value":     valMap["value"].(float64),
		"dimension": valMap["dimension"].(string),
	}
	return []map[string]interface{}{value}
}
