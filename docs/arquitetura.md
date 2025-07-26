# Arquitetura Atual do demo-signserver

Este documento descreve a arquitetura do projeto, fluxos principais e integrações.


## Diagrama de Arquitetura

```mermaid
graph TD
    subgraph Infraestrutura
        S3[S3 Bucket]
        EventBus[EventBus: Eventos Internos]
        Orquestrator[Orquestrator: Orquestração]
        WorkerPool[WorkerPool: Paralelismo]
        DynamoDB[DynamoDB: Logs/Resultados]
        POS[Fabricante POS: Assinatura APK]
        S3 --> Orquestrator
        Orquestrator --> EventBus
        EventBus --> WorkerPool
        WorkerPool --> POS
        WorkerPool --> DynamoDB
        WorkerPool --> S3
    end
```


## Diagrama de Sequência do Fluxo

```mermaid
sequenceDiagram
    participant Usuário
    participant S3
    participant Orquestrator
    participant EventBus
    participant WorkerPool
    participant POS as Fabricante POS
    participant DynamoDB

    Usuário->>S3: Upload de APK
    S3-->>Orquestrator: Evento de novo arquivo
    Orquestrator-->>EventBus: Publica eventos ("download", "sign", "upload")
    EventBus-->>WorkerPool: Executa handlers paralelos
    WorkerPool-->>POS: Solicita assinatura do APK
    POS-->>WorkerPool: Retorna APK assinado
    WorkerPool-->>DynamoDB: Salva resultado/log
    WorkerPool-->>S3: Salva APK assinado
    WorkerPool-->>Usuário: Notifica conclusão
```

## Fluxograma Geral

Consulte `docs/fluxogramas.md` para fluxogramas detalhados do fluxo ponta-a-ponta e observabilidade.

---

> **Observação:** O fluxo inicia após a criação de um arquivo no S3, disparando todo o processo de assinatura e orquestração descrito acima.
