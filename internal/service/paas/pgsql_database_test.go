package paas

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testPgsqlChildDatabaseID = "paas-db-12345678"

func TestPostgreSQLDatabaseParseResourceID(t *testing.T) {
	t.Parallel()

	for _, testCase := range []struct {
		id, serviceID, databaseID string
		expectErr                 bool
	}{
		{id: "fm-cluster-12345678/paas-db-12345678", serviceID: "fm-cluster-12345678", databaseID: "paas-db-12345678"},
		{id: "paas-db-12345678", expectErr: true},
		{id: "/paas-db-12345678", expectErr: true},
		{id: "fm-cluster-12345678/", expectErr: true},
		{id: "fm-cluster-12345678/paas-db-12345678/extra", expectErr: true},
	} {
		t.Run(testCase.id, func(t *testing.T) {
			serviceID, databaseID, err := PostgreSQLDatabaseParseResourceID(testCase.id)
			if testCase.expectErr {
				if err == nil {
					t.Fatalf("expected error for ID %q", testCase.id)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for ID %q: %s", testCase.id, err)
			}
			if serviceID != testCase.serviceID || databaseID != testCase.databaseID {
				t.Fatalf("got %q, %q, want %q, %q", serviceID, databaseID, testCase.serviceID, testCase.databaseID)
			}
			if id := PostgreSQLDatabaseCreateResourceID(serviceID, databaseID); id != testCase.id {
				t.Fatalf("CreateResourceID = %q, want %q", id, testCase.id)
			}
		})
	}
}

func TestExpandPostgreSQLDatabaseModifyParameters(t *testing.T) {
	t.Parallel()

	set := func(values ...interface{}) *schema.Set { return schema.NewSet(schema.HashString, values) }

	for _, testCase := range []struct {
		name                     string
		encoding, locale         string
		extensions, oldExtension *schema.Set
		expectExtensions         []interface{}
		expectNoExtensions       bool
	}{
		{name: "all parameters", encoding: "UTF8", locale: "ru_RU.UTF-8", extensions: set("pgcrypto"), oldExtension: set(), expectExtensions: []interface{}{"pgcrypto"}},
		{name: "removed extensions are sent as an empty array", extensions: set(), oldExtension: set("pgcrypto"), expectExtensions: []interface{}{}},
		{name: "extensions are omitted when never installed", extensions: set(), oldExtension: set(), expectNoExtensions: true},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			parameters := expandPostgreSQLDatabaseModifyParameters("app_owner", testCase.encoding, testCase.locale, testCase.extensions, testCase.oldExtension)

			if parameters["owner"] != "app_owner" {
				t.Fatalf("ModifyDatabase must always send the owner: %#v", parameters)
			}
			for key, value := range map[string]string{"encoding": testCase.encoding, "locale": testCase.locale} {
				got, ok := parameters[key]
				if value == "" && ok {
					t.Fatalf("empty %s must be omitted: %#v", key, parameters)
				}
				if value != "" && got != value {
					t.Fatalf("%s = %#v, want %q", key, got, value)
				}
			}

			extensions, ok := parameters["extensions"]
			if testCase.expectNoExtensions {
				if ok {
					t.Fatalf("extensions must be omitted: %#v", parameters)
				}
				return
			}
			if list, isList := extensions.([]interface{}); !ok || !isList || len(list) != len(testCase.expectExtensions) {
				t.Fatalf("extensions = %#v, want %#v", extensions, testCase.expectExtensions)
			}
		})
	}
}

func TestExpandPostgreSQLDatabaseModifyUsers(t *testing.T) {
	t.Parallel()

	set := func(values ...interface{}) *schema.Set { return schema.NewSet(schema.HashString, values) }

	if users := expandPostgreSQLDatabaseModifyUsers(set("app_reader"), set()); len(users) != 1 || aws.StringValue(users[0].Name) != "app_reader" {
		t.Fatalf("configured users must be sent: %#v", users)
	}
	// A non-nil empty list is serialized as "users":[] and revokes access of all users.
	if users := expandPostgreSQLDatabaseModifyUsers(set(), set("app_reader")); users == nil || len(users) != 0 {
		t.Fatalf("revoking all users must send an empty list, got %#v", users)
	}
	if users := expandPostgreSQLDatabaseModifyUsers(set(), set()); users != nil {
		t.Fatalf("users must be omitted when none were configured, got %#v", users)
	}
}

func TestExpandPostgreSQLDatabaseUserNames(t *testing.T) {
	t.Parallel()

	if users := expandPostgreSQLDatabaseUserNames(schema.NewSet(schema.HashString, nil)); users != nil {
		t.Fatalf("an empty set must produce no database users, got %#v", users)
	}

	users := expandPostgreSQLDatabaseUserNames(schema.NewSet(schema.HashString, []interface{}{"app_reader"}))
	if len(users) != 1 || users[0].(map[string]interface{})["name"] != "app_reader" {
		t.Fatalf("users = %#v", users)
	}
}

