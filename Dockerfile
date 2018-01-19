FROM docker.io/golang

WORKDIR /go/src/docker-image-info
COPY . .

RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run"] # ["docker-image-info"]
