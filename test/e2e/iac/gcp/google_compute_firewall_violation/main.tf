provider "google" {
 region      = "us-west1"
}

resource "google_compute_firewall" "portIsOpen" {
  provider = google
  name = "website-fw-2"
  network = "some-network"
  source_ranges = ["0.0.0.0/0"]
  allow {
    protocol = "tcp"
    ports = ["22"]
  }
  allow {
    protocol = "tcp"
    ports = ["3389"]
  }
  direction = "INGRESS"
}

resource "google_compute_firewall" "unrestrictedRdpAccess" {
  name    = "another-firewall"
  network = "another-network"
  direction = "INGRESS"

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["3389"]
  }

  source_tags = ["web"]
}