func TestPostgreSQLDatabaseOwnerNotInUsersValidation(t *testing.T) {
	t.Parallel()

	resource := ResourcePostgreSQLDatabase()

	for _, testCase := range []struct {
		name      string
		users     cty.Value
		errorText string
	}{
		{name: "owner in users", users: cty.SetVal([]cty.Value{cty.StringVal("app_owner")}), errorText: `users must not contain the database owner "app_owner"`},
		{name: "other users", users: cty.SetVal([]cty.Value{cty.StringVal("app_reader")})},
		{name: "unknown users", users: cty.UnknownVal(cty.Set(cty.String))},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			config := terraform.NewResourceConfigShimmed(cty.ObjectVal(map[string]cty.Value{
				"service_id": cty.StringVal(testPgsqlChildServiceID),
				"name":       cty.StringVal("app"),
				"owner":      cty.StringVal("app_owner"),
				"users":      testCase.users,
			}), resource.CoreConfigSchema())

			_, err := resource.Diff(context.Background(), nil, config, nil)
			if testCase.errorText == "" {
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), testCase.errorText) {
				t.Fatalf("expected %q, got %v", testCase.errorText, err)
			}
		})
	}
}

func TestPostgreSQLDatabaseDiff(t *testing.T) {
	t.Parallel()

	resource := ResourcePostgreSQLDatabase()
	initial := schema.TestResourceDataRaw(t, resource.Schema, testPgsqlChildDatabaseConfig(nil))
	initial.SetId(testPgsqlChildServiceID + "/" + testPgsqlChildDatabaseID)

	for _, testCase := range []struct {
		name          string
		change        map[string]interface{}
		expectNoDiff  bool
		expectReplace bool
	}{
		{name: "encoding differs only in case", change: map[string]interface{}{"encoding": "utf8"}, expectNoDiff: true},
		{name: "locale forces replacement", change: map[string]interface{}{"locale": "en_US.UTF-8"}, expectReplace: true},
		{name: "encoding forces replacement", change: map[string]interface{}{"encoding": "WIN1251"}, expectReplace: true},
		{name: "backup_id forces replacement", change: map[string]interface{}{"backup_id": "paas-backup-12345678"}, expectReplace: true},
		{name: "owner is updated in place", change: map[string]interface{}{"owner": "app_reader", "users": []interface{}{}}},
		{name: "extensions are updated in place", change: map[string]interface{}{"extensions": []interface{}{"pgcrypto", "pg_trgm"}}},
		{name: "backup_enabled is updated in place", change: map[string]interface{}{"backup_enabled": true}},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			diff, err := resource.Diff(context.Background(), initial.State(), terraform.NewResourceConfigRaw(testPgsqlChildDatabaseConfig(testCase.change)), nil)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if testCase.expectNoDiff {
				if diff != nil && len(diff.Attributes) > 0 {
					t.Fatalf("expected no diff, got %#v", diff.Attributes)
				}
				return
			}
			if diff == nil || len(diff.Attributes) == 0 {
				t.Fatal("expected a diff")
			}
			if diff.RequiresNew() != testCase.expectReplace {
				t.Fatalf("RequiresNew = %t, want %t", diff.RequiresNew(), testCase.expectReplace)
			}
		})
	}
}

