FROM docker.io/golang

WORKDIR /go/src/docker_image_info
COPY . .

RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run"] # ["docker_image_info"]
