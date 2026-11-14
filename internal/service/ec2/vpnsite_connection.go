package ec2

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"net"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

const vpnConnectionOperationInProgressTimeout = 10 * time.Minute

var vpnConnectionTunnelPrefixes = []string{"tunnel1_", "tunnel2_"}

// Configurable tunnel options, without tunnel prefix.
var vpnConnectionTunnelArguments = []string{
	"ike_versions",
	"inside_cidr",
	"phase1_dh_group_numbers",
	"phase1_encryption_algorithms",
	"phase1_integrity_algorithms",
	"phase1_lifetime_seconds",
	"phase2_dh_group_numbers",
	"phase2_encryption_algorithms",
	"phase2_integrity_algorithms",
	"phase2_lifetime_seconds",
	"preshared_key",
	"replay_window_size",
}

func ResourceVPNConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceVPNConnectionCreate,
		Read:   resourceVPNConnectionRead,
		Update: resourceVPNConnectionUpdate,
		Delete: resourceVPNConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		// TODO: change modify
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"customer_gateway_configuration": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"customer_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"local_ipv4_network_cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDRNetwork(0, 32),
			},
			"remote_ipv4_network_cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDRNetwork(0, 32),
			},
			"static_routes_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
			"tunnel1_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tunnel1_bgp_asn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tunnel1_bgp_holdtime": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tunnel1_cgw_inside_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tunnel1_ike_versions": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(VpnTunnelOptionsIKEVersion_Values(), false),
				},
			},
			"tunnel1_inside_cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validateVpnConnectionTunnelInsideCIDR,
			},
			"tunnel1_phase1_dh_group_numbers": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 6,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntInSlice([]int{5, 14, 15, 16, 17, 18, 19, 20, 21}),
				},
			},
			"tunnel1_phase1_encryption_algorithms": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(VpnTunnelOptionsPhase1EncryptionAlgorithm_Values(), false),
				},
			},
			"tunnel1_phase1_integrity_algorithms": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(VpnTunnelOptionsPhase1IntegrityAlgorithm_Values(), false),
				},
			},
			"tunnel1_phase1_lifetime_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(900, 28800),
			},
			"tunnel1_phase2_dh_group_numbers": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntInSlice([]int{0, 5, 14, 15, 16, 17, 18, 19, 20, 21}),
				},
			},
			"tunnel1_phase2_encryption_algorithms": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(VpnTunnelOptionsPhase2EncryptionAlgorithm_Values(), false),
				},
			},
			"tunnel1_phase2_integrity_algorithms": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(VpnTunnelOptionsPhase2IntegrityAlgorithm_Values(), false),
				},
			},
			"tunnel1_phase2_lifetime_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(900, 3600),
			},
			"tunnel1_preshared_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validateVpnConnectionTunnelPreSharedKey(),
			},
			"tunnel1_replay_window_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(32, 2048),
			},
			"tunnel1_vgw_inside_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tunnel2_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tunnel2_bgp_asn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tunnel2_bgp_holdtime": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tunnel2_cgw_inside_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tunnel2_ike_versions": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(VpnTunnelOptionsIKEVersion_Values(), false),
				},
			},
			"tunnel2_inside_cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validateVpnConnectionTunnelInsideCIDR,
			},
			"tunnel2_phase1_dh_group_numbers": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 6,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntInSlice([]int{5, 14, 15, 16, 17, 18, 19, 20, 21}),
				},
			},
			"tunnel2_phase1_encryption_algorithms": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(VpnTunnelOptionsPhase1EncryptionAlgorithm_Values(), false),
				},
			},
			"tunnel2_phase1_integrity_algorithms": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(VpnTunnelOptionsPhase1IntegrityAlgorithm_Values(), false),
				},
			},
			"tunnel2_phase1_lifetime_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(900, 28800),
			},
			"tunnel2_phase2_dh_group_numbers": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntInSlice([]int{0, 5, 14, 15, 16, 17, 18, 19, 20, 21}),
				},
			},
			"tunnel2_phase2_encryption_algorithms": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(VpnTunnelOptionsPhase2EncryptionAlgorithm_Values(), false),
				},
			},
			"tunnel2_phase2_integrity_algorithms": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(VpnTunnelOptionsPhase2IntegrityAlgorithm_Values(), false),
				},
			},
			"tunnel2_phase2_lifetime_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(900, 3600),
			},
			"tunnel2_preshared_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validateVpnConnectionTunnelPreSharedKey(),
			},
			"tunnel2_replay_window_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(32, 2048),
			},
			"tunnel2_vgw_inside_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tunnel_inside_ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      ec2.TunnelInsideIpVersionIpv4,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{ec2.TunnelInsideIpVersionIpv4}, false),
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(VpnConnectionType_Values(), false),
			},
			"vgw_telemetry": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accepted_route_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"last_status_change": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"outside_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			// ID of VPC is accepted as well.
			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(v interface{}) string {
					return NormalizeVPNGatewayID(v.(string))
				},
			},
		},

		CustomizeDiff: customdiff.Sequence(
			verify.SetTagsDiff,
			resourceVPNConnectionCustomizeDiff,
		),
	}
}

