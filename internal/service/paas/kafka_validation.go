package paas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/service/paas/services"
)

// kafkaMinDataVolumeSize is the minimum Kafka data volume size (GiB) accepted by K2 Cloud.
// The platform lowered it from 64 to 8 GiB on 2026-05-27 (K2 Cloud changelog, May 2026, "Fixes").
const kafkaMinDataVolumeSize = 8

type kafkaServiceConfig struct {
	isKafka                   bool
	highAvailability          bool
	hasCoordinatorBlock       bool
	hasCoordinatorRole        bool
	hasDataVolume             bool
	dataVolumeSize            int
	hasCoordinatorDataVolume  bool
	coordinatorDataVolumeSize int
}

// validateKafkaServiceConfiguration enforces Kafka CreateService rules from K2 Cloud:
// https://docs.k2.cloud/en/api/paas/actions/CreateService.html
//
// - additionalRoles: only "coordinator"; cannot be used together with coordinator (schema ConflictsWith).
// - coordinator: required when serviceType is kafka, highAvailability is true, and additionalRoles is not set.
//
// Non-HA Kafka on K2 Cloud uses additional_roles = ["coordinator"] to combine broker and coordinator;
// a dedicated coordinator block without additional_roles is rejected by the platform for single-node clusters.
func validateKafkaServiceConfiguration(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	// ClickHouse also uses additional_roles and coordinator, but with its own rules.
	if _, ok := d.GetOk(services.ClickHouse.ServiceType()); ok {
		return validateClickHouseServiceConfiguration(d)
	}

	config := kafkaServiceConfig{}

	if _, ok := d.GetOk(services.Kafka.ServiceType()); ok {
		config.isKafka = true
	}

	// Default schema value is false; GetOk is false when the attribute is false or unset.
	if v, ok := d.GetOk("high_availability"); ok {
		config.highAvailability = v.(bool)
	}

	if v, ok := d.GetOk("data_volume"); ok {
		if list, ok := v.([]interface{}); ok && len(list) > 0 && list[0] != nil {
			config.hasDataVolume = true
			config.dataVolumeSize = list[0].(map[string]interface{})["size"].(int)
		}
	}

	if v, ok := d.GetOk("coordinator"); ok {
		if list, ok := v.([]interface{}); ok && len(list) > 0 && list[0] != nil {
			config.hasCoordinatorBlock = true

			coordinatorMap := list[0].(map[string]interface{})
			if size, ok := coordinatorMap["data_volume_size"].(int); ok && size > 0 {
				config.hasCoordinatorDataVolume = true
				config.coordinatorDataVolumeSize = size
			}
		}
	}

	if v, ok := d.GetOk("additional_roles"); ok {
		for _, r := range v.([]interface{}) {
			if s, ok := r.(string); ok && s == "coordinator" {
				config.hasCoordinatorRole = true
				break
			}
		}
	}

	return validateKafkaServiceConfig(config)
}

func validateKafkaServiceConfig(config kafkaServiceConfig) error {
	if !config.isKafka {
		if config.hasCoordinatorRole {
			return fmt.Errorf("additional_roles is supported only for kafka services")
		}
		if config.hasCoordinatorBlock {
			return fmt.Errorf("coordinator is supported only for kafka services")
		}
		return nil
	}

	if config.hasDataVolume && config.dataVolumeSize < kafkaMinDataVolumeSize {
		return fmt.Errorf("kafka data_volume.size must be at least %d GiB", kafkaMinDataVolumeSize)
	}

	if config.hasCoordinatorBlock {
		if !config.hasCoordinatorDataVolume {
			return fmt.Errorf("kafka coordinator.data_volume_size must be set to at least %d GiB", kafkaMinDataVolumeSize)
		}
		if config.coordinatorDataVolumeSize < kafkaMinDataVolumeSize {
			return fmt.Errorf("kafka coordinator.data_volume_size must be at least %d GiB", kafkaMinDataVolumeSize)
		}
	}

	if config.highAvailability {
		if config.hasCoordinatorRole || config.hasCoordinatorBlock {
			return nil
		}
		return fmt.Errorf(
			"kafka with high_availability = true requires a \"coordinator\" block or additional_roles = [\"coordinator\"] " +
				"(CreateService: coordinator is required when highAvailability is true and additionalRoles is not set; " +
				"see https://docs.k2.cloud/en/api/paas/actions/CreateService.html)",
		)
	}

	if config.hasCoordinatorRole {
		return nil
	}
	if config.hasCoordinatorBlock {
		return fmt.Errorf(
			"kafka with high_availability = false requires additional_roles = [\"coordinator\"]; " +
				"a dedicated \"coordinator\" block is not supported for non-HA Kafka on K2 Cloud",
		)
	}
	return fmt.Errorf(
		"kafka with high_availability = false requires additional_roles = [\"coordinator\"]",
	)
}
