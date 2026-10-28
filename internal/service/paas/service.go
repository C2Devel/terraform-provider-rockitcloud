package paas

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/paas"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	"github.com/hashicorp/terraform-provider-aws/internal/service/paas/services"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceService() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceServiceCreate,
		ReadContext:   resourceServiceRead,
		UpdateContext: resourceServiceUpdate,
		DeleteContext: resourceServiceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		// TODO: change timeouts depending on the number of nodes when send to wait functions
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},

		CustomizeDiff: validateServiceConfiguration,

		Schema: map[string]*schema.Schema{
			"additional_roles": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					// K2 PaaS Kafka currently exposes only the coordinator role here.
					ValidateFunc: validation.StringInSlice([]string{"coordinator"}, false),
				},
				ConflictsWith: []string{"coordinator"},
			},
			"arbitrator_required": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"auto_created_security_group_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"auto_disk_expansion": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"condition_value": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(autoDiskExpansionMinConditionValue, autoDiskExpansionMaxConditionValue),
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"max_size_gb": {
							Type:     schema.TypeInt,
							Required: true,
							ValidateFunc: validation.All(
								validation.IntAtLeast(autoDiskExpansionSizeMultipleGb),
								validation.IntDivisibleBy(autoDiskExpansionSizeMultipleGb),
							),
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"available_environment_versions": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"backup_settings": {
				Type:             schema.TypeList,
				Optional:         true,
				MaxItems:         1,
				DiffSuppressFunc: verify.SuppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"expiration_days": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 3650),
						},
						"notification_email": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_login": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"coordinator": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				MaxItems:      1,
				ConflictsWith: []string{"additional_roles"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_volume_iops": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"data_volume_size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"data_volume_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"root_volume_iops": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"root_volume_size": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"root_volume_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"data_volume": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"iops": {
							Type:             schema.TypeInt,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: iopsDiffSuppressFunc,
						},
						"size": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          32,
							DiffSuppressFunc: dataVolumeSizeDiffSuppressFunc,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  ec2.VolumeTypeSt2,
						},
					},
				},
			},
			// When true, service databases are managed by aws_paas_pgsql_database resources:
			// the service neither reads databases into state nor sends them to the API.
			"databases_managed_externally": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"delete_interfaces_on_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"endpoints": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"addresses": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"environment_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"high_availability": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"instances": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoints": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"interface_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"instance_qty": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"main": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"coordinator": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(3, 20),
					validation.StringMatch(
						regexp.MustCompile(`^[a-z\d][a-z\d\-.]+[a-z\d]$`),
						"name must start and end with Latin letters or number "+
							"and can only contain lowercase Latin letters, numbers, periods (.) and hyphens (-)",
					),
				),
			},
			"network_interface_ids": {
				Type:         schema.TypeSet,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"network_interface_ids", "subnet_ids"},
				Elem:         &schema.Schema{Type: schema.TypeString},
			},
			"root_volume": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"iops": {
							Type:             schema.TypeInt,
							Optional:         true,
							ForceNew:         true,
							Computed:         true,
							DiffSuppressFunc: iopsDiffSuppressFunc,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Default:  32,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  ec2.VolumeTypeSt2,
						},
					},
				},
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"service_class": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssh_key_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_ids": {
				Type:         schema.TypeSet,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ExactlyOneOf: []string{"network_interface_ids", "subnet_ids"},
				Elem:         &schema.Schema{Type: schema.TypeString},
			},
			"supported_features": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"total_cpu_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"user_data": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"user_data_content_type"},
			},
			"user_data_content_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"user_data"},
				ValidateFunc: validation.StringInSlice([]string{"cloud-config", "x-shellscript"}, false),
			},
			// When true, service users are managed by aws_paas_pgsql_user resources:
			// the service neither reads users into state nor sends them to the API.
			"users_managed_externally": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			services.ClickHouse.ServiceType():    services.ClickHouse.ResourceSchema(),
			services.ELK.ServiceType():           services.ELK.ResourceSchema(),
			services.ElasticSearch.ServiceType(): services.ElasticSearch.ResourceSchema(),
			services.GitLab.ServiceType():        services.GitLab.ResourceSchema(),
			services.Kafka.ServiceType():         services.Kafka.ResourceSchema(),
			services.Memcached.ServiceType():     services.Memcached.ResourceSchema(),
			services.MongoDB.ServiceType():       services.MongoDB.ResourceSchema(),
			services.MySQL.ServiceType():         services.MySQL.ResourceSchema(),
			services.PostgreSQL.ServiceType():    services.PostgreSQL.ResourceSchema(),
			services.Prometheus.ServiceType():    services.Prometheus.ResourceSchema(),
			services.RabbitMQ.ServiceType():      services.RabbitMQ.ResourceSchema(),
			services.Redis.ServiceType():         services.Redis.ResourceSchema(),
		},
	}
}

func resourceServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn
	name := d.Get("name").(string)

	input := &paas.CreateServiceInput{
		HighAvailability: aws.Bool(d.Get("high_availability").(bool)),
		InstanceType:     aws.String(d.Get("instance_type").(string)),
		Name:             aws.String(d.Get("name").(string)),
		RootVolumeType:   aws.String(d.Get("root_volume.0.type").(string)),
		RootVolumeSize:   aws.Int64(int64(d.Get("root_volume.0.size").(int))),
		SecurityGroupIds: flex.ExpandStringSet(d.Get("security_group_ids").(*schema.Set)),
	}

	if aws.StringValue(input.RootVolumeType) == ec2.VolumeTypeIo2 {
		input.RootVolumeIops = aws.Int64(int64(d.Get("root_volume.0.iops").(int)))
	}

	// For services that support arbitrator_required, send an explicit boolean.
	// For unsupported services (e.g. kafka), omit the field entirely.
	manager := serviceManager(d)

	if manager == nil {
		return diag.Errorf("PaaS Service configuration error: unknown service")
	}

	if manager.Service().AllowArbitrator() {
		// Send explicit boolean only for services that support arbitrator_required.
		input.ArbitratorRequired = aws.Bool(d.Get("arbitrator_required").(bool))
	}

	if v, ok := d.GetOk("additional_roles"); ok {
		input.AdditionalRoles = flex.ExpandStringList(v.([]interface{}))
	}

	if v, ok := d.GetOk("coordinator"); ok {
		coordList := v.([]interface{})
		if len(coordList) > 0 && coordList[0] != nil {
			input.Coordinator = expandCoordinator(coordList[0].(map[string]interface{}))
		}
	}

	if v, ok := d.GetOk("backup_settings"); ok {
		backupSettingsMap := v.([]interface{})[0].(map[string]interface{})
		input.BackupSettings = expandBackupSettings(backupSettingsMap)
	}

	if v, ok := d.GetOk("instance_qty"); ok {
		input.InstanceQty = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("auto_disk_expansion"); ok && v.([]interface{})[0] != nil {
		input.AdeSettings = expandAutoDiskExpansionSettings(v.([]interface{})[0].(map[string]interface{}))
	}

	if _, ok := d.GetOk("data_volume"); ok {
		input.DataVolumeType = aws.String(d.Get("data_volume.0.type").(string))
		input.DataVolumeSize = aws.Int64(int64(d.Get("data_volume.0.size").(int)))

		if aws.StringValue(input.DataVolumeType) == ec2.VolumeTypeIo2 {
			input.DataVolumeIops = aws.Int64(int64(d.Get("data_volume.0.iops").(int)))
		}
	}

	if v, ok := d.GetOk("network_interface_ids"); ok {
		input.NetworkInterfaceIds = flex.ExpandStringSet(v.(*schema.Set))
	} else {
		input.SubnetIds = flex.ExpandStringSet(d.Get("subnet_ids").(*schema.Set))
	}

	input.ServiceType = aws.String(manager.ServiceType())

	parametersMap := serviceParametersFromResourceData(d, manager.ServiceType())
	input.ServiceClass = aws.String(parametersMap["class"].(string))

	input.Parameters = manager.ExpandServiceCreationParameters(parametersMap)

	if v, ok := d.GetOk("ssh_key_name"); ok && v != "" {
		input.SshKeyName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("user_data"); ok && v != "" {
		input.UserData = aws.String(v.(string))
		input.UserDataContentType = aws.String(d.Get("user_data_content_type").(string))
	}

	log.Printf("[DEBUG] Creating PaaS %s Service %q", manager.ServiceType(), name)
	output, err := conn.CreateServiceWithContext(
		ctx,
		input,
		request.WithLogLevel(aws.LogOff),
	)

	if err != nil {
		return diag.Errorf("error creating PaaS Service with name %s: %s", name, err)
	}
	if output == nil || output.Service == nil || aws.StringValue(output.Service.Id) == "" {
		return diag.Errorf("error creating PaaS Service with name %s: missing service ID in response", name)
	}

	id := aws.StringValue(output.Service.Id)
	d.SetId(id)
	d.Set("service_type", output.Service.ServiceType)

	_, err = waitServiceCreated(ctx, conn, id, d.Timeout(schema.TimeoutCreate))

	if err != nil {
		return diag.Errorf("error waiting for PaaS Service (%s) to create: %s", id, err)
	}

	// Update to apply changes for service users, databases
	// and parameters, that can't be set on service creation.
	return resourceServiceUpdate(ctx, d, meta)
}

func resourceServiceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	previousGitLabParameters, _ := d.GetOk(services.ServiceTypeGitLab)
	autoDiskExpansionConfigured := len(d.Get("auto_disk_expansion").([]interface{})) > 0

	diags := readService(ctx, d, meta)
	if diags.HasError() || d.Id() == "" {
		return diags
	}

	// Disabled auto disk expansion is equivalent to an omitted block.
	if !autoDiskExpansionConfigured && !autoDiskExpansionEnabled(d.Get("auto_disk_expansion").([]interface{})) {
		d.Set("auto_disk_expansion", nil)
	}

	serviceType := d.Get("service_type").(string)
	manager := services.Manager(serviceType)
	if manager == nil {
		return diag.Errorf("error reading PaaS Service (%s): unknown service type %q", d.Id(), serviceType)
	}

	if err := setUnsupportedArbitratorRequired(d, manager); err != nil {
		return diag.Errorf("error setting arbitrator_required: %s", err)
	}
	if serviceType == services.ServiceTypeGitLab {
		parameters := d.Get(serviceType).([]interface{})
		preserveGitLabInputOnlyParameters(parameters, previousGitLabParameters)
		if err := d.Set(serviceType, parameters); err != nil {
			return diag.Errorf("error preserving GitLab input-only parameters: %s", err)
		}
	}

	return diags
}

