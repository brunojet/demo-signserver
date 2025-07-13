# Módulo signer

Fluxo de assinatura desacoplado via EventBus:

1. **upload_received**: Recebe evento de upload, valida e publica evento de assinatura.
2. **sign_process**: Processa assinatura e publica evento de upload assinado.
3. **signed_upload**: Faz upload do arquivo assinado para o storage.

Cada handler é registrado no EventBus e pode ser escalado/testado separadamente.

Consulte os arquivos `*_handler.go` para detalhes e implemente a lógica de negócio conforme necessário.
