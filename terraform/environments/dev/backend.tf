# terraform/environments/dev/backend.tf
# SONA IoT Elderly Care Platform - Remote State Backend (DigitalOcean Spaces)
#
# This backend stores Terraform state in an S3-compatible DigitalOcean Spaces
# bucket. Before running `terraform init`, export the following environment
# variables with your DO Spaces access credentials:
#
#   export AWS_ACCESS_KEY_ID="<your-spaces-access-key>"
#   export AWS_SECRET_ACCESS_KEY="<your-spaces-secret-key>"
#
# These are the same credentials generated from the DigitalOcean control panel
# under API -> Spaces Keys.

terraform {
  backend "s3" {
    bucket = "sona-terraform-state"
    key    = "dev/terraform.tfstate"

    # DigitalOcean Spaces S3-compatible endpoint.
    # Format: https://<region>.digitaloceanspaces.com
    endpoint = "https://nyc3.digitaloceanspaces.com"
    region   = "us-east-1" # Required by the S3 backend but unused by DO Spaces

    # Disable AWS-specific validations that are not applicable to DO Spaces.
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    skip_requesting_account_id  = true
    skip_s3_checksum            = true
  }
}
