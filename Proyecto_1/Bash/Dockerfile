FROM alpine:latest

# Instalar bash y cliente docker
RUN apk add --no-cache bash docker-cli

# Copiar scripts
COPY crear_contenedores.sh /usr/local/bin/crear_contenedores.sh
COPY limpiar_contenedores.sh /usr/local/bin/limpiar_contenedores.sh

RUN chmod +x /usr/local/bin/crear_contenedores.sh \
    && chmod +x /usr/local/bin/limpiar_contenedores.sh

# Definir crontab
RUN echo "*/5 * * * * /bin/bash /usr/local/bin/crear_contenedores.sh >> /var/log/contenedores.log 2>&1" >> /etc/crontabs/root \
    && echo "*/6 * * * * /bin/bash /usr/local/bin/limpiar_contenedores.sh >> /var/log/limpiar.log 2>&1" >> /etc/crontabs/root

# Mantener cron en primer plano
CMD ["crond", "-f", "-l", "2"]
