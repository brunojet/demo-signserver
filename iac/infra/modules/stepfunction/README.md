# Módulo Step Function

Este módulo cria uma Step Function e sua role de execução.

- Cria a role de execução com trust policy mínima
- Permite receber uma role externa (opcional)

## Principais variáveis
- `name`: Nome base
- `role_arn`: (opcional) Role externa
- `definition`: Definição do workflow

## Outputs
- `state_machine_arn`, `stepfunction_exec_role_arn`, `stepfunction_exec_role_name`

---

Veja o orquestrador para exemplos de uso e como anexar policies adicionais.
