resource "google_cloud_run_service" "go_app" {
  name     = "go-background-worker"
  location = var.gcp_region

  template {
    spec {
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
  name        = "run-go-app-every-minute"
  description = "Triggers Go app every minute"
  schedule    = "* * * * *" # Every minute
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

    body = jsonencode({
      trigger = "scheduler"
    })
  }
}