func readService(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn
	id := d.Id()

	service, err := findServiceByID(ctx, conn, id)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] PaaS Service (%s) not found, removing from state", id)
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("error reading PaaS Service (%s): %+v", id, err)
	}

	if v := flattenBackupSettings(service.BackupSettings); len(v) > 0 {
		d.Set("backup_settings", v)
	}

	dataVolumeMap := []map[string]interface{}{
		{
			"type": aws.StringValue(service.DataVolumeType),
			"size": aws.Int64Value(service.DataVolumeSize),
			"iops": aws.Int64Value(service.DataVolumeIops),
		},
	}
	d.Set("available_environment_versions", service.AvailableEnvironmentVersions)

	if err := d.Set("auto_disk_expansion", flattenAutoDiskExpansionSettings(service.AutoDiskExpansionSettings)); err != nil {
		return diag.Errorf("error setting auto_disk_expansion: %s", err)
	}

	if aws.StringValue(service.DataVolumeType) != "" {
		d.Set("data_volume", dataVolumeMap)
	}

	d.Set("endpoints", flattenServiceEndpoints(service.Endpoints))

	d.Set("environment_version", service.EnvironmentVersion)

	d.Set("error_code", service.ErrorCode)
	d.Set("error_description", service.ErrorDescription)

	d.Set("high_availability", service.HighAvailability)

	d.Set("instances", flattenInstances(service.Instances))
	d.Set("instance_qty", service.InstanceQty)
	d.Set("instance_type", service.InstanceType)
	d.Set("nodes", flattenNodes(service.Nodes))

	d.Set("name", service.Name)
	if err := d.Set("additional_roles", flattenServiceAdditionalRoles(service)); err != nil {
		return diag.Errorf("error setting additional_roles: %s", err)
	}
	if err := d.Set("coordinator", flattenServiceCoordinator(service)); err != nil {
		return diag.Errorf("error setting coordinator: %s", err)
	}

	d.Set("network_interface_ids", service.NetworkInterfaceIds)

	serviceType := aws.StringValue(service.ServiceType)
	manager := services.Manager(serviceType)
	if manager == nil {
		return diag.Errorf("error reading PaaS Service (%s): unknown service type %q", id, serviceType)
	}

	parametersMap := manager.FlattenServiceParametersUsersDatabases(
		service.Parameters,
		service.Users,
		service.Databases,
	)
	if serviceType == services.ServiceTypeELK {
		preserveELKInputOnlyParameters(d, parametersMap)
	}
	// Externally managed users and databases belong to aws_paas_pgsql_user and
	// aws_paas_pgsql_database resources and must not appear as a diff here.
	// The keys are deleted rather than set to empty lists: d.Set writes missing blocks
	// as empty, and the read does not require the blocks to be in the schema.
	// aws_paas_pgsql has no user and database blocks, and d.Set rejects the whole
	// block with keys missing in the schema.
	if serviceUsersManagedExternally(d, serviceType) {
		delete(parametersMap, "user")
	}
	if serviceDatabasesManagedExternally(d, serviceType) {
		delete(parametersMap, "database")
	}
	parametersMap["class"] = service.ServiceClass
	d.Set(serviceType, []map[string]interface{}{parametersMap})

	rootVolumeMap := []map[string]interface{}{
		{
			"type": aws.StringValue(service.RootVolumeType),
			"size": aws.Int64Value(service.RootVolumeSize),
			"iops": aws.Int64Value(service.RootVolumeIops),
		},
	}
	d.Set("root_volume", rootVolumeMap)

	securityGroups := service.SecurityGroups
	var autoCreateSecurityGroupIds, nonAutoCreateSecurityGroupIds []*string
	for _, sg := range securityGroups {
		if aws.BoolValue(sg.CreatedAutomatically) {
			autoCreateSecurityGroupIds = append(autoCreateSecurityGroupIds, sg.Id)
		} else {
			nonAutoCreateSecurityGroupIds = append(nonAutoCreateSecurityGroupIds, sg.Id)
		}
	}
	d.Set("auto_created_security_group_ids", autoCreateSecurityGroupIds)
	d.Set("security_group_ids", nonAutoCreateSecurityGroupIds)

	d.Set("service_class", service.ServiceClass)
	d.Set("service_type", service.ServiceType)

	d.Set("ssh_key_name", service.SshKeyName)

	d.Set("status", service.Status)

	d.Set("subnet_ids", service.SubnetIds)

	d.Set("supported_features", service.SupportedFeatures)

	d.Set("total_cpu_count", service.TotalCpuCount)
	d.Set("total_memory", service.TotalMemory)

	d.Set("vpc_id", service.VpcId)

	return nil
}

