package eks

import "time"

const (
	ErrCodeClusterNotFound   = "InvalidKubernetesCluster.NotFound"
	ErrCodeNodegroupNotFound = "InvalidKubernetesNodegroup.NotFound"
	ErrCodeIPAddressInUse    = "InvalidIPAddress.InUse"
)

const (
	IdentityProviderConfigTypeOIDC = "oidc"
)

// clusterStatusModifying is reported by K2 while a cluster update is applied.
// It has no counterpart in the AWS SDK ClusterStatus enum.
const clusterStatusModifying = "MODIFYING"

const (
	ResourcesSecrets = "secrets"
)

func Resources_Values() []string {
	return []string{
		ResourcesSecrets,
	}
}

const (
	propagationTimeout = 2 * time.Minute
)
