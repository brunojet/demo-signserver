# Módulo Lambda

Este módulo cria funções Lambda e sua role de execução.

- Cria a role e policy básica de logs
- Permite receber uma role externa (opcional)
- Gera um código Hello World inline se não for informado um arquivo

## Principais variáveis
- `name`: Nome base da função
- `role_arn`: (opcional) Role externa

## Outputs
- `function_arn`, `role_arn`, `role_name`

---

Veja o orquestrador para exemplos de uso e como anexar policies adicionais.
