# Infraestrutura de Persistência - S3

Este módulo gerencia o bucket S3 para persistência do projeto, com suporte a múltiplos ambientes (dev, qa, prod).

## Como usar

1. Defina o ambiente desejado via variável `environment` (ex: dev, qa, prod).
2. O bucket será criado com versionamento ativado e proteção contra deleção acidental.
3. Recomenda-se usar um backend remoto para o state (ex: S3 + DynamoDB).

## Exemplo de uso

```hcl
module "persistencia_s3" {
  source     = "./infra/persistencia"
  environment = "dev"
  project     = "go-signserver"
}
```

## Variáveis principais
- `environment`: Ambiente (dev, qa, prod)
- `project`: Nome do projeto

## Proteções
- Versionamento ativado
- `prevent_destroy` habilitado

---

> **Atenção:** Nunca remova a proteção de deleção sem revisão!
