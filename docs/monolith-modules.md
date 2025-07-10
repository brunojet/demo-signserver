# Módulos e Fluxos — demo-signserver

## Visão Geral
O demo-signserver implementa um fluxo de assinatura de APKs Android, com arquitetura modular, observável e desacoplada. O núcleo da regra de negócio está em `internal/`, com orquestração baseada em eventos e handlers assíncronos.

## Módulos de Negócio (internal/)
- **request**: Recebe e valida pedidos de assinatura via API HTTP (handlers REST). Gera intents, valida perfis, retorna URL pré-assinada para upload.
- **signer**: Orquestra o fluxo de assinatura. Consome eventos de upload (S3), publica eventos no EventBus, executa handlers para processar assinatura e publicar resultados.
- **repository**: Implementa repositórios de domínio para persistência (DynamoDB, etc). Inclui RequestRepository, ProfileRepository, etc.
- **observability**: Middleware, tracing, logs estruturados e auditoria. Propaga traceID e requestID em todo o fluxo.
- **config**: Centraliza configuração e factories para adapters de storage, fila, etc.

## Estrutura de Diretórios

```
/demo-signserver
├── cmd/                # main.go (inicializa serviços e API)
├── internal/
│   ├── config/         # Factories, centralização de config
│   ├── observability/  # Tracing, logs, auditoria
│   ├── repository/     # Domínio, repositórios
│   ├── request/        # Handlers e serviços de intent/request
│   └── signer/         # Orquestração, handlers de assinatura
├── pkg/                # pacotes reutilizáveis (eventbus, workerpool, storage, etc)
├── docs/
├── iac/
└── README.md
```

## Fluxo de Trabalho (Regra de Negócio)
1. Usuário faz POST para criar intent de assinatura (request handler)
2. Serviço de request valida perfil, gera intent e retorna URL pré-assinada
3. Usuário faz upload do APK para S3 usando a URL
4. Evento S3 é capturado por um watcher/fila local (MessageQueueAdapter)
5. SignerService consome evento S3, publica no EventBus ("upload_received")
6. Handler de upload processa intent, busca dados, prepara para assinatura
7. Handler de assinatura executa lógica de assinatura (mock ou real), publica resultado
8. (Opcional) Notificação é enviada ao usuário

## Detalhes de Implementação
- Adapters de storage e fila são injetados via factories/configuração
- EventBus centraliza eventos assíncronos, com observabilidade automática e traceID
- Workerpool genérico para paralelismo seguro
- Mocks de infraestrutura em `pkg/*/mock` para testes
- Observabilidade: logs estruturados, tracing, métricas e auditoria em todos os handlers

## Referência
Consulte `docs/fluxogramas.md` e `docs/arquitetura.md` para diagramas e fluxos detalhados.
