package paas

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/service/paas/services"
)

func TestResourceServicePostgreSQLParametersDoNotRequireReplacement(t *testing.T) {
	t.Parallel()

	resource := ResourceService()
	resourceData := schema.TestResourceDataRaw(t, resource.Schema, testPostgreSQLServiceConfig(map[string]interface{}{
		"version":         "18",
		"work_mem":        "4194304",
		"max_connections": 100,
		"options": map[string]interface{}{
			"jit": "on",
		},
	}))
	resourceData.SetId("fm-cluster-12345678")

	diff, err := resource.Diff(
		context.Background(),
		resourceData.State(),
		terraform.NewResourceConfigRaw(testPostgreSQLServiceConfig(map[string]interface{}{
			"version":              "18",
			"work_mem":             "8388608",
			"max_connections":      200,
			"max_parallel_workers": 0,
			"options": map[string]interface{}{
				"jit": "off",
			},
		})),
		nil,
	)
	if err != nil {
		t.Fatalf("calculating PostgreSQL diff: %s", err)
	}
	if diff == nil || diff.Empty() {
		t.Fatal("parameters update unexpectedly produced an empty diff")
	}
	if diff.RequiresNew() {
		t.Fatalf("parameters update unexpectedly requires PostgreSQL replacement: %#v", diff.Attributes)
	}
}

func TestResourceServicePostgreSQLVersionRequiresReplacement(t *testing.T) {
	t.Parallel()

	resource := ResourceService()

	for name, parameters := range map[string][2]map[string]interface{}{
		"version": {
			{"version": "17"},
			{"version": "18"},
		},
		"replication_mode": {
			{"version": "18", "replication_mode": "asynchronous"},
			{"version": "18", "replication_mode": "synchronous"},
		},
	} {
		resourceData := schema.TestResourceDataRaw(t, resource.Schema, testPostgreSQLServiceConfig(parameters[0]))
		resourceData.SetId("fm-cluster-12345678")

		diff, err := resource.Diff(
			context.Background(),
			resourceData.State(),
			terraform.NewResourceConfigRaw(testPostgreSQLServiceConfig(parameters[1])),
			nil,
		)
		if err != nil {
			t.Fatalf("calculating PostgreSQL %s diff: %s", name, err)
		}
		if diff == nil || !diff.RequiresNew() {
			t.Fatalf("%s update must require PostgreSQL replacement: %#v", name, diff)
		}
	}
}

func TestResourceServicePostgreSQLRemovedParameterKeepsValue(t *testing.T) {
	t.Parallel()

	resource := ResourceService()
	resourceData := schema.TestResourceDataRaw(t, resource.Schema, testPostgreSQLServiceConfig(map[string]interface{}{
		"version":  "18",
		"work_mem": "8388608",
	}))
	resourceData.SetId("fm-cluster-12345678")

	diff, err := resource.Diff(
		context.Background(),
		resourceData.State(),
		terraform.NewResourceConfigRaw(testPostgreSQLServiceConfig(map[string]interface{}{
			"version": "18",
		})),
		nil,
	)
	if err != nil {
		t.Fatalf("calculating PostgreSQL diff: %s", err)
	}
	if diff == nil {
		return
	}
	for key, attrDiff := range diff.Attributes {
		if strings.HasPrefix(key, services.ServiceTypePostgreSQL+".") {
			t.Fatalf("removed parameter unexpectedly produced a diff: %s: %#v", key, attrDiff)
		}
	}
}

func TestChangedServiceParametersPostgreSQL(t *testing.T) {
	t.Parallel()

	got := testPostgreSQLChangedServiceParameters(t,
		map[string]interface{}{
			"version":              "18",
			"work_mem":             "4194304",
			"max_connections":      100,
			"max_parallel_workers": 4,
			"options": map[string]interface{}{
				"jit": "on",
			},
		},
		map[string]interface{}{
			"version":              "18",
			"work_mem":             "8388608",
			"max_connections":      100,
			"max_parallel_workers": 0,
			"options": map[string]interface{}{
				"jit": "on",
			},
		},
	)
	want := services.ServiceParameters{
		"work_mem": map[string]interface{}{
			"dimension": services.B,
			"value":     int64(8388608),
		},
		"max_parallel_workers": int64(0),
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected PostgreSQL patch parameters: got %#v want %#v", got, want)
	}
}

func TestChangedServiceParametersPostgreSQLUsersOnly(t *testing.T) {
	t.Parallel()

	user := []interface{}{
		map[string]interface{}{
			"name":     "user1",
			"password": "password1",
		},
	}

	got := testPostgreSQLChangedServiceParameters(t,
		map[string]interface{}{
			"version":         "18",
			"max_connections": 100,
		},
		map[string]interface{}{
			"version":         "18",
			"max_connections": 100,
			"user":            user,
		},
	)

	if len(got) != 0 {
		t.Fatalf("users change unexpectedly produced patch parameters: %#v", got)
	}
}

