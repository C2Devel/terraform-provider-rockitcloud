package eks

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tfec2 "github.com/hashicorp/terraform-provider-aws/internal/service/ec2"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

const (
	clusterCreateRetryTimeout = 10 * time.Minute
	clusterDeleteRetryTimeout = 60 * time.Minute
)

func ResourceCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceClusterCreate,
		Read:   resourceClusterRead,
		Update: resourceClusterUpdate,
		Delete: resourceClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		CustomizeDiff: verify.SetTagsDiff,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(clusterDeleteRetryTimeout),
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_authority": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled_cluster_log_types": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(eks.LogType_Values(), true),
				},
				Set: schema.HashString,
			},
			"encryption_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key_arn": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"resources": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice(Resources_Values(), false),
							},
						},
					},
				},
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"identity": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"oidc": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"issuer": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"kubernetes_network_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_family": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{eks.IpFamilyIpv4}, false),
						},
						"service_ipv4_cidr": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
							ValidateFunc: validation.All(
								validation.IsCIDRNetwork(12, 24),
								validation.StringMatch(regexp.MustCompile(`^(10|172\.(1[6-9]|2[0-9]|3[0-1])|192\.168)\..*`), "must be within 10.0.0.0/8, 172.16.0.0/12, or 192.168.0.0/16"),
							),
						},
						"pod_ipv4_cidr": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validation.IsCIDR,
						},
					},
				},
			},
			"legacy_cluster_params": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_autoscaler_config": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_autoscaler_required": {
										Type:     schema.TypeBool,
										Required: true,
										ForceNew: true,
									},
									"cluster_autoscaler_user": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						"docker_registry_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"volume_iops": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
									},
									"volume_size": {
										Type:     schema.TypeInt,
										Required: true,
										ForceNew: true,
									},
									"volume_type": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
						"ebs_provider_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ebs_user": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
						"ingress_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_type": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"public_ip": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"volume_iops": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
									},
									"volume_size": {
										Type:     schema.TypeInt,
										Required: true,
										ForceNew: true,
									},
									"volume_type": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
						"master_config": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"high_availability": {
										Type:     schema.TypeBool,
										Required: true,
										ForceNew: true,
									},
									"instance_type": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"public_ip": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"volume_iops": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
									},
									"volume_size": {
										Type:     schema.TypeInt,
										Required: true,
										ForceNew: true,
									},
									"volume_type": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
						"user_data_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_data": {
										Type:     schema.TypeString,
										Required: true,
									},
									"user_data_content_type": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice(
											[]string{"cloud-config", "x-shellscript"},
											false,
										),
									},
								},
							},
						},
						"nlb_provider_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nlb_user": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
						"placement_config": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"affinity": {
										Type:         schema.TypeString,
										ForceNew:     true,
										Computed:     true,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(eks.Affinity_Values(), false),
									},
									"tenancy": {
										Type:         schema.TypeString,
										Optional:     true,
										ForceNew:     true,
										Default:      eks.TenancyDefault,
										ValidateFunc: validation.StringInSlice(eks.Tenancy_Values(), false),
									},
									"host_id": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validClusterName,
			},
			"remote_access_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ec2_ssh_key": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"platform_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_arn": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true, // FIXME: Remove after UpdateClusterVersion is supported in C2 EKS API.
			},
			"vpc_config": {
				Type:     schema.TypeList,
				MinItems: 1,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_private_access": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"endpoint_public_access": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"public_access_cidrs": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: verify.ValidCIDRNetworkAddress,
							},
						},
						"security_group_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"subnet_ids": {
							Type:     schema.TypeSet,
							Required: true,
							ForceNew: true,
							MinItems: 1,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceClusterCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EKSConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))
	name := d.Get("name").(string)

	input := &eks.CreateClusterInput{
		Name:               aws.String(name),
		ResourcesVpcConfig: testAccClusterConfig_expandVPCRequest(d.Get("vpc_config").([]interface{})),
	}

	if v, ok := d.GetOk("role_arn"); ok && v.(string) != "" {
		input.RoleArn = aws.String(v.(string))
	}

	if _, ok := d.GetOk("kubernetes_network_config"); ok {
		input.KubernetesNetworkConfig = expandNetworkConfigRequest(d.Get("kubernetes_network_config").([]interface{}))
	}

	if v, ok := d.GetOk("legacy_cluster_params"); ok {
		input.LegacyClusterParams = expandLegacyClusterParams(v.([]interface{}))
	}

	if v, ok := d.GetOk("remote_access_config"); ok {
		input.RemoteAccessConfig = expandClusterRemoteAccessConfig(v.([]interface{}))
	}

	if v, ok := d.GetOk("version"); ok {
		input.Version = aws.String(v.(string))
	}

	if len(tags) > 0 {
		input.Tags = Tags(tags.IgnoreAWS())
	}

	clusterName := d.Get("name").(string)

	log.Printf("[DEBUG] Creating EKS Cluster: %s", input)
	var output *eks.CreateClusterOutput
	err := resource.Retry(clusterCreateRetryTimeout, func() *resource.RetryError {
		var err error

		output, err = conn.CreateCluster(input)

		// InvalidParameterException: roleArn, arn:aws:iam::123456789012:role/XXX, does not exist
		if tfawserr.ErrMessageContains(err, eks.ErrCodeInvalidParameterException, "does not exist") {
			return resource.RetryableError(err)
		}

		// InvalidParameterException: Error in role params
		if tfawserr.ErrMessageContains(err, eks.ErrCodeInvalidParameterException, "Error in role params") {
			return resource.RetryableError(err)
		}

		if tfawserr.ErrMessageContains(err, eks.ErrCodeInvalidParameterException, "Role could not be assumed because the trusted entity is not correct") {
			return resource.RetryableError(err)
		}

		// InvalidParameterException: The provided role doesn't have the Amazon EKS Managed Policies associated with it. Please ensure the following policy is attached: arn:aws:iam::aws:policy/AmazonEKSClusterPolicy
		if tfawserr.ErrMessageContains(err, eks.ErrCodeInvalidParameterException, "The provided role doesn't have the Amazon EKS Managed Policies associated with it") {
			return resource.RetryableError(err)
		}

		// InvalidParameterException: IAM role's policy must include the `ec2:DescribeSubnets` action
		if tfawserr.ErrMessageContains(err, eks.ErrCodeInvalidParameterException, "IAM role's policy must include") {
			return resource.RetryableError(err)
		}

		if tfawserr.ErrCodeEquals(err, ErrCodeIPAddressInUse) {
			log.Printf("[WARN] The specified IP address for EKS Cluster (%s) is in use, trying to create EKS Cluster for some time", clusterName)
			return resource.RetryableError(err)
		}

		cluster, _ := FindClusterByName(conn, clusterName)
		if cluster != nil && aws.StringValue(cluster.Status) == eks.ClusterStatusDeleted {
			log.Printf("[WARN] EKS Cluster with the same name (%s) found in the DELETED state, waiting for the specified name to become available", clusterName)
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		output, err = conn.CreateCluster(input)
	}

	if err != nil {
		return fmt.Errorf("error creating EKS Cluster (%s): %w", name, err)
	}

	d.SetId(aws.StringValue(output.Cluster.Name))

	_, err = waitClusterCreated(conn, d.Id(), d.Timeout(schema.TimeoutCreate))

	if err != nil {
		return fmt.Errorf("error waiting for EKS Cluster (%s) to create: %w", d.Id(), err)
	}

	return resourceClusterRead(d, meta)
}

func resourceClusterRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EKSConn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	cluster, err := FindClusterByName(conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] EKS Cluster (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading EKS Cluster (%s): %w", d.Id(), err)
	}

	d.Set("arn", cluster.Arn)

	if err := d.Set("certificate_authority", flattenCertificate(cluster.CertificateAuthority)); err != nil {
		return fmt.Errorf("error setting certificate_authority: %w", err)
	}

	d.Set("created_at", aws.TimeValue(cluster.CreatedAt).String())

	if err := d.Set("enabled_cluster_log_types", flattenEnabledLogTypes(cluster.Logging)); err != nil {
		return fmt.Errorf("error setting enabled_cluster_log_types: %w", err)
	}

	if err := d.Set("encryption_config", flattenEncryptionConfig(cluster.EncryptionConfig)); err != nil {
		return fmt.Errorf("error setting encryption_config: %w", err)
	}

	d.Set("endpoint", cluster.Endpoint)

	if err := d.Set("identity", flattenIdentity(cluster.Identity)); err != nil {
		return fmt.Errorf("error setting identity: %w", err)
	}

	if err := d.Set("kubernetes_network_config", flattenNetworkConfig(cluster.KubernetesNetworkConfig)); err != nil {
		return fmt.Errorf("error setting kubernetes_network_config: %w", err)
	}

	if err := d.Set("legacy_cluster_params", flattenLegacyClusterParams(cluster.LegacyClusterParams)); err != nil {
		return fmt.Errorf("error setting legacy_cluster_params: %w", err)
	}

	if err := d.Set("remote_access_config", flattenClusterRemoteAccessConfig(cluster.RemoteAccessConfig)); err != nil {
		return fmt.Errorf("error setting remote_access_config: %w", err)
	}

	d.Set("name", cluster.Name)
	d.Set("platform_version", cluster.PlatformVersion)
	d.Set("role_arn", cluster.RoleArn)
	d.Set("status", cluster.Status)
	d.Set("version", cluster.Version)

	if err := d.Set("vpc_config", flattenVPCConfigResponse(cluster.ResourcesVpcConfig)); err != nil {
		return fmt.Errorf("error setting vpc_config: %w", err)
	}

	tags := KeyValueTags(cluster.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	// lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}

	return nil
}

func resourceClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EKSConn
	connEc2 := meta.(*conns.AWSClient).EC2Conn

	if d.HasChange("vpc_config.0.security_group_ids") {
		input := &eks.UpdateClusterConfigInput{
			Name:               aws.String(d.Id()),
			ResourcesVpcConfig: expandVPCSecurityGroupUpdateRequest(d.Get("vpc_config.0.security_group_ids").(*schema.Set)),
		}

		log.Printf("[DEBUG] Updating EKS Cluster (%s) security groups: %s", d.Id(), input)
		output, err := conn.UpdateClusterConfig(input)

		if err != nil {
			return fmt.Errorf("error updating EKS Cluster (%s) security groups: %w", d.Id(), err)
		}

		updateID := aws.StringValue(output.Update.Id)

		_, err = waitClusterUpdateSuccessful(conn, d.Id(), updateID, d.Timeout(schema.TimeoutUpdate))

		if err != nil {
			return fmt.Errorf("error waiting for EKS Cluster (%s) security group update (%s): %w", d.Id(), updateID, err)
		}
	}

	if d.HasChange("legacy_cluster_params.0.user_data_config") {
		input := &eks.UpdateClusterUserDataInput{
			Name: aws.String(d.Id()),
			UserDataConfig: expandUserDataConfig(
				d.Get("legacy_cluster_params.0.user_data_config").([]interface{}),
			),
		}

		log.Printf("[DEBUG] Updating EKS Cluster (%s) user data: %s", d.Id(), input)
		output, err := conn.UpdateClusterUserData(input)

		if err != nil {
			return fmt.Errorf("error updating EKS Cluster (%s) user data: %w", d.Id(), err)
		}

		updateID := aws.StringValue(output.Update.Id)

		_, err = waitClusterUpdateSuccessful(conn, d.Id(), updateID, d.Timeout(schema.TimeoutUpdate))

		if err != nil {
			return fmt.Errorf("error waiting for EKS Cluster (%s) user data update (%s): %w", d.Id(), updateID, err)
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")
		// FIXME: Use eks.UpdateTags after TagResource and UntagResource are supported in C2 EKS API.
		// To use EC2 API arn contains the cluster id.
		if err := tfec2.UpdateTags(connEc2, d.Get("arn").(string), o, n); err != nil {
			return fmt.Errorf("error updating tags: %w", err)
		}
	}

	return resourceClusterRead(d, meta)
}

func resourceClusterDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EKSConn

	log.Printf("[DEBUG] Deleting EKS Cluster: %s", d.Id())

	input := &eks.DeleteClusterInput{
		Name: aws.String(d.Id()),
	}

	// If a cluster is scaling up due to load a delete request will fail
	// This is a temporary workaround until EKS supports multiple parallel mutating operations
	err := tfresource.RetryConfigContext(context.Background(), 0*time.Second, 1*time.Minute, 0*time.Second, 30*time.Second, clusterDeleteRetryTimeout, func() *resource.RetryError {
		var err error

		_, err = conn.DeleteCluster(input)

		if tfawserr.ErrMessageContains(err, eks.ErrCodeResourceInUseException, "in progress") {
			log.Printf("[DEBUG] eks cluster update in progress: %v", err)
			return resource.RetryableError(err)
		}

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		_, err = conn.DeleteCluster(input)
	}

	if tfawserr.ErrCodeEquals(err, ErrCodeClusterNotFound) {
		log.Printf("[WARN] EKS Cluster (%s) not found, removing from state", d.Id())
		return nil
	}

	// Sometimes the EKS API returns the ResourceNotFound error in this form:
	// ClientException: No cluster found for name: tf-acc-test-0o1f8
	if tfawserr.ErrMessageContains(err, eks.ErrCodeClientException, "No cluster found for name:") {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting EKS Cluster (%s): %w", d.Id(), err)
	}

	if _, err = waitClusterDeleted(conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return fmt.Errorf("error waiting for EKS Cluster (%s) to delete: %w", d.Id(), err)
	}

	return nil
}

func testAccClusterConfig_expandVPCRequest(l []interface{}) *eks.VpcConfigRequest {
	if len(l) == 0 {
		return nil
	}

	m := l[0].(map[string]interface{})

	vpcConfigRequest := &eks.VpcConfigRequest{
		SecurityGroupIds: flex.ExpandStringSet(m["security_group_ids"].(*schema.Set)),
		SubnetIds:        flex.ExpandStringSet(m["subnet_ids"].(*schema.Set)),
	}

	return vpcConfigRequest
}

func expandVPCSecurityGroupUpdateRequest(securityGroupIDs *schema.Set) *eks.VpcConfigRequest {
	return &eks.VpcConfigRequest{
		SecurityGroupIds: flex.ExpandStringSet(securityGroupIDs),
	}
}

func expandNetworkConfigRequest(tfList []interface{}) *eks.KubernetesNetworkConfigRequest {
	tfMap, ok := tfList[0].(map[string]interface{})

	if !ok {
		return nil
	}

	apiObject := &eks.KubernetesNetworkConfigRequest{}

	if v, ok := tfMap["service_ipv4_cidr"].(string); ok && v != "" {
		apiObject.ServiceIpv4Cidr = aws.String(v)
	}

	if v, ok := tfMap["pod_ipv4_cidr"].(string); ok && v != "" {
		apiObject.PodIpv4Cidr = aws.String(v)
	}

	if v, ok := tfMap["ip_family"].(string); ok && v != "" {
		apiObject.IpFamily = aws.String(v)
	}

	return apiObject
}

func expandLegacyClusterParams(tfList []interface{}) *eks.LegacyClusterParamsRequest {
	if len(tfList) == 0 {
		return nil
	}

	tfMap, ok := tfList[0].(map[string]interface{})
	if !ok || tfMap == nil {
		return nil
	}

	legacyParams := &eks.LegacyClusterParamsRequest{}

	if clusterAutoscalerConfig, ok := tfMap["cluster_autoscaler_config"].([]interface{}); ok && len(clusterAutoscalerConfig) > 0 {
		legacyParams.ClusterAutoscalerConfig = expandClusterAutoscalerConfig(clusterAutoscalerConfig)
	}

	if dockerRegistryConfig, ok := tfMap["docker_registry_config"].([]interface{}); ok && len(dockerRegistryConfig) > 0 {
		legacyParams.DockerRegistryConfig = expandDockerRegistryConfig(dockerRegistryConfig)
	}

	if ebsProviderConfig, ok := tfMap["ebs_provider_config"].([]interface{}); ok && len(ebsProviderConfig) > 0 {
		legacyParams.EbsProviderConfig = expandEbsProviderConfig(ebsProviderConfig)
	}

	if ingressConfig, ok := tfMap["ingress_config"].([]interface{}); ok && len(ingressConfig) > 0 {
		legacyParams.IngressConfig = expandIngressConfig(ingressConfig)
	}

	if masterConfig, ok := tfMap["master_config"].([]interface{}); ok && len(masterConfig) > 0 {
		legacyParams.MasterConfig = expandMasterConfig(masterConfig)
	}

	if userDataConfig, ok := tfMap["user_data_config"].([]interface{}); ok && len(userDataConfig) > 0 {
		legacyParams.UserDataConfig = expandUserDataConfig(userDataConfig)
	}

	if nlbProviderConfig, ok := tfMap["nlb_provider_config"].([]interface{}); ok && len(nlbProviderConfig) > 0 {
		legacyParams.NlbProviderConfig = expandNlbProviderConfig(nlbProviderConfig)
	}

	if placementConfig, ok := tfMap["placement_config"].([]interface{}); ok && len(placementConfig) > 0 {
		legacyParams.PlacementConfig = expandPlacementConfig(placementConfig)
	}

	return legacyParams
}

func expandClusterAutoscalerConfig(tfList []interface{}) *eks.ClusterAutoscalerConfig {
	if len(tfList) == 0 || tfList[0] == nil {
		return nil
	}

	tfMap := tfList[0].(map[string]interface{})
	config := &eks.ClusterAutoscalerConfig{
		ClusterAutoscalerRequired: aws.Bool(tfMap["cluster_autoscaler_required"].(bool)),
	}

	if v, ok := tfMap["cluster_autoscaler_user"].(string); ok && v != "" {
		config.ClusterAutoscalerUser = aws.String(v)
	}

	return config
}

func expandDockerRegistryConfig(tfList []interface{}) *eks.DockerRegistryConfig {
	if len(tfList) == 0 {
		return nil
	}

	tfMap, ok := tfList[0].(map[string]interface{})
	if !ok || tfMap == nil {
		return nil
	}

	dockerRegistryConfig := &eks.DockerRegistryConfig{
		DockerRegistryRequired: aws.Bool(true),
	}

	if v, ok := tfMap["volume_iops"].(int); ok && v != 0 {
		dockerRegistryConfig.DockerRegistryVolumeIops = aws.Int64(int64(v))
	}

	if v, ok := tfMap["volume_size"].(int); ok && v != 0 {
		dockerRegistryConfig.DockerRegistryVolumeSize = aws.Int64(int64(v))
	}

	if v, ok := tfMap["volume_type"].(string); ok && v != "" {
		dockerRegistryConfig.DockerRegistryVolumeType = aws.String(v)
	}

	return dockerRegistryConfig
}

func expandEbsProviderConfig(tfList []interface{}) *eks.EbsProviderConfigRequest {
	if len(tfList) == 0 {
		return nil
	}

	tfMap, ok := tfList[0].(map[string]interface{})
	if !ok || tfMap == nil {
		return nil
	}

	ebsProviderConfig := &eks.EbsProviderConfigRequest{
		EbsProviderRequired: aws.Bool(true),
	}

	if v, ok := tfMap["ebs_user"].(string); ok && v != "" {
		ebsProviderConfig.EbsUser = aws.String(v)
	}

	return ebsProviderConfig
}

func expandIngressConfig(tfList []interface{}) *eks.IngressConfig {
	if len(tfList) == 0 {
		return nil
	}

	tfMap, ok := tfList[0].(map[string]interface{})
	if !ok || tfMap == nil {
		return nil
	}

	ingressConfig := &eks.IngressConfig{
		IngressRequired: aws.Bool(true),
	}

	if v, ok := tfMap["instance_type"].(string); ok && v != "" {
		ingressConfig.IngressInstanceType = aws.String(v)
	}

	if v, ok := tfMap["public_ip"].(string); ok && v != "" {
		ingressConfig.IngressPublicIp = aws.String(v)
	}

	if v, ok := tfMap["volume_iops"].(int); ok && v != 0 {
		ingressConfig.IngressVolumeIops = aws.Int64(int64(v))
	}

	if v, ok := tfMap["volume_size"].(int); ok && v != 0 {
		ingressConfig.IngressVolumeSize = aws.Int64(int64(v))
	}

	if v, ok := tfMap["volume_type"].(string); ok && v != "" {
		ingressConfig.IngressVolumeType = aws.String(v)
	}

	return ingressConfig
}

func expandPlacementConfig(tfList []interface{}) *eks.PlacementConfig {
	if len(tfList) == 0 {
		return nil
	}

	tfMap, ok := tfList[0].(map[string]interface{})
	if !ok || tfMap == nil {
		return nil
	}
	placementConfig := &eks.PlacementConfig{}

	if v, ok := tfMap["affinity"].(string); ok && v != "" {
		placementConfig.Affinity = aws.String(v)
	}

	if v, ok := tfMap["tenancy"].(string); ok && v != "" {
		placementConfig.Tenancy = aws.String(v)
	}

	if v, ok := tfMap["host_id"].(string); ok && v != "" {
		placementConfig.HostId = aws.String(v)
	}

	return placementConfig
}

func expandMasterConfig(tfList []interface{}) *eks.MasterConfig {
	if len(tfList) == 0 {
		return nil
	}

	tfMap, ok := tfList[0].(map[string]interface{})
	if !ok || tfMap == nil {
		return nil
	}

	masterConfig := &eks.MasterConfig{}

	if v, ok := tfMap["high_availability"].(bool); ok {
		masterConfig.HighAvailability = aws.Bool(v)
	}

	if v, ok := tfMap["instance_type"].(string); ok && v != "" {
		masterConfig.MastersInstanceType = aws.String(v)
	}

	if v, ok := tfMap["public_ip"].(string); ok && v != "" {
		masterConfig.MasterPublicIp = aws.String(v)
	}

	if v, ok := tfMap["volume_iops"].(int); ok && v != 0 {
		masterConfig.MastersVolumeIops = aws.Int64(int64(v))
	}

	if v, ok := tfMap["volume_size"].(int); ok && v != 0 {
		masterConfig.MastersVolumeSize = aws.Int64(int64(v))
	}

	if v, ok := tfMap["volume_type"].(string); ok && v != "" {
		masterConfig.MastersVolumeType = aws.String(v)
	}

	return masterConfig
}

func expandUserDataConfig(tfList []interface{}) *eks.UserDataConfig {
	if len(tfList) == 0 {
		return nil
	}

	tfMap, ok := tfList[0].(map[string]interface{})
	if !ok || tfMap == nil {
		return nil
	}

	userDataConfig := &eks.UserDataConfig{}

	if v, ok := tfMap["user_data"].(string); ok && v != "" {
		userDataConfig.UserData = aws.String(v)
	}

	if v, ok := tfMap["user_data_content_type"].(string); ok && v != "" {
		userDataConfig.UserDataContentType = aws.String(v)
	}

	return userDataConfig
}

func expandNlbProviderConfig(tfList []interface{}) *eks.NlbProviderConfigRequest {
	if len(tfList) == 0 {
		return nil
	}

	tfMap, ok := tfList[0].(map[string]interface{})
	if !ok || tfMap == nil {
		return nil
	}

	nlbProviderConfig := &eks.NlbProviderConfigRequest{
		NlbProviderRequired: aws.Bool(true),
	}

	if v, ok := tfMap["nlb_user"].(string); ok && v != "" {
		nlbProviderConfig.NlbUser = aws.String(v)
	}

	return nlbProviderConfig
}

func expandClusterRemoteAccessConfig(tfList []interface{}) *eks.RemoteAccessConfig {
	if len(tfList) == 0 || tfList[0] == nil {
		return nil
	}

	tfMap := tfList[0].(map[string]interface{})
	config := &eks.RemoteAccessConfig{}

	if v, ok := tfMap["ec2_ssh_key"].(string); ok && v != "" {
		config.Ec2SshKey = aws.String(v)
	}

	return config
}

func flattenCertificate(certificate *eks.Certificate) []map[string]interface{} {
	if certificate == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"data": aws.StringValue(certificate.Data),
	}

	return []map[string]interface{}{m}
}

