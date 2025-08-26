variable "gcp_project" {
  description = "gcp Project ID"
  default     = "flatchecker"
}

variable "gcp_region" {
  description = "default gcp region"
  default     = "europe-north1"
}

variable "gcp_zone" {
  description = "default gcp zone"
  default     = "europe-north1-a"
}

variable "image_uri" {
  description = "Docker image in Artifact Registry"
  default = "europe-north1-docker.pkg.dev/flatchecker/flatchecker-containers/flatchecker-scheduler:latest"
}

variable "service_account_email" {
  description = "Service account email"
  default     = "db-service-account@flatchecker.iam.gserviceaccount.com"
}

variable "db_name" {
  description = "Database name"
  default     = "flatchecker"
}