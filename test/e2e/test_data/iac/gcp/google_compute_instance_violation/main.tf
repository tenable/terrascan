provider "google" {
  region      = "us-west1"
}

resource "google_compute_instance" "checkVM_NoFullCloudAccess" {
  name         = "test"
  machine_type = "n1-standard-1"
  zone         = "us-central1-a"

  tags = ["foo", "bar"]

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  // Local SSD disk
  scratch_disk {
    interface = "SCSI"
  }

  network_interface {
    network = "default"

    access_config {
      // Ephemeral IP
    }
  }

  metadata = {
    foo = "bar"
  }

  metadata_startup_script = "echo hi > /test.txt"

  service_account {
    scopes = ["cloud-platform"]
  }
}

resource "google_compute_instance" "defaultServiceAccountUsed" {
  name         = "test2"
  machine_type = "n1-standard-1"
  zone         = "us-central1-a"

  tags = ["foo", "bar"]

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  network_interface {
    network = "default"

    access_config {
      // Ephemeral IP
    }
  }

  service_account {
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
    email = "abcd1234-compute@developer.gserviceaccount.com."
  }
}


resource "google_compute_instance" "checkIpForward" {
  name         = "sample_name"
  machine_type = "n1-standard-1"
  zone         = "us-central1-a"

  tags = ["block-project-ssh-keys", "false"]

  can_ip_forward = true

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  // Local SSD disk
  scratch_disk {
    interface = "SCSI"
  }

  network_interface {
    network = "default"

    access_config {
      // Ephemeral IP
    }
  }

  metadata = {
    block-project-ssh-keys = "false"
  }

  metadata_startup_script = "echo hi > /test.txt"

  service_account {
    scopes = ["cloud-platform"]
  }
}

resource "google_compute_instance" "projectWideSshKeysUsed" {
  name         = "sample_name2"
  machine_type = "n1-standard-1"
  zone         = "us-central1-a"

  tags = ["block-project-ssh-keys", "false"]

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  // Local SSD disk
  scratch_disk {
    interface = "SCSI"
  }

  network_interface {
    network = "default"

    access_config {
      // Ephemeral IP
    }
  }

  metadata = {
    block-project-ssh-keys = "false"
  }

  metadata_startup_script = "echo hi > /test.txt"

  service_account {
    scopes = ["cloud-platform"]
  }
}

resource "google_compute_project_metadata" "projectWideSshKeysUsed" {
  metadata = {
    block-project-ssh-keys = "false"
  }
}

resource "google_compute_instance" "osLoginEnabled" {
  name         = "sample_name2"
  machine_type = "n1-standard-1"
  zone         = "us-central1-a"

  tags = ["block-project-ssh-keys", "false"]

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  // Local SSD disk
  scratch_disk {
    interface = "SCSI"
  }

  network_interface {
    network = "default"

    access_config {
      // Ephemeral IP
    }
  }

  metadata = {
    enable-oslogin = "false"
  }

  metadata_startup_script = "echo hi > /test.txt"

  service_account {
    scopes = ["cloud-platform"]
  }
}

resource "google_compute_project_metadata" "osLoginEnabled" {
  metadata = {
    enable-oslogin = "false"
  }
}

resource "google_compute_instance" "serialPortEnabled" {
  name         = "sample_name2"
  machine_type = "n1-standard-1"
  zone         = "us-central1-a"

  tags = ["block-project-ssh-keys", "false"]

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  // Local SSD disk
  scratch_disk {
    interface = "SCSI"
  }

  network_interface {
    network = "default"

    access_config {
      // Ephemeral IP
    }
  }

  metadata = {
    serial-port-enable = "false"
  }

  metadata_startup_script = "echo hi > /test.txt"

  service_account {
    scopes = ["cloud-platform"]
  }
}

resource "google_compute_project_metadata" "serialPortEnabled" {
  metadata = {
    serial-port-enable = "false"
  }
}

resource "google_compute_instance" "encryptedwithCsek" {
  name         = "sample_name3"
  machine_type = "n1-standard-1"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  // Local SSD disk
  scratch_disk {
    interface = "SCSI"
  }

  network_interface {
    network = "default"

    access_config {
      // Ephemeral IP
    }
  }

  service_account {
    scopes = ["cloud-platform"]
  }
}

resource "google_compute_instance" "shieldedVmEenabled" {
  name         = "sample_name4"
  machine_type = "n1-standard-1"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  // Local SSD disk
  scratch_disk {
    interface = "SCSI"
  }

  network_interface {
    network = "default"

    access_config {
      // Ephemeral IP
    }
  }

  shielded_instance_config {
    enable_vtpm = false
  }
  
  service_account {
    scopes = ["cloud-platform"]
  }
}