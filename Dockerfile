FROM golang:alpine AS builder

ARG GO_ENVIRONMENT
ARG LOG_LEVEL
ARG PNKL_CONSUL_ADDR 
#ARG PNKL_CONSUL_TOKEN


ENV GO_ENVIRONMENT=$GO_ENVIRONMENT
ENV LOG_LEVEL=$LOG_LEVEL
ENV PNKL_CONSUL_ADDR=$PNKL_CONSUL_ADDR 
ENV PNKL_CONSUL_TOKEN 5b0e9aba-79de-e63b-8a0f-9eb865769ae5
ENV GIN_MODE release

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