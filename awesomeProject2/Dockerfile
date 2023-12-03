FROM golang:1.20

WORKDIR /awesomeProject2

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

WORKDIR /awesomeProject2/cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /greenlight

EXPOSE 4000
CMD ["/greenlight"]