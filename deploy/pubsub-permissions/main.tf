resource "google_project_iam_member" "pubsub_publisher_binding" {
  project = var.gcp_project
  role    = "roles/pubsub.publisher"
  member  = "serviceAccount:${var.service_account_email}"
}

resource "google_project_iam_member" "pubsub_subscriber_binding" {
  project = var.gcp_project
  role    = "roles/pubsub.subscriber"
  member  = "serviceAccount:${var.service_account_email}"
}

resource "google_project_iam_member" "pubsub_editor_binding" {
  project = var.gcp_project
  role    = "roles/pubsub.editor"
  member  = "serviceAccount:${var.service_account_email}"
}