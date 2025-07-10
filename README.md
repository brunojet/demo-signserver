# Projeto: demo-signserver

Este projeto implementa uma infraestrutura robusta para assinatura de APKs Android, com arquitetura modular, observável e altamente testável.

## Estrutura do Projeto

- **cmd/**: Ponto de entrada da aplicação (main.go)
- **internal/**: Lógica de domínio e orquestração (intent, signing, app, repository)
- **pkg/**: Pacotes reutilizáveis (eventbus, workerpool, storage, repository, notification)
- **docs/**: Documentação, diagramas e fluxogramas
- **iac/**: Infraestrutura como código (Terraform)

## Fluxo Resumido
1. Usuário solicita assinatura de APK (intent)
2. Recebe URL para upload
3. Faz upload do APK para S3
4. Evento S3 dispara orquestração
5. Orchestrator processa intent, aciona workerpool
6. Workerpool executa assinatura, salva resultado e notifica usuário

## Observabilidade
- Logs estruturados, métricas e tracing distribuído (OpenTelemetry)
- EventBus com observabilidade automática em todos os handlers
- Propagação de traceID ponta-a-ponta

## Documentação e Fluxogramas
Consulte `docs/arquitetura.md`, `docs/fluxogramas.md` e `docs/eventbus.md` para diagramas, exemplos e detalhes de cada módulo.

---

Para dúvidas ou sugestões, abra uma issue ou entre em contato com o mantenedor.