// Returns first configured argument of tunnel, empty string if none.
func vpnConnectionTunnelConfiguredArgument(config cty.Value, prefix string) string {
	for _, k := range vpnConnectionTunnelArguments {
		v := config.GetAttr(prefix + k)

		if v.IsNull() {
			continue
		}

		if v.IsKnown() && v.CanIterateElements() && v.LengthInt() == 0 {
			continue
		}

		return prefix + k
	}

	return ""
}

func resourceVPNConnectionCustomizeDiff(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	config := d.GetRawConfig()

	if config.IsNull() || !config.IsKnown() {
		return nil
	}

	// Sets instead of MaxItems: argument types are kept same as in AWS provider.
	for _, prefix := range vpnConnectionTunnelPrefixes {
		for _, k := range []string{"ike_versions", "phase2_dh_group_numbers"} {
			if v := config.GetAttr(prefix + k); v.IsKnown() && !v.IsNull() && v.LengthInt() > 1 {
				return fmt.Errorf("%s%s: only one value is supported", prefix, k)
			}
		}
	}

	if d.Get("static_routes_only").(bool) {
		return fmt.Errorf("static_routes_only: only false is supported")
	}

	// Number of TunnelOptions entries is number of tunnels, second entry cannot be sent alone.
	tunnel1Configured := vpnConnectionTunnelConfiguredArgument(config, "tunnel1_") != ""
	tunnel2Argument := vpnConnectionTunnelConfiguredArgument(config, "tunnel2_")

	if tunnel2Argument != "" && !tunnel1Configured {
		return fmt.Errorf("%s requires at least one tunnel1_* argument", tunnel2Argument)
	}

	// Removed tunnel2_* arguments give no diff (computed), replacement is forced.
	if d.Id() != "" && tunnel1Configured && tunnel2Argument == "" && d.Get("tunnel2_address").(string) != "" {
		if err := d.SetNewComputed("tunnel2_address"); err != nil {
			return err
		}

		return d.ForceNew("tunnel2_address")
	}

	return nil
}

func resourceVPNConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	customerGatewayID := d.Get("customer_gateway_id").(string)
	vpnGatewayID := NormalizeVPNGatewayID(d.Get("vpn_gateway_id").(string))

	// Create returns existing VPN connection of same gateways instead of error.
	existing, err := FindVPNConnectionByCustomerGatewayIDAndVPNGatewayID(conn, customerGatewayID, vpnGatewayID)

	if err != nil && !tfresource.NotFound(err) {
		return fmt.Errorf(
			"error reading EC2 VPN Connections of Customer Gateway (%s) and VPN Gateway (%s): %w",
			customerGatewayID, vpnGatewayID, err,
		)
	}

	if existing != nil {
		return fmt.Errorf(
			"EC2 VPN Connection (%s) between Customer Gateway (%s) and VPN Gateway (%s) already exists",
			aws.StringValue(existing.VpnConnectionId), customerGatewayID, vpnGatewayID,
		)
	}

	options := expandVpnConnectionOptionsSpecification(d)

	// Values unknown during plan can turn out empty.
	if expandVpnTunnelOptionsSpecification(d, "tunnel1_") == nil && len(options.TunnelOptions) > 0 {
		return fmt.Errorf("tunnel2_* arguments require at least one tunnel1_* argument")
	}

	input := &ec2.CreateVpnConnectionInput{
		CustomerGatewayId: aws.String(customerGatewayID),
		Options:           options,
		TagSpecifications: ec2TagSpecificationsFromKeyValueTags(tags, ec2.ResourceTypeVpnConnection),
		Type:              aws.String(d.Get("type").(string)),
		VpnGatewayId:      aws.String(vpnGatewayID),
	}

	log.Printf("[DEBUG] Creating EC2 VPN Connection: %s", input)
	outputRaw, err := tfresource.RetryWhenAWSErrCodeEquals(
		vpnConnectionOperationInProgressTimeout,
		func() (interface{}, error) {
			return conn.CreateVpnConnection(input)
		},
		ErrCodeOperationInProgress,
	)

	if err != nil {
		return fmt.Errorf("error creating EC2 VPN Connection: %w", err)
	}

	output := outputRaw.(*ec2.CreateVpnConnectionOutput)
	d.SetId(aws.StringValue(output.VpnConnection.VpnConnectionId))

	if _, err := WaitVPNConnectionCreated(conn, d.Id()); err != nil {
		return fmt.Errorf("error waiting for EC2 VPN Connection (%s) create: %w", d.Id(), err)
	}

	return resourceVPNConnectionRead(d, meta)
}

func resourceVPNConnectionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	vpnConnection, err := FindVPNConnectionByID(conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] EC2 VPN Connection (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading EC2 VPN Connection (%s): %w", d.Id(), err)
	}

	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   ec2.ServiceName,
		Region:    meta.(*conns.AWSClient).Region,
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("vpn-connection/%s", d.Id()),
	}.String()
	d.Set("arn", arn)
	d.Set("customer_gateway_configuration", vpnConnection.CustomerGatewayConfiguration)
	d.Set("customer_gateway_id", vpnConnection.CustomerGatewayId)
	d.Set("type", vpnConnection.Type)
	d.Set("vpn_gateway_id", vpnConnection.VpnGatewayId)

	if err := d.Set("vgw_telemetry", flattenVgwTelemetries(vpnConnection.VgwTelemetry)); err != nil {
		return fmt.Errorf("error setting vgw_telemetry: %w", err)
	}

	tags := KeyValueTags(vpnConnection.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}

	// BGP parameters and inside addresses of tunnels are reported in XML only.
	var config XmlVpnConnectionConfig

	if v := aws.StringValue(vpnConnection.CustomerGatewayConfiguration); v != "" {
		if err := xml.Unmarshal([]byte(v), &config); err != nil {
			return fmt.Errorf("error parsing EC2 VPN Connection (%s) customer gateway configuration: %w", d.Id(), err)
		}
	}

	staticRoutesOnly := false
	tunnelInsideIpVersion := ec2.TunnelInsideIpVersionIpv4

	var tunnelOptions []*ec2.TunnelOption

	if v := vpnConnection.Options; v != nil {
		tunnelOptions = v.TunnelOptions

		d.Set("local_ipv4_network_cidr", v.LocalIpv4NetworkCidr)
		d.Set("remote_ipv4_network_cidr", v.RemoteIpv4NetworkCidr)

		if v.StaticRoutesOnly != nil {
			staticRoutesOnly = aws.BoolValue(v.StaticRoutesOnly)
		}

		if v.TunnelInsideIpVersion != nil {
			tunnelInsideIpVersion = aws.StringValue(v.TunnelInsideIpVersion)
		}
	}

	d.Set("static_routes_only", staticRoutesOnly)
	d.Set("tunnel_inside_ip_version", tunnelInsideIpVersion)

	for i, prefix := range vpnConnectionTunnelPrefixes {
		var tunnelOption *ec2.TunnelOption

		if i < len(tunnelOptions) {
			tunnelOption = tunnelOptions[i]
		}

		if err := flattenTunnelOption(d, prefix, tunnelOption, config.tunnel(i, tunnelOption)); err != nil {
			return err
		}
	}

	return nil
}

func resourceVPNConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	// TODO: change modify
	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Id(), o, n); err != nil {
			return fmt.Errorf("error updating EC2 VPN Connection (%s) tags: %w", d.Id(), err)
		}
	}

	return resourceVPNConnectionRead(d, meta)
}

func resourceVPNConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	log.Printf("[INFO] Deleting EC2 VPN Connection: %s", d.Id())
	_, err := tfresource.RetryWhenAWSErrCodeEquals(
		vpnConnectionOperationInProgressTimeout,
		func() (interface{}, error) {
			return conn.DeleteVpnConnection(&ec2.DeleteVpnConnectionInput{
				VpnConnectionId: aws.String(d.Id()),
			})
		},
		ErrCodeOperationInProgress,
	)

	if tfawserr.ErrCodeEquals(err, ErrCodeInvalidVpnConnectionIDNotFound) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting EC2 VPN Connection (%s): %w", d.Id(), err)
	}

	if _, err := WaitVPNConnectionDeleted(conn, d.Id()); err != nil {
		return fmt.Errorf("error waiting for EC2 VPN Connection (%s) delete: %w", d.Id(), err)
	}

	return nil
}

func expandVpnConnectionOptionsSpecification(d *schema.ResourceData) *ec2.VpnConnectionOptionsSpecification {
	apiObject := &ec2.VpnConnectionOptionsSpecification{}

	if v, ok := d.GetOk("local_ipv4_network_cidr"); ok {
		apiObject.LocalIpv4NetworkCidr = aws.String(v.(string))
	}

	if v, ok := d.GetOk("remote_ipv4_network_cidr"); ok {
		apiObject.RemoteIpv4NetworkCidr = aws.String(v.(string))
	}

	apiObject.StaticRoutesOnly = aws.Bool(d.Get("static_routes_only").(bool))
	apiObject.TunnelInsideIpVersion = aws.String(d.Get("tunnel_inside_ip_version").(string))

	// No entries or two give two tunnels, one entry gives single tunnel.
	for _, prefix := range vpnConnectionTunnelPrefixes {
		if v := expandVpnTunnelOptionsSpecification(d, prefix); v != nil {
			apiObject.TunnelOptions = append(apiObject.TunnelOptions, v)
		}
	}

	return apiObject
}

func expandVpnTunnelOptionsSpecification(d *schema.ResourceData, prefix string) *ec2.VpnTunnelOptionsSpecification {
	apiObject := &ec2.VpnTunnelOptionsSpecification{}
	configured := false

	if v, ok := d.GetOk(prefix + "ike_versions"); ok && v.(*schema.Set).Len() > 0 {
		for _, v := range v.(*schema.Set).List() {
			apiObject.IKEVersions = append(apiObject.IKEVersions, &ec2.IKEVersionsRequestListValue{
				Value: aws.String(v.(string)),
			})
		}

		configured = true
	}

	if v, ok := d.GetOk(prefix + "inside_cidr"); ok {
		apiObject.TunnelInsideCidr = aws.String(v.(string))
		configured = true
	}

	if v, ok := d.GetOk(prefix + "phase1_dh_group_numbers"); ok && v.(*schema.Set).Len() > 0 {
		for _, v := range v.(*schema.Set).List() {
			apiObject.Phase1DHGroupNumbers = append(apiObject.Phase1DHGroupNumbers, &ec2.Phase1DHGroupNumbersRequestListValue{
				Value: aws.Int64(int64(v.(int))),
			})
		}

		configured = true
	}

	if v, ok := d.GetOk(prefix + "phase1_encryption_algorithms"); ok && v.(*schema.Set).Len() > 0 {
		for _, v := range v.(*schema.Set).List() {
			apiObject.Phase1EncryptionAlgorithms = append(apiObject.Phase1EncryptionAlgorithms, &ec2.Phase1EncryptionAlgorithmsRequestListValue{
				Value: aws.String(v.(string)),
			})
		}

		configured = true
	}

	if v, ok := d.GetOk(prefix + "phase1_integrity_algorithms"); ok && v.(*schema.Set).Len() > 0 {
		for _, v := range v.(*schema.Set).List() {
			apiObject.Phase1IntegrityAlgorithms = append(apiObject.Phase1IntegrityAlgorithms, &ec2.Phase1IntegrityAlgorithmsRequestListValue{
				Value: aws.String(v.(string)),
			})
		}

		configured = true
	}

	if v, ok := d.GetOk(prefix + "phase1_lifetime_seconds"); ok {
		apiObject.Phase1LifetimeSeconds = aws.Int64(int64(v.(int)))
		configured = true
	}

	if v, ok := d.GetOk(prefix + "phase2_dh_group_numbers"); ok && v.(*schema.Set).Len() > 0 {
		for _, v := range v.(*schema.Set).List() {
			apiObject.Phase2DHGroupNumbers = append(apiObject.Phase2DHGroupNumbers, &ec2.Phase2DHGroupNumbersRequestListValue{
				Value: aws.Int64(int64(v.(int))),
			})
		}

		configured = true
	}

	if v, ok := d.GetOk(prefix + "phase2_encryption_algorithms"); ok && v.(*schema.Set).Len() > 0 {
		for _, v := range v.(*schema.Set).List() {
			apiObject.Phase2EncryptionAlgorithms = append(apiObject.Phase2EncryptionAlgorithms, &ec2.Phase2EncryptionAlgorithmsRequestListValue{
				Value: aws.String(v.(string)),
			})
		}

		configured = true
	}

	if v, ok := d.GetOk(prefix + "phase2_integrity_algorithms"); ok && v.(*schema.Set).Len() > 0 {
		for _, v := range v.(*schema.Set).List() {
			apiObject.Phase2IntegrityAlgorithms = append(apiObject.Phase2IntegrityAlgorithms, &ec2.Phase2IntegrityAlgorithmsRequestListValue{
				Value: aws.String(v.(string)),
			})
		}

		configured = true
	}

	if v, ok := d.GetOk(prefix + "phase2_lifetime_seconds"); ok {
		apiObject.Phase2LifetimeSeconds = aws.Int64(int64(v.(int)))
		configured = true
	}

	if v, ok := d.GetOk(prefix + "preshared_key"); ok {
		apiObject.PreSharedKey = aws.String(v.(string))
		configured = true
	}

	if v, ok := d.GetOk(prefix + "replay_window_size"); ok {
		apiObject.ReplayWindowSize = aws.Int64(int64(v.(int)))
		configured = true
	}

	if !configured {
		return nil
	}

	return apiObject
}

