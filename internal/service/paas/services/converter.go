package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/paas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/experimental/nullable"
	"log"
	"strconv"
)

func (s service) ExpandServiceParameters(tfMap map[string]interface{}) ServiceParameters {
	if tfMap == nil {
		return nil
	}

	serviceParameters := ServiceParameters{}

	if v, ok := tfMap["logging"].([]interface{}); ok && len(v) > 0 {
		loggingMap := v[0].(map[string]interface{})
		serviceParameters["logging"] = true

		if v, ok := loggingMap["log_to"].(string); ok && v != "" {
			serviceParameters["log_to"] = v
		}

		if v, ok := loggingMap["logging_tags"].(*schema.Set); ok && v.Len() > 0 {
			serviceParameters["logging_tags"] = v.List()
		}

		delete(tfMap, "logging")
	} else {
		serviceParameters["logging"] = false
	}

	if v, ok := tfMap["monitoring"].([]interface{}); ok && len(v) > 0 {
		monitoringMap := v[0].(map[string]interface{})
		serviceParameters["monitoring"] = true

		if v, ok := monitoringMap["monitor_by"].(string); ok && v != "" {
			serviceParameters["monitor_by"] = v
		}

		if v, ok := monitoringMap["monitoring_labels"].(map[string]interface{}); ok && len(v) > 0 {
			serviceParameters["monitoring_labels"] = v
		}

		delete(tfMap, "monitoring")
	} else {
		serviceParameters["monitoring"] = false
	}

	for k, v := range s.toInterface().expandServiceParameters(tfMap) {
		serviceParameters[k] = v
	}

	return serviceParameters
}

// ExpandUsers converts terraform representation of list of users to api representation.
func (s service) ExpandUsers(tfList []interface{}, forDatabase bool) []*paas.UserCreateRequest {
	if len(tfList) == 0 {
		return nil
	}

	var users []*paas.UserCreateRequest

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})
		if !ok {
			continue
		}

		user := s.toInterface().ExpandUser(tfMap, forDatabase)
		if user == nil {
			continue
		}

		users = append(users, user)
	}

	return users
}

// ExpandUser converts terraform representation of service user to api representation.
// If forDatabase is true, user is considered a database user.
func (s service) ExpandUser(tfMap map[string]interface{}, forDatabase bool) *paas.UserCreateRequest {
	if tfMap == nil {
		return nil
	}

	user := &paas.UserCreateRequest{}

	if v, ok := tfMap["name"].(string); ok && v != "" {
		user.Name = aws.String(v)
		delete(tfMap, "name")
	}

	if forDatabase {
		user.Parameters = s.toInterface().expandDatabaseUserParameters(tfMap)
	} else {
		user.Parameters = s.toInterface().expandUserParameters(tfMap)
	}

	return user
}

// ExpandDatabases converts terraform representation of list of databases to api representation.
func (s service) ExpandDatabases(tfList []interface{}) []*paas.DatabaseCreateRequest {
	if len(tfList) == 0 {
		return nil
	}

	var databases []*paas.DatabaseCreateRequest

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})
		if !ok {
			continue
		}

		database := s.toInterface().ExpandDatabase(tfMap)
		if database == nil {
			continue
		}

		databases = append(databases, database)
	}

	return databases
}

// ExpandDatabase converts terraform representation of database to api representation.
func (s service) ExpandDatabase(tfParameters map[string]interface{}) *paas.DatabaseCreateRequest {
	if tfParameters == nil {
		return nil
	}

	database := &paas.DatabaseCreateRequest{}

	if v, ok := tfParameters["backup_enabled"].(bool); ok && s.allowBackup {
		database.BackupEnabled = aws.Bool(v)
		delete(tfParameters, "backup_enabled")
	}

	if v, ok := tfParameters["name"].(string); ok && v != "" {
		database.Name = aws.String(v)
		delete(tfParameters, "name")
	}

	if v, ok := tfParameters["user"].([]interface{}); ok && len(v) > 0 {
		database.Users = s.toInterface().ExpandUsers(v, true)
		delete(tfParameters, "user")
	}

	database.Parameters = s.toInterface().expandDatabaseParameters(tfParameters)

	return database
}

