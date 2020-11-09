terraform {
  required_providers {
    aqua = {
      versions = ["0.0.1"]
      source = "aquasec.com/security/aqua"
    }
  }
}

provider "aqua" {
  user = "username"
  password = "password"
  host = "192.168.1.52"
  port = 443
  secure = true
  verify = false
}

resource "aqua_create_registry" "aws" {
  name = "AWS2"
  description =  "An AWS registry2"
  type = "AWS"
  url = "us-east-1"
  username = "aws_secret_key"
  password = "aws_acess_key"
  prefixes = []
  auto_pull = false
  auto_pull_max = 0
  auto_pull_time = ""


}