FROM golang:1.21-alpine as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

FROM gcr.io/distroless/static-debian12:nonroot
ENV TZ Asia/Tokyo
COPY --from=builder /app/app /
ENTRYPOINT [ "/app" ]
