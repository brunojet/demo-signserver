# Arquitetura Atual do demo-signserver

Este documento descreve a arquitetura do projeto, fluxos principais e integrações.

## Diagrama de Arquitetura AWS

```mermaid
graph TD
    subgraph AWS
        S3[S3 Bucket]
        Lambda[Lambda: Processa Pedido]
        StepFn[Step Functions: Orquestração]
        EventBridge[EventBridge: Eventos]
        DynamoDB[DynamoDB: Logs/Resultados]
        POS[Fabricante POS: Assinatura APK]
        S3 --> EventBridge
        EventBridge --> StepFn
        StepFn --> Lambda
        Lambda --> POS
        Lambda --> DynamoDB
        StepFn --> S3
    end
```

## Diagrama de Sequência do Fluxo

```mermaid
sequenceDiagram
    participant Usuário
    participant S3
    participant EventBridge
    participant StepFn as Step Functions
    participant Lambda
    participant POS as Fabricante POS
    participant DynamoDB

    Usuário->>S3: Upload de APK
    S3-->>EventBridge: Evento de novo arquivo
    EventBridge-->>StepFn: Inicia orquestração
    StepFn-->>Lambda: Executa função de assinatura
    Lambda-->>POS: Solicita assinatura do APK
    POS-->>Lambda: Retorna APK assinado
    Lambda-->>DynamoDB: Salva resultado/log
    Lambda-->>S3: Salva APK assinado
    StepFn-->>Usuário: Notifica conclusão
```

## Fluxograma Geral

Consulte `docs/fluxogramas.md` para fluxogramas detalhados do fluxo ponta-a-ponta e observabilidade.

---

> **Observação:** O fluxo inicia após a criação de um arquivo no S3, disparando todo o processo de assinatura e orquestração descrito acima.
