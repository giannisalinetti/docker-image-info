# Docker Image Info

A small tool to view Docker images details

## Getting Started

This program can be used as complementary in shell scripts to grab more detailed output.
To use it clone the repository and build it or run it as a Container.

### Prerequisites

Golang version 1.8 or higher.

To install Golang on RHEL/CentOS:

```
$ sudo yum install -y golang
```

To install Golang on Fedora

```
$ sudo dnf install -y golang
```

To install on Debian/Ubuntu

```
$ sudo apt-get install golang
```

### Installing

The GOPATH environment variable must be defined before building. A common
path is `~/go`.
To build the standalone executable and install it under $GOPATH/bin

```
$ go install
```

To build the Docker image:

```
$ docker build -t docker_image_info .

To run the docker image we need to map the /var/run directory as a container volume
to grant access to `/var/run/docker.sock`.

```
$ docker run -v /var/run:/var/run docker_image_info
```

To run che container with extra flags:

```
$ docker run -v /var/run:/var/run docker_image_info docker_image_info -yaml
```

## Authors

* **Giovan Battista Salinetti** - *Initial work* - [giannisalinetti]

## License

GPL License
