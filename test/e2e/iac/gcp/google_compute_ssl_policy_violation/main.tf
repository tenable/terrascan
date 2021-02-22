provider "google" {
    region = "us-west1"
}

resource "google_compute_ssl_policy" "weakCipherSuitesEnabled" {
  name    = "product-ssl-policy"
  profile = "MODERN"
}