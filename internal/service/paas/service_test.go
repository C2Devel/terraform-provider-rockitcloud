package paas_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/paas"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfpaas "github.com/hashicorp/terraform-provider-aws/internal/service/paas"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func init() {
	acctest.RegisterServiceErrorCheckFunc(paas.EndpointsID, testAccErrorCheckSkip)
}

func testAccErrorCheckSkip(t *testing.T) resource.ErrorCheckFunc {
	return acctest.ErrorCheckSkipMessagesContaining(
		t,
		"VmTypeNotFound",
		"UnknownVolumeType",
		"Parameter \"version\" has invalid value",
	)
}

func TestAccPaaSServiceElasticSearch_basic(t *testing.T) {
	resourceName := "aws_paas_service.test"
	vpcResourceName := "aws_vpc.test"

	randServiceName := fmt.Sprintf("terraform-test-%s", sdkacctest.RandString(5))
	randKeyName := fmt.Sprintf("terraform-test-%s", sdkacctest.RandString(5))
	publicKey, _, err := sdkacctest.RandSSHKeyPair(acctest.DefaultEmailAddress)
	if err != nil {
		t.Fatalf("error generating random SSH key: %s", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, paas.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceConfig(randServiceName, randKeyName, publicKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "arbitrator_required", "false"),
					resource.TestCheckResourceAttr(resourceName, "auto_created_security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "backup_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "data_volume.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "data_volume.*", map[string]string{
						"iops": "0",
						"size": "32",
						"type": "gp2",
					}),
					resource.TestCheckResourceAttr(resourceName, "delete_interfaces_on_destroy", "true"),
					resource.TestCheckResourceAttr(resourceName, "elasticsearch.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "elasticsearch.*", map[string]string{
						"allow_anonymous": "",
						"anonymous_role":  "",
						"class":           "search",
						"kibana":          "false",
						"logging.#":       "0",
						"monitoring.#":    "0",
						"options.%":       "0",
						"password":        "",
						"version":         "8.2.2",
					}),
					resource.TestCheckResourceAttr(resourceName, "endpoints.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "error_code", ""),
					resource.TestCheckResourceAttr(resourceName, "error_description", ""),
					resource.TestCheckResourceAttr(resourceName, "high_availability", "false"),
					resource.TestCheckResourceAttr(resourceName, "instances.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "instances.0.endpoints.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "instances.0.index"),
					resource.TestCheckResourceAttrSet(resourceName, "instances.0.instance_id"),
					resource.TestCheckResourceAttrSet(resourceName, "instances.0.interface_id"),
					resource.TestCheckResourceAttrSet(resourceName, "instances.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "instances.0.private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "instances.0.role"),
					resource.TestCheckResourceAttrSet(resourceName, "instances.0.status"),
					resource.TestCheckResourceAttr(resourceName, "instance_type", "c5.large"),
					resource.TestCheckResourceAttr(resourceName, "name", randServiceName),
					resource.TestCheckResourceAttr(resourceName, "network_interface_ids.#", "0"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "root_volume.*", map[string]string{
						"iops": "0",
						"size": "32",
						"type": "gp2",
					}),
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "service_class", "search"),
					resource.TestCheckResourceAttr(resourceName, "service_type", "elasticsearch"),
					resource.TestCheckResourceAttr(resourceName, "ssh_key_name", randKeyName),
					resource.TestCheckResourceAttr(resourceName, "status", tfpaas.ServiceStatusReady),
					resource.TestCheckResourceAttr(resourceName, "subnet_ids.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "supported_features.#"),
					resource.TestCheckResourceAttrSet(resourceName, "total_cpu_count"),
					resource.TestCheckResourceAttrSet(resourceName, "total_memory"),
					resource.TestCheckNoResourceAttr(resourceName, "user_data"),
					resource.TestCheckNoResourceAttr(resourceName, "user_data_content_type"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", vpcResourceName, "id"),
					// TODO: parametrize test with different service types
					resource.TestCheckResourceAttr(resourceName, "memcached.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "mongodb.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "mssql.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "pgsql.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "rabbitmq.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "redis.#", "0"),
				),
			},
		},
	})
}

func testAccCheckServiceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("PaaS Service %s is not found in state", resourceName)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).PaaSConn
		_, err := tfpaas.FindServiceByID(conn, rs.Primary.ID)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckServiceDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).PaaSConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_paas_service" {
			continue
		}

		service, err := tfpaas.FindServiceByID(conn, rs.Primary.ID)
		serviceDeleted := service != nil && aws.StringValue(service.Status) == tfpaas.ServiceStatusDeleted

		if tfresource.NotFound(err) || serviceDeleted {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("PaaS Service (%s) still exists", rs.Primary.ID)
	}

	return nil
}

func testAccServiceConfig(serviceName, keyName, publicKey string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  cidr_block = "10.0.1.0/24"
  vpc_id     = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_key_pair" "test" {
  key_name   = %[1]q
  public_key = %[2]q
}

resource "aws_paas_service" "test" {
  name          = %[3]q 
  instance_type = "c5.large"

  root_volume {
    type = "gp2"
  }

  data_volume {
    type = "gp2"
  }

  delete_interfaces_on_destroy = true
  security_group_ids           = [aws_vpc.test.default_security_group_id]
  subnet_ids                   = [aws_subnet.test.id]

  ssh_key_name = aws_key_pair.test.key_name

  elasticsearch {
    version = "8.12"
  }
}
`, keyName, publicKey, serviceName)
}
