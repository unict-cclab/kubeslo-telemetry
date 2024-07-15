FROM golang:1.21 AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o /kubeslo-telemetry -buildvcs=false

FROM debian:bookworm-slim
WORKDIR /
COPY --from=build /kubeslo-telemetry /kubeslo-telemetry
EXPOSE 8080
ENTRYPOINT ["/kubeslo-telemetry"]