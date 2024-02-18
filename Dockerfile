# database container
FROM postgres:alpine as database

# Set environment variables
ENV POSTGRES_USER=postgres
ENV POSTGRES_DB=cardboard_bank
ENV POSTGRES_PASSWORD=root

RUN [[ -d /docker-entrypoint-initdb.d ]] || mkdir /docker-entrypoint-initdb.d

# Copy script to container
COPY ./cmd/db/database.sql /docker-entrypoint-initdb.d/

EXPOSE 5432



# build stage container
FROM golang:alpine as build

WORKDIR /src

# Copy module files to download dependencies
COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

# Copy source to container
COPY . /src

RUN go build -o /bin/auth ./cmd/auth
# RUN go build -o /bin/api ./cmd/api



# api container
FROM alpine as authProd

COPY --from=build /bin/auth /bin/auth

EXPOSE 80

CMD ["auth"]