func flattenIdentity(identity *eks.Identity) []map[string]interface{} {
	if identity == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"oidc": flattenOIDC(identity.Oidc),
	}

	return []map[string]interface{}{m}
}

func flattenOIDC(oidc *eks.OIDC) []map[string]interface{} {
	if oidc == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"issuer": aws.StringValue(oidc.Issuer),
	}

	return []map[string]interface{}{m}
}

func flattenEncryptionConfig(apiObjects []*eks.EncryptionConfig) []interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []interface{}

	for _, apiObject := range apiObjects {
		tfMap := map[string]interface{}{
			"provider":  flattenProvider(apiObject.Provider),
			"resources": aws.StringValueSlice(apiObject.Resources),
		}

		tfList = append(tfList, tfMap)
	}

	return tfList
}

func flattenProvider(apiObject *eks.Provider) []interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{
		"key_arn": aws.StringValue(apiObject.KeyArn),
	}

	return []interface{}{tfMap}
}

func flattenVPCConfigResponse(vpcConfig *eks.VpcConfigResponse) []map[string]interface{} {
	if vpcConfig == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"cluster_security_group_id": aws.StringValue(vpcConfig.ClusterSecurityGroupId),
		"endpoint_private_access":   aws.BoolValue(vpcConfig.EndpointPrivateAccess),
		"endpoint_public_access":    aws.BoolValue(vpcConfig.EndpointPublicAccess),
		"security_group_ids":        flex.FlattenStringSet(vpcConfig.SecurityGroupIds),
		"subnet_ids":                flex.FlattenStringSet(vpcConfig.SubnetIds),
		"public_access_cidrs":       flex.FlattenStringSet(vpcConfig.PublicAccessCidrs),
		"vpc_id":                    aws.StringValue(vpcConfig.VpcId),
	}

	return []map[string]interface{}{m}
}

