resizer
=======

Golang-based image resizer using RabbitMQ and Amazon S3.

##Prerequisites##
- RabbitMQ
- Amazon S3 account with credentials (see `setenv.sh`)

##How it works##

Two components: `resizer.go` and `client.go`

- `resizer.go` is the image resizing worker application
- `client.go` is a bare bones client application

1. `client.go` sends a message containing a Amazon S3 object key ("filename") to a queue via RabbitMQ
2. `resizer.go` receives messages from the queue and downloads the file from S3 using the object key from the message
3. `resizer.go` resizes the image to a thumbnail with a predefined size and uploads it again to S3 with a suffix

