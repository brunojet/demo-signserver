
# Módulos e Fluxos — demo-signserver

## Visão Geral
O demo-signserver implementa um fluxo de assinatura de APKs Android, com arquitetura modular, observável e altamente desacoplada. O núcleo da regra de negócio está em `internal/`, com orquestração baseada em eventos, handlers assíncronos e integração com serviços externos.

## Módulos de Negócio (internal/)
- **request**: Recebe e valida pedidos de assinatura via API HTTP (handlers REST). Gera intents, valida perfis, retorna URL pré-assinada para upload. Inclui serviços e DTOs para requests e perfis.
- **orquestrator**: Orquestra o fluxo de eventos do sistema. Consome eventos de upload (S3), publica eventos no EventBus, executa handlers para download, assinatura e upload. Centraliza o fluxo assíncrono e paralelismo via workerpool.
- **ext_signer**: Integra com assinadores externos (ex: fabricantes POS), via adapters. Permite assinatura real ou mock de APKs.
- **repository**: Implementa repositórios de domínio para persistência (DynamoDB, etc). Inclui RequestRepository, ProfileRepository, objetos de domínio e mocks.
- **observability**: Middleware, tracing, logs estruturados, métricas e auditoria. Propaga traceID e requestID em todo o fluxo, integra com EventBus e workerpool.
- **config**: Centraliza configuração e factories para adapters de storage, fila, eventbus, etc.

## Estrutura de Diretórios

```
/demo-signserver
├── cmd/                # main.go (inicializa serviços e API)
├── internal/
│   ├── config/         # Factories, centralização de config
│   ├── observability/  # Tracing, logs, auditoria
│   ├── repository/     # Domínio, repositórios, mocks
│   ├── request/        # Handlers, serviços, DTOs de intent/request
│   ├── orquestrator/   # Orquestração, handlers de download/upload/sign
│   └── ext_signer/     # Integração com assinadores externos/adapters
├── pkg/
│   ├── eventbus/       # EventBus desacoplado, observabilidade
│   ├── workerpool/     # Paralelismo seguro, métricas
│   ├── storage/        # Adapters S3/local, abstração de storage
│   ├── observability/  # Logs, métricas, tracing
│   └── ...             # Outros pacotes reutilizáveis
├── docs/
├── iac/
└── README.md
```

## Fluxo de Trabalho (Regra de Negócio)
1. Usuário faz POST para criar intent de assinatura (request handler)
2. Serviço de request valida perfil, gera intent e retorna URL pré-assinada
3. Usuário faz upload do APK para S3 usando a URL
4. Evento S3 é capturado por um watcher/fila local (MessageQueueAdapter)
5. Orquestrator consome evento S3, publica no EventBus ("storage_download", "sign_process", "storage_upload")
6. Handlers de download/upload/sign processam intent, buscam dados, executam assinatura (mock ou real via ext_signer)
7. EventBus gerencia handlers assíncronos, paralelismo via workerpool, rastreabilidade e métricas
8. Observabilidade registra logs, métricas, tracing e auditoria em todos os handlers
9. (Opcional) Notificação é enviada ao usuário

## Detalhes de Implementação
- Adapters de storage (S3/local) e fila são injetados via factories/configuração
- EventBus centraliza eventos assíncronos, com observabilidade automática, traceID e integração com workerpool
- Workerpool genérico para paralelismo seguro, métricas e exportação periódica
- Integração com assinadores externos via ext_signer/adapters
- Mocks de infraestrutura em `pkg/*/mock` para testes
- Observabilidade: logs estruturados, tracing, métricas e auditoria em todos os handlers e eventos

## Referência
Consulte `docs/fluxogramas.md` e `docs/arquitetura.md` para diagramas e fluxos detalhados.
