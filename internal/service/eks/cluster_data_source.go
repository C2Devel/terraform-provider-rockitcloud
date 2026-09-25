package eks

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

func DataSourceCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceClusterRead,

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
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
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
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_family": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pod_ipv4_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_ipv4_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"legacy_cluster_params": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_autoscaler_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_autoscaler_required": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"cluster_autoscaler_user": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"docker_registry_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"volume_iops": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"volume_size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"volume_type": {
										Type:     schema.TypeString,
										Computed: true,
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
										Computed: true,
									},
								},
							},
						},
						"ingress_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"public_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"volume_iops": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"volume_size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"volume_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"master_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"high_availability": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"instance_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"public_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"volume_iops": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"volume_size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"volume_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"user_data_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_data": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"user_data_content_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"nlb_provider_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nlb_user": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"placement_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"affinity": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tenancy": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"host_id": {
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
			},
			"platform_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tftags.TagsSchemaComputed(),
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_config": {
				Type:     schema.TypeList,
				Computed: true,
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
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"security_group_ids": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"subnet_ids": {
							Type:     schema.TypeSet,
							Computed: true,
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

func dataSourceClusterRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EKSConn
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	name := d.Get("name").(string)
	cluster, err := FindClusterByName(conn, name)

	if err != nil {
		return fmt.Errorf("error reading EKS Cluster (%s): %w", name, err)
	}

	d.SetId(name)
	d.Set("arn", cluster.Arn)

	if err := d.Set("certificate_authority", flattenCertificate(cluster.CertificateAuthority)); err != nil {
		return fmt.Errorf("error setting certificate_authority: %w", err)
	}

	d.Set("created_at", aws.TimeValue(cluster.CreatedAt).String())

	if err := d.Set("enabled_cluster_log_types", flattenEnabledLogTypes(cluster.Logging)); err != nil {
		return fmt.Errorf("error setting enabled_cluster_log_types: %w", err)
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

	d.Set("name", cluster.Name)
	d.Set("platform_version", cluster.PlatformVersion)
	d.Set("role_arn", cluster.RoleArn)
	d.Set("status", cluster.Status)

	d.Set("version", cluster.Version)

	if err := d.Set("vpc_config", flattenVPCConfigResponse(cluster.ResourcesVpcConfig)); err != nil {
		return fmt.Errorf("error setting vpc_config: %w", err)
	}

	if err := d.Set("tags", KeyValueTags(cluster.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	return nil
}
