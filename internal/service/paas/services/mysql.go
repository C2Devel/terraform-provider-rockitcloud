package services

import (
	"math"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/experimental/nullable"
)

type mySQLManager struct {
	service
}

var MySQL = mySQLManager{
	service{
		name:               ServiceTypeMySQL,
		class:              []string{ServiceClassDatabase},
		defaultClass:       ServiceClassDatabase,
		allowArbitrator:    true,
		allowBackup:        true,
		dataVolumeRequired: true,
		usersEnabled:       true,
		databasesEnabled:   true,
		loggingEnabled:     true,
		monitoringEnabled:  true,
	},
}

func (s mySQLManager) serviceParametersSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connect_timeout": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(2, 31536000),
		},
		"galera_options": {
			Type:     schema.TypeMap,
			Optional: true,
			ForceNew: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"gcache_size": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(128 * Megabyte),
		},
		"gcs_fc_factor": {
			Type:         schema.TypeFloat,
			Optional:     true,
			ForceNew:     true,
			Default:      mySQLGcsFcFactorDefault,
			ValidateFunc: validation.FloatBetween(0.0, 1.0),
		},
		"gcs_fc_limit": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 2147483647),
		},
		"gcs_fc_master_slave": {
			Type:         nullable.TypeNullableBool,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: nullable.ValidateTypeStringNullableBool,
		},
		"gcs_fc_single_primary": {
			Type:         nullable.TypeNullableBool,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: nullable.ValidateTypeStringNullableBool,
		},
		"innodb_buffer_pool_instances": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 64),
		},
		"innodb_buffer_pool_size": {
			Type:     nullable.TypeNullableInt,
			Optional: true,
			ForceNew: true,
			// FIXME: uncomment the validation after all paas clusters are on v4.0
			// ValidateFunc: nullable.ValidateTypeStringNullableIntAtLeast(128*Megabyte),
		},
		"innodb_change_buffering": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"inserts",
				"deletes",
				"changes",
				"purges",
				"all",
				"none",
			}, false),
		},
		"innodb_flush_log_at_trx_commit": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(0, 2),
		},
		"innodb_io_capacity": {
			Type:     nullable.TypeNullableInt,
			Optional: true,
			ForceNew: true,
			// the actual max value is 2^64 - 1, but it doesn't fit in int64
			ValidateFunc: nullable.ValidateTypeStringNullableIntBetween(100, math.MaxInt64),
		},
		"innodb_io_capacity_max": {
			Type:     nullable.TypeNullableInt,
			Optional: true,
			ForceNew: true,
			// the actual max value is 2^64 - 1, but it doesn't fit in int64
			ValidateFunc: nullable.ValidateTypeStringNullableIntBetween(100, math.MaxInt64),
		},
		"innodb_log_file_size": {
			Type:     nullable.TypeNullableInt,
			Optional: true,
			ForceNew: true,
			// FIXME: uncomment the validation after all paas clusters are on v4.0
			// ValidateFunc: nullable.ValidateTypeStringNullableIntBetween(4*Megabyte, 4*Gigabyte),
		},
		"innodb_log_files_in_group": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(2, 100),
		},
		"innodb_purge_threads": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 32),
		},
		"innodb_thread_concurrency": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      mySQLInnodbThreadConcurrencyDefault,
			ValidateFunc: validation.IntBetween(0, 1000),
		},
		"innodb_strict_mode": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "OFF",
			ValidateFunc: validation.StringInSlice([]string{"ON", "OFF"}, false),
		},
		"innodb_sync_array_size": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 1024),
		},
		"max_allowed_packet": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      16 * Megabyte,
			ValidateFunc: validation.IntBetween(16*Megabyte, 1*Gigabyte),
		},
		"max_connect_errors": {
			Type:     nullable.TypeNullableInt,
			Optional: true,
			ForceNew: true,
			// the actual max value is 2^64 - 1, but it doesn't fit in int64
			ValidateFunc: nullable.ValidateTypeStringNullableIntBetween(1, math.MaxInt64),
		},
		"max_connections": {
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
			// FIXME: uncomment the validation after all paas clusters are on v4.0
			// ValidateFunc: validation.IntBetween(10, 100000),
		},
		"max_heap_table_size": {
			Type:         nullable.TypeNullableInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: nullable.ValidateTypeStringNullableIntBetween(16*Kilobyte, 4*Gigabyte),
		},
		"options": {
			Type:     schema.TypeMap,
			Optional: true,
			ForceNew: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"pxc_strict_mode": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"DISABLED", "PERMISSIVE", "ENFORCING", "MASTER"}, false),
		},
		"table_open_cache": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 1048576),
		},
		"thread_cache_size": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      mySQLThreadCacheSizeDefault,
			ValidateFunc: validation.IntBetween(0, 16*Kilobyte),
		},
		"tmp_table_size": {
			Type:         nullable.TypeNullableInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: nullable.ValidateTypeStringNullableIntBetween(1*Kilobyte, 4*Gigabyte),
		},
		"transaction_isolation": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"READ-UNCOMMITTED",
				"READ-COMMITTED",
				"REPEATABLE-READ",
				"SERIALIZABLE",
			}, false),
		},
		"vendor": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"mariadb", "percona", "mysql"}, false),
		},
		"version": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"wait_timeout": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 31536000),
		},
	}
}

