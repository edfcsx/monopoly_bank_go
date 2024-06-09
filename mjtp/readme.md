# MJTP - Message JSON Transfer Protocol

O MJTP é um protocolo que transfere mensagens de dados JSON.

Aqui está um exemplo de como uma mensagem MJTP pode ser estruturada:

```
/index MJTP/1.0 {"var":"foo"}\r\n\r\n
```
Neste exemplo, `/index` é o recurso que está sendo solicitado, `MJTP/1.0` indica que estamos usando a versão 1.0 do protocolo MJTP, e `{"var":"foo"}` é a mensagem JSON que está sendo transferida e \r\n\r\n indica o final da mensagem.

## Detalhes

O MJTP é projetado para ser simples e eficiente, permitindo a transferência de mensagens JSON de uma maneira estruturada e padronizada. Ele usa a estrutura de linha única para separar diferentes partes da mensagem, tornando-o fácil de analisar e processar.

A mensagem JSON é enviada como uma string, permitindo que qualquer tipo de dados JSON seja transferido. Isso inclui objetos JSON, arrays, strings, números, booleanos e null.

O MJTP é ideal para situações onde você precisa transferir dados JSON de forma eficiente e estruturada, como em aplicações web, APIs, e em qualquer lugar onde os dados JSON precisam ser transferidos entre diferentes partes de um sistema.


O MJTP é um protocolo que envia apenas mensagens, sem fornecer indicações de retorno de que a mensagem foi bem recebida. Por isso, ele foi projetado para funcionar em cima de um protocolo que já estabeleça essa camada de conexão entre cliente e servidor, como o protocolo WebSocket.