func TestResourcePostgreSQLDatabaseCreate(t *testing.T) {
	t.Parallel()

	created := false
	conn := testPaaSClientWithDebugBodyLogging(&recordingSDKLogger{}, roundTripFunc(func(request *http.Request) (*http.Response, error) {
		databasesPath := "/services/" + testPgsqlChildServiceID + "/databases"
		switch {
		case request.Method == http.MethodGet && request.URL.Path == "/services/"+testPgsqlChildServiceID:
			return jsonResponse(request, testPgsqlChildServiceResponse("pgsql")), nil
		case request.Method == http.MethodPost && request.URL.Path == databasesPath:
			var body struct {
				Database struct {
					Name          string                 `json:"name"`
					BackupEnabled *bool                  `json:"backupEnabled"`
					Parameters    map[string]interface{} `json:"parameters"`
					Users         []struct {
						Name string `json:"name"`
					} `json:"users"`
				} `json:"database"`
			}
			if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			database := body.Database
			if database.Name != "app" || database.BackupEnabled == nil || *database.BackupEnabled {
				t.Fatalf("unexpected CreateDatabase request: %#v", database)
			}
			if database.Parameters["owner"] != "app_owner" {
				t.Fatalf("CreateDatabase must send the owner: %#v", database.Parameters)
			}
			if extensions, ok := database.Parameters["extensions"].([]interface{}); !ok || len(extensions) != 1 || extensions[0] != "pgcrypto" {
				t.Fatalf("CreateDatabase must send extensions: %#v", database.Parameters)
			}
			for _, key := range []string{"encoding", "locale", "backup_id", "backup_db_name"} {
				if _, ok := database.Parameters[key]; ok {
					t.Fatalf("unset %s must not be sent: %#v", key, database.Parameters)
				}
			}
			if len(database.Users) != 1 || database.Users[0].Name != "app_reader" {
				t.Fatalf("CreateDatabase must send database users: %#v", database.Users)
			}
			created = true
			return jsonResponse(request, `{"database":{"id":"`+testPgsqlChildDatabaseID+`","name":"app","status":"available"}}`), nil
		case request.Method == http.MethodGet && request.URL.Path == databasesPath:
			return jsonResponse(request, `{"databases":[{"id":"`+testPgsqlChildDatabaseID+`","name":"app","status":"available"}]}`), nil
		case request.Method == http.MethodGet && request.URL.Path == databasesPath+"/"+testPgsqlChildDatabaseID:
			return jsonResponse(request, testPgsqlChildDatabaseResponse(`{"owner":"app_owner","extensions":["pgcrypto"]}`, `[{"name":"app_reader"}]`)), nil
		default:
			t.Fatalf("unexpected PaaS request: %s %s", request.Method, request.URL.Path)
			return nil, nil
		}
	}))

	data := schema.TestResourceDataRaw(t, ResourcePostgreSQLDatabase().Schema, map[string]interface{}{
		"service_id": testPgsqlChildServiceID,
		"name":       "app",
		"owner":      "app_owner",
		"extensions": []interface{}{"pgcrypto"},
		"users":      []interface{}{"app_reader"},
	})
	data.MarkNewResource()

	if diags := resourcePostgreSQLDatabaseCreate(context.Background(), data, testPaaSClientMeta(conn)); diags.HasError() {
		t.Fatalf("creating database: %#v", diags)
	}
	if !created {
		t.Fatal("CreateDatabase was not called")
	}
	if want := testPgsqlChildServiceID + "/" + testPgsqlChildDatabaseID; data.Id() != want {
		t.Fatalf("ID = %q, want %q", data.Id(), want)
	}
	if data.Get("database_id") != testPgsqlChildDatabaseID || data.Get("owner") != "app_owner" {
		t.Fatalf("database attributes were not read: database_id=%v owner=%v", data.Get("database_id"), data.Get("owner"))
	}
	if users := data.Get("users").(*schema.Set); users.Len() != 1 || !users.Contains("app_reader") {
		t.Fatalf("users = %#v", users.List())
	}
	// Encoding and locale are not returned by the API, so they stay unset instead of being cleared.
	if data.Get("encoding") != "" || data.Get("locale") != "" {
		t.Fatalf("encoding = %q, locale = %q, want empty", data.Get("encoding"), data.Get("locale"))
	}
}

func TestResourcePostgreSQLDatabaseUpdateSendsAllParameters(t *testing.T) {
	t.Parallel()

	resource := ResourcePostgreSQLDatabase()
	initial := schema.TestResourceDataRaw(t, resource.Schema, testPgsqlChildDatabaseConfig(nil))
	initial.SetId(testPgsqlChildServiceID + "/" + testPgsqlChildDatabaseID)

	// The new owner must be removed from users in the same change.
	config := testPgsqlChildDatabaseConfig(map[string]interface{}{"owner": "app_reader", "users": []interface{}{}})
	diff, err := resource.Diff(context.Background(), initial.State(), terraform.NewResourceConfigRaw(config), nil)
	if err != nil || diff == nil || diff.RequiresNew() {
		t.Fatalf("changing the owner must update the database in place: %v", err)
	}

	modified := false
	conn := testPaaSClientWithDebugBodyLogging(&recordingSDKLogger{}, roundTripFunc(func(request *http.Request) (*http.Response, error) {
		databasesPath := "/services/" + testPgsqlChildServiceID + "/databases"
		switch {
		case request.Method == http.MethodPut && request.URL.Path == databasesPath+"/"+testPgsqlChildDatabaseID:
			var body struct {
				Database map[string]interface{} `json:"database"`
			}
			if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			parameters, _ := body.Database["parameters"].(map[string]interface{})
			for key, want := range map[string]string{"owner": "app_reader", "encoding": "UTF8", "locale": "ru_RU.UTF-8"} {
				if parameters[key] != want {
					t.Fatalf("ModifyDatabase must send all current parameters, %s = %#v, want %q", key, parameters[key], want)
				}
			}
			if extensions, ok := parameters["extensions"].([]interface{}); !ok || len(extensions) != 1 {
				t.Fatalf("ModifyDatabase must send current extensions: %#v", parameters)
			}
			if users, ok := body.Database["users"].([]interface{}); !ok || len(users) != 0 {
				t.Fatalf("revoking all users must send users:[]: %#v", body.Database)
			}
			modified = true
			return jsonResponse(request, `{}`), nil
		case request.Method == http.MethodGet && request.URL.Path == "/services/"+testPgsqlChildServiceID:
			return jsonResponse(request, testPgsqlChildServiceResponse("pgsql")), nil
		case request.Method == http.MethodGet && request.URL.Path == databasesPath:
			return jsonResponse(request, `{"databases":[{"id":"`+testPgsqlChildDatabaseID+`","name":"app","status":"available"}]}`), nil
		case request.Method == http.MethodGet && request.URL.Path == databasesPath+"/"+testPgsqlChildDatabaseID:
			return jsonResponse(request, testPgsqlChildDatabaseResponse(
				`{"owner":"app_reader","encoding":"UTF8","locale":"ru_RU.UTF-8","extensions":["pgcrypto"]}`, `[]`)), nil
		default:
			t.Fatalf("unexpected PaaS request: %s %s", request.Method, request.URL.Path)
			return nil, nil
		}
	}))

	state, diags := resource.Apply(context.Background(), initial.State(), diff, testPaaSClientMeta(conn))
	if diags.HasError() || !modified {
		t.Fatalf("updating database failed: %#v", diags)
	}
	if state.Attributes["owner"] != "app_reader" || state.Attributes["users.#"] != "0" {
		t.Fatalf("owner = %q, users.# = %q", state.Attributes["owner"], state.Attributes["users.#"])
	}
}

