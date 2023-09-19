FROM golang:1.19

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o project ./cmd/main.go

EXPOSE 9092

CMD ["./project"]
