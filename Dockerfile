FROM golang:alpine AS builder

ARG GO_ENVIRONMENT
ARG LOG_LEVEL
ARG CONSUL_ADDR 
ARG PNKL_CONSUL_TOKEN


ENV GO_ENVIRONMENT=$GO_ENVIRONMENT
ENV LOG_LEVEL=$LOG_LEVEL
ENV PNKL_CONSUL_ADDR=$CONSUL_ADDR 
ENV PNKL_CONSUL_TOKEN=$PNKL_CONSUL_TOKEN

WORKDIR /
ADD go.mod .
COPY . .
RUN go build -o bin/service-api.exe -ldflags="-s -w"

FROM alpine
RUN apk update && \
    apk add --no-cache curl
WORKDIR /
COPY --from=builder /bin .
EXPOSE 8080
CMD ["./service-api.exe"]