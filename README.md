# Projeto: demo-signserver

Este projeto tem como objetivo criar uma infraestrutura para um serviço de assinatura de APKs Android. O fluxo principal consiste em receber pedidos de assinatura, processá-los e realizar a assinatura dos APKs utilizando o fabricante de POS (Point of Sale).

## Estrutura do Projeto

- **iac/**: Contém toda a infraestrutura como código (IaC) utilizando Terraform.
  - **modules/**: Módulos reutilizáveis do Terraform para recursos como DynamoDB, EventBridge, Lambda, S3, Step Functions, etc.
  - **orchestrator/** e **persistence/**: Pastas com stacks específicas para orquestração e persistência dos dados.

## Fluxo Resumido
1. Recebimento do pedido de assinatura de APK.
2. Orquestração do processo via Step Functions e EventBridge.
3. Execução da assinatura utilizando Lambda e integração com o fabricante de POS.
4. Armazenamento dos resultados e logs em DynamoDB e S3.

## Observações
- O código de infraestrutura foi migrado para este repositório.
- O projeto está em desenvolvimento inicial.

---

Para dúvidas ou sugestões, abra uma issue ou entre em contato com o mantenedor.