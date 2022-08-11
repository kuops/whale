package cri

import (
	"context"
	"google.golang.org/grpc"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"time"
)

type Client struct {
	RuntimeClient runtimeapi.RuntimeServiceClient
}

var (
	timeout = time.Second * 2
)

func NewClient(addr string) (*Client,*grpc.ClientConn,error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(timeout))
	if err != nil {
		return nil,nil,err
	}
	runtimeClient := runtimeapi.NewRuntimeServiceClient(conn)
	return &Client{RuntimeClient: runtimeClient},conn,nil
}

func (client *Client) ListAllPods() ([]*runtimeapi.PodSandbox,error) {
	response,err := client.RuntimeClient.ListPodSandbox(context.Background(),&runtimeapi.ListPodSandboxRequest{})
	return response.Items, err
}

func (client *Client) ListNamespacePods(namespace string) ([]*runtimeapi.PodSandbox,error) {
	response,err := client.RuntimeClient.ListPodSandbox(context.Background(),&runtimeapi.ListPodSandboxRequest{
		Filter: &runtimeapi.PodSandboxFilter{
			LabelSelector: map[string]string{
				"io.kubernetes.pod.namespace": namespace,
			},
		},
	})
	return response.Items, err
}

func (client *Client) ListRunningContainers() ([]*runtimeapi.Container,error) {
	response,err := client.RuntimeClient.ListContainers(context.Background(),&runtimeapi.ListContainersRequest{
		Filter: &runtimeapi.ContainerFilter{
			State: &runtimeapi.ContainerStateValue{State: runtimeapi.ContainerState_CONTAINER_RUNNING},
		},
	})
	return response.Containers, err
}

func (client *Client) ListPodContainers(podUID string) ([]*runtimeapi.Container,error) {
	response,err :=  client.RuntimeClient.ListContainers(context.Background(),&runtimeapi.ListContainersRequest{
		Filter: &runtimeapi.ContainerFilter{
			LabelSelector: map[string]string{
				"io.kubernetes.pod.uid": podUID,
			},
		},
	})
	return response.Containers, err
}

func (client *Client) ContainerStats(id string) (*runtimeapi.ContainerStatsResponse,error)  {
	return client.RuntimeClient.ContainerStats(context.Background(),&runtimeapi.ContainerStatsRequest{ContainerId: id})
}

func (client *Client) ContainerStatus(id string) (*runtimeapi.ContainerStatusResponse,error) {
	return client.RuntimeClient.ContainerStatus(context.Background(),&runtimeapi.ContainerStatusRequest{
		ContainerId: id,
		Verbose: true,
	})
}

func (client *Client) PodStatus(id string) (*runtimeapi.PodSandboxStatusResponse,error) {
	return client.RuntimeClient.PodSandboxStatus(context.Background(),&runtimeapi.PodSandboxStatusRequest{
		PodSandboxId: id,
		Verbose: true,
	})
}

