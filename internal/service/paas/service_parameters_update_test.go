package paas

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/service/paas/services"
)

type serviceParameterChange struct {
	parameter string
	old       interface{}
	new       interface{}
}

// The values mirror c2 (ps/bootstrap/paas/{pgsql,mysql}): parameters with
// allow_modify=True are updated in place through ModifyServiceParameters,
// the others can not be modified for an existing service.
var (
	postgreSQLEditableParameterChanges = []serviceParameterChange{
		{"autovacuum", "ON", "OFF"},
		{"autovacuum_max_workers", 3, 4},
		{"autovacuum_vacuum_cost_delay", 2, 5},
		{"autovacuum_vacuum_cost_limit", 200, 400},
		{"autovacuum_analyze_scale_factor", 0.1, 0.2},
		{"autovacuum_vacuum_scale_factor", 0.2, 0.3},
		{"effective_cache_size", "4294967296", "8589934592"},
		{"effective_io_concurrency", 1, 2},
		{"maintenance_work_mem", "67108864", "134217728"},
		{"max_connections", 100, 200},
		{"max_wal_size", "1073741824", "2147483648"},
		{"min_wal_size", "83886080", "167772160"},
		{"max_parallel_maintenance_workers", 2, 3},
		{"max_parallel_workers", 8, 16},
		{"max_parallel_workers_per_gather", 2, 4},
		{"max_worker_processes", 8, 16},
		{"shared_buffers", "134217728", "268435456"},
		{"wal_buffers", 16, 32},
		{"work_mem", "4194304", "8388608"},
		{"options", map[string]interface{}{"log_min_duration_statement": "100"}, map[string]interface{}{"log_min_duration_statement": "200"}},
	}
	postgreSQLImmutableParameterChanges = []serviceParameterChange{
		{"replication_mode", "asynchronous", "synchronous"},
		{"wal_keep_segments", 16, 64},
		{"version", "15", "16"},
	}

	mySQLEditableParameterChanges = []serviceParameterChange{
		{"connect_timeout", 10, 20},
		{"galera_options", map[string]interface{}{"gcache.page_size": "128M"}, map[string]interface{}{"gcache.page_size": "256M"}},
		{"gcs_fc_factor", 0.5, 0.6},
		{"gcs_fc_limit", 16, 32},
		{"innodb_buffer_pool_size", "134217728", "268435456"},
		{"innodb_change_buffering", "all", "none"},
		{"innodb_flush_log_at_trx_commit", 1, 2},
		{"innodb_io_capacity", "200", "400"},
		{"innodb_io_capacity_max", "2000", "4000"},
		{"innodb_log_file_size", "50331648", "100663296"},
		{"innodb_log_files_in_group", 2, 3},
		{"innodb_purge_threads", 4, 8},
		{"innodb_strict_mode", "OFF", "ON"},
		{"max_allowed_packet", 16 * services.Megabyte, 32 * services.Megabyte},
		{"max_connect_errors", "100", "200"},
		{"max_connections", 100, 200},
		{"max_heap_table_size", "16777216", "33554432"},
		{"options", map[string]interface{}{"slow_query_log": "ON"}, map[string]interface{}{"slow_query_log": "OFF"}},
		{"table_open_cache", 2000, 4000},
		{"thread_cache_size", 8, 16},
		{"tmp_table_size", "16777216", "33554432"},
	}
	mySQLImmutableParameterChanges = []serviceParameterChange{
		{"gcache_size", 128 * services.Megabyte, 256 * services.Megabyte},
		{"gcs_fc_master_slave", "true", "false"},
		{"gcs_fc_single_primary", "true", "false"},
		{"innodb_buffer_pool_instances", 1, 2},
		{"innodb_thread_concurrency", 0, 8},
		{"innodb_sync_array_size", 1, 2},
		{"pxc_strict_mode", "ENFORCING", "PERMISSIVE"},
		{"transaction_isolation", "REPEATABLE-READ", "READ-COMMITTED"},
		{"wait_timeout", 28800, 3600},
		{"vendor", "mysql", "percona"},
		{"version", "8.0.32", "8.0.28"},
	}
)