func (s mySQLManager) serviceParametersDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connect_timeout": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"galera_options": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"gcache_size": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"gcs_fc_factor": {
			Type:     schema.TypeFloat,
			Computed: true,
		},
		"gcs_fc_limit": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"gcs_fc_master_slave": {
			Type:     nullable.TypeNullableBool,
			Computed: true,
		},
		"gcs_fc_single_primary": {
			Type:     nullable.TypeNullableBool,
			Computed: true,
		},
		"innodb_buffer_pool_instances": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"innodb_buffer_pool_size": {
			Type:     nullable.TypeNullableInt,
			Computed: true,
		},
		"innodb_change_buffering": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"innodb_flush_log_at_trx_commit": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"innodb_io_capacity": {
			Type:     nullable.TypeNullableInt,
			Computed: true,
		},
		"innodb_io_capacity_max": {
			Type:     nullable.TypeNullableInt,
			Computed: true,
		},
		"innodb_log_file_size": {
			Type:     nullable.TypeNullableInt,
			Computed: true,
		},
		"innodb_log_files_in_group": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"innodb_purge_threads": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"innodb_thread_concurrency": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"innodb_strict_mode": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"innodb_sync_array_size": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"max_allowed_packet": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"max_connect_errors": {
			Type:     nullable.TypeNullableInt,
			Computed: true,
		},
		"max_connections": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"max_heap_table_size": {
			Type:     nullable.TypeNullableInt,
			Computed: true,
		},
		"options": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"pxc_strict_mode": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"table_open_cache": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"thread_cache_size": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"tmp_table_size": {
			Type:     nullable.TypeNullableInt,
			Computed: true,
		},
		"transaction_isolation": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vendor": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"version": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"wait_timeout": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func (s mySQLManager) userParametersSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"host": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 60),
		},
		"password": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
			ValidateFunc: validation.All(
				validation.StringIsNotEmpty,
				validation.StringDoesNotContainAny("`'\"\\"),
			),
		},
	}
}

func (s mySQLManager) userParametersDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"host": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"password": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (s mySQLManager) databaseParametersSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"backup_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"backup_db_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"charset": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"collate": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

func (s mySQLManager) databaseParametersDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"backup_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"backup_db_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"charset": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"collate": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (s mySQLManager) databaseUserParametersSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"options": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"GRANT", "NONE"}, false),
			},
		},
		"privileges": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(mySQLDatabaseUserPrivileges(), false),
			},
		},
	}
}

