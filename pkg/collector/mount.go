package collector

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"log"
	"whale/pkg/client/cri"
	"whale/pkg/util"
)

type mountCollector struct {
	dirSizeDesc *prometheus.Desc
	ctx         context.Context
}

// Describe implements the prometheus.Collector interface.
func (collector *mountCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.dirSizeDesc
}

// Collect implements the prometheus.Collector interface.
func (collector *mountCollector) Collect(ch chan<- prometheus.Metric) {
	err := collector.update(ch)
	if err != nil {
		log.Println(err)
	}
}

func newMountCollector(ctx context.Context) *mountCollector {
	const subsystem = "mount"
	return &mountCollector{
		dirSizeDesc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "dir_size"),
			"container mount host_path directory size", []string{"label_app", "pod", "pod_ip", "node_ip", "namespace", "host_path", "container_path"}, nil,
		),
		ctx: ctx,
	}
}

func (collector *mountCollector) update(ch chan<- prometheus.Metric) error {
	client := collector.ctx.Value("containerdClient").(*cri.Client)
	podNamespace := collector.ctx.Value("namespace").(string)
	mountPaths := collector.ctx.Value("mountPaths").([]string)
	allNamespaces := collector.ctx.Value("allNamespaces").(bool)
	nodeIP := collector.ctx.Value("nodeIP").(string)
	var pods []*runtimeapi.PodSandbox
	var err error
	if allNamespaces {
		pods, err = client.ListAllPods()
		if err != nil {
			return err
		}
	} else {
		pods, err = client.ListNamespacePods(podNamespace)
		if err != nil {
			return err
		}
	}

	for _, pod := range pods {
		podStatus, err := client.PodStatus(pod.Id)
		if err != nil {
			return err
		}

		appContainerStatus, err := collector.getAppContainerStatus(pod)
		if err != nil {
			return err
		}

		if appContainerStatus != nil {
			for _, mount := range appContainerStatus.Status.Mounts {
				if util.SliceContains(mountPaths, mount.ContainerPath) {
					size, err := util.DirSize(mount.HostPath)
					if err != nil {
						log.Println("err",err)
					}
					ch <- prometheus.MustNewConstMetric(collector.dirSizeDesc, prometheus.GaugeValue, float64(size),
						pod.Labels["app"], pod.Metadata.Name, podStatus.Status.Network.Ip, nodeIP, pod.Metadata.Namespace, mount.HostPath, mount.ContainerPath)
				}
			}
		}
	}
	return nil
}

func (collector *mountCollector) getAppContainerStatus(pod *runtimeapi.PodSandbox) (*runtimeapi.ContainerStatusResponse, error) {
	client := collector.ctx.Value("containerdClient").(*cri.Client)

	var containers []*runtimeapi.Container
	containers, err := client.ListPodContainers(pod.Metadata.Uid)
	if err != nil {
		return nil, err
	}

	var appContainerStatus *runtimeapi.ContainerStatusResponse
	for _, container := range containers {
		status, err := client.ContainerStatus(container.Id)
		if err != nil {
			return nil, err
		}
		if status.Status.Metadata.Name == pod.Labels["app"] {
			appContainerStatus = status
		}
	}
	return appContainerStatus, nil
}
