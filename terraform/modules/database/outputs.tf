# terraform/modules/database/outputs.tf
# SONA IoT Elderly Care Platform - Database Module Outputs

output "cluster_id" {
  description = "ID of the PostgreSQL database cluster"
  value       = digitalocean_database_cluster.postgres.id
}

output "cluster_urn" {
  description = "URN of the PostgreSQL database cluster"
  value       = digitalocean_database_cluster.postgres.urn
}

output "host" {
  description = "Public hostname of the database cluster"
  value       = digitalocean_database_cluster.postgres.host
}

output "private_host" {
  description = "Private hostname of the database cluster (VPC only)"
  value       = digitalocean_database_cluster.postgres.private_host
}

output "port" {
  description = "Port of the database cluster"
  value       = digitalocean_database_cluster.postgres.port
}

output "database_name" {
  description = "Name of the application database"
  value       = digitalocean_database_db.app_database.name
}

output "user" {
  description = "Application database user name"
  value       = digitalocean_database_user.app_user.name
}

output "password" {
  description = "Application database user password"
  value       = digitalocean_database_user.app_user.password
  sensitive   = true
}

output "connection_uri" {
  description = "Full connection URI for the application database (public)"
  value       = "postgresql://${digitalocean_database_user.app_user.name}:${digitalocean_database_user.app_user.password}@${digitalocean_database_cluster.postgres.host}:${digitalocean_database_cluster.postgres.port}/${digitalocean_database_db.app_database.name}?sslmode=require"
  sensitive   = true
}

output "private_connection_uri" {
  description = "Full connection URI for the application database (VPC private)"
  value       = "postgresql://${digitalocean_database_user.app_user.name}:${digitalocean_database_user.app_user.password}@${digitalocean_database_cluster.postgres.private_host}:${digitalocean_database_cluster.postgres.port}/${digitalocean_database_db.app_database.name}?sslmode=require"
  sensitive   = true
}

output "connection_pool_uri" {
  description = "Connection URI via PgBouncer connection pool"
  value       = digitalocean_database_connection_pool.app_pool.uri
  sensitive   = true
}

output "connection_pool_private_uri" {
  description = "Private connection URI via PgBouncer connection pool"
  value       = digitalocean_database_connection_pool.app_pool.private_uri
  sensitive   = true
}
