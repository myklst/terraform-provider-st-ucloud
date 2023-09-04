resource "st-ucloud_cdn_domain" "def" {
  depends_on = [
    st-ucloud_cert.def
  ]

  domain       = "xxxxx-cdn.com"
  test_url     = "http://xxxx-cdn.com/"
  area_code    = "cn"
  cdn_type     = "web"
  cdn_protocol = "http|https"
  cert_name    = "def"
  tag          = "Default"

  origin_conf = {
    origin_ip_list = ["origin.xxxx-cdn.com"]
    #    origin_host = "origin.xxxx-cdn.com"
    origin_port          = 80
    origin_protocol      = "http"
    origin_follow301     = 1
    backup_origin_enable = false
  }

  cache_conf = {
    #    cache_host = ""
    cache_list = [
      {
        path_pattern = "/*"
        ttl : 30
        cache_unit     = "day"
        cache_behavior = true
        #        follow_origin_rule = true
        description = "test"
      },
    ]
  }

  access_control_conf = {
    #    ip_blacklist = []
    #    ip_blacklist_empty = true
    refer_conf = {
      #      refer_type = 0
      #      null_refer = 0
      #      refer_list = []
    }
    #    enable_refer = false
  }

  advanced_conf = {
    #    http_client_header_list = []
    #    http_origin_header_list = []
    #    http_to_https = false
  }
}
