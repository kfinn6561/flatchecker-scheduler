module "db_service_account" {
  source  = "./db-iam-user"
  service_account_name = "scheduler-service-account"
  db_name = var.db_name
  db_password_secret_id = var.db_password_secret_id
}

resource "google_project_iam_member" "pubsub_publisher_binding" {
  project = var.gcp_project
  role    = "roles/pubsub.publisher"
  member  = "serviceAccount:${module.db_service_account.db_service_account_email}"
}

resource "google_project_iam_member" "pubsub_subscriber_binding" {
  project = var.gcp_project
  role    = "roles/pubsub.subscriber"
  member  = "serviceAccount:${module.db_service_account.db_service_account_email}"
}

resource "google_cloud_run_service" "go_app" {
  name     = "flatchecker-scheduler"
  location = var.gcp_region

  template {
    metadata {
      annotations = {
        "run.googleapis.com/cloudsql-instances" = "${var.gcp_project}:${var.gcp_region}:${var.db_name}"
      }
    }
    spec {
      service_account_name = module.db_service_account.db_service_account_email

      containers {
        image = var.image_uri
        resources {
          limits = {
            memory = "256Mi"
            cpu    = "1"
          }
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

# Allow unauthenticated invocations (optional: secure with IAM below)
resource "google_cloud_run_service_iam_member" "invoker" {
  service  = google_cloud_run_service.go_app.name
  location = var.gcp_region
  role     = "roles/run.invoker"
  member   = "allUsers"
}

# Cloud Scheduler job
resource "google_cloud_scheduler_job" "invoke_cloud_run" {
  name        = "run-scheduler-every-minute"
  description = "Triggers the scheduler every minute"
  schedule    = "* * * * *"    # Every minute
  region      = "europe-west1" # Cloud Scheduler is not available in europe-north1
  time_zone   = "Etc/UTC"

  http_target {
    http_method = "POST"
    uri         = google_cloud_run_service.go_app.status[0].url

    # If using authentication instead of allUsers:
    # oidc_token {
    #   service_account_email = google_service_account.scheduler_sa.email
    # }

    headers = {
      "Content-Type" = "application/json"
    }

    body = base64encode(jsonencode({
      trigger = "scheduler"
    }))
  }
}