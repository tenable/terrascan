resource "google_compute_disk" "default" {
  name = "test-disk"
}
resource "google_compute_ssl_policy" "default-ssl" {
  name = "test-ssl-policy"
  min_tls_version = "1.1"
}
resource "google_compute_firewall" "default_firewall" {
  name = "test-firewall"
  destination_ranges = ["0.0.0.0/0"]
}
