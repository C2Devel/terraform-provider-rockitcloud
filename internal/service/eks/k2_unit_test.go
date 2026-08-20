package eks

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceClusterK2Schema(t *testing.T) {
	resource := ResourceCluster()
	if err := resource.InternalValidate(nil, true); err != nil {
		t.Fatalf("unexpected schema validation error: %s", err)
	}

	for _, name := range []string{"enabled_cluster_log_types", "encryption_config"} {
		if field := resource.Schema[name]; !field.Computed || field.Optional {
			t.Fatalf("%s must be computed-only", name)
		}
	}

	vpc := resource.Schema["vpc_config"].Elem.(*schema.Resource).Schema
	if vpc["security_group_ids"].ForceNew {
		t.Fatal("security_group_ids must update in place")
	}
	for _, name := range []string{"endpoint_private_access", "endpoint_public_access", "public_access_cidrs"} {
		if field := vpc[name]; !field.Computed || field.Optional {
			t.Fatalf("vpc_config.%s must be computed-only", name)
		}
	}

	legacy := resource.Schema["legacy_cluster_params"].Elem.(*schema.Resource).Schema
	if !legacy["master_config"].Required {
		t.Fatal("master_config must be required when legacy_cluster_params is configured")
	}
	for blockName, flagName := range map[string]string{
		"docker_registry_config": "docker_registry_required",
		"ebs_provider_config":    "ebs_provider_required",
		"ingress_config":         "ingress_required",
		"nlb_provider_config":    "nlb_provider_required",
	} {
		block := legacy[blockName].Elem.(*schema.Resource).Schema
		if _, ok := block[flagName]; ok {
			t.Fatalf("%s must remain implicitly enabled for backward compatibility", flagName)
		}
	}
	userData := legacy["user_data_config"].Elem.(*schema.Resource).Schema
	if userData["user_data"].ForceNew || userData["user_data_content_type"].ForceNew {
		t.Fatal("user_data_config fields must update in place")
	}

	if got := resource.Timeouts.Delete; got == nil || *got != 60*time.Minute {
		t.Fatalf("cluster delete timeout must match the 60 minute delete retry window, got %v", got)
	}
}