func setUnsupportedArbitratorRequired(d *schema.ResourceData, manager services.ServiceManager) error {
	// DescribeService does not return arbitratorRequired. For service types that
	// cannot use an arbitrator, however, its value is unambiguously false and
	// must be restored during import. Otherwise Terraform sees the schema
	// default as removed and proposes replacing the service.
	if manager.Service().AllowArbitrator() {
		return nil
	}

	return d.Set("arbitrator_required", false)
}

func preserveELKInputOnlyParameters(d *schema.ResourceData, parametersMap map[string]interface{}) {
	if password, ok := parametersMap["password"].(string); !ok || password == "" {
		if password, ok := d.GetOk(services.ServiceTypeELK + ".0.password"); ok {
			parametersMap["password"] = password
		}
	}

	if options, ok := parametersMap["options"].(map[string]interface{}); !ok || len(options) == 0 {
		if options, ok := d.GetOk(services.ServiceTypeELK + ".0.options"); ok {
			parametersMap["options"] = options
		}
	}
}

func serviceParametersForUpdate(serviceType string, input services.ServiceParameters) services.ServiceParameters {
	if serviceType == services.ServiceTypeGitLab {
		// GitLab's root password is an initial credential, not an editable
		// service setting. Keep it on CreateService only, including when a
		// deferred integration update immediately follows creation.
		output := services.ServiceParameters{}
		for key, value := range input {
			if key != "password" {
				output[key] = value
			}
		}
		return output
	}
	if serviceType != services.ServiceTypeELK {
		return input
	}

	output := services.ServiceParameters{}
	for _, key := range []string{"monitoring", "monitor_by", "monitoring_labels"} {
		if value, ok := input[key]; ok {
			output[key] = value
		}
	}

	return output
}

func preserveGitLabInputOnlyParameters(parameters []interface{}, previous interface{}) {
	previousParameters, ok := previous.([]interface{})
	if !ok || len(previousParameters) == 0 || len(parameters) == 0 {
		return
	}
	current, currentOK := parameters[0].(map[string]interface{})
	old, oldOK := previousParameters[0].(map[string]interface{})
	if !currentOK || !oldOK {
		return
	}
	preserve := func(current, old map[string]interface{}, names []string) {
		for _, name := range names {
			if value, ok := current[name].(string); ok && value != "" {
				continue
			}
			if value, ok := old[name].(string); ok && value != "" {
				current[name] = value
			}
		}
	}
	preserve(current, old, []string{"password"})
	for block, names := range map[string][]string{
		"ldap": {"password", "client_cert", "client_key"},
		"oidc": {"client_secret"},
		"smtp": {"password"},
	} {
		currentBlocks, _ := current[block].([]interface{})
		oldBlocks, _ := old[block].([]interface{})
		// A missing remote or configured block must stay disabled. Preserve
		// only secrets of an integration that is present in both states.
		if len(currentBlocks) == 0 || len(oldBlocks) == 0 {
			continue
		}
		currentBlock, currentOK := currentBlocks[0].(map[string]interface{})
		oldBlock, oldOK := oldBlocks[0].(map[string]interface{})
		if currentOK && oldOK {
			preserve(currentBlock, oldBlock, names)
		}
	}
}

