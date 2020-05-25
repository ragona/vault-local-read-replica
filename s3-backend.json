storage "s3" {
    region = "us-west-2"
    bucket = "ragona-vault-test"
}

listener "tcp" {
    address     = "127.0.0.1:8200"
    tls_disable = 1       
}

seal "awskms" {
    region     = "us-west-2"
    kms_key_id = "d50ca385-b9af-4f65-99b8-695b426b78a0"
              
}

disable_mlock = true
