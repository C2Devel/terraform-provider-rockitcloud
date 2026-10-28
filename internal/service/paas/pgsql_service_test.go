package paas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/service/paas/services"
)

func TestResourcePostgreSQLServiceSchema(t *testing.T) {
	t.Parallel()

	serviceSchema := ResourcePostgreSQLService().Schema

	for _, serviceType := range services.ManagedServiceTypes() {
		if _, ok := serviceSchema[serviceType]; ok != (serviceType == services.ServiceTypePostgreSQL) {
			t.Errorf("block %s present = %t, want only the pgsql block", serviceType, ok)
		}
	}

	pgsql := serviceSchema[services.ServiceTypePostgreSQL]
	if !pgsql.Required || pgsql.Optional || len(pgsql.ExactlyOneOf) != 0 {
		t.Errorf("pgsql block must be required without ExactlyOneOf: %#v", pgsql)
	}
	for _, k := range []string{"user", "database"} {
		if _, ok := pgsql.Elem.(*schema.Resource).Schema[k]; ok {
			t.Errorf("pgsql.%s must not be supported", k)
		}
	}

	for _, k := range []string{"users_managed_externally", "databases_managed_externally", "additional_roles", "coordinator"} {
		if s := serviceSchema[k]; s == nil || !s.Computed || s.Optional || s.Required {
			t.Errorf("%s must be a read-only attribute: %#v", k, s)
		}
	}

	// New common arguments of aws_paas_service must also appear in aws_paas_pgsql.
	for k := range ResourceService().Schema {
		if _, ok := serviceSchema[k]; !ok && !containsString(services.ManagedServiceTypes(), k) {
			t.Errorf("aws_paas_service argument %s is missing in aws_paas_pgsql", k)
		}
	}
}

func TestResourcePostgreSQLServiceConfigValidation(t *testing.T) {
	t.Parallel()

	resource := ResourcePostgreSQLService()

	for name, tc := range map[string]struct {
		change    func(map[string]interface{})
		errorText string
	}{
		"valid configuration": {},
		"inline user block": {
			change: func(c map[string]interface{}) {
				c[services.ServiceTypePostgreSQL] = []interface{}{map[string]interface{}{
					"version": "18", "user": []interface{}{map[string]interface{}{"name": "user1", "password": "password1"}},
				}}
			},
			errorText: `"user"`,
		},
		"inline database block": {
			change: func(c map[string]interface{}) {
				c[services.ServiceTypePostgreSQL] = []interface{}{map[string]interface{}{
					"version": "18", "database": []interface{}{map[string]interface{}{"name": "db1", "owner": "user1"}},
				}}
			},
			errorText: `"database"`,
		},
		"block of another service type": {
			change: func(c map[string]interface{}) {
				c[services.ServiceTypeKafka] = []interface{}{map[string]interface{}{"version": "3.7.0"}}
			},
			errorText: `"kafka"`,
		},
		"read-only coordinator": {
			change: func(c map[string]interface{}) {
				c["coordinator"] = []interface{}{map[string]interface{}{"instance_type": "c5.large"}}
			},
			errorText: `"coordinator"`,
		},
		"read-only users_managed_externally": {
			change:    func(c map[string]interface{}) { c["users_managed_externally"] = false },
			errorText: `"users_managed_externally"`,
		},
		"missing pgsql block": {
			change:    func(c map[string]interface{}) { delete(c, services.ServiceTypePostgreSQL) },
			errorText: `"pgsql"`,
		},
	} {
		config := testPostgreSQLServiceConfig(map[string]interface{}{"version": "18"})
		config["security_group_ids"] = []interface{}{"sg-12345678"}
		if tc.change != nil {
			tc.change(config)
		}

		diags := resource.Validate(terraform.NewResourceConfigRaw(config))
		if tc.errorText == "" {
			if diags.HasError() {
				t.Errorf("%s: unexpected errors: %#v", name, diags)
			}
			continue
		}
		if !diags.HasError() || !strings.Contains(fmt.Sprintf("%#v", diags), tc.errorText) {
			t.Errorf("%s: expected an error about %s, got %#v", name, tc.errorText, diags)
		}
	}
}

