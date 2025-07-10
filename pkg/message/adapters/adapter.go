package adapters

// MessageQueueAdapterInterface é uma interface genérica para consumidores de eventos (ex: SQS, local, etc)
type MessageQueueAdapterInterface interface {
	// Start inicia o consumo de eventos e chama o callback para cada mensagem recebida
	Start(onMessage func(event any)) (stop func())
}
