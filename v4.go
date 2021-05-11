package metadata

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Based on https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-metadata-endpoint-v4.html
// and https://docs.aws.amazon.com/AmazonECS/latest/userguide/task-metadata-endpoint-v4-fargate.html

type ContainerMetadataV4 struct {
	DockerId       string `json:"DockerId"`
	Name       string `json:"Name"`
	DockerName       string `json:"DockerName"`
	Image       string `json:"Image"`
	ImageID       string `json:"ImageID"`
	Labels     struct {
		EcsCluster               string `json:"com.amazonaws.ecs.cluster"`
		EcsContainerName         string `json:"com.amazonaws.ecs.container-name"`
		EcsTaskArn               string `json:"com.amazonaws.ecs.task-arn"`
		EcsTaskDefinitionFamily  string `json:"com.amazonaws.ecs.task-definition-family"`
		EcsTaskDefinitionVersion string `json:"com.amazonaws.ecs.task-definition-version"`
	} `json:"Labels"`
	DesiredStatus string `json:"DesiredStatus"`
	KnownStatus   string `json:"KnownStatus"`
	Limits        struct {
		CPU    float64 `json:"CPU"`
		Memory int     `json:"Memory"`
	} `json:"Limits"`
	CreatedAt time.Time `json:"CreatedAt"`
	StartedAt time.Time `json:"StartedAt"`
	Type      string    `json:"Type"`
	Networks  []struct {
		NetworkMode              string   `json:"NetworkMode"`
		IPv4Addresses            []string `json:"IPv4Addresses"`
		AttachmentIndex          int      `json:"AttachmentIndex"`
		IPv4SubnetCIDRBlock      string   `json:"IPv4SubnetCIDRBlock"`
		MACAddress               string   `json:"MACAddress"`
		DomainNameServers        []string `json:"DomainNameServers"`
		DomainNameSearchList     []string `json:"DomainNameSearchList"`
		PrivateDNSName           string   `json:"PrivateDNSName"`
		SubnetGatewayIpv4Address string   `json:"SubnetGatewayIpv4Address"`
	} `json:"Networks"`
	ContainerARN   string `json:"ContainerARN"`
	LogOptions struct {
		AwslogsCreateGroup bool `json:"awslogs-create-group"`
		AwslogsGroup       string `json:"awslogs-group"`
		AwslogsRegion      string `json:"awslogs-region"`
		AwslogsStream      string `json:"awslogs-stream"`
	} `json:"LogOptions"`
	LogDriver   string `json:"LogDriver"`
}

// Retrieve ECS Container Metadata in V4 format
func GetContainerV4(ctx context.Context, client *http.Client) (*ContainerMetadataV4, error) {
	metadataUrl := os.Getenv(ecsMetadataUriEnvV4)
	if metadataUrl == "" {
		return nil, fmt.Errorf("missing metadata uri in environment (%s)", ecsMetadataUriEnvV3)
	}

	containerMetadata := &ContainerMetadataV4{}
	body, err := fetch(ctx, client, fmt.Sprintf("%s", metadataUrl))

	err = json.Unmarshal(body, &containerMetadata)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal into container metadata v4: %w", err)
	}

	return containerMetadata, nil
}
