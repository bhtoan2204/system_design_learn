pid_file = "/tmp/vault-agent.pid"

auto_auth {
  method "token" {
    config = {
      token = "root"
    }
  }

  sink "file" {
    config = {
      path = "/vault/secrets/agent-token"
    }
  }
}

template {
  source      = "/vault/templates/secret.tmpl"
  destination = "/vault/secrets/mysecret.txt"
}
