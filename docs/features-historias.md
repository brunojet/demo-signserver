# Features e Histórias — demo-signserver

## Feature 1: Intent de Assinatura
- História 1.1: Criar entidade e serviço de intent
- História 1.2: Implementar repositório de intents (interface + adapter)
- História 1.3: Expor endpoint para criar intent
- História 1.4: Validar e persistir intent no DynamoDB
- História 1.5: Testes unitários e integração do fluxo de intent

## Feature 2: Upload de APK
- História 2.1: Implementar serviço de upload
- História 2.2: Integrar com storage S3 (interface + adapter)
- História 2.3: Expor endpoint/canal para upload
- História 2.4: Validar e armazenar arquivo no S3
- História 2.5: Testes do fluxo de upload

## Feature 3: Assinatura de APK
- História 3.1: Criar serviço de assinatura (signing_service)
- História 3.2: Integrar com workerpool para processamento assíncrono
- História 3.3: Persistir resultado da assinatura
- História 3.4: Testes unitários e integração do fluxo de assinatura

## Feature 4: Notificação
- História 4.1: Implementar serviço de notificação (interface + adapter)
- História 4.2: Integrar com Webhook/Outro canal
- História 4.3: Enviar notificação ao finalizar assinatura
- História 4.4: Testes do fluxo de notificação

## Feature 5: Orquestração e Workers
- História 5.1: Implementar orquestrador central
- História 5.2: Criar e integrar workers (upload, signing, notification)
- História 5.3: Gerenciar filas/canais (SQS)
- História 5.4: Testes de orquestração e workers

## Feature 6: APIs e Integrações
- História 6.1: Definir contratos de API (OpenAPI/Swagger)
- História 6.2: Implementar handlers HTTP/gRPC
- História 6.3: Testes de integração das APIs

## Feature 7: Configuração, Scripts e Documentação
- História 7.1: Adicionar arquivos de configuração
- História 7.2: Criar scripts auxiliares (setup, seed, etc.)
- História 7.3: Atualizar README e documentação de uso

## Feature 8: Infraestrutura e Deploy
- História 8.1: Ajustar/implementar scripts de infraestrutura (se necessário)
- História 8.2: Preparar scripts de build/deploy