func (s mySQLManager) databaseUserParametersDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"options": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"privileges": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (s mySQLManager) expandServiceParameters(tfMap map[string]interface{}) ServiceParameters {
	if tfMap == nil {
		return nil
	}

	serviceParameters := ServiceParameters{}

	if v, ok := tfMap["connect_timeout"].(int); ok && v != 0 {
		serviceParameters["connect_timeout"] = int64(v)
	}

	if v, ok := tfMap["galera_options"].(map[string]interface{}); ok && len(v) > 0 {
		serviceParameters["galera_options"] = v
	}

	if v, ok := tfMap["gcache_size"].(int); ok && v != 0 {
		serviceParameters["gcache_size"] = map[string]interface{}{
			"dimension": B,
			"value":     int64(v),
		}
	}

	if v, ok := tfMap["gcs_fc_factor"].(float64); ok && v != mySQLGcsFcFactorDefault {
		serviceParameters["gcs_fc_factor"] = v
	}

	if v, ok := tfMap["gcs_fc_limit"].(int); ok && v != 0 {
		serviceParameters["gcs_fc_limit"] = int64(v)
	}

	if v, null, _ := nullable.Bool(tfMap["gcs_fc_master_slave"].(string)).Value(); !null {
		serviceParameters["gcs_fc_master_slave"] = v
	}

	if v, null, _ := nullable.Bool(tfMap["gcs_fc_single_primary"].(string)).Value(); !null {
		serviceParameters["gcs_fc_single_primary"] = v
	}

	if v, ok := tfMap["innodb_buffer_pool_instances"].(int); ok && v != 0 {
		serviceParameters["innodb_buffer_pool_instances"] = int64(v)
	}

	if v, _, _ := nullable.Int(tfMap["innodb_buffer_pool_size"].(string)).Value(); v != 0 {
		serviceParameters["innodb_buffer_pool_size"] = map[string]interface{}{
			"dimension": B,
			"value":     v,
		}
	}

	if v, ok := tfMap["innodb_change_buffering"].(string); ok && v != "" {
		serviceParameters["innodb_change_buffering"] = v
	}

	if v, ok := tfMap["innodb_flush_log_at_trx_commit"].(int); ok {
		serviceParameters["innodb_flush_log_at_trx_commit"] = int64(v)
	}

	if v, _, _ := nullable.Int(tfMap["innodb_io_capacity"].(string)).Value(); v != 0 {
		serviceParameters["innodb_io_capacity"] = v
	}

	if v, _, _ := nullable.Int(tfMap["innodb_io_capacity_max"].(string)).Value(); v != 0 {
		serviceParameters["innodb_io_capacity_max"] = v
	}

	if v, _, _ := nullable.Int(tfMap["innodb_log_file_size"].(string)).Value(); v != 0 {
		serviceParameters["innodb_log_file_size"] = map[string]interface{}{
			"dimension": B,
			"value":     v,
		}
	}

	if v, ok := tfMap["innodb_log_files_in_group"].(int); ok && v != 0 {
		serviceParameters["innodb_log_files_in_group"] = int64(v)
	}

	if v, ok := tfMap["innodb_purge_threads"].(int); ok && v != 0 {
		serviceParameters["innodb_purge_threads"] = int64(v)
	}

	if v, ok := tfMap["innodb_thread_concurrency"].(int); ok && v != mySQLInnodbThreadConcurrencyDefault {
		serviceParameters["innodb_thread_concurrency"] = int64(v)
	}

	if v, ok := tfMap["innodb_strict_mode"].(string); ok && v != "" {
		serviceParameters["innodb_strict_mode"] = v
	}

	if v, ok := tfMap["innodb_sync_array_size"].(int); ok && v != 0 {
		serviceParameters["innodb_sync_array_size"] = int64(v)
	}

	if v, ok := tfMap["max_allowed_packet"].(int); ok && v != 0 {
		serviceParameters["max_allowed_packet"] = map[string]interface{}{
			"dimension": B,
			"value":     int64(v),
		}
	}

	if v, _, _ := nullable.Int(tfMap["max_connect_errors"].(string)).Value(); v != 0 {
		serviceParameters["max_connect_errors"] = v
	}

	if v, ok := tfMap["max_connections"].(int); ok && v != 0 {
		serviceParameters["max_connections"] = int64(v)
	}

	if v, _, _ := nullable.Int(tfMap["max_heap_table_size"].(string)).Value(); v != 0 {
		serviceParameters["max_heap_table_size"] = map[string]interface{}{
			"dimension": B,
			"value":     v,
		}
	}

	if v, ok := tfMap["options"].(map[string]interface{}); ok && len(v) > 0 {
		serviceParameters["options"] = v
	}

	if v, ok := tfMap["pxc_strict_mode"].(string); ok && v != "" {
		serviceParameters["pxc_strict_mode"] = v
	}

	if v, ok := tfMap["table_open_cache"].(int); ok && v != 0 {
		serviceParameters["table_open_cache"] = int64(v)
	}

	if v, ok := tfMap["thread_cache_size"].(int); ok && v != mySQLThreadCacheSizeDefault {
		serviceParameters["thread_cache_size"] = int64(v)
	}

	if v, _, _ := nullable.Int(tfMap["tmp_table_size"].(string)).Value(); v != 0 {
		serviceParameters["tmp_table_size"] = map[string]interface{}{
			"dimension": B,
			"value":     v,
		}
	}

	if v, ok := tfMap["transaction_isolation"].(string); ok && v != "" {
		serviceParameters["transaction_isolation"] = v
	}

	if v, ok := tfMap["vendor"].(string); ok && v != "" {
		serviceParameters["vendor"] = v
	}

	if v, ok := tfMap["version"].(string); ok && v != "" {
		serviceParameters["version"] = v
	}

	if v, ok := tfMap["wait_timeout"].(int); ok && v != 0 {
		serviceParameters["wait_timeout"] = int64(v)
	}

	return serviceParameters
}

