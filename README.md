# RabbitMQ Listener

`rmq` is a simple RabbitMQ listener implemented as a CLI; it connects to the given RabbitMQ server and `vhost` and listens for messages, printing them and their metadata to STDOUT.

# Usage

`rmq` accepts the following parameters:

```bash
$> rmq --url amqp://username:password@server:5672/vhost --queue myQueue
