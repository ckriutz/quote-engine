FROM golang:1.12.9 AS builder
WORKDIR /go/src/quote-engine
COPY . .
WORKDIR /go/src/quote-engine
RUN go get
RUN CGO_ENABLED=0 go build -a -o quote-engine

FROM scratch
COPY --from=builder /go/src/quote-engine/quote-engine .
COPY --from=builder /go/src/quote-engine/quotes.json .
EXPOSE 8080
ENTRYPOINT ["./quote-engine"]