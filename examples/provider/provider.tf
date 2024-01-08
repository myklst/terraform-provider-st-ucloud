terraform {
  required_providers {
    st-ucloud = {
      source = "myklst/st-ucloud"
    }
  }
}

provider "st-ucloud" {
  region      = "cn-bj2"
  zone        = "cn-bj2-02"
  project_id  = "xxxx"
  public_key  = "xxxxx"
  private_key = "xxxx"
}
