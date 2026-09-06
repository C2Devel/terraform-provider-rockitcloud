package ec2

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceEIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceEIPCreate,
		Read:   resourceEIPRead,
		Update: resourceEIPUpdate,
		Delete: resourceEIPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		CustomizeDiff: verify.SetTagsDiff,

		Timeouts: &schema.ResourceTimeout{
			Read:   schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"allocation_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"associate_with_private_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"association_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"carrier_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"customer_owned_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"customer_owned_ipv4_pool": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network_border_group": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"network_interface": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_dns": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_dns": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ipv4_pool": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
			"vpc": {
				Type:       schema.TypeBool,
				Optional:   true,
				Default:    true,
				Deprecated: "The specified value is ignored.",
			},
		},
	}
}

func resourceEIPCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	// EC2-Classic is not supported; all Elastic IPs are VPC-domain.
	allocOpts := &ec2.AllocateAddressInput{
		Domain:            aws.String(ec2.DomainTypeVpc),
		TagSpecifications: ec2TagSpecificationsFromKeyValueTags(tags, ec2.ResourceTypeElasticIp),
	}

	if v, ok := d.GetOk("address"); ok {
		allocOpts.Address = aws.String(v.(string))
	}

	if v, ok := d.GetOk("public_ipv4_pool"); ok {
		allocOpts.PublicIpv4Pool = aws.String(v.(string))
	}

	if v, ok := d.GetOk("customer_owned_ipv4_pool"); ok {
		allocOpts.CustomerOwnedIpv4Pool = aws.String(v.(string))
	}

	if v, ok := d.GetOk("network_border_group"); ok {
		allocOpts.NetworkBorderGroup = aws.String(v.(string))
	}

	log.Printf("[DEBUG] EIP create configuration: %#v", allocOpts)
	allocResp, err := conn.AllocateAddress(allocOpts)
	if err != nil {
		return fmt.Errorf("Error creating EIP: %s", err)
	}

	d.SetId(aws.StringValue(allocResp.AllocationId))

	log.Printf("[INFO] EIP ID: %s (domain: %v)", d.Id(), aws.StringValue(allocResp.Domain))

	return resourceEIPUpdate(d, meta)
}

func resourceEIPRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	id := d.Id()

	// The ID is an allocation ID, unless the EIP was imported by its public IP.
	req := &ec2.DescribeAddressesInput{}
	if strings.Contains(id, "eipalloc") {
		req.AllocationIds = []*string{aws.String(id)}
	} else {
		req.PublicIps = []*string{aws.String(id)}
	}

	log.Printf("[DEBUG] EIP describe configuration: %s", req)

	var err error
	var describeAddresses *ec2.DescribeAddressesOutput

	if d.IsNewResource() {
		err := resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			describeAddresses, err = conn.DescribeAddresses(req)
			if err != nil {
				if tfawserr.ErrCodeEquals(err, ErrCodeInvalidAllocationIDNotFound, ErrCodeInvalidAddressNotFound) {
					return resource.RetryableError(err)
				}

				return resource.NonRetryableError(err)
			}
			return nil
		})
		if tfresource.TimedOut(err) {
			describeAddresses, err = conn.DescribeAddresses(req)
		}
		if err != nil {
			return fmt.Errorf("Error retrieving EIP: %s", err)
		}
	} else {
		describeAddresses, err = conn.DescribeAddresses(req)
		if err != nil {
			if tfawserr.ErrCodeEquals(err, ErrCodeInvalidAllocationIDNotFound, ErrCodeInvalidAddressNotFound) {
				log.Printf("[WARN] EIP not found, removing from state: %s", req)
				d.SetId("")
				return nil
			}
			return err
		}
	}

	var address *ec2.Address

	// In the case that AWS returns more EIPs than we intend it to, we loop
	// over the returned addresses to see if it's in the list of results
	for _, addr := range describeAddresses.Addresses {
		if aws.StringValue(addr.AllocationId) == id || aws.StringValue(addr.PublicIp) == id {
			address = addr
			break
		}
	}

	if address == nil {
		log.Printf("[WARN] EIP %q not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("association_id", address.AssociationId)
	if address.InstanceId != nil {
		d.Set("instance", address.InstanceId)
	} else {
		d.Set("instance", "")
	}
	if address.NetworkInterfaceId != nil {
		d.Set("network_interface", address.NetworkInterfaceId)
	} else {
		d.Set("network_interface", "")
	}

	d.Set("private_ip", address.PrivateIpAddress)
	d.Set("public_ip", address.PublicIp)

	privateDNS, publicDNS, err := eipDNSNames(conn, address.NetworkInterfaceId)
	if err != nil {
		return fmt.Errorf("error reading EC2 Network Interface (%s) for EIP (%s): %w", aws.StringValue(address.NetworkInterfaceId), d.Id(), err)
	}
	d.Set("private_dns", privateDNS)
	d.Set("public_dns", publicDNS)

	d.Set("allocation_id", address.AllocationId)
	d.Set("carrier_ip", address.CarrierIp)
	d.Set("customer_owned_ipv4_pool", address.CustomerOwnedIpv4Pool)
	d.Set("customer_owned_ip", address.CustomerOwnedIp)
	d.Set("network_border_group", address.NetworkBorderGroup)
	d.Set("public_ipv4_pool", address.PublicIpv4Pool)

	d.Set("vpc", aws.StringValue(address.Domain) == ec2.DomainTypeVpc)
	d.Set("domain", address.Domain)

	// The EIP can be imported by its IP address; store the allocation ID instead.
	if net.ParseIP(id) != nil {
		log.Printf("[DEBUG] Re-assigning EIP ID (%s) to its Allocation ID (%s)", d.Id(), aws.StringValue(address.AllocationId))
		d.SetId(aws.StringValue(address.AllocationId))
	}

	tags := KeyValueTags(address.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}

	return nil
}

func resourceEIPUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	// If we are updating an EIP that is not newly created, and we are attached to
	// an instance or interface, detach first.
	disassociate := false
	if !d.IsNewResource() && d.HasChanges("instance", "network_interface", "associate_with_private_ip") {
		oldInstance, _ := d.GetChange("instance")

		if oldInstance.(string) != "" || d.Get("association_id").(string) != "" {
			disassociate = true
		}
	}
	if disassociate {
		if err := disassociateEip(d, meta); err != nil {
			return err
		}
	}

	v_instance, ok_instance := d.GetOk("instance")
	v_interface, ok_interface := d.GetOk("network_interface")

	associateByInstance := d.HasChange("instance") && ok_instance
	// "network_interface" is Optional+Computed, so it can hold a stale ENI ID
	// from a prior instance-based association; don't treat that as a trigger.
	associateByInterface := !associateByInstance &&
		d.HasChanges("network_interface", "associate_with_private_ip") && ok_interface
	associate := associateByInstance || associateByInterface
	if associate {
		instanceId := v_instance.(string)
		networkInterfaceId := v_interface.(string)

		var privateIpAddress *string
		if v := d.Get("associate_with_private_ip").(string); v != "" {
			privateIpAddress = aws.String(v)
		}
		assocOpts := &ec2.AssociateAddressInput{
			AllocationId:     aws.String(d.Id()),
			PrivateIpAddress: privateIpAddress,
		}
		// Specifying both InstanceId and NetworkInterfaceId is rejected.
		if associateByInstance {
			assocOpts.InstanceId = aws.String(instanceId)
		} else {
			assocOpts.NetworkInterfaceId = aws.String(networkInterfaceId)
		}

		log.Printf("[DEBUG] EIP associate configuration: %s", assocOpts)

		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			_, err := conn.AssociateAddress(assocOpts)
			if err != nil {
				if tfawserr.ErrCodeEquals(err, ErrCodeInvalidAllocationIDNotFound) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if tfresource.TimedOut(err) {
			_, err = conn.AssociateAddress(assocOpts)
		}
		if err != nil {
			// Prevent saving instance if association failed
			// e.g. missing internet gateway in VPC
			d.Set("instance", "")
			d.Set("network_interface", "")
			return fmt.Errorf("Failure associating EIP: %s", err)
		}
	}

	if d.HasChange("tags_all") && !d.IsNewResource() {
		o, n := d.GetChange("tags_all")
		if err := UpdateTags(conn, d.Id(), o, n); err != nil {
			return fmt.Errorf("error updating EIP (%s) tags: %s", d.Id(), err)
		}
	}

	return resourceEIPRead(d, meta)
}

func resourceEIPDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	// If we are attached to an instance or interface, detach first.
	if d.Get("instance").(string) != "" || d.Get("association_id").(string) != "" {
		if err := disassociateEip(d, meta); err != nil {
			return err
		}
	}

	log.Printf("[DEBUG] EIP release (destroy) address allocation: %v", d.Id())
	input := &ec2.ReleaseAddressInput{
		AllocationId: aws.String(d.Id()),
	}
	if v := d.Get("network_border_group").(string); v != "" {
		input.NetworkBorderGroup = aws.String(v)
	}

	// AuthFailure is the only error worth waiting out here
	// Every other refusal is permanent, so retrying it just hide reason timeout
	_, err := tfresource.RetryWhenAWSErrCodeEquals(d.Timeout(schema.TimeoutDelete), func() (interface{}, error) {
		return conn.ReleaseAddress(input)
	}, ErrCodeAuthFailure)

	if tfawserr.ErrCodeEquals(err, ErrCodeInvalidAllocationIDNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("error releasing EC2 EIP (%s): %w", d.Id(), err)
	}
	return nil
}

func disassociateEip(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	log.Printf("[DEBUG] Disassociating EIP: %s", d.Id())

	associationID := d.Get("association_id").(string)
	if associationID == "" {
		return nil
	}

	_, err := conn.DisassociateAddress(&ec2.DisassociateAddressInput{
		AssociationId: aws.String(associationID),
	})

	// First check if the association ID is not found. If this
	// is the case, then it was already disassociated somehow,
	// and that is okay. The most common reason for this is that
	// the instance or ENI it was attached to was destroyed.
	if tfawserr.ErrCodeEquals(err, ErrCodeInvalidAssociationIDNotFound) {
		err = nil
	}

	return err
}

// eipDNSNames returns the DNS names of the network interface the EIP is attached
// to. They are read from the interface rather than from its instance, since an
// instance may own several interfaces and only this one matches the EIP.
func eipDNSNames(conn *ec2.EC2, networkInterfaceID *string) (*string, *string, error) {
	if aws.StringValue(networkInterfaceID) == "" {
		return nil, nil, nil
	}

	eni, err := FindNetworkInterfaceByID(conn, aws.StringValue(networkInterfaceID))
	if err != nil {
		if tfresource.NotFound(err) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	var publicDNS *string
	if eni.Association != nil {
		publicDNS = eni.Association.PublicDnsName
	}

	return eni.PrivateDnsName, publicDNS, nil
}

func ConvertIPToDashIP(ip string) string {
	return strings.Replace(ip, ".", "-", -1)
}
