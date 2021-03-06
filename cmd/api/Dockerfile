FROM golang:alpine AS build

RUN apk update && \
    apk add curl \
            git \
            bash \
            make \
            ca-certificates && \
    rm -rf /var/cache/apk/*

WORKDIR /app

# copy module files first so that they don't need to be downloaded again if no change
COPY go.* ./
RUN go mod download
RUN go mod verify

# copy source files and build the binary
COPY . .
RUN make build


FROM alpine:latest
RUN apk --no-cache add ca-certificates bash curl
RUN mkdir -p /var/log/app
WORKDIR /app/

COPY --from=build /app/api .
COPY --from=build /app/cmd/api/entrypoint.sh .
COPY config/*.yml ./config/
RUN ls -la
ENTRYPOINT ["./entrypoint.sh"]

