terraform {
  required_providers {
    st-ucloud = {
      source = "example.local/myklst/st-ucloud"
    }
  }
}

provider "st-ucloud" {
  api_url = "http://api.ucloud.cn"
  public_key = "xxxxx"
  private_key = "xxxx"
  project_id = "xxxx"
  region = "cn-bj2"
  zone = "cn-bj2-02"
}