// Both arguments are nil for missing second tunnel.
func flattenTunnelOption(d *schema.ResourceData, prefix string, apiObject *ec2.TunnelOption, config *XmlIpsecTunnel) error {
	if apiObject == nil {
		apiObject = &ec2.TunnelOption{}
	}

	if config == nil {
		config = &XmlIpsecTunnel{}
	}

	var s []*string
	var i []*int64

	for _, v := range apiObject.IkeVersions {
		s = append(s, v.Value)
	}
	if err := d.Set(prefix+"ike_versions", aws.StringValueSlice(s)); err != nil {
		return fmt.Errorf("error setting %sike_versions: %w", prefix, err)
	}
	s = nil

	for _, v := range apiObject.Phase1DHGroupNumbers {
		i = append(i, v.Value)
	}
	if err := d.Set(prefix+"phase1_dh_group_numbers", aws.Int64ValueSlice(i)); err != nil {
		return fmt.Errorf("error setting %sphase1_dh_group_numbers: %w", prefix, err)
	}
	i = nil

	for _, v := range apiObject.Phase1EncryptionAlgorithms {
		s = append(s, v.Value)
	}
	if err := d.Set(prefix+"phase1_encryption_algorithms", aws.StringValueSlice(s)); err != nil {
		return fmt.Errorf("error setting %sphase1_encryption_algorithms: %w", prefix, err)
	}
	s = nil

	for _, v := range apiObject.Phase1IntegrityAlgorithms {
		s = append(s, v.Value)
	}
	if err := d.Set(prefix+"phase1_integrity_algorithms", aws.StringValueSlice(s)); err != nil {
		return fmt.Errorf("error setting %sphase1_integrity_algorithms: %w", prefix, err)
	}
	s = nil

	for _, v := range apiObject.Phase2DHGroupNumbers {
		i = append(i, v.Value)
	}
	if err := d.Set(prefix+"phase2_dh_group_numbers", aws.Int64ValueSlice(i)); err != nil {
		return fmt.Errorf("error setting %sphase2_dh_group_numbers: %w", prefix, err)
	}

	for _, v := range apiObject.Phase2EncryptionAlgorithms {
		s = append(s, v.Value)
	}
	if err := d.Set(prefix+"phase2_encryption_algorithms", aws.StringValueSlice(s)); err != nil {
		return fmt.Errorf("error setting %sphase2_encryption_algorithms: %w", prefix, err)
	}
	s = nil

	for _, v := range apiObject.Phase2IntegrityAlgorithms {
		s = append(s, v.Value)
	}
	if err := d.Set(prefix+"phase2_integrity_algorithms", aws.StringValueSlice(s)); err != nil {
		return fmt.Errorf("error setting %sphase2_integrity_algorithms: %w", prefix, err)
	}

	d.Set(prefix+"phase1_lifetime_seconds", apiObject.Phase1LifetimeSeconds)
	d.Set(prefix+"phase2_lifetime_seconds", apiObject.Phase2LifetimeSeconds)
	d.Set(prefix+"replay_window_size", apiObject.ReplayWindowSize)

	d.Set(prefix+"inside_cidr", apiObject.TunnelInsideCidr)

	d.Set(prefix+"address", apiObject.OutsideIpAddress)
	d.Set(prefix+"preshared_key", apiObject.PreSharedKey)

	d.Set(prefix+"bgp_asn", config.BGPASN)
	d.Set(prefix+"bgp_holdtime", config.BGPHoldTime)
	d.Set(prefix+"cgw_inside_address", config.CgwInsideAddress)
	d.Set(prefix+"vgw_inside_address", config.VgwInsideAddress)

	return nil
}

