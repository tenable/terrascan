terraform {
  required_providers {
    foo = {
      source = "hashicorp/foo"
      configuration_aliases = [ foo.bar ]
    }
  }
}
