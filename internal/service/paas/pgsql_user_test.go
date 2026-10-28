package paas

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	testPgsqlChildServiceID = "fm-cluster-12345678"
	testPgsqlChildUserID    = "paas-user-12345678"
)

func TestPostgreSQLUserParseResourceID(t *testing.T) {
	t.Parallel()

	for _, testCase := range []struct {
		id, serviceID, userID string
		expectErr             bool
	}{
		{id: "fm-cluster-12345678/paas-user-12345678", serviceID: "fm-cluster-12345678", userID: "paas-user-12345678"},
		{id: "paas-user-12345678", expectErr: true},
		{id: "/paas-user-12345678", expectErr: true},
		{id: "fm-cluster-12345678/", expectErr: true},
		{id: "fm-cluster-12345678/paas-user-12345678/extra", expectErr: true},
	} {
		t.Run(testCase.id, func(t *testing.T) {
			serviceID, userID, err := PostgreSQLUserParseResourceID(testCase.id)
			if testCase.expectErr {
				if err == nil {
					t.Fatalf("expected error for ID %q", testCase.id)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for ID %q: %s", testCase.id, err)
			}
			if serviceID != testCase.serviceID || userID != testCase.userID {
				t.Fatalf("got %q, %q, want %q, %q", serviceID, userID, testCase.serviceID, testCase.userID)
			}
			if id := PostgreSQLUserCreateResourceID(serviceID, userID); id != testCase.id {
				t.Fatalf("CreateResourceID = %q, want %q", id, testCase.id)
			}
		})
	}
}

func TestExpandPostgreSQLUserModifyParameters(t *testing.T) {
	t.Parallel()

	set := func(values ...interface{}) *schema.Set { return schema.NewSet(schema.HashString, values) }

	for _, testCase := range []struct {
		name            string
		roles, oldRoles *schema.Set
		expectRoles     []interface{}
		expectNoRoles   bool
	}{
		{name: "roles are sent", roles: set("pg_monitor"), oldRoles: set(), expectRoles: []interface{}{"pg_monitor"}},
		{name: "removed roles are sent as an empty array", roles: set(), oldRoles: set("pg_monitor"), expectRoles: []interface{}{}},
		// Older PaaS environments reject roles, so they must be omitted when never assigned.
		{name: "roles are omitted when never assigned", roles: set(), oldRoles: set(), expectNoRoles: true},
		{name: "nil sets omit roles", expectNoRoles: true},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			parameters := expandPostgreSQLUserModifyParameters("pgsql-test-secret", testCase.roles, testCase.oldRoles)

			if parameters["password"] != "pgsql-test-secret" {
				t.Fatalf("ModifyUser must always send the password: %#v", parameters)
			}

			roles, ok := parameters["roles"]
			if testCase.expectNoRoles {
				if ok {
					t.Fatalf("roles must be omitted: %#v", parameters)
				}
				return
			}
			list, isList := roles.([]interface{})
			if !ok || !isList || len(list) != len(testCase.expectRoles) {
				t.Fatalf("roles = %#v, want %#v", roles, testCase.expectRoles)
			}
			for i := range list {
				if list[i] != testCase.expectRoles[i] {
					t.Fatalf("roles = %#v, want %#v", list, testCase.expectRoles)
				}
			}
		})
	}
}

