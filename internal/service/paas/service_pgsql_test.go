package paas

import (
	"context"
	"encoding/json"
	"net/http"
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
		"version":                          "18",
		"work_mem":                         "8388608",
		"max_parallel_maintenance_workers": 3,
		"wal_keep_segments":                5,
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

func TestResourceServicePostgreSQLUpdateRequests(t *testing.T) {
	t.Parallel()

	user := []interface{}{map[string]interface{}{"name": "user1", "password": "password1"}}
	servicePath := "/services/" + testPostgreSQLServiceID

	for name, tc := range map[string]struct {
		new        map[string]interface{}
		wantEvents []string
		wantPatch  map[string]interface{}
	}{
		"parameter changed": {
			new:        map[string]interface{}{"version": "18", "work_mem": "8388608"},
			wantEvents: []string{"PATCH " + servicePath + "/parameters", "GET " + servicePath, "GET " + servicePath},
			wantPatch:  map[string]interface{}{"work_mem": map[string]interface{}{"dimension": "B", "value": 8388608.0}},
		},
		"users only": {
			new:        map[string]interface{}{"version": "18", "work_mem": "4194304", "user": user},
			wantEvents: []string{"PUT " + servicePath, "GET " + servicePath, "GET " + servicePath},
		},
	} {
		resource := ResourceService()
		initial := schema.TestResourceDataRaw(t, resource.Schema, testPostgreSQLServiceConfig(map[string]interface{}{"version": "18", "work_mem": "4194304"}))
		initial.SetId(testPostgreSQLServiceID)
		diff, err := resource.Diff(context.Background(), initial.State(), terraform.NewResourceConfigRaw(testPostgreSQLServiceConfig(tc.new)), nil)
		if err != nil {
			t.Fatalf("%s: calculating PostgreSQL diff: %s", name, err)
		}
		data, err := schema.InternalMap(resource.Schema).Data(initial.State(), diff)
		if err != nil {
			t.Fatalf("%s: building PostgreSQL resource data: %s", name, err)
		}

		var events []string
		var patch map[string]interface{}
		conn := testPaaSClientWithDebugBodyLogging(&recordingSDKLogger{}, roundTripFunc(func(request *http.Request) (*http.Response, error) {
			events = append(events, request.Method+" "+request.URL.Path)
			if request.Method == http.MethodPatch {
				patch = testGitLabInternalRequestParameters(t, request)
			}
			return jsonResponse(request, testPostgreSQLServiceResponse(t)), nil
		}))
		if diags := resourceServiceUpdate(context.Background(), data, testPaaSClientMeta(conn)); diags.HasError() {
			t.Fatalf("%s: updating PostgreSQL service: %#v", name, diags)
		}
		if !reflect.DeepEqual(events, tc.wantEvents) {
			t.Errorf("%s: unexpected PaaS requests: got %#v want %#v", name, events, tc.wantEvents)
		}
		if !reflect.DeepEqual(patch, tc.wantPatch) {
			t.Errorf("%s: unexpected PostgreSQL patch parameters: got %#v want %#v", name, patch, tc.wantPatch)
		}
	}
}

func testPostgreSQLServiceResponse(t *testing.T) string {
	t.Helper()

	body, err := json.Marshal(map[string]interface{}{"service": map[string]interface{}{
		"id": testPostgreSQLServiceID, "name": "tf-pgsql-test", "serviceType": services.ServiceTypePostgreSQL, "serviceClass": "database", "status": "READY",
		"instanceType": "c5.large", "rootVolumeType": "gp2", "rootVolumeSize": 32, "dataVolumeType": "gp2", "dataVolumeSize": 32,
		"subnetIds": []string{"subnet-12345678"}, "parameters": map[string]interface{}{"version": "18"},
	}})
	if err != nil {
		t.Fatalf("encoding PostgreSQL HTTP response fixture: %s", err)
	}
	return string(body)
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

const testPostgreSQLServiceID = "fm-cluster-12345678"

func TestResourceServicePostgreSQLManagedExternallyValidation(t *testing.T) {
	t.Parallel()

	user := []interface{}{map[string]interface{}{"name": "user1", "password": "password1"}}
	database := []interface{}{map[string]interface{}{"name": "db1", "owner": "user1"}}

	for name, tc := range map[string]struct {
		flags     map[string]interface{}
		pgsql     map[string]interface{}
		redis     bool
		errorText string
	}{
		"user blocks with users_managed_externally": {
			flags:     map[string]interface{}{"users_managed_externally": true},
			pgsql:     map[string]interface{}{"version": "18", "user": user},
			errorText: "pgsql.0.user must not be set when users_managed_externally is true",
		},
		"database blocks with databases_managed_externally": {
			flags:     map[string]interface{}{"databases_managed_externally": true},
			pgsql:     map[string]interface{}{"version": "18", "user": user, "database": database},
			errorText: "pgsql.0.database must not be set when databases_managed_externally is true",
		},
		"users_managed_externally for another service type": {
			flags:     map[string]interface{}{"users_managed_externally": true},
			redis:     true,
			errorText: "users_managed_externally is supported only for pgsql services",
		},
		"both flags without blocks": {
			flags: map[string]interface{}{"users_managed_externally": true, "databases_managed_externally": true},
			pgsql: map[string]interface{}{"version": "18"},
		},
		"external users with database blocks": {
			flags: map[string]interface{}{"users_managed_externally": true},
			pgsql: map[string]interface{}{"version": "18", "database": database},
		},
		"blocks without flags": {
			pgsql: map[string]interface{}{"version": "18", "user": user, "database": database},
		},
	} {
		config := testPostgreSQLServiceConfig(tc.pgsql)
		if tc.redis {
			delete(config, services.ServiceTypePostgreSQL)
			config[services.ServiceTypeRedis] = []interface{}{map[string]interface{}{"version": "7.2"}}
		}
		for k, v := range tc.flags {
			config[k] = v
		}

		_, err := ResourceService().Diff(context.Background(), nil, terraform.NewResourceConfigRaw(config), nil)
		if tc.errorText == "" {
			if err != nil {
				t.Errorf("%s: unexpected error: %s", name, err)
			}
			continue
		}
		if err == nil || !strings.Contains(err.Error(), tc.errorText) {
			t.Errorf("%s: expected %q, got %v", name, tc.errorText, err)
		}
	}
}

// The API replaces the whole lists of users and databases with ModifyService, so externally
// managed users and databases must never be sent, otherwise the cloud deletes them.
func TestResourceServicePostgreSQLManagedExternallyRequests(t *testing.T) {
	t.Parallel()

	user := []interface{}{map[string]interface{}{"name": "user1", "password": "password1"}}
	database := []interface{}{map[string]interface{}{"name": "db1", "owner": "user1"}}
	servicePath := "/services/" + testPostgreSQLServiceID

	for name, tc := range map[string]struct {
		newFlags         map[string]interface{}
		newPgsql         map[string]interface{}
		wantModify       bool
		wantUsersSent    bool
		wantDatabaseSent bool
	}{
		"both flags enabled and blocks removed": {
			newFlags: map[string]interface{}{"users_managed_externally": true, "databases_managed_externally": true},
			newPgsql: map[string]interface{}{"version": "18"},
		},
		"users flag enabled and database removed": {
			newFlags:         map[string]interface{}{"users_managed_externally": true},
			newPgsql:         map[string]interface{}{"version": "18"},
			wantModify:       true,
			wantDatabaseSent: true,
		},
		"flags disabled and blocks removed": {
			newPgsql:         map[string]interface{}{"version": "18"},
			wantModify:       true,
			wantUsersSent:    true,
			wantDatabaseSent: true,
		},
	} {
		resource := ResourceService()
		initial := schema.TestResourceDataRaw(t, resource.Schema, testPostgreSQLServiceConfig(map[string]interface{}{
			"version": "18", "user": user, "database": database,
		}))
		initial.SetId(testPostgreSQLServiceID)

		config := testPostgreSQLServiceConfig(tc.newPgsql)
		for k, v := range tc.newFlags {
			config[k] = v
		}
		diff, err := resource.Diff(context.Background(), initial.State(), terraform.NewResourceConfigRaw(config), nil)
		if err != nil {
			t.Fatalf("%s: calculating PostgreSQL diff: %s", name, err)
		}
		if diff.RequiresNew() {
			t.Fatalf("%s: switching flags must not replace the service", name)
		}
		data, err := schema.InternalMap(resource.Schema).Data(initial.State(), diff)
		if err != nil {
			t.Fatalf("%s: building PostgreSQL resource data: %s", name, err)
		}

		var modifyBody map[string]interface{}
		var mutations []string
		conn := testPaaSClientWithDebugBodyLogging(&recordingSDKLogger{}, roundTripFunc(func(request *http.Request) (*http.Response, error) {
			if request.Method != http.MethodGet {
				mutations = append(mutations, request.Method+" "+request.URL.Path)
			}
			if request.Method == http.MethodPut && request.URL.Path == servicePath {
				if err := json.NewDecoder(request.Body).Decode(&modifyBody); err != nil {
					t.Fatal(err)
				}
			}
			return jsonResponse(request, testPostgreSQLServiceWithUsersAndDatabasesResponse(t)), nil
		}))
		if diags := resourceServiceUpdate(context.Background(), data, testPaaSClientMeta(conn)); diags.HasError() {
			t.Fatalf("%s: updating PostgreSQL service: %#v", name, diags)
		}

		if !tc.wantModify {
			if len(mutations) > 0 {
				t.Errorf("%s: no requests must change the service, got %#v", name, mutations)
			}
		} else if modifyBody == nil {
			t.Errorf("%s: ModifyService was not called, requests: %#v", name, mutations)
		}
		if _, ok := modifyBody["users"]; ok != tc.wantUsersSent {
			t.Errorf("%s: users sent = %t, want %t: %#v", name, ok, tc.wantUsersSent, modifyBody)
		}
		if _, ok := modifyBody["databases"]; ok != tc.wantDatabaseSent {
			t.Errorf("%s: databases sent = %t, want %t: %#v", name, ok, tc.wantDatabaseSent, modifyBody)
		}

		// Externally managed users and databases returned by the API must not appear in the service state.
		for flag, key := range map[string]string{"users_managed_externally": "user", "databases_managed_externally": "database"} {
			blocks := data.Get(services.ServiceTypePostgreSQL + ".0." + key).([]interface{})
			if data.Get(flag).(bool) && len(blocks) != 0 {
				t.Errorf("%s: %s blocks must not be read into state when %s is true: %#v", name, key, flag, blocks)
			}
			if !data.Get(flag).(bool) && len(blocks) == 0 {
				t.Errorf("%s: %s blocks must be read into state when %s is false", name, key, flag)
			}
		}
	}
}

func testPostgreSQLServiceWithUsersAndDatabasesResponse(t *testing.T) string {
	t.Helper()

	body, err := json.Marshal(map[string]interface{}{"service": map[string]interface{}{
		"id": testPostgreSQLServiceID, "name": "tf-pgsql-test", "serviceType": services.ServiceTypePostgreSQL, "serviceClass": "database", "status": "READY",
		"instanceType": "c5.large", "rootVolumeType": "gp2", "rootVolumeSize": 32, "dataVolumeType": "gp2", "dataVolumeSize": 32,
		"subnetIds": []string{"subnet-12345678"}, "parameters": map[string]interface{}{"version": "18"},
		"users": []interface{}{
			map[string]interface{}{"id": "paas-user-1", "name": "user1", "parameters": map[string]interface{}{"password": "password1"}},
			map[string]interface{}{"id": "paas-user-2", "name": "external", "parameters": map[string]interface{}{"password": "password2"}},
		},
		"databases": []interface{}{
			map[string]interface{}{"id": "paas-db-1", "name": "db1", "parameters": map[string]interface{}{"owner": "user1"}},
		},
	}})
	if err != nil {
		t.Fatalf("encoding PostgreSQL HTTP response fixture: %s", err)
	}
	return string(body)
}
