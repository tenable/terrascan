provider "google" {
    region = "us-west1"
}

resource "github_repository" "privateRepoEnabled" {
  name        = "sample_repository"
}