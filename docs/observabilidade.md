# Observabilidade no Projeto demo-signserver

## Visão Geral
O sistema implementa observabilidade ponta-a-ponta, com logs estruturados, métricas, tracing distribuído e rastreio de eventos via traceID.

## Itens Implementados
- **Logs estruturados**: Todos os handlers e serviços usam logs key-value, incluindo traceID, erro, duração e contexto.
- **Correlações**: traceID único por evento, propagado em todo o fluxo (EventBus, handlers, storage, etc).
- **Métricas**: Métricas de entrada, sucesso, erro e panic por handler, expostas via MetricsService.
- **Tracing distribuído**: Pronto para integração com OpenTelemetry (spans, parentSpanID).
- **Dashboards/Alertas**: Estrutura pronta para integração com Prometheus, Grafana, Datadog.
- **Auditoria**: Operações sensíveis logadas com contexto completo.

## Fluxo de Observabilidade
Consulte `docs/fluxogramas.md` para o ciclo de observabilidade dos eventos.

## Checklist de Observabilidade
- [x] Logging estruturado implementado
- [x] traceID propagado e logado
- [x] Métricas expostas e monitoradas
- [ ] Tracing distribuído ativo (OpenTelemetry)
- [ ] Dashboards e alertas configurados
- [x] Auditoria de operações sensíveis

## Ferramentas Sugeridas
- Logging: zap, logrus, zerolog
- Métricas: Prometheus, Grafana
- Tracing: OpenTelemetry, Jaeger
- Dashboards/Alertas: Grafana, Kibana, Datadog

---

> **Observação:** O EventBus centraliza a observabilidade dos fluxos assíncronos, garantindo rastreio e métricas automáticas em todos os eventos.
