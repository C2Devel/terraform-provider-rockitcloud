package services

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type postgreSQLManager struct {
	service
}

var PostgreSQL = postgreSQLManager{
	service{
		name:               ServiceTypePostgreSQL,
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

func (s postgreSQLManager) serviceParametersSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"autovacuum": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "ON",
			ValidateFunc: validation.StringInSlice([]string{"ON", "OFF"}, false),
		},
		"autovacuum_max_workers": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      3,
			ValidateFunc: validation.IntBetween(1, 262143),
		},
		"autovacuum_vacuum_cost_delay": {
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.Any(
				validation.IntInSlice([]int{-1}),
				validation.IntBetween(1, 100),
			),
		},
		"autovacuum_vacuum_cost_limit": {
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
			Default:  -1,
			ValidateFunc: validation.Any(
				validation.IntInSlice([]int{-1}),
				validation.IntBetween(1, 10000),
			),
		},
		"autovacuum_analyze_scale_factor": {
			Type:         schema.TypeFloat,
			Optional:     true,
			ForceNew:     true,
			Default:      0.1,
			ValidateFunc: validation.FloatBetween(0, 100),
		},
		"autovacuum_vacuum_scale_factor": {
			Type:         schema.TypeFloat,
			Optional:     true,
			ForceNew:     true,
			Default:      0.2,
			ValidateFunc: validation.FloatBetween(0, 100),
		},
		"effective_cache_size": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      524288,
			ValidateFunc: validation.IntBetween(1, 2147483647),
		},
		"effective_io_concurrency": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      1,
			ValidateFunc: validation.IntBetween(0, 1000),
		},
		"maintenance_work_mem": {
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
			Default:  64 * Megabyte,
			ValidateFunc: validation.All(
				validation.IntAtLeast(1*Megabyte),
				validation.IntDivisibleBy(Kilobyte),
			),
		},
		"max_connections": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      100,
			ValidateFunc: validation.IntBetween(1, 262143),
		},
		"max_wal_size": {
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
			Default:  1 * Gigabyte,
			ValidateFunc: validation.All(
				validation.IntAtLeast(2*Megabyte),
				validation.IntDivisibleBy(Megabyte),
			),
		},
		// TODO: add validation that depends on version value
		"max_parallel_maintenance_workers": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      postgreSQLMaxParallelMaintenanceWorkersDefault,
			ValidateFunc: validation.IntBetween(0, 1024),
		},
		"max_parallel_workers": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      8,
			ValidateFunc: validation.IntBetween(0, 1024),
		},
		"max_parallel_workers_per_gather": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      2,
			ValidateFunc: validation.IntBetween(0, 1024),
		},
		"max_worker_processes": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      8,
			ValidateFunc: validation.IntBetween(0, 262143),
		},
		"min_wal_size": {
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
			Default:  80 * Megabyte,
			ValidateFunc: validation.All(
				validation.IntAtLeast(32*Megabyte),
				validation.IntDivisibleBy(Megabyte),
			),
		},
		"options": {
			Type:     schema.TypeMap,
			Optional: true,
			ForceNew: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"replication_mode": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"asynchronous",
				"synchronous",
				"synchronous_strict",
			}, false),
		},
		"shared_buffers": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      1024,
			ValidateFunc: validation.IntBetween(16, 1073741823),
		},
		"version": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"10.21",
				"11.16",
				"12.11",
				"13.7",
				"14.4",
				"15.2",
			}, false),
		},
		// TODO: add validation that depends on version value
		"wal_keep_segments": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      postgreSQLWalKeepSegmentsDefault,
			ValidateFunc: validation.IntBetween(0, 2147483647),
		},
		"wal_buffers": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(8, 262143),
		},
		"work_mem": {
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
			Default:  4 * Megabyte,
			ValidateFunc: validation.All(
				validation.IntAtLeast(64*Kilobyte),
				validation.IntDivisibleBy(Kilobyte),
			),
		},
	}
}