func TestResourcePostgreSQLDatabaseReadRemovesMissingDatabase(t *testing.T) {
	t.Parallel()

	conn := testPaaSClientWithDebugBodyLogging(&recordingSDKLogger{}, roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.Method == http.MethodGet && request.URL.Path == "/services/"+testPgsqlChildServiceID+"/databases" {
			return jsonResponse(request, `{"databases":[]}`), nil
		}
		t.Fatalf("unexpected PaaS request: %s %s", request.Method, request.URL.Path)
		return nil, nil
	}))

	data := schema.TestResourceDataRaw(t, ResourcePostgreSQLDatabase().Schema, testPgsqlChildDatabaseConfig(nil))
	data.SetId(testPgsqlChildServiceID + "/" + testPgsqlChildDatabaseID)

	if diags := resourcePostgreSQLDatabaseRead(context.Background(), data, testPaaSClientMeta(conn)); diags.HasError() {
		t.Fatalf("reading database: %#v", diags)
	}
	if data.Id() != "" {
		t.Fatalf("a database missing in the service must be removed from state, ID = %q", data.Id())
	}
}

func TestResourcePostgreSQLDatabaseReadRemovesDatabaseOfDeletedService(t *testing.T) {
	t.Parallel()

	conn := testPaaSClientWithDebugBodyLogging(&recordingSDKLogger{}, roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.Method == http.MethodGet && request.URL.Path == "/services/"+testPgsqlChildServiceID+"/databases" {
			return testPaaSErrorResponse(request, http.StatusNotFound, ServiceNotFoundCode, "service not found"), nil
		}
		t.Fatalf("unexpected PaaS request: %s %s", request.Method, request.URL.Path)
		return nil, nil
	}))

	data := schema.TestResourceDataRaw(t, ResourcePostgreSQLDatabase().Schema, testPgsqlChildDatabaseConfig(nil))
	data.SetId(testPgsqlChildServiceID + "/" + testPgsqlChildDatabaseID)

	if diags := resourcePostgreSQLDatabaseRead(context.Background(), data, testPaaSClientMeta(conn)); diags.HasError() {
		t.Fatalf("reading database: %#v", diags)
	}
	if data.Id() != "" {
		t.Fatalf("a database of a deleted service must be removed from state, ID = %q", data.Id())
	}
}

// testPgsqlChildDatabaseConfig returns a database configuration with the given changes applied.
func testPgsqlChildDatabaseConfig(changes map[string]interface{}) map[string]interface{} {
	config := map[string]interface{}{
		"service_id": testPgsqlChildServiceID,
		"name":       "app",
		"owner":      "app_owner",
		"encoding":   "UTF8",
		"locale":     "ru_RU.UTF-8",
		"extensions": []interface{}{"pgcrypto"},
		"users":      []interface{}{"app_reader"},
	}
	for k, v := range changes {
		config[k] = v
	}
	return config
}

func testPgsqlChildDatabaseResponse(parameters, users string) string {
	return `{"database":{"id":"` + testPgsqlChildDatabaseID + `","name":"app","status":"available","backupEnabled":false,` +
		`"parameters":` + parameters + `,"users":` + users + `}}`
}
