resource "google_alloydb_cluster" "coams_cluster" {
  cluster_id = "coams-cluster"
  location   = var.region
  network_config {
    network = "default" 
    # Assumes standard VPC peering setup for AlloyDB
  }
  
  initial_users {
    user     = "coams_admin"
    password = "change_me_in_secrets_manager"
  }
}

resource "google_alloydb_instance" "coams_primary" {
  cluster       = google_alloydb_cluster.coams_cluster.name
  instance_id   = "coams-primary"
  instance_type = "PRIMARY"

  machine_config {
    cpu_count = 2
  }
}

# The actual pgvector extension creation and database schema setup 
# will be managed out-of-band by the COAMS application migrations
# to avoid state locking issues, similar to how Spanner is handled in GERP.
