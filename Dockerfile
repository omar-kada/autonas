# ----------- Builder Stage -----------
FROM golang:1.25.1-bookworm AS builder


RUN mkdir /autonas
WORKDIR /autonas

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o autonas autonas.go

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
RUN mkdir /autonas
RUN mkdir /autonas/config
WORKDIR /autonas

COPY --from=builder /autonas/autonas /autonas/
RUN chmod 755 /autonas/autonas

ARG CONFIG_FILES
ARG CONFIG_REPO

WORKDIR /autonas/config
# Start the application
CMD ["sh", "-c", "/autonas/autonas run -c ${CONFIG_FILES} -r ${CONFIG_REPO}"]