func TestResourceServiceParameterChangesReplacement(t *testing.T) {
	t.Parallel()

	for _, testCase := range []struct {
		serviceType string
		base        map[string]interface{}
		editable    []serviceParameterChange
		immutable   []serviceParameterChange
	}{
		{
			serviceType: services.ServiceTypePostgreSQL,
			base:        map[string]interface{}{"version": "15"},
			editable:    postgreSQLEditableParameterChanges,
			immutable:   postgreSQLImmutableParameterChanges,
		},
		{
			serviceType: services.ServiceTypeMySQL,
			base:        map[string]interface{}{"vendor": "mysql", "version": "8.0.32"},
			editable:    mySQLEditableParameterChanges,
			immutable:   mySQLImmutableParameterChanges,
		},
	} {
		testCase := testCase

		for _, wantReplacement := range []bool{false, true} {
			wantReplacement := wantReplacement
			changes := testCase.editable
			if wantReplacement {
				changes = testCase.immutable
			}

			for _, change := range changes {
				change := change
				name := testCase.serviceType + "/" + change.parameter
				if wantReplacement {
					name += "/replaces"
				} else {
					name += "/in-place"
				}

				t.Run(name, func(t *testing.T) {
					t.Parallel()

					resource := ResourceService()
					resourceData := schema.TestResourceDataRaw(t, resource.Schema,
						testServiceParametersConfig(testCase.serviceType, testCase.base, change.parameter, change.old))
					resourceData.SetId("fm-cluster-12345678")

					diff, err := resource.Diff(
						context.Background(),
						resourceData.State(),
						terraform.NewResourceConfigRaw(
							testServiceParametersConfig(testCase.serviceType, testCase.base, change.parameter, change.new)),
						nil,
					)
					if err != nil {
						t.Fatalf("calculating diff: %s", err)
					}
					if diff == nil || diff.Empty() {
						t.Fatal("changing the parameter produced an empty diff")
					}

					prefix := testCase.serviceType + ".0." + change.parameter
					found := false
					for key := range diff.Attributes {
						if key == prefix || strings.HasPrefix(key, prefix+".") {
							found = true
						}
					}
					if !found {
						t.Fatalf("diff does not contain %s: %#v", prefix, diff.Attributes)
					}

					if got := diff.RequiresNew(); got != wantReplacement {
						t.Fatalf("RequiresNew = %t, want %t: %#v", got, wantReplacement, diff.Attributes)
					}
				})
			}
		}
	}
}

// Creating a service with any of the parameters must keep working: the
// parameters stay in the plan of a new resource.
func TestResourceServiceCreateWithEditableParameters(t *testing.T) {
	t.Parallel()

	for _, testCase := range []struct {
		serviceType string
		base        map[string]interface{}
		changes     []serviceParameterChange
	}{
		{services.ServiceTypePostgreSQL, map[string]interface{}{"version": "15"}, postgreSQLEditableParameterChanges},
		{services.ServiceTypeMySQL, map[string]interface{}{"vendor": "mysql", "version": "8.0.32"}, mySQLEditableParameterChanges},
	} {
		testCase := testCase

		t.Run(testCase.serviceType, func(t *testing.T) {
			t.Parallel()

			configuration := testServiceParametersConfig(testCase.serviceType, testCase.base, "", nil)
			block := configuration[testCase.serviceType].([]interface{})[0].(map[string]interface{})
			for _, change := range testCase.changes {
				block[change.parameter] = change.new
			}

			resource := ResourceService()
			diff, err := resource.Diff(
				context.Background(),
				nil,
				terraform.NewResourceConfigRaw(configuration),
				nil,
			)
			if err != nil {
				t.Fatalf("calculating diff: %s", err)
			}
			if diff == nil {
				t.Fatal("create produced no diff")
			}

			for _, change := range testCase.changes {
				prefix := testCase.serviceType + ".0." + change.parameter
				found := false
				for key := range diff.Attributes {
					if key == prefix || strings.HasPrefix(key, prefix+".") {
						found = true
					}
				}
				if !found {
					t.Errorf("create diff does not contain %s", prefix)
				}
			}
		})
	}
}

func testServiceParametersConfig(
	serviceType string,
	base map[string]interface{},
	parameter string,
	value interface{},
) map[string]interface{} {
	block := map[string]interface{}{}
	for key, v := range base {
		block[key] = v
	}
	if parameter != "" {
		block[parameter] = value
	}

	return map[string]interface{}{
		"name":          "tf-" + serviceType + "-test",
		"instance_type": "c5.large",
		"root_volume": []interface{}{
			map[string]interface{}{
				"type": "gp2",
				"size": 32,
			},
		},
		"data_volume": []interface{}{
			map[string]interface{}{
				"type": "gp2",
				"size": 32,
			},
		},
		"subnet_ids": []interface{}{"subnet-12345678"},
		serviceType: []interface{}{
			block,
		},
	}
}
