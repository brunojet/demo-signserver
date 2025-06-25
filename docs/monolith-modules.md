## Contexto

Após uma avaliação do projeto, decidimos que a melhor abordagem seria manter uma estrutura monolítica para o serviço de assinatura de APKs. Essa decisão foi baseada na simplicidade e facilidade de manutenção que um monólito oferece, especialmente em estágios iniciais de desenvolvimento.

## Módulos

A estrutura monolítica será dividida em módulos, cada um responsável por uma parte específica do serviço. Os módulos identificados até agora são:

- **Módulo de Intenção**: Responsável pela registro do pedido de assinatura, encaminhamento do perfil de assinatura e retorno da URL assinada para upload do aplicativo.
- **Módulo de Assinatura**: Responsável pelo processo de assinatura de APKs, desde o recebimento do evento do S3 até o download do APK para processamento da assinatura e posterior upload do APK assinado de volta ao S3.
- **Módulo de Notificações**: Responsável pelo envio de notificações aos usuários, como webhooks ou mensagens de status via SNS (Simple Notification Service).

Cada módulo será desenvolvido de forma independente, mas todos estarão integrados na mesma base de código.

## Estrutura do Projeto

Abaixo está um exemplo de como pode ser organizada a estrutura de diretórios e arquivos para um projeto monolítico em Go seguindo o padrão DDD (Domain-Driven Design), considerando que os itens compartilhados (modelos e repositórios) ficam em `pkg/` para permitir reutilização futura por outros projetos:

```text
/demo-signserver
├── cmd/                        # Ponto de entrada da aplicação (main.go)
│   └── main.go
├── internal/                   # Código interno da aplicação (não exportado)
│   ├── app/                    # Orquestração de fluxos, channels, etc.
│   │   ├── orchestrator.go         # Coordena o fluxo principal do sistema
│   │   ├── upload_worker.go        # Worker dedicado ao upload do APK assinado
│   │   ├── download_worker.go      # Worker dedicado ao download do APK do S3 para disco
│   │   ├── event.go                # Definição de eventos internos (structs para assinatura, upload, notificação)
│   │   └── channels.go             # Definição e inicialização dos canais usados entre os workers
│   ├── intent/                 # Contexto de Intenção
│   │   ├── application/
│   │   │   └── intent_service.go   # Casos de uso: criar intent, gerar URL, atualizar status
│   │   └── interfaces/
│   │       └── http_handler.go     # Handler HTTP para pedidos de intenção
│   ├── signing/                # Contexto de Assinatura
│   │   ├── application/
│   │   │   └── signing_service.go   # Lógica de assinatura, integração com workerpool
│   │   ├── adapters/                # Adaptadores para diferentes assinadores/fabricantes
│   │   │   ├── signing_adapter_positivo.go   # Integração com assinador Positivo
│   │   │   ├── signing_adapter_fabricanteX.go# Integração com outro fabricante
│   │   │   └── ...                          # Outros adaptadores conforme necessário
│   │   └── interfaces/
│   │       ├── worker_handler.go        # Handler para eventos do workerpool
│   │       └── signing_adapter_interface.go # Interface que define o contrato dos adaptadores de assinadores
├── pkg/                        # Código reutilizável/exportável por outros projetos
│   ├── repository/             # Repositórios compartilhados
│   │   ├── domain/                 # Interfaces dos repositórios (contratos)
│   │   │   ├── device_profile.go
│   │   │   └── intent.go
│   │   ├── service/                # Interfaces/contratos de serviços auxiliares
│   │   │   ├── device_profile.go
│   │   │   └── intent.go
│   │   ├── adapters/               # Implementações concretas dos repositórios
│   │   |   └── dynamodb.go         # Adapter para DynamoDB (pode haver outros, ex: postgres.go)
│   |   └── interfaces/             # Interfaces públicas do serviço
│   |       └── adapters.go
│   ├── storage/                    # Abstração para acesso a storages
│   │   ├── service/                # Interface/contrato para operações de storage
│   │   │   └── storage.go
│   │   ├── adapters/               # Implementações concretas dos storages
│   │   |    └── s3.go               # Integração com AWS S3
│   |   └── interfaces/             # Interfaces públicas do serviço
│   |       └── adapters.go
│   ├── workerpool/                 # Implementação genérica do pool de workers
│   │   └── service/
│   │       └── workerpool.go
│   └── notification/               # Serviço de notificação genérico e reutilizável
│       ├── domain/                 # Entidades/valores de notificação genéricos
│       │   └── notification.go
│       ├── service/                # Interface/contrato para operações de notificação
│       │   └── notification.go
│       ├── adapters/               # Adapters para SNS, webhooks, etc.
│       │   └── webhook.go
│       ├── interfaces/             # Interfaces públicas do serviço
│       │   └── adapters.go
│       └── notification_service.go # Serviço que recebe objeto abstrato e parâmetros de envio
├── api/                        # Definições de APIs (OpenAPI/Swagger, protos, etc)
├── configs/                    # Arquivos de configuração
├── scripts/                    # Scripts auxiliares
├── go.mod
├── go.sum
└── README.md
```

