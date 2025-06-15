#!/bin/bash
# terraform_common.sh - Funções utilitárias para workflows Terraform

# Carrega variáveis de ambiente e define o workspace do Terraform
function setup_terraform_env() {
  # Detecta evento e branch
  if [[ "$GITHUB_EVENT_NAME" == "pull_request" ]]; then
    BRANCH_NAME="$GITHUB_BASE_REF"
  else
    BRANCH_NAME="${GITHUB_REF##*/}"
  fi

  # Se branch for feature/*, bugfix/* ou hotfix/*, valida o padrão do sufixo
  if [[ "$BRANCH_NAME" =~ ^(feature|bugfix|hotfix)/(.+)$ ]]; then
    PREFIX="${BASH_REMATCH[1]}"
    SUFIXO="${BASH_REMATCH[2]}"
    # Aceita sufixo simples (string) ou com hífens (string-string...)
    if [[ ! "$SUFIXO" =~ ^[a-zA-Z0-9]+(-[a-zA-Z0-9]+)*$ ]]; then
      echo "Erro: O nome do branch após o prefixo deve estar no padrão string ou string-string (ex: feature/foo ou feature/foo-bar)"
      exit 1
    fi
    MY_WORKSPACE="${PREFIX}-$(map_branch_to_env "$SUFIXO")"
  else
    MY_WORKSPACE="$(map_branch_to_env "$BRANCH_NAME")"
  fi

  echo "Branch detectado: $BRANCH_NAME"
  echo "Workspace do Terraform definido para: $MY_WORKSPACE"
  echo "MY_WORKSPACE=$MY_WORKSPACE" >> $GITHUB_ENV
}

# Seleciona ou cria o workspace do Terraform
function select_or_create_workspace() {
  terraform workspace select "$MY_WORKSPACE" || terraform workspace new "$MY_WORKSPACE"
}

# Função para mapear branch para ambiente lógico
function map_branch_to_env() {
  case "$1" in
    develop)
      echo "dev" ;;
    release)
      echo "hml" ;;
    master)
      echo "prod" ;;
    *)
      echo "$1" ;;
  esac
}
