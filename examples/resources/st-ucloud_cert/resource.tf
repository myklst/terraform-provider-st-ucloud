resource "st-ucloud_cert" "def" {
  cert_name = "def"
  user_cert = "-----BEGIN CERTIFICATE-----\nxxxx\n-----END CERTIFICATE-----\n"
  private_key = "-----BEGIN RSA PRIVATE KEY-----\nxxxx\n-----END RSA PRIVATE KEY-----\n"
  ca_cert = "-----BEGIN CERTIFICATE-----\nxxxx\n-----END CERTIFICATE-----\n"
}

