package ec2

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func ResourceNATGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceNATGatewayCreate,
		Read:   resourceNATGatewayRead,
		Update: resourceNATGatewayUpdate,
		Delete: resourceNATGatewayDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"auto_provision_zones": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      ec2.AvailabilityModeRegional,
				ValidateFunc: validation.StringInSlice([]string{ec2.AvailabilityModeRegional}, false),
			},
			"availability_zone_addresses": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allocation_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"connectivity_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      ec2.ConnectivityTypePublic,
				ValidateFunc: validation.StringInSlice([]string{ec2.ConnectivityTypePublic}, false),
			},
			"nat_gateway_addresses": natGatewayAddressesSchema(),
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: customdiff.All(
			verify.SetTagsDiff,
			// Addresses of NAT gateway are changed by EIP (dis)association so they are recomputed
			customdiff.ComputedIf("nat_gateway_addresses", func(_ context.Context, diff *schema.ResourceDiff, _ interface{}) bool {
				return diff.HasChange("availability_zone_addresses")
			}),
			// EIPs of NAT gateway created with auto allocated addresses cannot be managed manually
			func(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
				if diff.Id() == "" {
					return nil
				}

				if diff.Get("auto_provision_zones").(string) == ec2.AutoProvisionZonesStateEnabled && diff.HasChange("availability_zone_addresses") {
					return diff.ForceNew("availability_zone_addresses")
				}

				return nil
			},
		),
	}
}

func natGatewayAddressesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"allocation_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"association_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"availability_zone": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"is_primary": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"private_ip": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"public_ip": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"status": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func resourceNATGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	input := &ec2.CreateNatGatewayInput{
		AvailabilityMode:  aws.String(d.Get("availability_mode").(string)),
		ConnectivityType:  aws.String(d.Get("connectivity_type").(string)),
		TagSpecifications: ec2TagSpecificationsFromKeyValueTags(tags, ec2.ResourceTypeNatgateway),
		VpcId:             aws.String(d.Get("vpc_id").(string)),
	}

	if v, ok := d.GetOk("availability_zone_addresses"); ok && v.(*schema.Set).Len() > 0 {
		tfList := v.(*schema.Set).List()

		if err := validateAvailabilityZoneAddresses(tfList); err != nil {
			return fmt.Errorf("error creating EC2 NAT Gateway: %w", err)
		}

		input.AvailabilityZoneAddresses = expandAvailabilityZoneAddresses(tfList)
	}

	log.Printf("[DEBUG] Creating EC2 NAT Gateway: %s", input)
	output, err := conn.CreateNatGateway(input)

	if err != nil {
		return fmt.Errorf("error creating EC2 NAT Gateway: %w", err)
	}

	d.SetId(aws.StringValue(output.NatGateway.NatGatewayId))

	if _, err := WaitNATGatewayCreated(conn, d.Id()); err != nil {
		return fmt.Errorf("error waiting for EC2 NAT Gateway (%s) create: %w", d.Id(), err)
	}

	return resourceNATGatewayRead(d, meta)
}

func resourceNATGatewayRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	ng, err := FindNATGatewayByID(conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] EC2 NAT Gateway (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading EC2 NAT Gateway (%s): %w", d.Id(), err)
	}

	d.Set("auto_provision_zones", ng.AutoProvisionZones)
	d.Set("availability_mode", ng.AvailabilityMode)
	d.Set("connectivity_type", ng.ConnectivityType)
	d.Set("vpc_id", ng.VpcId)

	if aws.StringValue(ng.AutoProvisionZones) == ec2.AutoProvisionZonesStateEnabled {
		d.Set("availability_zone_addresses", nil)
	} else {
		if err := d.Set("availability_zone_addresses", flattenAvailabilityZoneAddresses(ng.NatGatewayAddresses)); err != nil {
			return fmt.Errorf("error setting availability_zone_addresses: %w", err)
		}
	}

	if err := d.Set("nat_gateway_addresses", flattenNATGatewayAddresses(ng.NatGatewayAddresses)); err != nil {
		return fmt.Errorf("error setting nat_gateway_addresses: %w", err)
	}

	tags := KeyValueTags(ng.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}

	return nil
}

func resourceNATGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	if d.HasChange("availability_zone_addresses") {
		_, n := d.GetChange("availability_zone_addresses")

		if err := validateAvailabilityZoneAddresses(n.(*schema.Set).List()); err != nil {
			return fmt.Errorf("error updating EC2 NAT Gateway (%s): %w", d.Id(), err)
		}

		// only one EIP is allowed per availability zone
		newByAZ := make(map[string]string)
		for _, tfMapRaw := range n.(*schema.Set).List() {
			tfMap := tfMapRaw.(map[string]interface{})
			newByAZ[tfMap["availability_zone"].(string)] = tfMap["allocation_id"].(string)
		}

		ng, err := FindNATGatewayByID(conn, d.Id())

		if err != nil {
			return fmt.Errorf("error reading EC2 NAT Gateway (%s): %w", d.Id(), err)
		}

		liveByAZ := make(map[string]string)
		for _, address := range ng.NatGatewayAddresses {
			liveByAZ[aws.StringValue(address.AvailabilityZone)] = aws.StringValue(address.AllocationId)
		}

		for _, address := range ng.NatGatewayAddresses {
			az := aws.StringValue(address.AvailabilityZone)

			if allocationID, ok := newByAZ[az]; ok && allocationID == aws.StringValue(address.AllocationId) {
				continue
			}

			input := &ec2.DisassociateNatGatewayAddressInput{
				NatGatewayId:   aws.String(d.Id()),
				AssociationIds: aws.StringSlice([]string{aws.StringValue(address.AssociationId)}),
			}

			log.Printf("[DEBUG] Disassociating EC2 NAT Gateway address: %s", input)
			if _, err := conn.DisassociateNatGatewayAddress(input); err != nil {
				return fmt.Errorf("error disassociating EC2 NAT Gateway (%s) address in %s: %w", d.Id(), az, err)
			}
		}

		for az, allocationID := range newByAZ {
			if liveByAZ[az] == allocationID {
				continue
			}

			input := &ec2.AssociateNatGatewayAddressInput{
				NatGatewayId:     aws.String(d.Id()),
				AvailabilityZone: aws.String(az),
				AllocationIds:    aws.StringSlice([]string{allocationID}),
			}

			log.Printf("[DEBUG] Associating EC2 NAT Gateway address: %s", input)
			if _, err := conn.AssociateNatGatewayAddress(input); err != nil {
				return fmt.Errorf("error associating EC2 NAT Gateway (%s) address in %s: %w", d.Id(), az, err)
			}
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, d.Id(), o, n); err != nil {
			return fmt.Errorf("error updating EC2 NAT Gateway (%s) tags: %w", d.Id(), err)
		}
	}

	return resourceNATGatewayRead(d, meta)
}

func resourceNATGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	log.Printf("[INFO] Deleting EC2 NAT Gateway: %s", d.Id())
	_, err := conn.DeleteNatGateway(&ec2.DeleteNatGatewayInput{
		NatGatewayId: aws.String(d.Id()),
	})

	if err != nil && !tfawserr.ErrCodeEquals(err, ErrCodeInvalidNatGatewayIDNotFound) {
		return fmt.Errorf("error deleting EC2 NAT Gateway (%s): %w", d.Id(), err)
	}

	// SDK retry may get NotFound while gateway still deleting
	if _, err := WaitNATGatewayDeleted(conn, d.Id()); err != nil {
		return fmt.Errorf("error waiting for EC2 NAT Gateway (%s) delete: %w", d.Id(), err)
	}

	return nil
}

func validateAvailabilityZoneAddresses(tfList []interface{}) error {
	seenAZs := make(map[string]struct{}, len(tfList))
	azsByAllocationID := make(map[string]string, len(tfList))

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})

		if !ok {
			continue
		}

		availabilityZone, allocationID := tfMap["availability_zone"].(string), tfMap["allocation_id"].(string)

		if _, ok := seenAZs[availabilityZone]; ok {
			return fmt.Errorf("availability zone %s is specified more than once in availability_zone_addresses: only one Elastic IP per availability zone is allowed", availabilityZone)
		}

		seenAZs[availabilityZone] = struct{}{}

		if otherAvailabilityZone, ok := azsByAllocationID[allocationID]; ok {
			return fmt.Errorf("allocation %s is specified in availability_zone_addresses for both %s and %s: an Elastic IP can serve only one availability zone", allocationID, otherAvailabilityZone, availabilityZone)
		}

		azsByAllocationID[allocationID] = availabilityZone
	}

	return nil
}

func expandAvailabilityZoneAddresses(tfList []interface{}) []*ec2.AvailabilityZoneAddress {
	apiObjects := make([]*ec2.AvailabilityZoneAddress, 0, len(tfList))

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})

		if !ok {
			continue
		}

		apiObjects = append(apiObjects, &ec2.AvailabilityZoneAddress{
			AvailabilityZone: aws.String(tfMap["availability_zone"].(string)),
			AllocationIds:    aws.StringSlice([]string{tfMap["allocation_id"].(string)}),
		})
	}

	return apiObjects
}

func flattenAvailabilityZoneAddresses(apiObjects []*ec2.NatGatewayAddress) []interface{} {
	tfList := make([]interface{}, 0, len(apiObjects))

	for _, apiObject := range apiObjects {
		tfList = append(tfList, map[string]interface{}{
			"allocation_id":     aws.StringValue(apiObject.AllocationId),
			"availability_zone": aws.StringValue(apiObject.AvailabilityZone),
		})
	}

	return tfList
}

func flattenNATGatewayAddresses(apiObjects []*ec2.NatGatewayAddress) []interface{} {
	tfList := make([]interface{}, 0, len(apiObjects))

	for _, apiObject := range apiObjects {
		tfList = append(tfList, map[string]interface{}{
			"allocation_id":     aws.StringValue(apiObject.AllocationId),
			"association_id":    aws.StringValue(apiObject.AssociationId),
			"availability_zone": aws.StringValue(apiObject.AvailabilityZone),
			"is_primary":        aws.BoolValue(apiObject.IsPrimary),
			"private_ip":        aws.StringValue(apiObject.PrivateIp),
			"public_ip":         aws.StringValue(apiObject.PublicIp),
			"status":            aws.StringValue(apiObject.Status),
		})
	}

	return tfList
}