// changedServiceParameters returns only the expanded parameters whose Terraform attributes have changed.
func changedServiceParameters(d *schema.ResourceData, serviceType string, parameters services.ServiceParameters) services.ServiceParameters {
	prefix := serviceType + ".0."
	// log_to and monitor_by are required while logging and monitoring are enabled, so they are passed with any change of the block.
	attribute := map[string]string{
		"log_to":            "logging",
		"logging_tags":      "logging.0.logging_tags",
		"monitor_by":        "monitoring",
		"monitoring_labels": "monitoring.0.monitoring_labels",
	}

	output := services.ServiceParameters{}
	for key, value := range parameters {
		attr := key
		if v, ok := attribute[key]; ok {
			attr = v
		}

		if d.HasChange(prefix + attr) {
			output[key] = value
		}
	}

	// Expanders omit zero numeric values, so a changed attribute explicitly set to zero is passed separately.
	if v, ok := d.Get(serviceType).([]interface{}); ok && len(v) > 0 && v[0] != nil {
		for key, value := range v[0].(map[string]interface{}) {
			if _, ok := output[key]; ok || !d.HasChange(prefix+key) {
				continue
			}

			switch value {
			case 0:
				output[key] = int64(0)
			case 0.0:
				output[key] = 0.0
			}
		}
	}

	// PatchServiceParameters keeps omitted values, so removed tags and labels are passed as empty values.
	for key, empty := range map[string]interface{}{
		"logging_tags":      []interface{}{},
		"monitoring_labels": map[string]interface{}{},
	} {
		block, _, _ := strings.Cut(attribute[key], ".")
		if _, ok := output[key]; !ok && output[block] == true && d.HasChange(prefix+attribute[key]) {
			output[key] = empty
		}
	}

	return output
}

func resourceServiceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn
	id := d.Id()

	manager := serviceManager(d)

	if manager == nil {
		return diag.Errorf("PaaS Service (%s) configuration error: unknown service", id)
	}

	serviceType := manager.ServiceType()

	if d.HasChange("data_volume.0.size") && !d.IsNewResource() {
		oldRaw, newRaw := d.GetChange("data_volume.0.size")
		oldSize, newSize := oldRaw.(int), newRaw.(int)

		if err := validateIncreaseOnly("data_volume.size", oldSize, newSize); err != nil {
			return diag.FromErr(err)
		}

		input := &paas.ModifyInstanceVolumeSizeInput{
			ServiceId: aws.String(id),
			Size:      aws.Int64(int64(newSize)),
		}

		log.Printf("[DEBUG] Modifying PaaS Service data volume size: %+v", input)
		_, err := conn.ModifyInstanceVolumeSizeWithContext(ctx, input)

		if err != nil {
			return diag.Errorf("error modifying PaaS Service (%s) data volume size: %s", id, err)
		}

		_, err = waitServiceUpdated(ctx, conn, id, d.Timeout(schema.TimeoutUpdate))

		if err != nil {
			return diag.Errorf("error waiting for PaaS Service (%s) data volume size to update: %s", id, err)
		}
	}

	if d.HasChange("data_volume.0.iops") && !d.IsNewResource() {
		if volumeType := d.Get("data_volume.0.type").(string); strings.ToLower(volumeType) != ec2.VolumeTypeIo2 {
			return diag.Errorf("data_volume.iops can only be updated when data_volume.type is %q", ec2.VolumeTypeIo2)
		}

		input := &paas.ModifyInstanceVolumeIopsInput{
			ServiceId: aws.String(id),
			Iops:      aws.Int64(int64(d.Get("data_volume.0.iops").(int))),
		}

		log.Printf("[DEBUG] Modifying PaaS Service data volume IOPS: %+v", input)
		_, err := conn.ModifyInstanceVolumeIopsWithContext(ctx, input)

		if err != nil {
			return diag.Errorf("error modifying PaaS Service (%s) data volume IOPS: %s", id, err)
		}

		_, err = waitServiceUpdated(ctx, conn, id, d.Timeout(schema.TimeoutUpdate))

		if err != nil {
			return diag.Errorf("error waiting for PaaS Service (%s) data volume IOPS to update: %s", id, err)
		}
	}

	if d.HasChange("auto_disk_expansion") && !d.IsNewResource() {
		oldRaw, newRaw := d.GetChange("auto_disk_expansion")

		// UpdateAutoDiskExpansionSettings only modifies the service settings and does not start a task.
		if input := autoDiskExpansionUpdateInput(id, oldRaw.([]interface{}), newRaw.([]interface{})); input != nil {
			if err := updateAutoDiskExpansionSettings(ctx, conn, input); err != nil {
				return diag.Errorf("error modifying PaaS Service (%s) auto disk expansion settings: %s", id, err)
			}
		}
	}

	if d.HasChange("instance_type") && !d.IsNewResource() {
		nodeRoles, ok := serviceInstanceTypeNodeRolesFromState(d)

		if !ok {
			return diag.Errorf("error determining PaaS Service (%s) node roles for instance type update: nodes not in state, run terraform refresh first", id)
		}

		for _, nodeRole := range nodeRoles {
			input := &paas.ModifyInstanceTypeInput{
				ServiceId:    aws.String(id),
				InstanceType: aws.String(d.Get("instance_type").(string)),
				NodeRole:     aws.String(nodeRole),
			}

			log.Printf("[DEBUG] Modifying PaaS Service instance type: %+v", input)
			_, err := conn.ModifyInstanceTypeWithContext(ctx, input)

			if err != nil {
				return diag.Errorf("error modifying PaaS Service (%s) instance type for node role %s: %s", id, nodeRole, err)
			}

			_, err = waitServiceUpdated(ctx, conn, id, d.Timeout(schema.TimeoutUpdate))

			if err != nil {
				return diag.Errorf("error waiting for PaaS Service (%s) instance type for node role %s to update: %s", id, nodeRole, err)
			}
		}
	}

	parametersMap := serviceParametersFromResourceData(d, serviceType)

	if d.HasChange(serviceType) && (!d.IsNewResource() || manager.HasDeferredParameters(parametersMap)) {
		input := &paas.ModifyServiceParametersInput{
			ServiceId: aws.String(id),
		}

		input.Parameters = serviceParametersForUpdate(
			serviceType,
			manager.ExpandServiceParameters(parametersMap),
		)

		var err error
		modified := true
		if serviceType == services.ServiceTypePostgreSQL {
			// PatchServiceParameters changes only the passed parameters, so values set outside
			// of Terraform (for example, PgBouncer parameters) are preserved.
			patchInput := &paas.PatchServiceParametersInput{
				ServiceId:  aws.String(id),
				Parameters: changedServiceParameters(d, serviceType, input.Parameters),
			}

			// Users and databases are updated with ModifyService below.
			modified = len(patchInput.Parameters) > 0
			if modified {
				log.Printf("[DEBUG] Patching PaaS Service parameters: %s", patchInput)
				_, err = conn.PatchServiceParametersWithContext(ctx, patchInput)
			}
		} else {
			log.Printf("[DEBUG] Modifying PaaS %s Service (%s) parameters", serviceType, id)
			_, err = conn.ModifyServiceParametersWithContext(
				ctx,
				input,
				request.WithLogLevel(aws.LogOff),
			)
		}

		if err != nil {
			return diag.Errorf("error modifying PaaS Service (%s) parameters: %s", id, err)
		}

		if modified {
			_, err = waitServiceUpdated(ctx, conn, id, d.Timeout(schema.TimeoutUpdate))

			if err != nil {
				return diag.Errorf("error waiting for PaaS Service (%s) parameters to update: %s", id, err)
			}
		}
	}

	userKey := fmt.Sprintf("%s.0.user", serviceType)
	databaseKey := fmt.Sprintf("%s.0.database", serviceType)

	// Externally managed users and databases are never sent with ModifyService, because the API
	// replaces the whole lists and would delete objects created by aws_paas_pgsql_user
	// and aws_paas_pgsql_database.
	usersChanged := d.HasChange(userKey) && !serviceUsersManagedExternally(d, serviceType)
	databasesChanged := d.HasChange(databaseKey) && !serviceDatabasesManagedExternally(d, serviceType)

	if usersChanged || databasesChanged || d.HasChange("backup_settings") {
		input := &paas.ModifyServiceInput{
			ServiceId: aws.String(id),
		}

		if d.HasChange("backup_settings") {
			if v, ok := d.GetOk("backup_settings"); ok {
				backupSettingsMap := v.([]interface{})[0].(map[string]interface{})
				input.BackupSettings = expandBackupSettings(backupSettingsMap)
			}
		}

		if usersChanged {
			if v, ok := d.GetOk(userKey); ok {
				input.Users = manager.ExpandUsers(v.([]interface{}), false)
				if serviceType == services.ServiceTypePostgreSQL {
					previous, _ := d.GetChange(userKey)
					previousUsers, _ := previous.([]interface{})
					clearPostgreSQLUserRoles(input.Users, previousUsers)
				}
			} else {
				input.Users = []*paas.UserCreateRequest{}
			}
		}

		if databasesChanged {
			if v, ok := d.GetOk(databaseKey); ok {
				input.Databases = manager.ExpandDatabases(v.([]interface{}))
			} else {
				input.Databases = []*paas.DatabaseCreateRequest{}
			}
		}

		log.Printf("[DEBUG] Modifying PaaS %s Service (%s) users, databases and backup settings", serviceType, id)
		_, err := conn.ModifyServiceWithContext(
			ctx,
			input,
			request.WithLogLevel(aws.LogOff),
		)

		if err != nil {
			return diag.Errorf("error modifying PaaS Service (%s): %s", id, err)
		}

		_, err = waitServiceUpdated(ctx, conn, id, d.Timeout(schema.TimeoutUpdate))

		if err != nil {
			return diag.Errorf("error waiting for PaaS Service (%s) to update: %s", id, err)
		}

	}

	return resourceServiceRead(ctx, d, meta)
}

func resourceServiceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).PaaSConn
	id := d.Id()

	input := &paas.DeleteServiceInput{
		ServiceId:        aws.String(id),
		DeleteInterfaces: aws.Bool(d.Get("delete_interfaces_on_destroy").(bool)),
	}

	log.Printf("[DEBUG] Deleting PaaS Service: %s", input)
	_, err := conn.DeleteServiceWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, ServiceNotFoundCode) {
		log.Printf("[WARN] PaaS Service (%s) not found, removing from state", id)
		return nil
	}

	if err != nil {
		return diag.Errorf("error deleting PaaS Service (%s): %s", id, err)
	}

	_, err = waitServiceDeleted(ctx, conn, id, d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return diag.Errorf("error waiting for PaaS Service (%s) to delete: %s", d.Id(), err)
	}

	return nil
}

// iopsDiffSuppressFunc suppress diff if volume type is not io2 and iops is unset or configured as 0.
func iopsDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	volumeType := d.Get(strings.Replace(k, "iops", "type", 1)).(string)
	return strings.ToLower(volumeType) != ec2.VolumeTypeIo2 && new == "0"
}

func validateIncreaseOnly(name string, oldValue, newValue int) error {
	if newValue < oldValue {
		return fmt.Errorf("%s can only be increased in-place, got decrease from %d to %d", name, oldValue, newValue)
	}

	return nil
}

func serviceInstanceTypeNodeRolesFromState(d *schema.ResourceData) ([]string, bool) {
	v, ok := d.GetOk("nodes.0.main.0.role")
	if !ok {
		return nil, false
	}

	// A dedicated coordinator has its own instance_type configuration. Changing
	// the common instance_type must resize only the main nodes; a combined role
	// is kept intact so combined server/coordinator nodes are resized together.
	return []string{v.(string)}, true
}

func serviceManager(d *schema.ResourceData) services.ServiceManager {
	for _, serviceType := range services.ManagedServiceTypes() {
		_, exists := d.GetOk(serviceType)

		if exists {
			return services.Manager(serviceType)
		}
	}

	log.Printf("[WARN] There is no service specified in configuration.")
	return nil
}

func expandBackupSettings(tfMap map[string]interface{}) *paas.BackupSettingsRequest {
	if tfMap == nil {
		return nil
	}

	backupSettings := &paas.BackupSettingsRequest{}

	if v, ok := tfMap["bucket_name"].(string); ok && v != "" {
		backupSettings.BucketName = aws.String(v)
	}

	if v, ok := tfMap["enabled"].(bool); ok {
		backupSettings.Enabled = aws.Bool(v)
	}

	if v, ok := tfMap["expiration_days"].(int); ok && v != 0 {
		backupSettings.BackupExpirationDays = aws.Int64(int64(v))
	}

	if v, ok := tfMap["notification_email"].(string); ok && v != "" {
		backupSettings.NotificationEmail = aws.String(v)
	}

	if v, ok := tfMap["start_time"].(string); ok && v != "" {
		backupSettings.StartTime = aws.String(v)
	}

	if v, ok := tfMap["user_login"].(string); ok && v != "" {
		backupSettings.UserLogin = aws.String(v)
	}

	return backupSettings
}

func flattenBackupSettings(backupSettings *paas.BackupSettingsResponse) []map[string]interface{} {
	if backupSettings == nil {
		return []map[string]interface{}{}
	}

	tfMap := map[string]interface{}{}

	if v := backupSettings.BucketName; v != nil {
		tfMap["bucket_name"] = v
	}

	if v := backupSettings.Enabled; v != nil {
		tfMap["enabled"] = v
	}

	if v := backupSettings.BackupExpirationDays; v != nil {
		tfMap["expiration_days"] = v
	}

	if v := backupSettings.NotificationEmail; v != nil {
		tfMap["notification_email"] = v
	}

	if v := backupSettings.StartTime; v != nil {
		tfMap["start_time"] = v
	}

	if v := backupSettings.UserId; v != nil {
		tfMap["user_id"] = v
	}

	if v := backupSettings.UserLogin; v != nil {
		tfMap["user_login"] = v
	}

	// ignore when api returns `"backupSettings": {}` (block is omitted in config)
	if len(tfMap) == 0 {
		return []map[string]interface{}{}
	}

	return []map[string]interface{}{tfMap}
}

func flattenServiceEndpoints(endpoints []*paas.ServiceEndpoint) []map[string]interface{} {
	if endpoints == nil {
		return []map[string]interface{}{}
	}

	var tfList []map[string]interface{}
	for _, endpoint := range endpoints {
		if endpoint == nil {
			continue
		}

		tfMap := map[string]interface{}{}

		if v := endpoint.Addresses; v != nil {
			tfMap["addresses"] = v
		}

		if v := endpoint.Name; v != nil {
			tfMap["name"] = v
		}

		tfList = append(tfList, tfMap)
	}

	return tfList
}

