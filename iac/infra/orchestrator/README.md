# Orchestrator Terraform

Este diretório contém o orquestrador principal da infraestrutura, responsável por:
- Instanciar os módulos de Lambda, Step Function e EventBridge
- Gerenciar as policies e permissões entre os recursos
- Garantir deploy limpo, modular e sem dependências cíclicas

## Estrutura
- `main.tf`: Orquestra todos os recursos e faz os attachments de permissões
- `variables.tf`: Variáveis de ambiente e parâmetros globais
- `backend.tf`: Configuração do backend S3 para o state

## Observações
- As policies são criadas aqui para garantir que dependam de ARNs já existentes
- Não há mais dependência de módulos IAM externos
- Use workspaces para múltiplos ambientes (dev, qa, prod)

---

Dúvidas? Consulte os comentários no `main.tf` ou o README dos módulos.
