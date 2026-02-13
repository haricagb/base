# terraform/modules/compute/outputs.tf
# SONA IoT Elderly Care Platform - Compute Module Outputs

output "cluster_id" {
  description = "ID of the Kubernetes cluster"
  value       = digitalocean_kubernetes_cluster.sona.id
}

output "cluster_urn" {
  description = "URN of the Kubernetes cluster"
  value       = digitalocean_kubernetes_cluster.sona.urn
}

output "cluster_name" {
  description = "Name of the Kubernetes cluster"
  value       = digitalocean_kubernetes_cluster.sona.name
}

output "cluster_endpoint" {
  description = "API server endpoint of the Kubernetes cluster"
  value       = digitalocean_kubernetes_cluster.sona.endpoint
}

output "cluster_ipv4_address" {
  description = "Public IPv4 address of the Kubernetes cluster API server"
  value       = digitalocean_kubernetes_cluster.sona.ipv4_address
}

output "k8s_version" {
  description = "Running Kubernetes version"
  value       = digitalocean_kubernetes_cluster.sona.version
}

output "cluster_status" {
  description = "Status of the Kubernetes cluster"
  value       = digitalocean_kubernetes_cluster.sona.status
}

output "kubeconfig" {
  description = "Raw kubeconfig for cluster access"
  value       = digitalocean_kubernetes_cluster.sona.kube_config[0].raw_config
  sensitive   = true
}

output "kubeconfig_host" {
  description = "Kubernetes API server host from kubeconfig"
  value       = digitalocean_kubernetes_cluster.sona.kube_config[0].host
}

output "kubeconfig_token" {
  description = "Authentication token from kubeconfig"
  value       = digitalocean_kubernetes_cluster.sona.kube_config[0].token
  sensitive   = true
}

output "kubeconfig_cluster_ca_certificate" {
  description = "Cluster CA certificate (base64 encoded)"
  value       = digitalocean_kubernetes_cluster.sona.kube_config[0].cluster_ca_certificate
  sensitive   = true
}

output "default_node_pool_id" {
  description = "ID of the default node pool"
  value       = digitalocean_kubernetes_cluster.sona.node_pool[0].id
}

output "iot_worker_pool_id" {
  description = "ID of the IoT worker node pool (if created)"
  value       = var.enable_iot_worker_pool ? digitalocean_kubernetes_node_pool.iot_workers[0].id : null
}
