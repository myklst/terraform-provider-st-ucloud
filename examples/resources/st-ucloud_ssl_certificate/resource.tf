resource "st-ucloud_ssl_certificate" "test" {
  cert_name   = "test"
  user_cert   = "-----BEGIN CERTIFICATE-----\nxxxx\n-----END CERTIFICATE-----\n"
  private_key = "-----BEGIN RSA PRIVATE KEY-----\nxxxx\n-----END RSA PRIVATE KEY-----\n"
  ca_cert     = "-----BEGIN CERTIFICATE-----\nxxxx\n-----END CERTIFICATE-----\n"
}
