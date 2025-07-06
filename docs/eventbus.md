# EventBus: Arquitetura e Documentação

## Objetivo
O `EventBus` é um barramento de eventos local, thread-safe, que permite registrar handlers para diferentes tipos de eventos. Para cada tipo de evento registrado, uma goroutine dedicada é criada para consumir eventos da fila e executar o handler correspondente. O EventBus é desacoplado do mecanismo de fila, usando um adapter (interface) para a fila.

## Estrutura

```go
// EventBus mantém o controle dos handlers e goroutines por tipo de evento.
type EventBus struct {
    queueAdapter QueueInterface
    handlers    map[string]Handler
    stopChans   map[string]chan struct{}
    mu          sync.Mutex
}
```

- `queueAdapter`: Adapter para a fila de mensagens (interface).
- `handlers`: Mapeia eventType para handler.
- `stopChans`: Mapeia eventType para canal de parada da goroutine.
- `mu`: Mutex para garantir thread-safety.

## Métodos

### Register
```go
func (b *EventBus) Register(eventType string, handler Handler)
```
Registra um handler para um tipo de evento. Cria uma goroutine dedicada que consome eventos desse tipo da fila e executa o handler.

### Unregister
```go
func (b *EventBus) Unregister(eventType string)
```
Remove o handler e para a goroutine associada ao tipo de evento.

### Stop
```go
func (b *EventBus) Stop()
```
Para todas as goroutines e encerra o consumo de eventos.

## Exemplo de Uso

```go
bus := NewEventBus(queue)
bus.Register("upload", uploadHandler)
// ...
bus.Unregister("upload")
bus.Stop()
```

## Vantagens
- Isolamento por tipo de evento (cada um com sua goroutine)
- Fácil de registrar/desregistrar handlers dinamicamente
- Desacoplado do mecanismo de fila (pode ser local, SQS, etc)
- Thread-safe

## Possíveis Extensões
- Suporte a middlewares
- Retry/backoff por tipo de evento
- Métricas e logging por handler
