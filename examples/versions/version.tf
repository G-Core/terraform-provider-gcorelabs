terraform {
  required_version = ">= 0.13.0"
  required_providers {
    gcore = {
      source  = "local.gcorelabs.com/repo/gcore"
      version = "~>0.0.8"
    }
  }
}