# Observabilidade no Projeto demo-signserver

## Objetivo

Garantir que o sistema seja monitorável, rastreável e auditável, facilitando a detecção de problemas, análise de performance e auditoria de operações.

## Itens a serem implementados

### 1. Logs Estruturados
- Utilizar logs estruturados (JSON ou key-value) em todos os handlers, serviços e pontos críticos do sistema.
- Incluir informações relevantes: timestamp, nível (info, warn, error), requestID, usuário (se aplicável), payloads, erros e stacktrace.
- Padronizar o formato dos logs para facilitar integração com ferramentas de análise (ex: ELK, Datadog, CloudWatch).

### 2. Correlação de Requisições
- Gerar e propagar um identificador único de requisição (requestID) em todo o ciclo de vida da request HTTP.
- Incluir o requestID nos logs e nas respostas HTTP (header).
- Permitir rastrear uma requisição ponta-a-ponta.

### 3. Métricas
- Expor métricas de negócio e técnicas via endpoint /metrics (Prometheus).
- Métricas sugeridas:
  - Total de requisições por endpoint
  - Latência média e percentis
  - Taxa de erro por endpoint
  - Métricas de uso de recursos (memória, CPU, goroutines)

### 4. Tracing Distribuído
- Integrar tracing distribuído (OpenTelemetry ou Jaeger) para rastrear chamadas entre serviços e dependências externas.
- Instrumentar handlers, serviços e integrações externas (DynamoDB, S3, etc).
- Permitir análise de gargalos e dependências.

### 5. Alertas e Dashboards
- Configurar alertas automáticos para erros críticos, alta latência ou indisponibilidade.
- Criar dashboards para visualização de métricas e logs em tempo real.

### 6. Auditoria
- Registrar operações sensíveis (criação, alteração, deleção) com informações de quem fez, quando e o que foi alterado.
- Garantir rastreabilidade para fins de compliance.

## Ferramentas Sugeridas
- Logging: zap, logrus, zerolog
- Métricas: Prometheus, Grafana
- Tracing: OpenTelemetry, Jaeger
- Dashboards/Alertas: Grafana, Kibana, Datadog

## Plano de Execução
1. Definir padrão de logs e implementar logging estruturado.
2. Adicionar requestID e correlação em todo o fluxo.
3. Instrumentar métricas e expor endpoint /metrics.
4. Integrar tracing distribuído.
5. Configurar dashboards e alertas.
6. Implementar auditoria de operações sensíveis.

## Checklist de Observabilidade
- [ ] Logging estruturado implementado
- [ ] requestID propagado e logado
- [ ] Métricas expostas e monitoradas
- [ ] Tracing distribuído ativo
- [ ] Dashboards e alertas configurados
- [ ] Auditoria de operações sensíveis

---

Este documento deve ser revisado e atualizado conforme a evolução do projeto e das necessidades de monitoramento.
