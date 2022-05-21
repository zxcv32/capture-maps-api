FROM golang:1.18 AS build
WORKDIR /capture/
COPY go.mod /capture/
COPY go.sum /capture/
RUN go mod download
COPY src/ /capture/src
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -tags=nomsgpack -a -o app ./src

FROM alpine:3
LABEL NAME="capture-maps-api"
WORKDIR /capture
ENV GIN_MODE=release
EXPOSE 8090
COPY --from=build /capture/app ./
# Please specify at least INFLUXDB_ORG and INFLUXDB_TOKEN at runtime
CMD ["./app"]
