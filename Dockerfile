FROM golang:1.8

WORKDIR /go/src/app
COPY . .

RUN go-wrapper download
RUN go-wrapper install

# Adjust the options as necessary using environment.
# The root of storage unless specified will be '/go/src/app/downloads'.
# One should change the advertised URL suit the deployment of this container.
#
#  -l, --listen string   Listen address (default "127.0.0.1")
#  -p, --port string     Port number (default "8080")
#  -r, --root string     Storage directory
#  -u, --url string      Filedrop server URL to advertise

CMD ["go-wrapper", "run", "--listen", "0.0.0.0", "--root", "/"]
