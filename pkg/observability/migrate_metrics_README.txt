MIGRAÇÃO DE MÉTRICAS:

- O conteúdo de internal/observability/metrics.go foi migrado para pkg/observability/metrics.go.
- Use agora o serviço centralizado MetricsService para registrar e exportar métricas.
- Funções antigas como IncRequestCount devem ser adaptadas para usar MetricsService.Inc("request_count", ...).
- Remova ou arquive o arquivo antigo após atualizar todos os usos no projeto.