// Capturing the process-wide provider logger requires a serial test.
func TestResourcePostgreSQLUserCreate(t *testing.T) {
	var providerLogs bytes.Buffer
	originalLogOutput := log.Writer()
	log.SetOutput(&providerLogs)
	defer log.SetOutput(originalLogOutput)

	const password = "pgsql-user-create-secret"
	sdkLogger := &recordingSDKLogger{}
	created := false

	conn := testPaaSClientWithDebugBodyLogging(sdkLogger, roundTripFunc(func(request *http.Request) (*http.Response, error) {
		switch {
		case request.Method == http.MethodGet && request.URL.Path == "/services/"+testPgsqlChildServiceID:
			return jsonResponse(request, testPgsqlChildServiceResponse("pgsql")), nil
		case request.Method == http.MethodPost && request.URL.Path == "/services/"+testPgsqlChildServiceID+"/users":
			var body struct {
				User struct {
					Name       string                 `json:"name"`
					Parameters map[string]interface{} `json:"parameters"`
				} `json:"user"`
			}
			if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if body.User.Name != "app_reader" || body.User.Parameters["password"] != password {
				t.Fatalf("unexpected CreateUser request: %#v", body)
			}
			if roles, ok := body.User.Parameters["roles"].([]interface{}); !ok || len(roles) != 1 || roles[0] != "pg_monitor" {
				t.Fatalf("CreateUser must send roles: %#v", body.User.Parameters)
			}
			created = true
			return jsonResponse(request, `{"user":{"id":"`+testPgsqlChildUserID+`","name":"app_reader"}}`), nil
		case request.Method == http.MethodGet && request.URL.Path == "/services/"+testPgsqlChildServiceID+"/users":
			return jsonResponse(request, `{"users":[{"id":"`+testPgsqlChildUserID+`","name":"app_reader"}]}`), nil
		case request.Method == http.MethodGet && request.URL.Path == "/services/"+testPgsqlChildServiceID+"/users/"+testPgsqlChildUserID:
			return jsonResponse(request, testPgsqlChildUserResponse(password, `["pg_monitor"]`)), nil
		default:
			t.Fatalf("unexpected PaaS request: %s %s", request.Method, request.URL.Path)
			return nil, nil
		}
	}))

	data := schema.TestResourceDataRaw(t, ResourcePostgreSQLUser().Schema, testPgsqlChildUserConfig(password, []interface{}{"pg_monitor"}))
	data.MarkNewResource()

	if diags := resourcePostgreSQLUserCreate(context.Background(), data, testPaaSClientMeta(conn)); diags.HasError() {
		t.Fatalf("creating user: %#v", diags)
	}
	if !created {
		t.Fatal("CreateUser was not called")
	}
	if want := testPgsqlChildServiceID + "/" + testPgsqlChildUserID; data.Id() != want {
		t.Fatalf("ID = %q, want %q", data.Id(), want)
	}
	if data.Get("user_id") != testPgsqlChildUserID || data.Get("password") != password {
		t.Fatalf("user attributes were not read: user_id=%v", data.Get("user_id"))
	}
	if roles := data.Get("roles").(*schema.Set); roles.Len() != 1 || !roles.Contains("pg_monitor") {
		t.Fatalf("roles = %#v", roles.List())
	}
	if strings.Contains(providerLogs.String(), password) || strings.Contains(sdkLogger.String(), password) {
		t.Fatal("user password was written to provider or SDK HTTP debug logs")
	}
}

func TestResourcePostgreSQLUserCreateRejectsOtherServiceTypes(t *testing.T) {
	t.Parallel()

	conn := testPaaSClientWithDebugBodyLogging(&recordingSDKLogger{}, roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.Method == http.MethodGet && request.URL.Path == "/services/"+testPgsqlChildServiceID {
			return jsonResponse(request, testPgsqlChildServiceResponse("redis")), nil
		}
		t.Fatalf("no user requests must be sent for a non-PostgreSQL service: %s %s", request.Method, request.URL.Path)
		return nil, nil
	}))

	data := schema.TestResourceDataRaw(t, ResourcePostgreSQLUser().Schema, testPgsqlChildUserConfig("pgsql-test-secret", nil))
	data.MarkNewResource()

	diags := resourcePostgreSQLUserCreate(context.Background(), data, testPaaSClientMeta(conn))
	if !diags.HasError() || !strings.Contains(diags[0].Summary, `supports only "pgsql" services`) {
		t.Fatalf("expected service type error, got %#v", diags)
	}
}

