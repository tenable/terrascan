provider "google"{
    region = "us-west1"
}

resource "google_compute_disk" "vmEncryptedwithCsek" {
  name  = "sample-disk"
  type  = "pd-ssd"
  zone  = "us-central1-a"
  image = "debian-8-jessie-v20170523"
  labels = {
    environment = "dev"
  }
  physical_block_size_bytes = 4096
}