func TestResourcePostgreSQLServiceReplicationModeValidation(t *testing.T) {
	t.Parallel()

	resource := ResourcePostgreSQLService()

	for name, tc := range map[string]struct {
		highAvailability bool
		replicationMode  interface{}
		wantError        bool
	}{
		"high availability without replication mode": {highAvailability: true, wantError: true},
		"high availability with replication mode":    {highAvailability: true, replicationMode: "synchronous"},
		"high availability with unknown mode":        {highAvailability: true, replicationMode: testUnknownConfigValue},
		"single node without replication mode":       {},
	} {
		pgsqlParameters := map[string]interface{}{"version": "16"}
		if tc.replicationMode != nil {
			pgsqlParameters["replication_mode"] = tc.replicationMode
		}
		config := testPostgreSQLServiceConfig(pgsqlParameters)
		config["security_group_ids"] = []interface{}{"sg-12345678"}
		config["high_availability"] = tc.highAvailability

		_, err := resource.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(config), nil)
		if tc.wantError {
			if err == nil || !strings.Contains(err.Error(), "pgsql.0.replication_mode must be set if high_availability is true") {
				t.Errorf("%s: expected a replication_mode error, got %v", name, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s: unexpected error: %s", name, err)
		}
	}
}

// testUnknownConfigValue marks a value that is known only during apply in a raw configuration.
const testUnknownConfigValue = "74D93920-ED26-11E3-AC10-0800200C9A66"

func TestResourcePostgreSQLServiceCreate(t *testing.T) {
	t.Parallel()

	var mutations []string
	conn := testPaaSClientWithDebugBodyLogging(&recordingSDKLogger{}, roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.Method != http.MethodGet {
			mutations = append(mutations, request.Method+" "+request.URL.Path)
		}
		if request.Method == http.MethodPost && request.URL.Path == "/services" {
			var body map[string]interface{}
			if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if body["serviceType"] != services.ServiceTypePostgreSQL {
				t.Errorf("CreateService must create a pgsql service: %#v", body["serviceType"])
			}
			for _, k := range []string{"users", "databases"} {
				if _, ok := body[k]; ok {
					t.Errorf("CreateService must not send %s: %#v", k, body)
				}
			}
			return jsonResponse(request, `{"service":{"id":"`+testPostgreSQLServiceID+`","serviceType":"pgsql"}}`), nil
		}
		return jsonResponse(request, testPostgreSQLServiceWithUsersAndDatabasesResponse(t)), nil
	}))

	data := schema.TestResourceDataRaw(t, ResourcePostgreSQLService().Schema, testPostgreSQLServiceConfig(map[string]interface{}{"version": "18"}))
	data.MarkNewResource()

	if diags := resourcePostgreSQLServiceCreate(context.Background(), data, testPaaSClientMeta(conn)); diags.HasError() {
		t.Fatalf("creating PostgreSQL service: %#v", diags)
	}
	if !reflect.DeepEqual(mutations, []string{"POST /services"}) {
		t.Errorf("unexpected mutating requests: %#v", mutations)
	}
	if data.Id() != testPostgreSQLServiceID || data.Get("service_type") != services.ServiceTypePostgreSQL {
		t.Errorf("ID = %q, service_type = %v", data.Id(), data.Get("service_type"))
	}
	// The pgsql block read from the API must be written to the state.
	if data.Get("pgsql.0.version") != "18" {
		t.Errorf("pgsql block was not read into the state: version = %q", data.Get("pgsql.0.version"))
	}
	for _, k := range []string{"users_managed_externally", "databases_managed_externally"} {
		if !data.Get(k).(bool) {
			t.Errorf("%s must be true", k)
		}
	}
}

