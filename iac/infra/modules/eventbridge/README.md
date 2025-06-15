# Módulo EventBridge

Este módulo cria regras e targets do EventBridge.

- Cria a role de invoke para o target
- Permite customizar o padrão de evento e o target

## Principais variáveis
- `name`: Nome base
- `event_pattern`: Padrão do evento
- `target_arn`: ARN do target
- `invoke_role_arn`: Role para invocar o target

## Outputs
- `eventbridge_rule_name`, `eventbridge_rule_arn`

---

Veja o orquestrador para exemplos de uso e como passar a role de invoke.