func flattenNodes(nodes *paas.Nodes) []map[string]interface{} {
	if nodes == nil {
		return []map[string]interface{}{}
	}

	tfMap := map[string]interface{}{}

	if nodes.Main != nil {
		if role := aws.StringValue(nodes.Main.Role); role != "" {
			tfMap["main"] = []map[string]interface{}{{"role": role}}
		}
	}

	if nodes.Coordinator != nil {
		if role := aws.StringValue(nodes.Coordinator.Role); role != "" {
			tfMap["coordinator"] = []map[string]interface{}{{"role": role}}
		}
	}

	if len(tfMap) == 0 {
		return []map[string]interface{}{}
	}

	return []map[string]interface{}{tfMap}
}

func flattenInstances(instances []*paas.Instance) []map[string]interface{} {
	if instances == nil {
		return []map[string]interface{}{}
	}

	var tfList []map[string]interface{}
	for _, instance := range instances {
		if instance == nil {
			continue
		}

		tfMap := map[string]interface{}{}

		if v := instance.Endpoints; v != nil {
			tfMap["endpoints"] = flattenInstanceEndpoints(v)
		}

		if v := instance.Index; v != nil {
			tfMap["index"] = v
		}

		if v := instance.InstanceId; v != nil {
			tfMap["instance_id"] = v
		}

		if v := instance.InterfaceId; v != nil {
			tfMap["interface_id"] = v
		}

		if v := instance.Name; v != nil {
			tfMap["name"] = v
		}

		if v := instance.PrivateIp; v != nil {
			tfMap["private_ip"] = v
		}

		if v := instance.Role; v != nil {
			tfMap["role"] = v
		}

		if v := instance.Status; v != nil {
			tfMap["status"] = v
		}

		tfList = append(tfList, tfMap)
	}

	return tfList
}

func flattenInstanceEndpoints(endpoints []*paas.InstanceEndpoint) []map[string]interface{} {
	if endpoints == nil {
		return []map[string]interface{}{}
	}

	var tfList []map[string]interface{}
	for _, endpoint := range endpoints {
		if endpoint == nil {
			continue
		}

		tfMap := map[string]interface{}{}

		if v := endpoint.Address; v != nil {
			tfMap["address"] = v
		}

		if v := endpoint.Name; v != nil {
			tfMap["name"] = v
		}

		tfList = append(tfList, tfMap)
	}

	return tfList
}

func expandCoordinator(tfMap map[string]interface{}) *paas.NodeRequest {
	if tfMap == nil {
		return nil
	}

	nodeRequest := &paas.NodeRequest{
		InstanceType:   aws.String(tfMap["instance_type"].(string)),
		RootVolumeType: aws.String(tfMap["root_volume_type"].(string)),
		RootVolumeSize: aws.Int64(int64(tfMap["root_volume_size"].(int))),
	}

	if v, ok := tfMap["root_volume_iops"].(int); ok && v > 0 {
		nodeRequest.RootVolumeIops = aws.Int64(int64(v))
	}

	if v, ok := tfMap["data_volume_type"].(string); ok && v != "" {
		nodeRequest.DataVolumeType = aws.String(v)
	}

	if v, ok := tfMap["data_volume_size"].(int); ok && v > 0 {
		nodeRequest.DataVolumeSize = aws.Int64(int64(v))
	}

	if v, ok := tfMap["data_volume_iops"].(int); ok && v > 0 {
		nodeRequest.DataVolumeIops = aws.Int64(int64(v))
	}

	return nodeRequest
}

func flattenServiceAdditionalRoles(service *paas.Service) []string {
	if service == nil || service.Nodes == nil || service.Nodes.Main == nil {
		return nil
	}

	// Dedicated coordinator nodes are represented explicitly by nodes.coordinator.
	if service.Nodes.Coordinator != nil {
		return nil
	}

	for _, role := range strings.Split(aws.StringValue(service.Nodes.Main.Role), ",") {
		if strings.TrimSpace(role) == "coordinator" {
			return []string{"coordinator"}
		}
	}

	return nil
}

func flattenServiceCoordinator(service *paas.Service) []map[string]interface{} {
	if service == nil || service.Nodes == nil || service.Nodes.Coordinator == nil {
		return nil
	}

	node := service.Nodes.Coordinator

	return []map[string]interface{}{
		{
			"data_volume_iops": aws.Int64Value(node.DataVolumeIops),
			"data_volume_size": aws.Int64Value(node.DataVolumeSize),
			"data_volume_type": aws.StringValue(node.DataVolumeType),
			"instance_type":    aws.StringValue(node.InstanceType),
			"root_volume_iops": aws.Int64Value(node.RootVolumeIops),
			"root_volume_size": aws.Int64Value(node.RootVolumeSize),
			"root_volume_type": aws.StringValue(node.RootVolumeType),
		},
	}
}
