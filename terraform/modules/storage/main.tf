# terraform/modules/storage/main.tf
# SONA IoT Elderly Care Platform - Storage Module
# DigitalOcean Spaces bucket for static assets, Svelte builds, and media uploads.

resource "digitalocean_spaces_bucket" "assets" {
  name   = "${var.project_name}-assets-${var.environment}"
  region = var.spaces_region

  acl = var.bucket_acl

  versioning {
    enabled = var.versioning_enabled
  }

  # -------------------------------------------------------------------------
  # CORS Rules
  # Allow the SONA web application to access bucket resources directly.
  # -------------------------------------------------------------------------
  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET", "HEAD"]
    allowed_origins = var.cors_allowed_origins
    max_age_seconds = 3600
  }

  # Allow uploads from the application (presigned URLs)
  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET", "PUT", "POST", "DELETE", "HEAD"]
    allowed_origins = var.cors_upload_origins
    max_age_seconds = 3600
  }

  # Lifecycle rule: expire incomplete multipart uploads
  lifecycle_rule {
    id      = "abort-incomplete-uploads"
    enabled = true

    abort_incomplete_multipart_upload_days = 7
  }

  # Lifecycle rule: transition old media uploads to cheaper storage (optional)
  dynamic "lifecycle_rule" {
    for_each = var.archive_after_days > 0 ? [1] : []
    content {
      id      = "archive-old-media"
      enabled = true
      prefix  = "uploads/"

      expiration {
        days = var.archive_after_days
      }
    }
  }

  force_destroy = var.force_destroy
}

# -----------------------------------------------------------------------------
# CDN Endpoint
# Serve static assets (Svelte builds, images) via DigitalOcean CDN.
# -----------------------------------------------------------------------------
resource "digitalocean_cdn" "assets" {
  origin           = digitalocean_spaces_bucket.assets.bucket_domain_name
  ttl              = var.cdn_ttl
  certificate_name = var.cdn_custom_domain != "" ? var.cdn_certificate_name : null
  custom_domain    = var.cdn_custom_domain != "" ? var.cdn_custom_domain : null
}