func (s mySQLManager) expandUserParameters(tfMap map[string]interface{}) UserParameters {
	if tfMap == nil {
		return nil
	}

	userParameters := UserParameters{}

	if v, ok := tfMap["host"].(string); ok && v != "" {
		userParameters["host"] = v
	}

	if v, ok := tfMap["password"].(string); ok && v != "" {
		userParameters["password"] = v
	}

	return userParameters
}

func (s mySQLManager) expandDatabaseParameters(tfMap map[string]interface{}) DatabaseParameters {
	if tfMap == nil {
		return nil
	}

	databaseParameters := DatabaseParameters{}

	if v, ok := tfMap["backup_id"].(string); ok && v != "" {
		databaseParameters["backup_id"] = v
	}

	if v, ok := tfMap["backup_db_name"].(string); ok && v != "" {
		databaseParameters["backup_db_name"] = v
	}

	if v, ok := tfMap["charset"].(string); ok && v != "" {
		databaseParameters["charset"] = v
	}

	if v, ok := tfMap["collate"].(string); ok && v != "" {
		databaseParameters["collate"] = v
	}

	return databaseParameters
}

func (s mySQLManager) expandDatabaseUserParameters(tfMap map[string]interface{}) DatabaseUserParameters {
	if tfMap == nil {
		return nil
	}

	databaseUserParameters := DatabaseUserParameters{}

	if v, ok := tfMap["options"].(*schema.Set); ok && v.Len() > 0 {
		databaseUserParameters["options"] = v.List()
	}

	if v, ok := tfMap["privileges"].(*schema.Set); ok && v.Len() > 0 {
		databaseUserParameters["privileges"] = v.List()
	}

	return databaseUserParameters
}

