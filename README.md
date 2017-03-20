[![Build Status](https://travis-ci.org/neckhair/owntracks-eventr.svg?branch=master)](https://travis-ci.org/neckhair/owntracks-eventr)

# Owntracks Eventr

Listens for Owntracks events on MQTT and writes them into a file for further processing.

## Test on local machine

    # run a local pre-configured MQTT server in Docker
    docker-compose up

    # compile and run the daemon
    export MQTT_PASSWORD=secretpassword
    go install
    owntracks-eventr listen --ca-cert docker/ca.crt --username eventr

## Knowhow

### How do I create CA and server certificates?

It's all described in Mosquittos' man page: https://mosquitto.org/man/mosquitto-tls-7.html

    # Generate CA Certificate (with -nodes for password-less key)
    openssl req -new -x509 -extensions v3_ca -nodes -keyout ca.key -out ca.cr

    # Generate un-encrypted server key
    openssl genrsa -out server.key 2048

    # Generate CSR
    openssl req -out server.csr -key server.key -new

    # Sign certificate
    openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt
