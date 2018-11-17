#!/bin/bash

docker run -d --hostname rabbitmq --name rabbitmq -e RABBITMQ_ERLANG_COOKIE='cookie' -e RABBITMQ_DEFAULT_USER=admin -e RABBITMQ_DEFAULT_PASS=password  -e RABBITMQ_HIPE_COMPILE=1 -p 8080:15672 -p 5672:5672 rabbitmq:3-management

# then run 
#    rmq publish -a amqp://admin:password@localhost:5672 -q myQueue
# and
#    rmq subscribe -a amqp://admin:password@localhost:5672 -q myQueue