func (s mySQLManager) flattenServiceParameters(serviceParameters ServiceParameters) map[string]interface{} {
	if serviceParameters == nil {
		return map[string]interface{}{}
	}

	tfMap := map[string]interface{}{}

	if v, ok := serviceParameters["connectTimeout"].(int64); ok {
		tfMap["connect_timeout"] = v
	}

	if v, ok := serviceParameters["galeraOptions"].(map[string]interface{}); ok {
		tfMap["galera_options"] = v
	}

	if vMap, okMap := serviceParameters["gcacheSize"].(map[string]interface{}); okMap {
		bytes, err := parseBytes(vMap["value"].(int64), vMap["dimension"].(string))

		if err == nil {
			tfMap["gcache_size"] = bytes
		}
	}

	if v, ok := serviceParameters["gcsFcFactor"].(float64); ok {
		tfMap["gcs_fc_factor"] = v
	} else {
		tfMap["gcs_fc_factor"] = mySQLGcsFcFactorDefault
	}

	if v, ok := serviceParameters["gcsFcLimit"].(int64); ok {
		tfMap["gcs_fc_limit"] = v
	}

	if v, ok := serviceParameters["gcsFcMasterSlave"].(bool); ok {
		tfMap["gcs_fc_master_slave"] = strconv.FormatBool(v)
	}

	if v, ok := serviceParameters["gcsFcSinglePrimary"].(bool); ok {
		tfMap["gcs_fc_single_primary"] = strconv.FormatBool(v)
	}

	if v, ok := serviceParameters["innodbBufferPoolInstances"].(int64); ok {
		tfMap["innodb_buffer_pool_instances"] = v
	}

	if vMap, okMap := serviceParameters["innodbBufferPoolSize"].(map[string]interface{}); okMap {
		bytes, err := parseBytes(vMap["value"].(int64), vMap["dimension"].(string))

		if err == nil {
			tfMap["innodb_buffer_pool_size"] = strconv.FormatInt(bytes, 10)
		}
	}

	if v, ok := serviceParameters["innodbChangeBuffering"].(string); ok {
		tfMap["innodb_change_buffering"] = v
	}

	if v, ok := serviceParameters["innodbFlushLogAtTrxCommit"].(int64); ok {
		tfMap["innodb_flush_log_at_trx_commit"] = v
	}

	if v, ok := serviceParameters["innodbIoCapacity"].(int64); ok {
		tfMap["innodb_io_capacity"] = strconv.FormatInt(v, 10)
	}

	if v, ok := serviceParameters["innodbIoCapacityMax"].(int64); ok {
		tfMap["innodb_io_capacity_max"] = strconv.FormatInt(v, 10)
	}

	if vMap, okMap := serviceParameters["innodbLogFileSize"].(map[string]interface{}); okMap {
		bytes, err := parseBytes(vMap["value"].(int64), vMap["dimension"].(string))

		if err == nil {
			tfMap["innodb_log_file_size"] = strconv.FormatInt(bytes, 10)
		}
	}

	if v, ok := serviceParameters["innodbLogFilesInGroup"].(int64); ok {
		tfMap["innodb_log_files_in_group"] = v
	}

	if v, ok := serviceParameters["innodbPurgeThreads"].(int64); ok {
		tfMap["innodb_purge_threads"] = v
	}

	if v, ok := serviceParameters["innodbThreadConcurrency"].(int64); ok {
		tfMap["innodb_thread_concurrency"] = v
	} else {
		tfMap["innodb_thread_concurrency"] = mySQLInnodbThreadConcurrencyDefault
	}

	if v, ok := serviceParameters["innodbStrictMode"].(string); ok {
		tfMap["innodb_strict_mode"] = v
	}

	if v, ok := serviceParameters["innodbSyncArraySize"].(int64); ok {
		tfMap["innodb_sync_array_size"] = v
	}

	if vMap, okMap := serviceParameters["maxAllowedPacket"].(map[string]interface{}); okMap {
		bytes, err := parseBytes(vMap["value"].(int64), vMap["dimension"].(string))

		if err == nil {
			tfMap["max_allowed_packet"] = bytes
		}
	}

	if v, ok := serviceParameters["maxConnectErrors"].(int64); ok {
		tfMap["max_connect_errors"] = strconv.FormatInt(v, 10)
	}

	if v, ok := serviceParameters["maxConnections"].(int64); ok {
		tfMap["max_connections"] = v
	}

	if vMap, okMap := serviceParameters["maxHeapTableSize"].(map[string]interface{}); okMap {
		bytes, err := parseBytes(vMap["value"].(int64), vMap["dimension"].(string))

		if err == nil {
			tfMap["max_heap_table_size"] = strconv.FormatInt(bytes, 10)
		}
	}

	if v, ok := serviceParameters["options"].(map[string]interface{}); ok {
		tfMap["options"] = v
	}

	if v, ok := serviceParameters["pxcStrictMode"].(string); ok {
		tfMap["pxc_strict_mode"] = v
	}

	if v, ok := serviceParameters["tableOpenCache"].(int64); ok {
		tfMap["table_open_cache"] = v
	}

	if v, ok := serviceParameters["threadCacheSize"].(int64); ok {
		tfMap["thread_cache_size"] = v
	} else {
		tfMap["thread_cache_size"] = mySQLThreadCacheSizeDefault
	}

	if vMap, okMap := serviceParameters["tmpTableSize"].(map[string]interface{}); okMap {
		bytes, err := parseBytes(vMap["value"].(int64), vMap["dimension"].(string))

		if err == nil {
			tfMap["tmp_table_size"] = strconv.FormatInt(bytes, 10)
		}
	}

	if v, ok := serviceParameters["transactionIsolation"].(string); ok {
		tfMap["transaction_isolation"] = v
	}

	if v, ok := serviceParameters["vendor"].(string); ok {
		tfMap["vendor"] = v
	}

	if v, ok := serviceParameters["version"].(string); ok {
		tfMap["version"] = v
	}

	if v, ok := serviceParameters["waitTimeout"].(int64); ok {
		tfMap["wait_timeout"] = v
	}

	return tfMap
}

