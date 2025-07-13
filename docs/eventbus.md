# EventBus: Arquitetura e Documentação

## Objetivo
O `EventBus` é um barramento de eventos thread-safe, observável e robusto, com worker pool, shutdown seguro, rastreio ponta-a-ponta (traceID) e integração com métricas/logs customizados.

## Estrutura

```go
// EventBus mantém o controle dos handlers e worker pools por tipo de evento.
type EventBus struct {
    workers    map[HandlerName]*eventWorker
    ObsHandler ObservableHandler
    // ...
}
```

- **workers**: Mapeia eventType para worker pool.
- **ObsHandler**: Handler de observabilidade (logs, métricas, tracing).

## Principais Métodos

- `Register(eventType, handler, numWorkers, backlog)`
- `Unregister(eventType)`
- `Publish(eventType, data)`
- `PublishWithContext(ctx, eventType, data)`
- `Stop()`

## Observabilidade
- Todos os handlers são automaticamente envolvidos por middleware de observabilidade.
- Cada evento recebe um traceID único propagado em todo o fluxo.
- Logs estruturados, métricas e tracing distribuído (OpenTelemetry-ready).

## Exemplo de Uso
```go
bus := eventbus.NewEventBus()
bus.Register("upload_received", handler, 2, 10)
bus.Publish("upload_received", payload)
bus.Unregister("upload_received")
bus.Stop()
```

## Fluxograma do ciclo de eventos
Consulte `docs/fluxogramas.md` para o ciclo de observabilidade do EventBus.

## Vantagens
- Isolamento por tipo de evento (worker pool)
- Observabilidade automática
- Fácil integração com tracing e métricas
- Thread-safe e robusto

---
Consulte os testes unitários para exemplos avançados e integração com métricas/logs.
