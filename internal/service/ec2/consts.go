package ec2

const (
	// https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_CreditSpecificationRequest.html#API_CreditSpecificationRequest_Contents
	CPUCreditsStandard  = "standard"
	CPUCreditsUnlimited = "unlimited"
)

func CPUCredits_Values() []string {
	return []string{
		CPUCreditsStandard,
		CPUCreditsUnlimited,
	}
}

const (
	// https://docs.aws.amazon.com/vpc/latest/privatelink/vpce-interface.html#vpce-interface-lifecycle
	VpcEndpointStateAvailable         = "available"
	VpcEndpointStateDeleted           = "deleted"
	VpcEndpointStateDeleting          = "deleting"
	VpcEndpointStateFailed            = "failed"
	VpcEndpointStatePending           = "pending"
	VpcEndpointStatePendingAcceptance = "pendingAcceptance"
	VpcEndpointStateRejected          = "rejected"
)

// See https://docs.aws.amazon.com/vm-import/latest/userguide/vmimport-image-import.html#check-import-task-status
const (
	EBSSnapshotImportStateActive     = "active"
	EBSSnapshotImportStateDeleting   = "deleting"
	EBSSnapshotImportStateDeleted    = "deleted"
	EBSSnapshotImportStateUpdating   = "updating"
	EBSSnapshotImportStateValidating = "validating"
	EBSSnapshotImportStateValidated  = "validated"
	EBSSnapshotImportStateConverting = "converting"
	EBSSnapshotImportStateCompleted  = "completed"
)

// See https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_CreateNetworkInterface.html#API_CreateNetworkInterface_Example_2_Response.
const (
	NetworkInterfaceStatusPending = "pending"
)

// See https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeInternetGateways.html#API_DescribeInternetGateways_Example_1_Response.
const (
	InternetGatewayAttachmentStateAvailable = "available"
)

// See https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_CustomerGateway.html#API_CustomerGateway_Contents.
const (
	CustomerGatewayStateAvailable = "available"
	CustomerGatewayStateDeleted   = "deleted"
	CustomerGatewayStateDeleting  = "deleting"
	CustomerGatewayStatePending   = "pending"
)

const (
	VpnTunnelOptionsIKEVersion1 = "ikev1"
	VpnTunnelOptionsIKEVersion2 = "ikev2"
)

func VpnTunnelOptionsIKEVersion_Values() []string {
	return []string{
		VpnTunnelOptionsIKEVersion1,
		VpnTunnelOptionsIKEVersion2,
	}
}

// Supported tunnel option values:
// https://docs.k2.cloud/en/api/ec2/datatypes/VpnTunnelOptionsSpecification.html
const (
	VpnTunnelOptionsPhase1EncryptionAlgorithmAES128           = "aes128"
	VpnTunnelOptionsPhase1EncryptionAlgorithmAES256           = "aes256"
	VpnTunnelOptionsPhase1EncryptionAlgorithmAES_CTR128       = "aes_ctr128"
	VpnTunnelOptionsPhase1EncryptionAlgorithmAES_CTR256       = "aes_ctr256"
	VpnTunnelOptionsPhase1EncryptionAlgorithmAES_GCM128       = "aes_gcm128"
	VpnTunnelOptionsPhase1EncryptionAlgorithmAES_GCM256       = "aes_gcm256"
	VpnTunnelOptionsPhase1EncryptionAlgorithmCamellia128      = "camellia128"
	VpnTunnelOptionsPhase1EncryptionAlgorithmCamellia256      = "camellia256"
	VpnTunnelOptionsPhase1EncryptionAlgorithmChaCha20Poly1305 = "chacha20poly1305"
)

func VpnTunnelOptionsPhase1EncryptionAlgorithm_Values() []string {
	return []string{
		VpnTunnelOptionsPhase1EncryptionAlgorithmAES128,
		VpnTunnelOptionsPhase1EncryptionAlgorithmAES256,
		VpnTunnelOptionsPhase1EncryptionAlgorithmAES_CTR128,
		VpnTunnelOptionsPhase1EncryptionAlgorithmAES_CTR256,
		VpnTunnelOptionsPhase1EncryptionAlgorithmAES_GCM128,
		VpnTunnelOptionsPhase1EncryptionAlgorithmAES_GCM256,
		VpnTunnelOptionsPhase1EncryptionAlgorithmCamellia128,
		VpnTunnelOptionsPhase1EncryptionAlgorithmCamellia256,
		VpnTunnelOptionsPhase1EncryptionAlgorithmChaCha20Poly1305,
	}
}

