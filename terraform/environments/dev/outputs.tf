# terraform/environments/dev/outputs.tf
# SONA IoT Elderly Care Platform - Dev Environment Outputs

# ---------------------------------------------------------------------------
# Compute (Kubernetes)
# ---------------------------------------------------------------------------
output "k8s_cluster_endpoint" {
  description = "The endpoint URL of the DigitalOcean Kubernetes cluster"
  value       = module.compute.k8s_cluster_endpoint
}

output "k8s_cluster_id" {
  description = "The unique identifier of the Kubernetes cluster"
  value       = module.compute.k8s_cluster_id
}

# ---------------------------------------------------------------------------
# Database (PostgreSQL)
# ---------------------------------------------------------------------------
output "database_uri" {
  description = "Full connection URI for the managed PostgreSQL cluster"
  value       = module.database.database_uri
  sensitive   = true
}

output "database_host" {
  description = "Hostname of the managed PostgreSQL cluster"
  value       = module.database.database_host
}

# ---------------------------------------------------------------------------
# Storage (Spaces)
# ---------------------------------------------------------------------------
output "spaces_bucket_name" {
  description = "Name of the DigitalOcean Spaces bucket for IoT data and assets"
  value       = module.storage.spaces_bucket_name
}

output "spaces_endpoint" {
  description = "Regional endpoint URL for the Spaces bucket"
  value       = module.storage.spaces_endpoint
}

# ---------------------------------------------------------------------------
# Cache (Redis)
# ---------------------------------------------------------------------------
output "redis_uri" {
  description = "Full connection URI for the managed Redis cluster"
  value       = module.cache.redis_uri
  sensitive   = true
}

# ---------------------------------------------------------------------------
# Networking (VPC)
# ---------------------------------------------------------------------------
output "vpc_id" {
  description = "ID of the VPC containing all SONA platform resources"
  value       = module.networking.vpc_id
}
