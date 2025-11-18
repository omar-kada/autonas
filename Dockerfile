# ----------- Builder Stage -----------
FROM golang:1.25.1-bookworm AS builder


RUN mkdir /autonas
WORKDIR /autonas

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o autonas /autonas/cmd/autonas/main.go

# ----------- Production Stage -----------
FROM debian:bookworm-slim AS production

RUN apt update && apt install --yes --no-install-recommends \
    curl \
    ca-certificates \
    gnupg \
    unzip \
    dumb-init \
    && install -m 0755 -d /etc/apt/keyrings \
    && curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg \
    && chmod a+r /etc/apt/keyrings/docker.gpg \
    && echo \
         "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian \
         "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
         tee /etc/apt/sources.list.d/docker.list > /dev/null \
    && apt update \
    && apt --yes --no-install-recommends install \
         docker-ce-cli \
         docker-compose-plugin \
    && rm -rf /var/lib/apt/lists/*

# Move to working directory /build

ARG UID=1000
ARG GID=1000

RUN mkdir /app
RUN mkdir /app/config

WORKDIR /app

COPY --from=builder /autonas/autonas /app/
COPY frontend /app/frontend

RUN chmod -R 744 /app

ENV AUTONAS_WORKING_DIR="/app/config"
EXPOSE 8080

# Start the application
CMD ["sh", "-c", "/app/autonas run -d ${AUTONAS_WORKING_DIR} --add-write-perm ${AUTONAS_ADD_WRITE_PERM}"]