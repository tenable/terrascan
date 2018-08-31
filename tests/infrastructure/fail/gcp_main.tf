resource "google_compute_disk" "default" {
  name = "test-disk"
}
resource "google_compute_ssl_policy" "default-ssl" {
  name = "test-ssl-policy"
  min_tls_version = "1.1"
}
