FROM golang:alpine AS builder

WORKDIR /
ADD go.mod .
COPY . .
RUN go build -o bin/service-api.exe -ldflags="-s -w"

FROM alpine

ARG PNKL_VAULT_TOKEN
ARG GO_ENVIRONMENT
ARG LOG_LEVEL

ENV PNKL_VAULT_TOKEN=$PNKL_VAULT_TOKEN
ENV GO_ENVIRONMENT=$GO_ENVIRONMENT
ENV LOG_LEVEL=$LOG_LEVEL

RUN apk update && \
    apk add --no-cache curl
WORKDIR /
COPY --from=builder /bin .
COPY --from=builder /var /var
EXPOSE 8080 
EXPOSE 9000
CMD ["./service-api.exe"]
