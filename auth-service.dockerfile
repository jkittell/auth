FROM --platform=$BUILDPLATFORM golang:alpine AS builder
WORKDIR /build
COPY go.mod ./
RUN go mod download
COPY . ./
ARG TARGETOS
ARG TARGETARCH
ENV GOOS $TARGETOS
ENV GOARCH $TARGETARCH
RUN go build -o authApp ./cmd

FROM alpine:latest
RUN apk update
RUN apk upgrade
COPY --from=builder ["/build/authApp", "/"]
EXPOSE 80
RUN chmod +x /authApp
ENTRYPOINT [ "/authApp" ]