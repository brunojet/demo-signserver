## Contexto

Após uma avaliação do projeto, decidimos que a melhor abordagem seria manter uma estrutura monolítica para o serviço de assinatura de APKs. Essa decisão foi baseada na simplicidade e facilidade de manutenção que um monólito oferece, especialmente em estágios iniciais de desenvolvimento.

## Módulos

A estrutura monolítica será dividida em módulos, cada um responsável por uma parte específica do serviço. Os módulos identificados até agora são:

- **Módulo de Intenção**: Responsável pelo registro do pedido de assinatura, encaminhamento do perfil de assinatura e retorno da URL assinada para upload do aplicativo.
- **Módulo de Assinatura**: Responsável pelo processo de assinatura de APKs, desde o recebimento do evento do S3 até o download do APK para processamento da assinatura e posterior upload do APK assinado de volta ao S3.
- **Módulo de Notificações**: Responsável pelo envio de notificações aos usuários, como webhooks ou mensagens de status via SNS (Simple Notification Service).

Cada módulo será desenvolvido de forma independente, mas todos estarão integrados na mesma base de código.

## Estrutura do Projeto

A estrutura de diretórios foi atualizada para refletir a separação entre código de produção e mocks de teste, seguindo o padrão:

- Código de produção: `pkg/repository/services`, `pkg/storage/services`, etc.
- Mocks de teste: `pkg/repository/mock`, `pkg/storage/mock`, etc.

Exemplo de estrutura:

```text
/demo-signserver
├── cmd/                        # Ponto de entrada da aplicação (main.go)
│   └── main.go
├── internal/                   # Código interno da aplicação (não exportado)
│   ├── app/                    # Orquestração de fluxos, channels, etc.
│   ├── intent/                 # Contexto de Intenção
│   ├── signing/                # Contexto de Assinatura
│   └── repository/             # Interfaces de domínio (contratos)
│       └── domain/             # Interfaces dos repositórios
├── pkg/                        # Código reutilizável/exportável por outros projetos
│   ├── repository/
│   │   ├── services/           # Implementações concretas dos repositórios (DynamoDB, etc.)
│   │   ├── mock/               # Mocks para testes (ex: MockDynamoDBClient)
│   ├── storage/
│   │   ├── services/           # Serviços de storage (S3, STS, etc.)
│   │   ├── mock/               # Mocks de storage para testes
│   └── ...
├── api/                        # Definições de APIs (OpenAPI/Swagger, protos, etc)
├── configs/                    # Arquivos de configuração
├── scripts/                    # Scripts auxiliares
├── go.mod
├── go.sum
└── README.md
```

- O diretório `pkg/` contém código que pode ser importado por outros projetos Go, facilitando a reutilização caso o monólito evolua para microsserviços ou bibliotecas compartilhadas.
- Os mocks de infraestrutura (ex: MockDynamoDBClient, MockPresignClient) ficam em `pkg/repository/mock` e `pkg/storage/mock`, garantindo que apenas código de teste dependa deles.
- O domínio e as interfaces de repositório continuam em `internal/repository/domain`, mantendo o core do domínio desacoplado da infraestrutura.
- O diretório `internal/app/` centraliza a orquestração de fluxos que envolvem múltiplos contextos, como workers, channels e coordenação entre assinatura e notificação.

Essa abordagem prepara o projeto para uma possível evolução futura, mantendo a estrutura alinhada ao DDD, à filosofia Go de reutilização de código e às melhores práticas de isolamento de testes.

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
