# database container
FROM postgres:alpine as database
# Set environment variables
ENV POSTGRES_USER=postgres
ENV POSTGRES_DB=cardboard_bank
ENV POSTGRES_PASSWORD=root
RUN [[ -d "/docker-entrypoint-initdb.d" ]] || mkdir /docker-entrypoint-initdb.d ;
EXPOSE 5432

# Copy script to container
COPY ./database.sql /docker-entrypoint-initdb.d/

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
    go build -o /bin/auth ./web/auth && \
    go build -o /bin/api ./web/api ;

# auth container
FROM alpine AS authProd
COPY --from=build /bin/auth /bin/auth
EXPOSE 80
ENTRYPOINT /bin/auth ;
# api container
FROM alpine AS apiProd
COPY --from=build /bin/api /bin/api
EXPOSE 80
ENTRYPOINT /bin/api ;
