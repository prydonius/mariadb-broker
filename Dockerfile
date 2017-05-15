FROM golang:1.8
COPY . /go/src/github.com/prydonius/mariadb-broker
WORKDIR /go/src/github.com/prydonius/mariadb-broker
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /mariadb-broker .

FROM bitnami/minideb:latest
COPY --from=0 /mariadb-broker /mariadb-broker
ADD https://kubernetes-charts.storage.googleapis.com/mariadb-0.6.1.tgz /mariadb-0.6.1.tgz
CMD ["/mariadb-broker", "-logtostderr"]
