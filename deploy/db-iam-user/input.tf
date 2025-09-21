variable "gcp_project" {
  description = "gcp Project ID"
  default     = "flatchecker"
}

variable "db_name" {
  description = "name of the database"
}

variable "db_password_secret_id" {
  description = "Secret id for the database password"
}