// The same guarantees must hold for aws_paas_service with the flags enabled and for
// aws_paas_pgsql, because both resources share the CRUD functions.
func TestPostgreSQLServicesDoNotSendExternalUsersAndDatabases(t *testing.T) {
	t.Parallel()

	servicePath := "/services/" + testPostgreSQLServiceID

	for name, tc := range map[string]struct {
		resource *schema.Resource
		flags    map[string]interface{}
	}{
		"aws_paas_service with flags": {
			resource: ResourceService(),
			flags:    map[string]interface{}{"users_managed_externally": true, "databases_managed_externally": true},
		},
		"aws_paas_pgsql": {
			resource: ResourcePostgreSQLService(),
		},
	} {
		oldConfig := testPostgreSQLServiceConfig(map[string]interface{}{"version": "18", "work_mem": "4194304"})
		newConfig := testPostgreSQLServiceConfig(map[string]interface{}{"version": "18", "work_mem": "8388608"})
		for k, v := range tc.flags {
			oldConfig[k] = v
			newConfig[k] = v
		}

		initial := schema.TestResourceDataRaw(t, tc.resource.Schema, oldConfig)
		initial.SetId(testPostgreSQLServiceID)
		if tc.flags == nil {
			setPostgreSQLServiceExternalManagement(initial)
		}

		diff, err := tc.resource.Diff(context.Background(), initial.State(), terraform.NewResourceConfigRaw(newConfig), nil)
		if err != nil || diff == nil || diff.RequiresNew() {
			t.Fatalf("%s: changing work_mem must update the service in place: %v", name, err)
		}
		data, err := schema.InternalMap(tc.resource.Schema).Data(initial.State(), diff)
		if err != nil {
			t.Fatalf("%s: building resource data: %s", name, err)
		}

		var mutations []string
		conn := testPaaSClientWithDebugBodyLogging(&recordingSDKLogger{}, roundTripFunc(func(request *http.Request) (*http.Response, error) {
			if request.Method != http.MethodGet {
				mutations = append(mutations, request.Method+" "+request.URL.Path)
			}
			return jsonResponse(request, testPostgreSQLServiceWithUsersAndDatabasesResponse(t)), nil
		}))

		if diags := tc.resource.UpdateContext(context.Background(), data, testPaaSClientMeta(conn)); diags.HasError() {
			t.Fatalf("%s: updating service: %#v", name, diags)
		}
		if !reflect.DeepEqual(mutations, []string{"PATCH " + servicePath + "/parameters"}) {
			t.Errorf("%s: only the parameters must be patched, got %#v", name, mutations)
		}
		if data.Get(services.ServiceTypePostgreSQL+".0.version") != "18" {
			t.Errorf("%s: pgsql block was not read into the state", name)
		}
		if _, ok := data.GetOk(services.ServiceTypePostgreSQL + ".0.user.0.name"); ok {
			t.Errorf("%s: users returned by the API must not be read into the service state", name)
		}
		if _, ok := data.GetOk(services.ServiceTypePostgreSQL + ".0.database.0.name"); ok {
			t.Errorf("%s: databases returned by the API must not be read into the service state", name)
		}
	}
}

func TestResourcePostgreSQLServiceImport(t *testing.T) {
	t.Parallel()

	for _, serviceType := range []string{services.ServiceTypePostgreSQL, services.ServiceTypeRedis, services.ServiceTypeGitLab} {
		conn := testPaaSClientWithDebugBodyLogging(&recordingSDKLogger{}, roundTripFunc(func(request *http.Request) (*http.Response, error) {
			if request.Method != http.MethodGet {
				t.Fatalf("import must not change the service: %s %s", request.Method, request.URL.Path)
			}
			if serviceType == services.ServiceTypePostgreSQL {
				return jsonResponse(request, testPostgreSQLServiceWithUsersAndDatabasesResponse(t)), nil
			}
			return jsonResponse(request, `{"service":{"id":"`+testPostgreSQLServiceID+`","serviceType":"`+serviceType+`","status":"READY"}}`), nil
		}))

		data := schema.TestResourceDataRaw(t, ResourcePostgreSQLService().Schema, map[string]interface{}{})
		data.SetId(testPostgreSQLServiceID)

		diags := resourcePostgreSQLServiceRead(context.Background(), data, testPaaSClientMeta(conn))

		if serviceType != services.ServiceTypePostgreSQL {
			if !diags.HasError() || !strings.Contains(diags[0].Summary, `supports only "pgsql" services`) {
				t.Errorf("%s: importing another service type must fail, got %#v", serviceType, diags)
			}
			continue
		}
		if diags.HasError() {
			t.Fatalf("importing PostgreSQL service: %#v", diags)
		}
		if data.Get("service_type") != services.ServiceTypePostgreSQL || data.Get("pgsql.0.version") != "18" {
			t.Errorf("service was not read: service_type = %v, version = %v", data.Get("service_type"), data.Get("pgsql.0.version"))
		}
		if !data.Get("users_managed_externally").(bool) || !data.Get("databases_managed_externally").(bool) {
			t.Error("imported service must manage users and databases externally")
		}
	}
}

