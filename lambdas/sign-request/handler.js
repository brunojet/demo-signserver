// Lambda: sign-request
// Responsável por receber a intenção de assinatura e encaminhar o perfil de assinatura

exports.handler = async (event) => {
    // TODO: Implementar lógica de recebimento de intenção de assinatura
    // e envio do perfil de assinatura
    return {
        statusCode: 200,
        body: JSON.stringify({ message: 'Intenção de assinatura recebida.' })
    };
};
