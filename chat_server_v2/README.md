# Chat Server com Dispatcher e Thread Pool

## Protocolo

- Cliente faz conexão TCP com servidor na porta 7777
- Cliente envia JSON com tipo de conexão
    - JSON para conexão que envia mensagens no chat.
    ```
    {
      "type: "TALK",
      "name": "username"
    }
    ```
    - JSON para conexão que recebe mensagens do chat a partir de determinado momento.
    ```
    {
      "type: "VIEW",
      "timestamp" : 394103910
    }
    ```
- Servidor responde com JSON de resposta.
    ```
    {
      "status": "OK"
    }
    ```
    ```
    {
      "status": "ERROR",
      "message": "Error message"
    }
    ```
- Para enviar mensagens, o cliente deve enviar:
    ```
    ```
-