func TestExpandAndFlattenK2ClusterConfiguration(t *testing.T) {
	network := expandNetworkConfigRequest([]interface{}{map[string]interface{}{
		"ip_family":         "ipv4",
		"pod_ipv4_cidr":     "10.10.0.0/16",
		"service_ipv4_cidr": "10.20.0.0/16",
	}})
	if got := aws.StringValue(network.PodIpv4Cidr); got != "10.10.0.0/16" {
		t.Fatalf("unexpected pod CIDR: %q", got)
	}
	flattenedNetwork := flattenNetworkConfig(&eks.KubernetesNetworkConfigResponse{
		PodIpv4Cidr: aws.String("10.10.0.0/16"),
	})
	if got := flattenedNetwork[0].(map[string]interface{})["pod_ipv4_cidr"]; got != "10.10.0.0/16" {
		t.Fatalf("unexpected flattened pod CIDR: %q", got)
	}

	securityGroups := schema.NewSet(schema.HashString, []interface{}{"sg-1", "sg-2"})
	vpcUpdate := expandVPCSecurityGroupUpdateRequest(securityGroups)
	if len(vpcUpdate.SecurityGroupIds) != 2 || vpcUpdate.SubnetIds != nil ||
		vpcUpdate.EndpointPrivateAccess != nil || vpcUpdate.EndpointPublicAccess != nil ||
		vpcUpdate.PublicAccessCidrs != nil {
		t.Fatalf("security group update contains unsupported fields: %#v", vpcUpdate)
	}

	remoteAccess := expandClusterRemoteAccessConfig([]interface{}{map[string]interface{}{
		"ec2_ssh_key": "key-name",
	}})
	if got := aws.StringValue(remoteAccess.Ec2SshKey); got != "key-name" {
		t.Fatalf("unexpected SSH key: %q", got)
	}

	userData := expandUserDataConfig([]interface{}{map[string]interface{}{
		"user_data":              "#cloud-config",
		"user_data_content_type": "cloud-config",
	}})
	if got := aws.StringValue(userData.UserData); got != "#cloud-config" {
		t.Fatalf("unexpected user data: %q", got)
	}
	if got := aws.StringValue(userData.UserDataContentType); got != "cloud-config" {
		t.Fatalf("unexpected user data content type: %q", got)
	}
	flattenedUserData := flattenUserDataConfig(userData)
	if got := flattenedUserData[0].(map[string]interface{})["user_data"]; got != "#cloud-config" {
		t.Fatalf("unexpected flattened user data: %q", got)
	}

	legacy := expandLegacyClusterParams([]interface{}{map[string]interface{}{
		"cluster_autoscaler_config": []interface{}{map[string]interface{}{
			"cluster_autoscaler_required": false,
			"cluster_autoscaler_user":     "autoscaler",
		}},
		"docker_registry_config": []interface{}{map[string]interface{}{
			"volume_size": 10,
			"volume_type": "ssd",
		}},
		"ebs_provider_config": []interface{}{map[string]interface{}{
			"ebs_user": "ebs",
		}},
		"ingress_config": []interface{}{map[string]interface{}{
			"instance_type": "small",
			"volume_size":   10,
			"volume_type":   "ssd",
		}},
		"master_config": []interface{}{map[string]interface{}{
			"high_availability": false,
			"instance_type":     "small",
			"volume_size":       10,
			"volume_type":       "ssd",
		}},
		"nlb_provider_config": []interface{}{map[string]interface{}{
			"nlb_user": "nlb",
		}},
	}})

	if aws.BoolValue(legacy.ClusterAutoscalerConfig.ClusterAutoscalerRequired) {
		t.Fatal("explicit false Cluster Autoscaler flag was not preserved")
	}
	if !aws.BoolValue(legacy.DockerRegistryConfig.DockerRegistryRequired) ||
		!aws.BoolValue(legacy.EbsProviderConfig.EbsProviderRequired) ||
		!aws.BoolValue(legacy.IngressConfig.IngressRequired) ||
		!aws.BoolValue(legacy.NlbProviderConfig.NlbProviderRequired) {
		t.Fatal("existing integrated services must remain implicitly enabled")
	}

	flattened := flattenLegacyClusterParams(&eks.LegacyClusterParamsResponse{
		ClusterAutoscalerConfig: &eks.ClusterAutoscalerConfig{
			ClusterAutoscalerRequired: aws.Bool(true),
			ClusterAutoscalerUserName: aws.String("managed-autoscaler"),
		},
		EbsProviderConfig: &eks.EbsProviderConfigResponse{
			EbsProviderRequired: aws.Bool(true),
			EbsUser:             aws.String("configured-ebs"),
			EbsUserName:         aws.String("managed-ebs"),
		},
	})
	values := flattened[0].(map[string]interface{})
	autoscaler := values["cluster_autoscaler_config"].([]interface{})[0].(map[string]interface{})
	if got := autoscaler["cluster_autoscaler_user"]; got != "managed-autoscaler" {
		t.Fatalf("unexpected autoscaler user: %q", got)
	}
	ebs := values["ebs_provider_config"].([]interface{})[0].(map[string]interface{})
	if got := ebs["ebs_user"]; got != "configured-ebs" {
		t.Fatalf("unexpected EBS user: %q", got)
	}

	emptyLegacy := flattenLegacyClusterParams(&eks.LegacyClusterParamsResponse{
		ClusterAutoscalerConfig: &eks.ClusterAutoscalerConfig{
			ClusterAutoscalerRequired: aws.Bool(false),
		},
		DockerRegistryConfig: &eks.DockerRegistryConfig{
			DockerRegistryRequired: aws.Bool(false),
		},
		EbsProviderConfig: &eks.EbsProviderConfigResponse{
			EbsProviderRequired: aws.Bool(false),
		},
		IngressConfig: &eks.IngressConfig{
			IngressRequired: aws.Bool(false),
		},
		NlbProviderConfig: &eks.NlbProviderConfigResponse{
			NlbProviderRequired: aws.Bool(false),
		},
		UserDataConfig: &eks.UserDataConfig{},
		PlacementConfig: &eks.PlacementConfig{
			Tenancy: aws.String("default"),
		},
	})
	if len(emptyLegacy) != 0 {
		t.Fatalf("empty API legacy blocks must not create Terraform drift: %#v", emptyLegacy)
	}
	if got := flattenClusterRemoteAccessConfig(&eks.RemoteAccessConfig{}); len(got) != 0 {
		t.Fatalf("empty remote access block must not create Terraform drift: %#v", got)
	}
}

func TestResourceNodeGroupK2SchemaAndLabelValidation(t *testing.T) {
	resource := ResourceNodeGroup()
	if err := resource.InternalValidate(nil, true); err != nil {
		t.Fatalf("unexpected schema validation error: %s", err)
	}

	for _, name := range []string{"force_update_version", "launch_template", "release_version", "version"} {
		if field := resource.Schema[name]; !field.Computed || field.Optional {
			t.Fatalf("%s must be computed-only", name)
		}
	}

	minSize := resource.Schema["scaling_config"].Elem.(*schema.Resource).Schema["min_size"]
	if _, errors := minSize.ValidateFunc(0, "min_size"); len(errors) == 0 {
		t.Fatal("min_size=0 must be rejected")
	}

	valid := map[string]interface{}{
		"example.com/name": "value_1",
		"empty":            "",
	}
	if _, errors := validateKubernetesLabels(valid, "labels"); len(errors) != 0 {
		t.Fatalf("valid labels rejected: %v", errors)
	}

	invalid := map[string]interface{}{
		"Bad Prefix/name":  strings.Repeat("x", 64),
		"bad..prefix/name": "value",
	}
	if _, errors := validateKubernetesLabels(invalid, "labels"); len(errors) < 3 {
		t.Fatalf("invalid labels were not fully rejected: %v", errors)
	}
}