func flattenEnabledLogTypes(logging *eks.Logging) *schema.Set {
	enabledLogTypes := []*string{}

	if logging != nil {
		logSetups := logging.ClusterLogging
		for _, logSetup := range logSetups {
			if logSetup == nil || !aws.BoolValue(logSetup.Enabled) {
				continue
			}

			enabledLogTypes = append(enabledLogTypes, logSetup.Types...)
		}
	}

	return flex.FlattenStringSet(enabledLogTypes)
}

func flattenNetworkConfig(apiObject *eks.KubernetesNetworkConfigResponse) []interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{
		"pod_ipv4_cidr":     aws.StringValue(apiObject.PodIpv4Cidr),
		"service_ipv4_cidr": aws.StringValue(apiObject.ServiceIpv4Cidr),
		"ip_family":         aws.StringValue(apiObject.IpFamily),
	}

	return []interface{}{tfMap}
}

func flattenLegacyClusterParams(legacyParams *eks.LegacyClusterParamsResponse) []interface{} {
	if legacyParams == nil {
		return nil
	}

	tfMap := map[string]interface{}{
		"cluster_autoscaler_config": flattenClusterAutoscalerConfig(legacyParams.ClusterAutoscalerConfig),
		"docker_registry_config":    flattenDockerRegistryConfig(legacyParams.DockerRegistryConfig),
		"ebs_provider_config":       flattenEbsProviderConfig(legacyParams.EbsProviderConfig),
		"ingress_config":            flattenIngressConfig(legacyParams.IngressConfig),
		"master_config":             flattenMasterConfig(legacyParams.MasterConfig),
		"user_data_config":          flattenUserDataConfig(legacyParams.UserDataConfig),
		"placement_config":          flattenPlacementConfig(legacyParams.PlacementConfig),
		"nlb_provider_config":       flattenNlbProviderConfig(legacyParams.NlbProviderConfig),
	}

	hasConfig := false
	for _, value := range tfMap {
		if block, ok := value.([]interface{}); ok && len(block) > 0 {
			hasConfig = true
			break
		}
	}
	if !hasConfig {
		return nil
	}

	return []interface{}{tfMap}
}