func TestResourcePostgreSQLServiceImportArbitratorRequired(t *testing.T) {
	t.Parallel()

	resource := ResourcePostgreSQLService()

	instances := func(roles ...string) []interface{} {
		list := make([]interface{}, 0, len(roles))
		for index, role := range roles {
			list = append(list, map[string]interface{}{"index": index, "role": role, "status": "READY"})
		}
		return list
	}

	// The instance roles match the services created in K2 Cloud.
	for name, tc := range map[string]struct {
		highAvailability bool
		instances        []interface{}
		want             string
	}{
		"single node":                         {instances: instances("node"), want: "false"},
		"two nodes and an arbitrator":         {highAvailability: true, instances: instances("node", "node", "arbitrator"), want: "true"},
		"three nodes":                         {highAvailability: true, instances: instances("node", "node", "node"), want: "false"},
		"high availability without instances": {highAvailability: true},
	} {
		var response map[string]interface{}
		if err := json.Unmarshal([]byte(testPostgreSQLServiceWithUsersAndDatabasesResponse(t)), &response); err != nil {
			t.Fatal(err)
		}
		service := response["service"].(map[string]interface{})
		service["highAvailability"] = tc.highAvailability
		if tc.highAvailability {
			service["parameters"].(map[string]interface{})["replicationMode"] = "synchronous"
		}
		if tc.instances != nil {
			service["instances"] = tc.instances
		}
		body, err := json.Marshal(response)
		if err != nil {
			t.Fatal(err)
		}

		conn := testPaaSClientWithDebugBodyLogging(&recordingSDKLogger{}, roundTripFunc(func(request *http.Request) (*http.Response, error) {
			return jsonResponse(request, string(body)), nil
		}))

		// Import starts from an empty state without schema defaults.
		data := resource.Data(&terraform.InstanceState{ID: testPostgreSQLServiceID})

		if diags := resourcePostgreSQLServiceRead(context.Background(), data, testPaaSClientMeta(conn)); diags.HasError() {
			t.Fatalf("%s: importing PostgreSQL service: %#v", name, diags)
		}

		value, ok := data.State().Attributes["arbitrator_required"]
		if tc.want == "" {
			// Without instances the arbitrator cannot be determined, so the value is not guessed.
			if ok && value != "" {
				t.Errorf("%s: arbitrator_required must not be guessed, got %q", name, value)
			}
			continue
		}
		if !ok || value != tc.want {
			t.Errorf("%s: arbitrator_required must be restored as %s, got %q (set: %t)", name, tc.want, value, ok)
			continue
		}

		pgsqlParameters := map[string]interface{}{"version": "18"}
		config := testPostgreSQLServiceConfig(pgsqlParameters)
		if tc.highAvailability {
			pgsqlParameters["replication_mode"] = "synchronous"
			config["high_availability"] = true
			config["arbitrator_required"] = tc.want == "true"
		}
		diff, err := resource.Diff(context.Background(), data.State(), terraform.NewResourceConfigRaw(config), nil)
		if err != nil {
			t.Fatal(err)
		}
		if diff != nil && diff.RequiresNew() {
			var replaced []string
			for k, attr := range diff.Attributes {
				if attr.RequiresNew {
					replaced = append(replaced, k)
				}
			}
			t.Errorf("%s: an imported service must not be replaced, forced by %v", name, replaced)
		}
	}
}

func TestResourcePostgreSQLServiceReadRemovesMissingService(t *testing.T) {
	t.Parallel()

	conn := testPaaSClientWithDebugBodyLogging(&recordingSDKLogger{}, roundTripFunc(func(request *http.Request) (*http.Response, error) {
		return testPaaSErrorResponse(request, http.StatusNotFound, ServiceNotFoundCode, "service not found"), nil
	}))

	data := schema.TestResourceDataRaw(t, ResourcePostgreSQLService().Schema, map[string]interface{}{})
	data.SetId(testPostgreSQLServiceID)

	if diags := resourcePostgreSQLServiceRead(context.Background(), data, testPaaSClientMeta(conn)); diags.HasError() {
		t.Fatalf("reading a missing service: %#v", diags)
	}
	if data.Id() != "" {
		t.Fatalf("a missing service must be removed from state, ID = %q", data.Id())
	}
}