func (s postgreSQLManager) serviceParametersDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"autovacuum": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"autovacuum_max_workers": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"autovacuum_vacuum_cost_delay": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"autovacuum_vacuum_cost_limit": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"autovacuum_analyze_scale_factor": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"autovacuum_vacuum_scale_factor": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"effective_cache_size": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"effective_io_concurrency": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"maintenance_work_mem": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"max_connections": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"max_wal_size": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"max_parallel_maintenance_workers": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"max_parallel_workers": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"max_parallel_workers_per_gather": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"max_worker_processes": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"min_wal_size": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"options": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"replication_mode": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"shared_buffers": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"version": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"wal_keep_segments": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"wal_buffers": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"work_mem": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func (s postgreSQLManager) userParametersSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"password": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(8, 128),
				validation.StringDoesNotContainAny("`'\"\\"),
			),
		},
	}
}

func (s postgreSQLManager) userParametersDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"password": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (s postgreSQLManager) databaseParametersSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"backup_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"backup_db_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		// TODO: add validation that depends on locale value
		"encoding": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "UTF8",
		},
		"extensions": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(postgreSQLDatabaseExtensions(), false),
			},
		},
		"locale": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "ru_RU.UTF-8",
			ValidateFunc: validation.StringInSlice(postgreSQLDatabaseLocales(), false),
		},
		"owner": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

func (s postgreSQLManager) databaseParametersDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"backup_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"backup_db_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"encoding": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"extensions": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"locale": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"owner": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (s postgreSQLManager) expandServiceParameters(tfMap map[string]interface{}) ServiceParameters {
	if tfMap == nil {
		return nil
	}

	serviceParameters := ServiceParameters{}

	if v, ok := tfMap["autovacuum"].(string); ok && v != "" {
		serviceParameters["autovacuum"] = v
	}

	if v, ok := tfMap["autovacuum_max_workers"].(int); ok && v != 0 {
		serviceParameters["autovacuum_max_workers"] = int64(v)
	}

	if v, ok := tfMap["autovacuum_vacuum_cost_delay"].(int); ok && v != 0 {
		serviceParameters["autovacuum_vacuum_cost_delay"] = int64(v)
	}

	if v, ok := tfMap["autovacuum_vacuum_cost_limit"].(int); ok && v != 0 {
		serviceParameters["autovacuum_vacuum_cost_limit"] = int64(v)
	}

	if v, ok := tfMap["autovacuum_analyze_scale_factor"].(float64); ok && v != 0.0 {
		serviceParameters["autovacuum_analyze_scale_factor"] = v
	}

	if v, ok := tfMap["autovacuum_vacuum_scale_factor"].(float64); ok && v != 0.0 {
		serviceParameters["autovacuum_vacuum_scale_factor"] = v
	}

	if v, ok := tfMap["effective_cache_size"].(int); ok && v != 0 {
		serviceParameters["effective_cache_size"] = int64(v)
	}

	if v, ok := tfMap["effective_io_concurrency"].(int); ok {
		serviceParameters["effective_io_concurrency"] = int64(v)
	}

	if v, ok := tfMap["maintenance_work_mem"].(int); ok && v != 0 {
		serviceParameters["maintenance_work_mem"] = map[string]interface{}{
			"dimension": B,
			"value":     int64(v),
		}
	}

	if v, ok := tfMap["max_connections"].(int); ok && v != 0 {
		serviceParameters["max_connections"] = int64(v)
	}

	if v, ok := tfMap["max_wal_size"].(int); ok && v != 0 {
		serviceParameters["max_wal_size"] = map[string]interface{}{
			"dimension": B,
			"value":     int64(v),
		}
	}

	if v, ok := tfMap["max_parallel_maintenance_workers"].(int); ok && v != postgreSQLMaxParallelMaintenanceWorkersDefault {
		serviceParameters["max_parallel_maintenance_workers"] = int64(v)
	}

	if v, ok := tfMap["max_parallel_workers"].(int); ok {
		serviceParameters["max_parallel_workers"] = int64(v)
	}

	if v, ok := tfMap["max_parallel_workers_per_gather"].(int); ok {
		serviceParameters["max_parallel_workers_per_gather"] = int64(v)
	}

	if v, ok := tfMap["max_worker_processes"].(int); ok {
		serviceParameters["max_worker_processes"] = int64(v)
	}

	if v, ok := tfMap["min_wal_size"].(int); ok && v != 0 {
		serviceParameters["min_wal_size"] = map[string]interface{}{
			"dimension": B,
			"value":     int64(v),
		}
	}

	if v, ok := tfMap["options"].(map[string]interface{}); ok && len(v) > 0 {
		serviceParameters["options"] = v
	}

	if v, ok := tfMap["replication_mode"].(string); ok && v != "" {
		serviceParameters["replication_mode"] = v
	}

	if v, ok := tfMap["shared_buffers"].(int); ok && v != 0 {
		serviceParameters["shared_buffers"] = int64(v)
	}

	if v, ok := tfMap["version"].(string); ok && v != "" {
		serviceParameters["version"] = v
	}

	if v, ok := tfMap["wal_keep_segments"].(int); ok && v != postgreSQLWalKeepSegmentsDefault {
		serviceParameters["wal_keep_segments"] = int64(v)
	}

	if v, ok := tfMap["wal_buffers"].(int); ok && v != 0 {
		serviceParameters["wal_buffers"] = int64(v)
	}

	if v, ok := tfMap["work_mem"].(int); ok && v != 0 {
		serviceParameters["work_mem"] = map[string]interface{}{
			"dimension": B,
			"value":     int64(v),
		}
	}

	return serviceParameters
}

