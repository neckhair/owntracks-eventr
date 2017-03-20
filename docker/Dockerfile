FROM eclipse-mosquitto

RUN mkdir /etc/mosquitto/tls
COPY ca.crt server.crt server.key /etc/mosquitto/tls/

RUN touch /etc/mosquitto/passwd && \
    mosquitto_passwd -b /etc/mosquitto/passwd eventr secretpassword
