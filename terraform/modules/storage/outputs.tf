# terraform/modules/storage/outputs.tf
# SONA IoT Elderly Care Platform - Storage Module Outputs

output "bucket_name" {
  description = "Name of the Spaces bucket"
  value       = digitalocean_spaces_bucket.assets.name
}

output "bucket_urn" {
  description = "URN of the Spaces bucket"
  value       = digitalocean_spaces_bucket.assets.urn
}

output "bucket_domain_name" {
  description = "Domain name of the Spaces bucket (for direct access)"
  value       = digitalocean_spaces_bucket.assets.bucket_domain_name
}

output "bucket_region" {
  description = "Region of the Spaces bucket"
  value       = digitalocean_spaces_bucket.assets.region
}

output "bucket_endpoint" {
  description = "Full endpoint URL for the Spaces bucket"
  value       = "https://${digitalocean_spaces_bucket.assets.bucket_domain_name}"
}

output "cdn_id" {
  description = "ID of the CDN endpoint"
  value       = digitalocean_cdn.assets.id
}

output "cdn_endpoint" {
  description = "CDN endpoint URL for serving static assets"
  value       = "https://${digitalocean_cdn.assets.endpoint}"
}

output "cdn_custom_domain" {
  description = "Custom domain for the CDN (if configured)"
  value       = digitalocean_cdn.assets.custom_domain
}