func TestChangedServiceParametersPostgreSQLLoggingMonitoring(t *testing.T) {
	t.Parallel()

	logging := func(logTo string, tags ...interface{}) []interface{} {
		return []interface{}{map[string]interface{}{"log_to": logTo, "logging_tags": tags}}
	}
	monitoring := func(monitorBy string, labels map[string]interface{}) []interface{} {
		return []interface{}{map[string]interface{}{"monitor_by": monitorBy, "monitoring_labels": labels}}
	}

	for name, tc := range map[string]struct {
		old, new map[string]interface{}
		want     services.ServiceParameters
	}{
		"logging tags removed": {
			old: map[string]interface{}{"version": "18", "logging": logging("fm-cluster-1", "a", "b")},
			new: map[string]interface{}{"version": "18", "logging": logging("fm-cluster-1")},
			want: services.ServiceParameters{
				"logging":      true,
				"log_to":       "fm-cluster-1",
				"logging_tags": []interface{}{},
			},
		},
		"logging tags unchanged": {
			old: map[string]interface{}{"version": "18", "logging": logging("fm-cluster-1", "a")},
			new: map[string]interface{}{"version": "18", "logging": logging("fm-cluster-2", "a")},
			want: services.ServiceParameters{
				"logging": true,
				"log_to":  "fm-cluster-2",
			},
		},
		"logging disabled": {
			old: map[string]interface{}{"version": "18", "logging": logging("fm-cluster-1", "a")},
			new: map[string]interface{}{"version": "18"},
			want: services.ServiceParameters{
				"logging": false,
			},
		},
		"monitoring labels removed": {
			old: map[string]interface{}{"version": "18", "monitoring": monitoring("fm-cluster-1", map[string]interface{}{"env": "prod"})},
			new: map[string]interface{}{"version": "18", "monitoring": monitoring("fm-cluster-1", nil)},
			want: services.ServiceParameters{
				"monitoring":        true,
				"monitor_by":        "fm-cluster-1",
				"monitoring_labels": map[string]interface{}{},
			},
		},
	} {
		if got := testPostgreSQLChangedServiceParameters(t, tc.old, tc.new); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s: unexpected PostgreSQL patch parameters: got %#v want %#v", name, got, tc.want)
		}
	}
}

func TestExpandServiceParametersPostgreSQLOmitsUnsetParameters(t *testing.T) {
	t.Parallel()

	resource := ResourceService()
	resourceData := schema.TestResourceDataRaw(t, resource.Schema, testPostgreSQLServiceConfig(map[string]interface{}{
		"version":              "18",
		"effective_cache_size": "4294967296",
		"shared_buffers":       "1073741824",
	}))

	got := services.Manager(services.ServiceTypePostgreSQL).ExpandServiceParameters(
		resourceData.Get(services.ServiceTypePostgreSQL).([]interface{})[0].(map[string]interface{}),
	)
	want := services.ServiceParameters{
		"logging":    false,
		"monitoring": false,
		"version":    "18",
		"effective_cache_size": map[string]interface{}{
			"dimension": services.B,
			"value":     int64(4294967296),
		},
		"shared_buffers": map[string]interface{}{
			"dimension": services.B,
			"value":     int64(1073741824),
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected PostgreSQL create parameters: got %#v want %#v", got, want)
	}
}

func testPostgreSQLChangedServiceParameters(t *testing.T, oldParameters, newParameters map[string]interface{}) services.ServiceParameters {
	t.Helper()

	resource := ResourceService()
	resourceData := schema.TestResourceDataRaw(t, resource.Schema, testPostgreSQLServiceConfig(oldParameters))
	resourceData.SetId("fm-cluster-12345678")
	state := resourceData.State()

	diff, err := resource.Diff(
		context.Background(),
		state,
		terraform.NewResourceConfigRaw(testPostgreSQLServiceConfig(newParameters)),
		nil,
	)
	if err != nil {
		t.Fatalf("calculating PostgreSQL diff: %s", err)
	}

	d, err := schema.InternalMap(resource.Schema).Data(state, diff)
	if err != nil {
		t.Fatalf("building PostgreSQL resource data: %s", err)
	}

	parameters := services.Manager(services.ServiceTypePostgreSQL).ExpandServiceParameters(
		d.Get(services.ServiceTypePostgreSQL).([]interface{})[0].(map[string]interface{}),
	)

	return changedServiceParameters(d, services.ServiceTypePostgreSQL, parameters)
}

func testPostgreSQLServiceConfig(pgsqlParameters map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"name":          "tf-pgsql-test",
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
		services.ServiceTypePostgreSQL: []interface{}{
			pgsqlParameters,
		},
	}
}