func TestResourcePostgreSQLUserCreateRetriesWhenServiceIsBusy(t *testing.T) {
	t.Parallel()

	attempts := 0
	conn := testPaaSClientWithDebugBodyLogging(&recordingSDKLogger{}, roundTripFunc(func(request *http.Request) (*http.Response, error) {
		switch {
		case request.Method == http.MethodGet && request.URL.Path == "/services/"+testPgsqlChildServiceID:
			return jsonResponse(request, testPgsqlChildServiceResponse("pgsql")), nil
		case request.Method == http.MethodPost && request.URL.Path == "/services/"+testPgsqlChildServiceID+"/users":
			attempts++
			if attempts == 1 {
				return testPaaSErrorResponse(request, http.StatusBadRequest, OperationNotAllowedCode, "service is updating"), nil
			}
			return jsonResponse(request, `{"user":{"id":"`+testPgsqlChildUserID+`","name":"app_reader"}}`), nil
		case request.Method == http.MethodGet && request.URL.Path == "/services/"+testPgsqlChildServiceID+"/users":
			return jsonResponse(request, `{"users":[{"id":"`+testPgsqlChildUserID+`","name":"app_reader"}]}`), nil
		case request.Method == http.MethodGet && request.URL.Path == "/services/"+testPgsqlChildServiceID+"/users/"+testPgsqlChildUserID:
			return jsonResponse(request, testPgsqlChildUserResponse("pgsql-test-secret", `[]`)), nil
		default:
			t.Fatalf("unexpected PaaS request: %s %s", request.Method, request.URL.Path)
			return nil, nil
		}
	}))

	data := schema.TestResourceDataRaw(t, ResourcePostgreSQLUser().Schema, testPgsqlChildUserConfig("pgsql-test-secret", nil))
	data.MarkNewResource()

	if diags := resourcePostgreSQLUserCreate(context.Background(), data, testPaaSClientMeta(conn)); diags.HasError() {
		t.Fatalf("creating user: %#v", diags)
	}
	if attempts != 2 {
		t.Fatalf("CreateUser attempts = %d, want 2", attempts)
	}
}

func TestResourcePostgreSQLUserUpdateSendsAllParameters(t *testing.T) {
	t.Parallel()

	resource := ResourcePostgreSQLUser()
	initial := schema.TestResourceDataRaw(t, resource.Schema, testPgsqlChildUserConfig("pgsql-test-secret", []interface{}{"pg_monitor"}))
	initial.SetId(testPgsqlChildServiceID + "/" + testPgsqlChildUserID)

	diff, err := resource.Diff(context.Background(), initial.State(), terraform.NewResourceConfigRaw(testPgsqlChildUserConfig("pgsql-test-secret", nil)), nil)
	if err != nil || diff == nil || diff.RequiresNew() {
		t.Fatalf("removing roles must update the user in place: %v", err)
	}

	modified := false
	conn := testPaaSClientWithDebugBodyLogging(&recordingSDKLogger{}, roundTripFunc(func(request *http.Request) (*http.Response, error) {
		userPath := "/services/" + testPgsqlChildServiceID + "/users/" + testPgsqlChildUserID
		switch {
		case request.Method == http.MethodPut && request.URL.Path == userPath:
			var body struct {
				User struct {
					Parameters map[string]interface{} `json:"parameters"`
				} `json:"user"`
			}
			if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if body.User.Parameters["password"] != "pgsql-test-secret" {
				t.Fatalf("ModifyUser must send the current password: %#v", body.User.Parameters)
			}
			if roles, ok := body.User.Parameters["roles"].([]interface{}); !ok || len(roles) != 0 {
				t.Fatalf("revocation must send roles:[]: %#v", body.User.Parameters)
			}
			modified = true
			return jsonResponse(request, `{}`), nil
		case request.Method == http.MethodGet && request.URL.Path == "/services/"+testPgsqlChildServiceID:
			return jsonResponse(request, testPgsqlChildServiceResponse("pgsql")), nil
		case request.Method == http.MethodGet && request.URL.Path == "/services/"+testPgsqlChildServiceID+"/users":
			return jsonResponse(request, `{"users":[{"id":"`+testPgsqlChildUserID+`","name":"app_reader"}]}`), nil
		case request.Method == http.MethodGet && request.URL.Path == userPath:
			return jsonResponse(request, testPgsqlChildUserResponse("pgsql-test-secret", `[]`)), nil
		default:
			t.Fatalf("unexpected PaaS request: %s %s", request.Method, request.URL.Path)
			return nil, nil
		}
	}))

	state, diags := resource.Apply(context.Background(), initial.State(), diff, testPaaSClientMeta(conn))
	if diags.HasError() || !modified {
		t.Fatalf("updating user failed: %#v", diags)
	}
	if state.Attributes["roles.#"] != "0" {
		t.Fatalf("roles.# = %q, want 0", state.Attributes["roles.#"])
	}
}

