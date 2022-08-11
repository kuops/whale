# whale

whale is a pod hostpath directory size metrics exporter, support cri-api runtime.

```bash
❯ go run cmd/main.go --help
  -A, --all-namespaces          all namespace pods.
  -h, --help                    display help for whale.
  -l, --listen-address string   listen address. (default ":8080")
      --max-requests int        max http requests. (default 40)
  -p, --mount-paths strings     collector container mount paths. (default [/data/logs])
  -n, --namespace string        collector namespace pods. (default "default")
      --node-ip string          running node ip.
      --socket-path string      container runtime interface socket path. (default "unix:///run/containerd/containerd.sock")
```

example metrics:

```bash
container_mount_dir_size{app="whale", cluster_name="xxx-prod", container_path="/data/logs", controller_revision_hash="5945b78", host_path="/var/lib/kubelet/pods/050a5238-dca3-40e1-82b2-60677cdb04c5/volume-subpaths/data-logs/app/1", instance="10.3.248.29:8080", job="tke-retailcloud-prod-kubernetes-pods", kubernetes_namespace="sre-system", kubernetes_pod_name="whale-4zbqd", label_app="data-dev-web-saas", namespace="prod", node_ip="10.3.254.81", pod="data-dev-web-saas-1", pod_ip="10.3.248.105", pod_template_generation="7"} 30156194915
```
