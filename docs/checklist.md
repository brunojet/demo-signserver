# Checklist de Desenvolvimento — demo-signserver

## 1. Inicialização e Base
- [x] Inicializar projeto Go (`go mod init demo-signserver`)
- [x] Criar estrutura inicial de diretórios conforme documentação
- [x] Implementar `cmd/main.go` com Hello World

## 2. Implementação dos Pacotes Genéricos (`pkg/`)
- [ ] Definir e implementar as interfaces/contratos em:
  - [ ] `pkg/shared/domain/` (entidades centrais)
  - [ ] `pkg/repository/services/` (implementações de repositório)
  - [ ] `pkg/repository/mock/` (mocks de repositório para testes)
  - [ ] `pkg/storage/services/` (implementações de storage)
  - [ ] `pkg/storage/mock/` (mocks de storage para testes)
  - [ ] `pkg/notification/domain/` e `service/` (notificação genérica)
- [ ] Implementar os adapters genéricos:
  - [ ] `pkg/repository/services/dynamodb_service.go` (e outros, se necessário)
  - [ ] `pkg/storage/services/s3_service.go`
  - [ ] `pkg/notification/adapters/webhook.go` (e outros canais)
- [ ] Implementar o pool de workers em `pkg/workerpool/service/workerpool.go`

## 3. Testes Unitários dos Pacotes Genéricos
- [ ] Escrever testes unitários para cada pacote em `pkg/`
  - [ ] Testes para entidades/valores
  - [ ] Testes para interfaces (mocks/fakes em `mock/`)
  - [ ] Testes para adapters (usando mocks ou recursos locais)
  - [ ] Testes para workerpool

## 4. Módulos de Domínio (internal/)
- [ ] Implementar camada de aplicação de cada contexto:
  - [ ] `internal/intent/application/intent_service.go`
  - [ ] `internal/signing/application/signing_service.go`
  - [ ] `internal/notification` (se necessário, para lógica específica)
- [ ] Implementar interfaces (handlers, etc.) de cada contexto

## 5. Orquestração e Fluxos
- [ ] Implementar orquestrador em `internal/app/orchestrator.go`
- [ ] Implementar workers e canais (`upload_worker.go`, `download_worker.go`, `channels.go`, etc.)
- [ ] Integrar módulos de domínio com os pacotes genéricos (`pkg/`)

## 6. Testes de Integração e E2E
- [ ] Escrever testes de integração para fluxos principais
- [ ] Escrever testes end-to-end simulando o ciclo completo (intent → upload → assinatura → notificação)

## 7. Configuração, Scripts e Documentação
- [ ] Adicionar arquivos de configuração em `configs/`
- [ ] Criar scripts auxiliares em `scripts/` (ex: setup, seed, etc.)
- [ ] Atualizar e detalhar o `README.md` com instruções de uso, build e testes

## 8. APIs e Integrações
- [ ] Definir contratos de API em `api/` (OpenAPI/Swagger, protos, etc.)
- [ ] Implementar handlers HTTP/gRPC conforme necessário

## 9. Infraestrutura e Deploy
- [ ] Ajustar/implementar infraestrutura em `iac/` (Terraform, etc.)
- [ ] Preparar scripts de build/deploy