// expandServiceParameters converts terraform representation of service-specific parameters
// to api representation.
//
// If PaaS service has specific parameters, it should override this method.
func (s service) expandServiceParameters(_ map[string]interface{}) ServiceParameters {
	return nil
}

// expandUserParameters converts terraform representation of service-specific user parameters
// to api representation.
//
// If PaaS service has specific user parameters, it should override this method.
func (s service) expandUserParameters(_ map[string]interface{}) UserParameters {
	return nil
}

// expandDatabaseParameters converts terraform representation of service-specific database parameters
// to api representation.
//
// If PaaS service has specific database parameters, it should override this method.
func (s service) expandDatabaseParameters(_ map[string]interface{}) DatabaseParameters {
	return nil
}

// expandDatabaseUserParameters converts terraform representation of service-specific database user parameters
// to api representation.
//
// If PaaS service has specific database user parameters, it should override this method.
func (s service) expandDatabaseUserParameters(_ map[string]interface{}) DatabaseUserParameters {
	return nil
}

// todo: doc
func (s service) expandIntByteParameter(tfByteParameter interface{}) map[string]interface{} {
	if vList, okList := tfByteParameter.([]interface{}); okList && len(vList) > 0 {
		var vMap = vList[0].(map[string]interface{})

		// todo: check if not ok
		if v, _, _ := nullable.Int(vMap["value"].(string)).Value(); v != 0 {
			return map[string]interface{}{
				"value":     v,
				"dimension": vMap["dimension"].(string),
			}
		}
	}

	return nil
}

// todo: doc
func (s service) expandFloatByteParameter(tfByteParameter interface{}) map[string]interface{} {
	if vList, okList := tfByteParameter.([]interface{}); okList && len(vList) > 0 {
		var vMap = vList[0].(map[string]interface{})

		if v, ok := vMap["value"].(float64); ok && v != 0.0 {
			return map[string]interface{}{
				"value":     v,
				"dimension": vMap["dimension"].(string),
			}
		}
	}

	return nil
}

// FlattenServiceParametersUsersDatabases converts all blocks of service-specific parameters
// from api to terraform representation.
//
// It's Expand analogue is represented by three separate methods:
// ExpandServiceParameters, ExpandUsers and ExpandDatabases,
// because these blocks are separated in api representation.
func (s service) FlattenServiceParametersUsersDatabases(
	serviceParameters ServiceParameters,
	users []*paas.UserResponse,
	databases []*paas.DatabaseResponse,
) map[string]interface{} {
	tfMap := map[string]interface{}{}

	for k, v := range s.toInterface().flattenServiceParameters(serviceParameters) {
		tfMap[k] = v
	}

	if s.usersEnabled {
		tfMap["user"] = s.toInterface().FlattenUsers(users, false)
	}

	if s.databasesEnabled {
		tfMap["database"] = s.toInterface().FlattenDatabases(databases)
	}

	if s.loggingEnabled {
		if v, ok := serviceParameters["logging"].(bool); ok && v {
			loggingMap := map[string]interface{}{}

			if v, ok := serviceParameters["logTo"].(string); ok {
				loggingMap["log_to"] = v
			}

			if v, ok := serviceParameters["loggingTags"].([]interface{}); ok {
				loggingMap["logging_tags"] = v
			}

			tfMap["logging"] = []map[string]interface{}{loggingMap}
		}
	}

	if s.monitoringEnabled {
		if v, ok := serviceParameters["monitoring"].(bool); ok && v {
			monitoringMap := map[string]interface{}{}

			if v, ok := serviceParameters["monitorBy"].(string); ok {
				monitoringMap["monitor_by"] = v
			}

			if v, ok := serviceParameters["monitoringLabels"].(map[string]interface{}); ok {
				monitoringMap["monitoring_labels"] = v
			}

			tfMap["monitoring"] = []map[string]interface{}{monitoringMap}
		}
	}

	return tfMap
}

// FlattenUsers converts api representation of list of users to terraform representation.
func (s service) FlattenUsers(users []*paas.UserResponse, forDatabase bool) []interface{} {
	if len(users) == 0 {
		return nil
	}

	var tfList []interface{}

	for _, user := range users {
		if user == nil {
			continue
		}

		tfList = append(tfList, s.toInterface().FlattenUser(user, forDatabase))
	}

	return tfList
}

