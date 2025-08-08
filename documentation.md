# Documentación para entender gRCP

En gRPC, hay dos lados:

- Servidor gRPC: implementa las funciones (por ejemplo, la lógica de sumar dos números).
- Cliente gRPC: es quien llama a esas funciones.

```yaml
Cliente gRPC
└── Llama a las funciones del servidor como si fueran locales

Servidor gRPC
└── Implementa las funciones y responde al cliente
```

## Archivo `.proto`

Es como un manual de instrucciones que dice:

- Qué funciones puede hacer el servidor.
- Qué información necesita cada función.
- Qué información devuelve cada función.

Como el `cliente` y el `servidor` van a hablar.
Contrato `cliente` y `servidor`.

- En el servidor, sí: debes implementar todas las funciones definidas en el `.proto`.
- En el cliente, solo necesitas consumir esas funciones; el código se genera a partir del `.proto`.

## Servicios

Define las `funciones` que puede hacer el servidor.

## Funciones

```protobuf
rpc CreateTask(CreateTaskRequest) returns (TaskResponse);
     ↑              ↑                       ↑
   nombre        qué envías            qué recibes
```

## Mensajes

Podemos verlos como estructuras de datos u objetos que definen la forma que tendran los datos que se van a enviar o recibir.

- El mensaje principal no suele tener ningun prefijo o sufijo en su nombre.
Normalmente llamamos a los mensajes con sufijos como `Request` (lo que enviamos al servidor) o `Response` (lo que el servidor te devuelve).

## Tags o field number

Son identificadores binarios unicos para cada campo.

- Se usan internamente por Protobuf para serializar datos.
- Son obligatorios y permanentes: nunca deben cambiarse ni reutilizarse.

## Tipo de comunicación Unary

1. Unary (uno a uno):

- Cliente envía una solicitud, servidor responde una vez.
