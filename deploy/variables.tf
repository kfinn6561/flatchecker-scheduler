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
  default = "us-docker.pkg.dev/flatchecker/flatchecker-scheduler/flatchecker-scheduler:latest"
}