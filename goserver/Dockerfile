#FROM golang:1.20-alpine AS build
#
#RUN apk --no-cache add go
#
#WORKDIR /app
#
#COPY go.mod .
#COPY go.sum .
#
#RUN go mod download && go mod verify
#
#COPY main.go ./
#
#
#RUN CGO_ENABLED=0 go build -o /goserver
#
#FROM gcr.io/distroless/base-debian11
#
#COPY --from=build /goserver /goserver
#
#ENTRYPOINT ["/goserver" ]

FROM golang:1.19.5-buster

RUN mkdir -p /go/bin /go/logs /go/shared /app

WORKDIR /app

RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /goserver
ENTRYPOINT ["/goserver" ]

EXPOSE 8080
EXPOSE 8081

