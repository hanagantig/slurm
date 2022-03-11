FROM golang:1.17.2 AS build

WORKDIR /go/src/slurm

COPY . ./

ENV GOPATH /go/
ENV GO111MODULE on

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o slurmapi

FROM alpine:3.12.0

COPY --from=build /go/src/slurm/slurmapi /opt/slurm/slurmapi

WORKDIR /opt/slurm

EXPOSE 8080

CMD ./slurmapi
