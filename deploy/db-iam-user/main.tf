
resource "google_service_account" "db_service_account" {
  account_id  = "db-service-account"
  description = "Service account for the database"
}

resource "google_sql_user" "iam_service_account_user" {
  name     = google_service_account.db_service_account.email
  instance = var.db_name
  type     = "CLOUD_IAM_SERVICE_ACCOUNT"
}

resource "google_project_iam_binding" "cloud_sql_user" {
  project = var.gcp_project
  role    = "roles/cloudsql.instanceUser"
  members = [
    "serviceAccount:${google_service_account.db_service_account.email}"
  ]
}

resource "google_project_iam_binding" "cloud_sql_client" {
  project = var.gcp_project
  role    = "roles/cloudsql.client"
  members = [
    "serviceAccount:${google_service_account.db_service_account.email}"
  ]
}

resource "google_secret_manager_secret_iam_member" "db_password_accessor" {
  secret_id = var.db_password_secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.db_service_account.email}"
}