path "secret/data/marketplace/*" {
  capabilities = ["read"]
}

path "transit/encrypt/marketplace" {
  capabilities = ["update"]
}