- O diretório `pkg/` contém código que pode ser importado por outros projetos Go, facilitando a reutilização caso o monólito evolua para microsserviços ou bibliotecas compartilhadas.
- As entidades de domínio compartilhadas entre dois ou mais módulos ficam em `pkg/shared/domain/`, deixando claro que são centrais e reutilizáveis no projeto.
- Os repositórios compartilhados ficam em `pkg/repository/`, separados em `domain/` (interfaces/contratos) e `adapters/` (implementações concretas, como DynamoDB, Postgres, etc). Isso garante desacoplamento e facilita testes e manutenção.
- Os módulos de domínio (`intent`, `signing`, `notification`) podem importar esses pacotes normalmente.
- O diretório `internal/app/` centraliza a orquestração de fluxos que envolvem múltiplos contextos, como workers, channels e coordenação entre assinatura e notificação.
- O arquivo `orchestrator.go` pode conter a lógica que utiliza intent, signing e notification, mantendo cada domínio isolado e a lógica de coordenação separada.

Essa abordagem prepara o projeto para uma possível evolução futura, mantendo a estrutura alinhada ao DDD e à filosofia Go de reutilização de código.

## Fluxo de trabalho

1. O usuário envia um pedido de assinatura através do módulo de Intenção, informando o perfil de assinatura desejado.
2. O módulo de Intenção gera uma URL pré-assinada para upload do APK e cria um registro na tabela intent com status "pendente". O nome do APK é um UUID, que também é a chave primária do registro intent.
3. O usuário faz o upload do APK para o S3 usando a URL fornecida.
4. O S3 dispara um evento, que é capturado pelo orchestrator (internal/app/orchestrator.go).
5. O orchestrator localiza o registro intent pelo UUID, baixa o APK localmente e envia um evento para o workerpool de assinatura. O download não é paralelizado, mas cada intent é processada em paralelo pelo workerpool.
6. O workerpool de assinatura recebe o perfil e o local do APK, seleciona o assinador correspondente e executa o processo de assinatura. O processo pode envolver upload/download em polling até obter o APK assinado, que é salvo em disco.
7. O workerpool envia um evento para o workerpool de upload, que faz o upload do APK assinado para o S3 e atualiza o registro intent para "concluído", incluindo a URL do APK assinado.
8. Por fim, o módulo de Notificações envia uma mensagem de status ao usuário, informando que o APK foi assinado e está disponível para download.

Essa estrutura e fluxo garantem clareza, desacoplamento e escalabilidade, facilitando a manutenção e evolução do sistema.

- O diretório `adapters/` dentro de `internal/signing` centraliza os adaptadores para diferentes assinadores/fabricantes, como `signing_adapter_positivo.go`, `signing_adapter_fabricanteX.go` etc. Isso deixa explícito o papel de cada implementação e facilita a manutenção e expansão do sistema.

- O diretório `pkg/workerpool/` contém a implementação genérica do pool de workers, podendo ser reutilizado por diferentes módulos, contextos ou até outros projetos Go.

- O arquivo `download_worker.go` em `internal/app/` é responsável por controlar o fluxo de downloads do S3 para disco, permitindo limitar a concorrência e centralizar a lógica de gravação local dos APKs. Isso facilita o controle, a manutenção e a evolução do fluxo de downloads.

- O arquivo `signing_adapter_interface.go` em `internal/signing/interfaces/` define a interface que todos os adaptadores de assinadores devem implementar, garantindo padronização e desacoplamento entre a lógica de assinatura e as integrações com diferentes fabricantes.

- O diretório `internal/storage/` foi removido, pois toda a lógica de acesso a storage (S3, etc.) está centralizada em `pkg/storage`, tornando o código reutilizável e evitando duplicidade.

- O serviço de notificação genérico foi movido para `pkg/notification/`, incluindo suas camadas domain, infrastructure, interfaces e service, seguindo o padrão DDD e facilitando a reutilização.
