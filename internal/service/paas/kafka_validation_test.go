package paas

import "testing"

func TestValidateKafkaServiceConfig(t *testing.T) {
	testCases := []struct {
		name      string
		config    kafkaServiceConfig
		expectErr bool
	}{
		{
			name:   "non-kafka without kafka-only fields is allowed",
			config: kafkaServiceConfig{},
		},
		{
			name: "non-kafka rejects additional roles",
			config: kafkaServiceConfig{
				hasCoordinatorRole: true,
			},
			expectErr: true,
		},
		{
			name: "non-kafka rejects coordinator block",
			config: kafkaServiceConfig{
				hasCoordinatorBlock: true,
			},
			expectErr: true,
		},
		{
			name: "kafka rejects too-small service data volume",
			config: kafkaServiceConfig{
				isKafka:        true,
				hasDataVolume:  true,
				dataVolumeSize: 7,
			},
			expectErr: true,
		},
		{
			name: "kafka rejects missing coordinator data volume size",
			config: kafkaServiceConfig{
				isKafka:             true,
				highAvailability:    true,
				hasCoordinatorBlock: true,
			},
			expectErr: true,
		},
		{
			name: "kafka rejects too-small coordinator data volume",
			config: kafkaServiceConfig{
				isKafka:                   true,
				highAvailability:          true,
				hasCoordinatorBlock:       true,
				hasCoordinatorDataVolume:  true,
				coordinatorDataVolumeSize: 7,
			},
			expectErr: true,
		},
		{
			name: "ha kafka accepts dedicated coordinator",
			config: kafkaServiceConfig{
				isKafka:                   true,
				highAvailability:          true,
				hasDataVolume:             true,
				dataVolumeSize:            64,
				hasCoordinatorBlock:       true,
				hasCoordinatorDataVolume:  true,
				coordinatorDataVolumeSize: 64,
			},
		},
		{
			name: "kafka accepts minimum service data volume",
			config: kafkaServiceConfig{
				isKafka:            true,
				hasDataVolume:      true,
				dataVolumeSize:     8,
				hasCoordinatorRole: true,
			},
		},
		{
			name: "ha kafka accepts minimum coordinator data volume",
			config: kafkaServiceConfig{
				isKafka:                   true,
				highAvailability:          true,
				hasDataVolume:             true,
				dataVolumeSize:            8,
				hasCoordinatorBlock:       true,
				hasCoordinatorDataVolume:  true,
				coordinatorDataVolumeSize: 8,
			},
		},
		{
			name: "ha kafka accepts combined coordinator role",
			config: kafkaServiceConfig{
				isKafka:            true,
				highAvailability:   true,
				hasDataVolume:      true,
				dataVolumeSize:     64,
				hasCoordinatorRole: true,
			},
		},
		{
			name: "ha kafka requires coordinator placement",
			config: kafkaServiceConfig{
				isKafka:          true,
				highAvailability: true,
				hasDataVolume:    true,
				dataVolumeSize:   64,
			},
			expectErr: true,
		},
		{
			name: "non-ha kafka requires combined coordinator role",
			config: kafkaServiceConfig{
				isKafka:        true,
				hasDataVolume:  true,
				dataVolumeSize: 64,
			},
			expectErr: true,
		},
		{
			name: "non-ha kafka accepts combined coordinator role",
			config: kafkaServiceConfig{
				isKafka:            true,
				hasDataVolume:      true,
				dataVolumeSize:     64,
				hasCoordinatorRole: true,
			},
		},
		{
			name: "non-ha kafka rejects dedicated coordinator block",
			config: kafkaServiceConfig{
				isKafka:                   true,
				hasDataVolume:             true,
				dataVolumeSize:            64,
				hasCoordinatorBlock:       true,
				hasCoordinatorDataVolume:  true,
				coordinatorDataVolumeSize: 64,
			},
			expectErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := validateKafkaServiceConfig(testCase.config)

			if testCase.expectErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %s", err)
			}
		})
	}
}