// FlattenUser converts api representation of service user to terraform representation.
// If forDatabase is true, user is considered a database user.
func (s service) FlattenUser(user *paas.UserResponse, forDatabase bool) map[string]interface{} {
	if user == nil {
		return map[string]interface{}{}
	}

	tfMap := map[string]interface{}{}

	if v := user.Id; v != nil {
		tfMap["id"] = v
	}

	if v := user.Name; v != nil {
		tfMap["name"] = v
	}

	var parameters map[string]interface{}
	if forDatabase {
		parameters = s.toInterface().flattenDatabaseUserParameters(user.Parameters)
	} else {
		parameters = s.toInterface().flattenUserParameters(user.Parameters)
	}

	for k, v := range parameters {
		tfMap[k] = v
	}

	return tfMap
}

// FlattenDatabases converts api representation of list of databases to terraform representation.
func (s service) FlattenDatabases(databases []*paas.DatabaseResponse) []interface{} {
	if len(databases) == 0 {
		return nil
	}

	var tfList []interface{}

	for _, database := range databases {
		if database == nil {
			continue
		}

		tfList = append(tfList, s.toInterface().FlattenDatabase(database))
	}

	return tfList
}

// FlattenDatabase converts api representation of database to terraform representation.
func (s service) FlattenDatabase(database *paas.DatabaseResponse) map[string]interface{} {
	if database == nil {
		return map[string]interface{}{}
	}

	tfMap := map[string]interface{}{}

	if v := database.BackupEnabled; v != nil && s.allowBackup {
		tfMap["backup_enabled"] = v
	}

	if v := database.Id; v != nil {
		tfMap["id"] = v
	}

	if v := database.Name; v != nil {
		tfMap["name"] = v
	}

	if v := database.Users; v != nil {
		tfMap["user"] = s.toInterface().FlattenUsers(database.Users, true)
	}

	for k, v := range s.toInterface().flattenDatabaseParameters(database.Parameters) {
		tfMap[k] = v
	}

	return tfMap
}

// flattenServiceParameters converts api representation of service-specific parameters
// to terraform representation.
//
// If PaaS service has specific parameters, it should override this method.
func (s service) flattenServiceParameters(_ ServiceParameters) map[string]interface{} {
	return nil
}

// flattenUserParameters converts api representation of service-specific user parameters
// to terraform representation.
//
// If PaaS service has specific user parameters, it should override this method.
func (s service) flattenUserParameters(_ UserParameters) map[string]interface{} {
	return nil
}

// flattenDatabaseParameters converts api representation of service-specific database parameters
// to terraform representation.
//
// If PaaS service has specific database parameters, it should override this method.
func (s service) flattenDatabaseParameters(_ DatabaseParameters) map[string]interface{} {
	return nil
}

// flattenDatabaseUserParameters converts api representation of service-specific database user parameters
// to terraform representation.
//
// If PaaS service has specific database user parameters, it should override this method.
func (s service) flattenDatabaseUserParameters(_ DatabaseUserParameters) map[string]interface{} {
	return nil
}

// todo: doc
func (s service) flattenIntByteParameter(byteParameter interface{}) []map[string]interface{} {
	if vMap, okMap := byteParameter.(map[string]interface{}); okMap {
		if v, ok := vMap["value"].(int64); ok {
			return []map[string]interface{}{
				{
					"value":     strconv.FormatInt(v, 10),
					"dimension": vMap["dimension"].(string), // todo: optional?
				},
			}
		}
	}

	return nil
}

// todo: doc
func (s service) flattenFloatByteParameter(byteParameter interface{}) []map[string]interface{} {
	if vMap, okMap := byteParameter.(map[string]interface{}); okMap {
		var res map[string]interface{}

		if int64Val, int64Ok := vMap["value"].(int64); int64Ok {
			res["value"] = float64(int64Val)
		} else if floatVal, floatOk := vMap["value"].(float64); floatOk {
			res["value"] = floatVal
		} else {
			log.Printf("[ERROR] Unknown type '%T' for float parameter", vMap["value"])
			return nil
		}

		res["dimension"] = vMap["dimension"].(string)

		return []map[string]interface{}{res}
	}

	return nil
}
