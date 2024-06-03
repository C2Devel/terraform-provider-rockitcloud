package services

import (
	"github.com/aws/aws-sdk-go/service/paas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ServiceParameters = map[string]interface{}
type UserParameters = map[string]interface{}
type DatabaseParameters = map[string]interface{}
type DatabaseUserParameters = map[string]interface{}

// Type service represents PaaS service and provides default implementations for ServiceManager methods.
type service struct {
	name               string
	class              []string
	defaultClass       string
	allowArbitrator    bool
	allowBackup        bool
	dataVolumeRequired bool
	usersEnabled       bool
	databasesEnabled   bool
	loggingEnabled     bool
	monitoringEnabled  bool
}

// Interface schemaBuilder provides methods for building resource and datasource schemas for service.
//
// Private methods represent schema for different types of service-specific parameters.
type schemaBuilder interface {
	ResourceSchema() *schema.Schema
	serviceParametersSchema() map[string]*schema.Schema
	userParametersSchema() map[string]*schema.Schema
	databaseParametersSchema() map[string]*schema.Schema
	databaseUserParametersSchema() map[string]*schema.Schema

	intByteParameterSchema() *schema.Schema
	floatByteParameterSchema() *schema.Schema

	DataSourceSchema() *schema.Schema
	serviceParametersDataSourceSchema() map[string]*schema.Schema
	userParametersDataSourceSchema() map[string]*schema.Schema
	databaseParametersDataSourceSchema() map[string]*schema.Schema
	databaseUserParametersDataSourceSchema() map[string]*schema.Schema

	intByteParameterDataSourceSchema() *schema.Schema
	floatByteParameterDataSourceSchema() *schema.Schema
}

// Interface converter provides methods for mapping service parameters
// between terraform and api representation.
//
// Expand: terraform -> api. Flatten: api -> terraform.
//
// Private methods and ExpandServiceParameters are used to convert different types of service-specific parameters.
type converter interface {
	ExpandServiceParameters(map[string]interface{}) ServiceParameters
	ExpandUsers([]interface{}, bool) []*paas.UserCreateRequest
	ExpandUser(map[string]interface{}, bool) *paas.UserCreateRequest
	ExpandDatabases([]interface{}) []*paas.DatabaseCreateRequest
	ExpandDatabase(map[string]interface{}) *paas.DatabaseCreateRequest
	expandServiceParameters(map[string]interface{}) ServiceParameters
	expandUserParameters(map[string]interface{}) UserParameters
	expandDatabaseParameters(map[string]interface{}) DatabaseParameters
	expandDatabaseUserParameters(map[string]interface{}) DatabaseUserParameters

	expandIntByteParameter(interface{}) map[string]interface{}
	expandFloatByteParameter(interface{}) map[string]interface{}

	FlattenServiceParametersUsersDatabases(
		ServiceParameters, []*paas.UserResponse, []*paas.DatabaseResponse,
	) map[string]interface{}
	FlattenUsers([]*paas.UserResponse, bool) []interface{}
	FlattenUser(*paas.UserResponse, bool) map[string]interface{}
	FlattenDatabases([]*paas.DatabaseResponse) []interface{}
	FlattenDatabase(response *paas.DatabaseResponse) map[string]interface{}
	flattenServiceParameters(ServiceParameters) map[string]interface{}
	flattenUserParameters(UserParameters) map[string]interface{}
	flattenDatabaseParameters(DatabaseParameters) map[string]interface{}
	flattenDatabaseUserParameters(DatabaseUserParameters) map[string]interface{}

	flattenIntByteParameter(interface{}) []map[string]interface{}
	flattenFloatByteParameter(interface{}) []map[string]interface{}
}

// ServiceManager is a main public interface that includes all methods,
// that need to be implemented to use new PaaS services in terraform configuration.
type ServiceManager interface {
	ServiceType() string
	Service() service

	schemaBuilder
	converter
}

// Service returns itself. It is used to access service parameters from ServiceManager implementation.
func (s service) Service() service {
	return s
}

func (s service) ServiceType() string {
	return s.name
}

// toInterface returns implementation of ServiceManager for specified service type.
// It is used to call overridden methods in ServiceManager implementations.
func (s service) toInterface() ServiceManager {
	return Manager(s.ServiceType())
}
