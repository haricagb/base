# terraform/modules/networking/outputs.tf
# SONA IoT Elderly Care Platform - Networking Module Outputs

output "vpc_id" {
  description = "ID of the created VPC"
  value       = digitalocean_vpc.sona.id
}

output "vpc_urn" {
  description = "URN of the created VPC"
  value       = digitalocean_vpc.sona.urn
}

output "vpc_ip_range" {
  description = "IP range of the VPC"
  value       = digitalocean_vpc.sona.ip_range
}

output "firewall_web_id" {
  description = "ID of the web firewall"
  value       = digitalocean_firewall.web.id
}

output "firewall_database_id" {
  description = "ID of the database firewall"
  value       = digitalocean_firewall.database.id
}

output "firewall_k8s_id" {
  description = "ID of the Kubernetes internal firewall"
  value       = digitalocean_firewall.k8s_internal.id
}
