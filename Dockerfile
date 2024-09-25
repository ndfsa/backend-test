# database container
FROM postgres:alpine AS database
# Set environment variables
ENV POSTGRES_USER=postgres
ENV POSTGRES_DB=cardboard_bank
ENV POSTGRES_PASSWORD=root
RUN [[ -d "/docker-entrypoint-initdb.d" ]] || mkdir /docker-entrypoint-initdb.d ;
# Copy script to container
COPY ./database.sql /docker-entrypoint-initdb.d/
EXPOSE 5432

# nginx
FROM nginx:alpine AS gateway
COPY ./nginx.conf /etc/nginx/nginx.conf
EXPOSE 80

# build stage container
FROM golang:alpine AS build
WORKDIR /src
# Copy module files to download dependencies
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download ;
# Copy source to container
COPY . /src
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" \
    go build -o /bin/tui ./tui && \
    go build -o /bin/auth ./web/auth && \
    go build -o /bin/api ./web/api ;

# auth container
FROM alpine AS authprod
COPY --from=build /bin/auth /bin/auth
EXPOSE 80
ENTRYPOINT /bin/auth ;

# api container
FROM alpine AS apiprod
COPY --from=build /bin/api /bin/api
EXPOSE 80
ENTRYPOINT /bin/api ;

# tui container
FROM alpine AS tui
RUN apk add openssh gettext moreutils ;
RUN ssh-keygen -A ;
COPY --from=build /bin/tui /bin/tui
COPY ./tui/tui.conf /etc/ssh/ssh_config.d/tui.conf
COPY ./tui/wrapper.sh /bin/tui_wrapper
ARG DB_URL
RUN envsubst < /bin/tui_wrapper | sponge /bin/tui_wrapper && \
    adduser tui -s /bin/tui_wrapper -D && \
    echo "tui:tui" | chpasswd && \
    echo "" > /etc/motd ;
EXPOSE 22
ENTRYPOINT /usr/sbin/sshd -D -e ;

# prometheus metrics
FROM prom/prometheus AS prometheus
COPY ./prometheus.yaml /etc/prometheus/prometheus.yml
EXPOSE 9090

# grafana analytics
FROM grafana/grafana AS grafana
COPY ./datasource.yml /etc/grafana/provisioning/datasources/datasource.yml
EXPOSE 3000

# cadvisor analytics
FROM gcr.io/cadvisor/cadvisor:latest AS cadvisor
EXPOSE 8080