func flattenClusterAutoscalerConfig(config *eks.ClusterAutoscalerConfig) []interface{} {
	if config == nil {
		return nil
	}

	user := aws.StringValue(config.ClusterAutoscalerUser)
	if user == "" {
		user = aws.StringValue(config.ClusterAutoscalerUserName)
	}

	return []interface{}{map[string]interface{}{
		"cluster_autoscaler_required": aws.BoolValue(config.ClusterAutoscalerRequired),
		"cluster_autoscaler_user":     user,
	}}
}

func flattenDockerRegistryConfig(dockerRegistryConfig *eks.DockerRegistryConfig) []interface{} {
	if dockerRegistryConfig == nil {
		return nil
	}
	if !aws.BoolValue(dockerRegistryConfig.DockerRegistryRequired) &&
		aws.Int64Value(dockerRegistryConfig.DockerRegistryVolumeIops) == 0 &&
		aws.Int64Value(dockerRegistryConfig.DockerRegistryVolumeSize) == 0 &&
		aws.StringValue(dockerRegistryConfig.DockerRegistryVolumeType) == "" {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := dockerRegistryConfig.DockerRegistryVolumeIops; v != nil {
		tfMap["volume_iops"] = aws.Int64Value(v)
	}

	if v := dockerRegistryConfig.DockerRegistryVolumeSize; v != nil {
		tfMap["volume_size"] = aws.Int64Value(v)
	}

	if v := dockerRegistryConfig.DockerRegistryVolumeType; v != nil {
		tfMap["volume_type"] = aws.StringValue(v)
	}

	if len(tfMap) == 0 {
		return nil
	}

	return []interface{}{tfMap}
}

func flattenEbsProviderConfig(ebsProviderConfig *eks.EbsProviderConfigResponse) []interface{} {
	if ebsProviderConfig == nil {
		return nil
	}
	if !aws.BoolValue(ebsProviderConfig.EbsProviderRequired) &&
		aws.StringValue(ebsProviderConfig.EbsUser) == "" &&
		aws.StringValue(ebsProviderConfig.EbsUserName) == "" {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := ebsProviderConfig.EbsUser; v != nil {
		tfMap["ebs_user"] = aws.StringValue(v)
	} else if v := ebsProviderConfig.EbsUserName; v != nil {
		tfMap["ebs_user"] = aws.StringValue(v)
	}

	if len(tfMap) == 0 {
		return nil
	}

	return []interface{}{tfMap}
}

func flattenIngressConfig(ingressConfig *eks.IngressConfig) []interface{} {
	if ingressConfig == nil {
		return nil
	}
	if !aws.BoolValue(ingressConfig.IngressRequired) &&
		aws.StringValue(ingressConfig.IngressInstanceType) == "" &&
		aws.StringValue(ingressConfig.IngressPublicIp) == "" &&
		aws.Int64Value(ingressConfig.IngressVolumeIops) == 0 &&
		aws.Int64Value(ingressConfig.IngressVolumeSize) == 0 &&
		aws.StringValue(ingressConfig.IngressVolumeType) == "" {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := ingressConfig.IngressInstanceType; v != nil {
		tfMap["instance_type"] = aws.StringValue(v)
	}

	if v := ingressConfig.IngressPublicIp; v != nil {
		tfMap["public_ip"] = aws.StringValue(v)
	}

	if v := ingressConfig.IngressVolumeIops; v != nil {
		tfMap["volume_iops"] = aws.Int64Value(v)
	}

	if v := ingressConfig.IngressVolumeSize; v != nil {
		tfMap["volume_size"] = aws.Int64Value(v)
	}

	if v := ingressConfig.IngressVolumeType; v != nil {
		tfMap["volume_type"] = aws.StringValue(v)
	}

	if len(tfMap) == 0 {
		return nil
	}

	return []interface{}{tfMap}
}

func flattenMasterConfig(masterConfig *eks.MasterConfig) []interface{} {
	if masterConfig == nil {
		return nil
	}

	tfMap := map[string]interface{}{
		"high_availability": aws.BoolValue(masterConfig.HighAvailability),
		"instance_type":     aws.StringValue(masterConfig.MastersInstanceType),
		"public_ip":         aws.StringValue(masterConfig.MasterPublicIp),
		"volume_iops":       aws.Int64Value(masterConfig.MastersVolumeIops),
		"volume_size":       aws.Int64Value(masterConfig.MastersVolumeSize),
		"volume_type":       aws.StringValue(masterConfig.MastersVolumeType),
	}

	return []interface{}{tfMap}
}

func flattenUserDataConfig(userDataConfig *eks.UserDataConfig) []interface{} {
	if userDataConfig == nil {
		return nil
	}
	if aws.StringValue(userDataConfig.UserData) == "" &&
		aws.StringValue(userDataConfig.UserDataContentType) == "" {
		return nil
	}

	tfMap := map[string]interface{}{
		"user_data":              aws.StringValue(userDataConfig.UserData),
		"user_data_content_type": aws.StringValue(userDataConfig.UserDataContentType),
	}

	return []interface{}{tfMap}
}

func flattenNlbProviderConfig(nlbProviderConfig *eks.NlbProviderConfigResponse) []interface{} {
	if nlbProviderConfig == nil {
		return nil
	}
	if !aws.BoolValue(nlbProviderConfig.NlbProviderRequired) &&
		aws.StringValue(nlbProviderConfig.NlbUser) == "" &&
		aws.StringValue(nlbProviderConfig.NlbUserName) == "" {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := nlbProviderConfig.NlbUser; v != nil {
		tfMap["nlb_user"] = aws.StringValue(v)
	} else if v := nlbProviderConfig.NlbUserName; v != nil {
		tfMap["nlb_user"] = aws.StringValue(v)
	}

	if len(tfMap) == 0 {
		return nil
	}

	return []interface{}{tfMap}
}

func flattenClusterRemoteAccessConfig(config *eks.RemoteAccessConfig) []interface{} {
	if config == nil || aws.StringValue(config.Ec2SshKey) == "" {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"ec2_ssh_key": aws.StringValue(config.Ec2SshKey),
	}}
}

func flattenPlacementConfig(placementConfig *eks.PlacementConfig) []interface{} {
	if placementConfig == nil {
		return nil
	}
	if aws.StringValue(placementConfig.Affinity) == "" &&
		aws.StringValue(placementConfig.HostId) == "" &&
		(aws.StringValue(placementConfig.Tenancy) == "" || aws.StringValue(placementConfig.Tenancy) == "default") {
		return nil
	}

	tfMap := map[string]interface{}{
		"affinity": aws.StringValue(placementConfig.Affinity),
		"tenancy":  aws.StringValue(placementConfig.Tenancy),
		"host_id":  aws.StringValue(placementConfig.HostId),
	}

	return []interface{}{tfMap}
}