func TestClusterAuthIsExplicitlyUnsupported(t *testing.T) {
	err := dataSourceClusterAuthRead(nil, nil)
	if err == nil || !strings.Contains(err.Error(), "aws_eks_cluster_kubeconfig") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestK2FailureDetailsAreActionable(t *testing.T) {
	clusterErr := ClusterIssuesError([]*eks.ClusterIssue{{
		Code:        aws.String("InsufficientCapacity"),
		Message:     aws.String("control plane could not be provisioned"),
		ResourceIds: aws.StringSlice([]string{"cluster/test"}),
	}})
	if clusterErr == nil ||
		!strings.Contains(clusterErr.Error(), "InsufficientCapacity") ||
		!strings.Contains(clusterErr.Error(), "control plane could not be provisioned") {
		t.Fatalf("cluster health details were lost: %v", clusterErr)
	}

	updateErr := ErrorDetailsError([]*eks.ErrorDetail{{
		ErrorCode:    aws.String("InvalidParameter"),
		ErrorMessage: aws.String("security group update was rejected"),
		ResourceIds:  aws.StringSlice([]string{"sg-test"}),
	}})
	if updateErr == nil ||
		!strings.Contains(updateErr.Error(), "InvalidParameter") ||
		!strings.Contains(updateErr.Error(), "security group update was rejected") {
		t.Fatalf("update failure details were lost: %v", updateErr)
	}
}

func TestClusterUpdateFallsBackWhenDescribeUpdateIsUnavailable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/clusters/test/updates/update-1":
			w.Header().Set("X-Amzn-Errortype", "PathNotFoundError")
			http.Error(w, `{"message":"Specified path does not exist."}`, http.StatusBadRequest)
		case "/clusters/test":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"cluster":{"name":"test","status":"READY"}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	conn := eks.New(session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("test", "test", ""),
		Endpoint:    aws.String(server.URL),
		Region:      aws.String("ru-msk"),
	})))

	update, err := waitClusterUpdateSuccessful(conn, "test", "update-1", 5*time.Second)
	if err != nil {
		t.Fatalf("unexpected fallback waiter error: %s", err)
	}
	if got := aws.StringValue(update.Status); got != eks.UpdateStatusSuccessful {
		t.Fatalf("unexpected fallback update status: %q", got)
	}

	update, err = waitClusterUpdateSuccessful(conn, "test", "", 5*time.Second)
	if err != nil {
		t.Fatalf("unexpected empty update ID fallback error: %s", err)
	}
	if got := aws.StringValue(update.Status); got != eks.UpdateStatusSuccessful {
		t.Fatalf("unexpected empty ID fallback update status: %q", got)
	}
}

func TestNodegroupUpdateFallsBackWhenDescribeUpdateIsUnavailable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/clusters/test/updates/update-1":
			w.Header().Set("X-Amzn-Errortype", "PathNotFoundError")
			http.Error(w, `{"message":"Specified path does not exist."}`, http.StatusBadRequest)
		case "/clusters/test/node-groups/general":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"nodegroup":{"clusterName":"test","nodegroupName":"general","status":"ACTIVE"}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	conn := eks.New(session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("test", "test", ""),
		Endpoint:    aws.String(server.URL),
		Region:      aws.String("ru-msk"),
	})))

	ctx := context.Background()
	update, err := waitNodegroupUpdateSuccessful(ctx, conn, "test", "general", "update-1", 5*time.Second)
	if err != nil {
		t.Fatalf("unexpected fallback waiter error: %s", err)
	}
	if got := aws.StringValue(update.Status); got != eks.UpdateStatusSuccessful {
		t.Fatalf("unexpected fallback update status: %q", got)
	}

	update, err = waitNodegroupUpdateSuccessful(ctx, conn, "test", "general", "", 5*time.Second)
	if err != nil {
		t.Fatalf("unexpected empty update ID fallback error: %s", err)
	}
	if got := aws.StringValue(update.Status); got != eks.UpdateStatusSuccessful {
		t.Fatalf("unexpected empty ID fallback update status: %q", got)
	}
}