func flattenVgwTelemetry(apiObject *ec2.VgwTelemetry) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.AcceptedRouteCount; v != nil {
		tfMap["accepted_route_count"] = aws.Int64Value(v)
	}

	if v := apiObject.LastStatusChange; v != nil {
		tfMap["last_status_change"] = aws.TimeValue(v).Format(time.RFC3339)
	}

	if v := apiObject.OutsideIpAddress; v != nil {
		tfMap["outside_ip_address"] = aws.StringValue(v)
	}

	if v := apiObject.Status; v != nil {
		tfMap["status"] = aws.StringValue(v)
	}

	if v := apiObject.StatusMessage; v != nil {
		tfMap["status_message"] = aws.StringValue(v)
	}

	return tfMap
}

func flattenVgwTelemetries(apiObjects []*ec2.VgwTelemetry) []interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []interface{}

	for _, apiObject := range apiObjects {
		if apiObject == nil {
			continue
		}

		tfList = append(tfList, flattenVgwTelemetry(apiObject))
	}

	return tfList
}

// Customer gateway configuration of VPN connection in native XML format.
type XmlVpnConnectionConfig struct {
	Tunnels []XmlIpsecTunnel `xml:"ipsec_tunnel"`
}

type XmlIpsecTunnel struct {
	BGPASN           string `xml:"vpn_gateway>bgp>asn"`
	BGPHoldTime      int    `xml:"vpn_gateway>bgp>hold_time"`
	CgwInsideAddress string `xml:"customer_gateway>tunnel_inside_address>ip_address"`
	OutsideAddress   string `xml:"vpn_gateway>tunnel_outside_address>ip_address"`
	VgwInsideAddress string `xml:"vpn_gateway>tunnel_inside_address>ip_address"`
}

// Tunnels are matched by outside IP address, then by order.
func (c XmlVpnConnectionConfig) tunnel(index int, tunnelOption *ec2.TunnelOption) *XmlIpsecTunnel {
	if tunnelOption != nil {
		if address := aws.StringValue(tunnelOption.OutsideIpAddress); address != "" {
			for i := range c.Tunnels {
				if c.Tunnels[i].OutsideAddress == address {
					return &c.Tunnels[i]
				}
			}
		}
	}

	if index < len(c.Tunnels) {
		return &c.Tunnels[index]
	}

	return nil
}

func validateVpnConnectionTunnelPreSharedKey() schema.SchemaValidateFunc {
	return validation.All(
		validation.StringLenBetween(8, 64),
		validation.StringDoesNotMatch(regexp.MustCompile(`^0`), "cannot start with zero character"),
		validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z_.]+$`), "can only contain alphanumeric, period and underscore characters"),
	)
}

// Address inside /30 block is allowed, it is assigned to VPN gateway side of tunnel.
func validateVpnConnectionTunnelInsideCIDR(v interface{}, k string) (ws []string, errors []error) {
	const allowedCIDR = "169.254.252.0/22"

	value := v.(string)
	ip, ipNet, err := net.ParseCIDR(value)

	if err != nil {
		errors = append(errors, fmt.Errorf("%q must be valid CIDR block: %w", k, err))
		return
	}

	if ones, bits := ipNet.Mask.Size(); ones != 30 || bits != 32 {
		errors = append(errors, fmt.Errorf("%q must be /30 IPv4 CIDR block, got %q", k, value))
		return
	}

	if _, allowed, _ := net.ParseCIDR(allowedCIDR); !allowed.Contains(ip) {
		errors = append(errors, fmt.Errorf("%q must be within %s, got %q", k, allowedCIDR, value))
	}

	return
}