func (s postgreSQLManager) expandUserParameters(tfMap map[string]interface{}) UserParameters {
	if tfMap == nil {
		return nil
	}

	userParameters := UserParameters{}

	if v, ok := tfMap["password"].(string); ok && v != "" {
		userParameters["password"] = v
	}

	return userParameters
}

func (s postgreSQLManager) expandDatabaseParameters(tfMap map[string]interface{}) DatabaseParameters {
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

	if v, ok := tfMap["encoding"].(string); ok && v != "" {
		databaseParameters["encoding"] = v
	}

	if v, ok := tfMap["extensions"].(*schema.Set); ok && v.Len() > 0 {
		databaseParameters["extensions"] = v.List()
	}

	if v, ok := tfMap["locale"].(string); ok && v != "" {
		databaseParameters["locale"] = v
	}

	if v, ok := tfMap["owner"].(string); ok && v != "" {
		databaseParameters["owner"] = v
	}

	return databaseParameters
}

func (s postgreSQLManager) flattenServiceParameters(serviceParameters ServiceParameters) map[string]interface{} {
	if serviceParameters == nil {
		return map[string]interface{}{}
	}

	tfMap := map[string]interface{}{}

	if v, ok := serviceParameters["autovacuum"].(string); ok {
		tfMap["autovacuum"] = v
	}

	if v, ok := serviceParameters["autovacuumMaxWorkers"].(int64); ok {
		tfMap["autovacuum_max_workers"] = v
	}

	if v, ok := serviceParameters["autovacuumVacuumCostDelay"].(int64); ok {
		tfMap["autovacuum_vacuum_cost_delay"] = v
	}

	if v, ok := serviceParameters["autovacuumVacuumCostLimit"].(int64); ok {
		tfMap["autovacuum_vacuum_cost_limit"] = v
	}

	if v, ok := serviceParameters["autovacuumAnalyzeScaleFactor"].(float64); ok {
		tfMap["autovacuum_analyze_scale_factor"] = v
	}

	if v, ok := serviceParameters["autovacuumVacuumScaleFactor"].(float64); ok {
		tfMap["autovacuum_vacuum_scale_factor"] = v
	}

	if v, ok := serviceParameters["effectiveCacheSize"].(int64); ok {
		tfMap["effective_cache_size"] = v
	}

	if v, ok := serviceParameters["effectiveIoConcurrency"].(int64); ok {
		tfMap["effective_io_concurrency"] = v
	}

	if vMap, okMap := serviceParameters["maintenanceWorkMem"].(map[string]interface{}); okMap {
		bytes, err := parseBytes(vMap["value"].(int64), vMap["dimension"].(string))

		if err == nil {
			tfMap["maintenance_work_mem"] = bytes
		}
	}

	if v, ok := serviceParameters["maxConnections"].(int64); ok {
		tfMap["max_connections"] = v
	}

	if vMap, okMap := serviceParameters["maxWalSize"].(map[string]interface{}); okMap {
		bytes, err := parseBytes(vMap["value"].(int64), vMap["dimension"].(string))

		if err == nil {
			tfMap["max_wal_size"] = bytes
		}
	}

	if v, ok := serviceParameters["maxParallelMaintenanceWorkers"].(int64); ok {
		tfMap["max_parallel_maintenance_workers"] = v
	} else {
		tfMap["max_parallel_maintenance_workers"] = postgreSQLMaxParallelMaintenanceWorkersDefault
	}

	if v, ok := serviceParameters["maxParallelWorkers"].(int64); ok {
		tfMap["max_parallel_workers"] = v
	}

	if v, ok := serviceParameters["maxParallelWorkersPerGather"].(int64); ok {
		tfMap["max_parallel_workers_per_gather"] = v
	}

	if v, ok := serviceParameters["maxWorkerProcesses"].(int64); ok {
		tfMap["max_worker_processes"] = v
	}

	if vMap, okMap := serviceParameters["minWalSize"].(map[string]interface{}); okMap {
		bytes, err := parseBytes(vMap["value"].(int64), vMap["dimension"].(string))

		if err == nil {
			tfMap["min_wal_size"] = bytes
		}
	}

	if v, ok := serviceParameters["options"].(map[string]interface{}); ok {
		tfMap["options"] = v
	}

	if v, ok := serviceParameters["replicationMode"].(string); ok {
		tfMap["replication_mode"] = v
	}

	if v, ok := serviceParameters["sharedBuffers"].(int64); ok {
		tfMap["shared_buffers"] = v
	}

	if v, ok := serviceParameters["version"].(string); ok {
		tfMap["version"] = v
	}

	if v, ok := serviceParameters["walKeepSegments"].(int64); ok {
		tfMap["wal_keep_segments"] = v
	} else {
		tfMap["wal_keep_segments"] = postgreSQLWalKeepSegmentsDefault
	}

	if v, ok := serviceParameters["walBuffers"].(int64); ok {
		tfMap["wal_buffers"] = v
	}

	if vMap, okMap := serviceParameters["workMem"].(map[string]interface{}); okMap {
		bytes, err := parseBytes(vMap["value"].(int64), vMap["dimension"].(string))

		if err == nil {
			tfMap["work_mem"] = bytes
		}
	}

	return tfMap
}

func (s postgreSQLManager) flattenUserParameters(userParameters UserParameters) map[string]interface{} {
	if userParameters == nil {
		return map[string]interface{}{}
	}

	tfMap := map[string]interface{}{}

	if v, ok := userParameters["password"].(string); ok {
		tfMap["password"] = v
	}

	return tfMap
}

func (s postgreSQLManager) flattenDatabaseParameters(databaseParameters DatabaseParameters) map[string]interface{} {
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

	if v, ok := databaseParameters["encoding"].(string); ok {
		tfMap["encoding"] = v
	}

	if v, ok := databaseParameters["extensions"].([]interface{}); ok {
		tfMap["extensions"] = v
	}

	if v, ok := databaseParameters["locale"].(string); ok {
		tfMap["locale"] = v
	}

	if v, ok := databaseParameters["owner"].(string); ok {
		tfMap["owner"] = v
	}

	return tfMap
}