func TestResourcePostgreSQLUserReadRemovesMissingUser(t *testing.T) {
	t.Parallel()

	conn := testPaaSClientWithDebugBodyLogging(&recordingSDKLogger{}, roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.Method == http.MethodGet && request.URL.Path == "/services/"+testPgsqlChildServiceID+"/users" {
			return jsonResponse(request, `{"users":[{"id":"paas-user-other","name":"other"}]}`), nil
		}
		t.Fatalf("unexpected PaaS request: %s %s", request.Method, request.URL.Path)
		return nil, nil
	}))

	data := schema.TestResourceDataRaw(t, ResourcePostgreSQLUser().Schema, testPgsqlChildUserConfig("pgsql-test-secret", nil))
	data.SetId(testPgsqlChildServiceID + "/" + testPgsqlChildUserID)

	if diags := resourcePostgreSQLUserRead(context.Background(), data, testPaaSClientMeta(conn)); diags.HasError() {
		t.Fatalf("reading user: %#v", diags)
	}
	if data.Id() != "" {
		t.Fatalf("a user missing in the service must be removed from state, ID = %q", data.Id())
	}
}

func TestResourcePostgreSQLUserDeleteIgnoresMissingUser(t *testing.T) {
	t.Parallel()

	conn := testPaaSClientWithDebugBodyLogging(&recordingSDKLogger{}, roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.Method == http.MethodDelete && request.URL.Path == "/services/"+testPgsqlChildServiceID+"/users/"+testPgsqlChildUserID {
			return testPaaSErrorResponse(request, http.StatusNotFound, ResourceDoesNotExistCode, "user not found"), nil
		}
		t.Fatalf("unexpected PaaS request: %s %s", request.Method, request.URL.Path)
		return nil, nil
	}))

	data := schema.TestResourceDataRaw(t, ResourcePostgreSQLUser().Schema, testPgsqlChildUserConfig("pgsql-test-secret", nil))
	data.SetId(testPgsqlChildServiceID + "/" + testPgsqlChildUserID)

	if diags := resourcePostgreSQLUserDelete(context.Background(), data, testPaaSClientMeta(conn)); diags.HasError() {
		t.Fatalf("deleting a missing user must succeed: %#v", diags)
	}
}

func TestResourcePostgreSQLUserDeleteReturnsDependencyError(t *testing.T) {
	t.Parallel()

	conn := testPaaSClientWithDebugBodyLogging(&recordingSDKLogger{}, roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.Method == http.MethodDelete && request.URL.Path == "/services/"+testPgsqlChildServiceID+"/users/"+testPgsqlChildUserID {
			return testPaaSErrorResponse(request, http.StatusBadRequest, "CannotDeleteUser",
				"The user app_reader cannot be deleted. First, remove it from the databases: app"), nil
		}
		t.Fatalf("unexpected PaaS request: %s %s", request.Method, request.URL.Path)
		return nil, nil
	}))

	data := schema.TestResourceDataRaw(t, ResourcePostgreSQLUser().Schema, testPgsqlChildUserConfig("pgsql-test-secret", nil))
	data.SetId(testPgsqlChildServiceID + "/" + testPgsqlChildUserID)

	diags := resourcePostgreSQLUserDelete(context.Background(), data, testPaaSClientMeta(conn))
	if !diags.HasError() || !strings.Contains(diags[0].Summary, "CannotDeleteUser") {
		t.Fatalf("expected CannotDeleteUser error, got %#v", diags)
	}
}

func testPgsqlChildUserConfig(password string, roles []interface{}) map[string]interface{} {
	config := map[string]interface{}{
		"service_id": testPgsqlChildServiceID,
		"name":       "app_reader",
		"password":   password,
	}
	if roles != nil {
		config["roles"] = roles
	}
	return config
}

func testPgsqlChildServiceResponse(serviceType string) string {
	return `{"service":{"id":"` + testPgsqlChildServiceID + `","serviceType":"` + serviceType + `","status":"READY"}}`
}

func testPgsqlChildUserResponse(password, roles string) string {
	return `{"user":{"id":"` + testPgsqlChildUserID + `","name":"app_reader","parameters":{"password":"` + password + `","roles":` + roles + `}}}`
}

func testPaaSErrorResponse(request *http.Request, statusCode int, code, message string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body:    io.NopCloser(strings.NewReader(fmt.Sprintf(`{"code":%q,"message":%q}`, code, message))),
		Request: request,
	}
}
