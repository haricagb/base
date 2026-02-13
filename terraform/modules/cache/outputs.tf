# terraform/modules/cache/outputs.tf
# SONA IoT Elderly Care Platform - Cache Module Outputs

output "cluster_id" {
  description = "ID of the Redis database cluster"
  value       = digitalocean_database_cluster.redis.id
}

output "cluster_urn" {
  description = "URN of the Redis database cluster"
  value       = digitalocean_database_cluster.redis.urn
}

output "host" {
  description = "Public hostname of the Redis cluster"
  value       = digitalocean_database_cluster.redis.host
}

output "private_host" {
  description = "Private hostname of the Redis cluster (VPC only)"
  value       = digitalocean_database_cluster.redis.private_host
}

output "port" {
  description = "Port of the Redis cluster"
  value       = digitalocean_database_cluster.redis.port
}

output "password" {
  description = "Redis cluster password"
  value       = digitalocean_database_cluster.redis.password
  sensitive   = true
}

output "connection_uri" {
  description = "Full connection URI for the Redis cluster (public, TLS)"
  value       = digitalocean_database_cluster.redis.uri
  sensitive   = true
}

output "private_connection_uri" {
  description = "Full connection URI for the Redis cluster (VPC private, TLS)"
  value       = digitalocean_database_cluster.redis.private_uri
  sensitive   = true
}

output "redis_url" {
  description = "Redis URL in standard format (rediss://user:pass@host:port) for application config"
  value       = "rediss://default:${digitalocean_database_cluster.redis.password}@${digitalocean_database_cluster.redis.private_host}:${digitalocean_database_cluster.redis.port}"
  sensitive   = true
}