const (
	VpnTunnelOptionsPhase1IntegrityAlgorithmSHA1     = "sha1"
	VpnTunnelOptionsPhase1IntegrityAlgorithmSHA2_256 = "sha256"
	VpnTunnelOptionsPhase1IntegrityAlgorithmSHA2_384 = "sha384"
	VpnTunnelOptionsPhase1IntegrityAlgorithmSHA2_512 = "sha512"
)

func VpnTunnelOptionsPhase1IntegrityAlgorithm_Values() []string {
	return []string{
		VpnTunnelOptionsPhase1IntegrityAlgorithmSHA1,
		VpnTunnelOptionsPhase1IntegrityAlgorithmSHA2_256,
		VpnTunnelOptionsPhase1IntegrityAlgorithmSHA2_384,
		VpnTunnelOptionsPhase1IntegrityAlgorithmSHA2_512,
	}
}

const (
	VpnTunnelOptionsPhase2EncryptionAlgorithmAES128           = "aes128"
	VpnTunnelOptionsPhase2EncryptionAlgorithmAES256           = "aes256"
	VpnTunnelOptionsPhase2EncryptionAlgorithmAES_CCM128       = "aes_ccm128"
	VpnTunnelOptionsPhase2EncryptionAlgorithmAES_CCM256       = "aes_ccm256"
	VpnTunnelOptionsPhase2EncryptionAlgorithmAES_CTR128       = "aes_ctr128"
	VpnTunnelOptionsPhase2EncryptionAlgorithmAES_CTR256       = "aes_ctr256"
	VpnTunnelOptionsPhase2EncryptionAlgorithmAES_GCM128       = "aes_gcm128"
	VpnTunnelOptionsPhase2EncryptionAlgorithmAES_GCM256       = "aes_gcm256"
	VpnTunnelOptionsPhase2EncryptionAlgorithmCamellia128      = "camellia128"
	VpnTunnelOptionsPhase2EncryptionAlgorithmCamellia256      = "camellia256"
	VpnTunnelOptionsPhase2EncryptionAlgorithmChaCha20Poly1305 = "chacha20poly1305"
)

func VpnTunnelOptionsPhase2EncryptionAlgorithm_Values() []string {
	return []string{
		VpnTunnelOptionsPhase2EncryptionAlgorithmAES128,
		VpnTunnelOptionsPhase2EncryptionAlgorithmAES256,
		VpnTunnelOptionsPhase2EncryptionAlgorithmAES_CCM128,
		VpnTunnelOptionsPhase2EncryptionAlgorithmAES_CCM256,
		VpnTunnelOptionsPhase2EncryptionAlgorithmAES_CTR128,
		VpnTunnelOptionsPhase2EncryptionAlgorithmAES_CTR256,
		VpnTunnelOptionsPhase2EncryptionAlgorithmAES_GCM128,
		VpnTunnelOptionsPhase2EncryptionAlgorithmAES_GCM256,
		VpnTunnelOptionsPhase2EncryptionAlgorithmCamellia128,
		VpnTunnelOptionsPhase2EncryptionAlgorithmCamellia256,
		VpnTunnelOptionsPhase2EncryptionAlgorithmChaCha20Poly1305,
	}
}

const (
	VpnTunnelOptionsPhase2IntegrityAlgorithmSHA1     = "sha1"
	VpnTunnelOptionsPhase2IntegrityAlgorithmSHA2_256 = "sha256"
	VpnTunnelOptionsPhase2IntegrityAlgorithmSHA2_384 = "sha384"
	VpnTunnelOptionsPhase2IntegrityAlgorithmSHA2_512 = "sha512"
)

func VpnTunnelOptionsPhase2IntegrityAlgorithm_Values() []string {
	return []string{
		VpnTunnelOptionsPhase2IntegrityAlgorithmSHA1,
		VpnTunnelOptionsPhase2IntegrityAlgorithmSHA2_256,
		VpnTunnelOptionsPhase2IntegrityAlgorithmSHA2_384,
		VpnTunnelOptionsPhase2IntegrityAlgorithmSHA2_512,
	}
}

const (
	VpnConnectionTypeIpsec1      = "ipsec.1"
	VpnConnectionTypeIpsecLegacy = "ipsec.legacy"
)

func VpnConnectionType_Values() []string {
	return []string{
		VpnConnectionTypeIpsec1,
		VpnConnectionTypeIpsecLegacy,
	}
}

const (
	AmazonIPv6PoolID = "Amazon"
)

const (
	DefaultDHCPOptionsID = "default"
)

const (
	DefaultSecurityGroupName = "default"
)

const (
	LaunchTemplateVersionDefault = "$Default"
	LaunchTemplateVersionLatest  = "$Latest"
)