func (s mySQLManager) flattenUserParameters(userParameters UserParameters) map[string]interface{} {
	if userParameters == nil {
		return map[string]interface{}{}
	}

	tfMap := map[string]interface{}{}

	if v, ok := userParameters["host"].(string); ok {
		tfMap["host"] = v
	}

	if v, ok := userParameters["password"].(string); ok {
		tfMap["password"] = v
	}

	return tfMap
}

func (s mySQLManager) flattenDatabaseParameters(databaseParameters DatabaseParameters) map[string]interface{} {
	if databaseParameters == nil {
		return map[string]interface{}{}
	}

	tfMap := map[string]interface{}{}

	if v, ok := databaseParameters["backupId"].(string); ok {
		tfMap["backup_id"] = v
	}

	if v, ok := databaseParameters["backupDbName"].(string); ok {
		tfMap["backup_db_name"] = v
	}

	if v, ok := databaseParameters["charset"].(string); ok {
		tfMap["charset"] = v
	}

	if v, ok := databaseParameters["collate"].(string); ok {
		tfMap["collate"] = v
	}

	return tfMap
}

func (s mySQLManager) flattenDatabaseUserParameters(databaseUserParameters DatabaseUserParameters) map[string]interface{} {
	if databaseUserParameters == nil {
		return map[string]interface{}{}
	}

	tfMap := map[string]interface{}{}

	if v, ok := databaseUserParameters["options"].([]interface{}); ok {
		tfMap["options"] = v
	}

	if v, ok := databaseUserParameters["privileges"].([]interface{}); ok {
		tfMap["privileges"] = v
	}

	return tfMap
}
