package eks

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func statusAddon(ctx context.Context, conn *eks.EKS, clusterName, addonName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindAddonByClusterNameAndAddonName(ctx, conn, clusterName, addonName)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Status), nil
	}
}

func statusAddonUpdate(ctx context.Context, conn *eks.EKS, clusterName, addonName, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindAddonUpdateByClusterNameAddonNameAndID(ctx, conn, clusterName, addonName, id)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Status), nil
	}
}

func statusCluster(conn *eks.EKS, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindClusterByName(conn, name)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Status), nil
	}
}

func statusClusterUpdate(conn *eks.EKS, name, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindClusterUpdateByNameAndID(conn, name, id)

		if tfresource.NotFound(err) {
			return nil, eks.UpdateStatusInProgress, nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Status), nil
	}
}

func statusFargateProfile(conn *eks.EKS, clusterName, fargateProfileName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindFargateProfileByClusterNameAndFargateProfileName(conn, clusterName, fargateProfileName)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Status), nil
	}
}

func statusNodegroup(conn *eks.EKS, clusterName, nodeGroupName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindNodegroupByClusterNameAndNodegroupName(conn, clusterName, nodeGroupName)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Status), nil
	}
}

// statusNodegroupAfterUpdate reports an ACTIVE node group as still updating until the
// requested desired size is visible: K2 keeps the node group ACTIVE for a short time
// after UpdateNodegroupConfig is accepted.
func statusNodegroupAfterUpdate(conn *eks.EKS, clusterName, nodeGroupName string, desiredSize *int64) resource.StateRefreshFunc {
	refresh := statusNodegroup(conn, clusterName, nodeGroupName)

	return func() (interface{}, string, error) {
		output, status, err := refresh()

		if err != nil || desiredSize == nil || status != eks.NodegroupStatusActive {
			return output, status, err
		}

		if nodeGroup, ok := output.(*eks.Nodegroup); ok && nodeGroup.ScalingConfig != nil &&
			aws.Int64Value(nodeGroup.ScalingConfig.DesiredSize) != aws.Int64Value(desiredSize) {
			return output, eks.NodegroupStatusUpdating, nil
		}

		return output, status, nil
	}
}

func statusNodegroupUpdate(conn *eks.EKS, clusterName, nodeGroupName, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindNodegroupUpdateByClusterNameNodegroupNameAndID(conn, clusterName, nodeGroupName, id)

		if tfresource.NotFound(err) {
			return nil, eks.UpdateStatusInProgress, nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Status), nil
	}
}

func statusOIDCIdentityProviderConfig(ctx context.Context, conn *eks.EKS, clusterName, configName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := FindOIDCIdentityProviderConfigByClusterNameAndConfigName(ctx, conn, clusterName, configName)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.StringValue(output.Status), nil
	}
}
