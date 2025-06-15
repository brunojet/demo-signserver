terraform {
  backend "s3" {
    bucket       = "brunojet-tfstate"               # Substitua pelo nome do bucket S3 para o state
    key          = "orchestrator/terraform.tfstate" # Caminho/arquivo do state
    region       = "us-east-1"                      # Ajuste para sua região
    use_lockfile = true
    encrypt      = true
  }
}
