# terraform/modules/storage/variables.tf
# SONA IoT Elderly Care Platform - Storage Module Variables

variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
  default     = "sona"
}

variable "environment" {
  description = "Deployment environment (dev, staging, production)"
  type        = string
  default     = "production"
}

variable "spaces_region" {
  description = "DigitalOcean Spaces region (Spaces availability differs from Droplet regions)"
  type        = string
  default     = "nyc3"
}

variable "bucket_acl" {
  description = "Access control for the bucket (private or public-read)"
  type        = string
  default     = "private"

  validation {
    condition     = contains(["private", "public-read"], var.bucket_acl)
    error_message = "Bucket ACL must be 'private' or 'public-read'."
  }
}

variable "versioning_enabled" {
  description = "Enable versioning on the Spaces bucket"
  type        = bool
  default     = true
}

variable "cors_allowed_origins" {
  description = "Origins allowed to read from the bucket (e.g., Svelte app domains)"
  type        = list(string)
  default     = ["*"]
}

variable "cors_upload_origins" {
  description = "Origins allowed to upload to the bucket (presigned URL uploads)"
  type        = list(string)
  default     = ["*"]
}

variable "archive_after_days" {
  description = "Number of days after which to expire old uploads (0 to disable)"
  type        = number
  default     = 0
}

variable "force_destroy" {
  description = "Allow Terraform to destroy the bucket even if it contains objects"
  type        = bool
  default     = false
}

# ---- CDN ----

variable "cdn_ttl" {
  description = "CDN cache TTL in seconds"
  type        = number
  default     = 3600
}

variable "cdn_custom_domain" {
  description = "Custom domain for the CDN endpoint (leave empty to skip)"
  type        = string
  default     = ""
}

variable "cdn_certificate_name" {
  description = "Name of the SSL certificate for CDN custom domain"
  type        = string
  default     = ""
}
