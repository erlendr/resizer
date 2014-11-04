package main

import (
	"bytes"
	"github.com/erlendr/store"
	"github.com/nfnt/resize"
	"github.com/streadway/amqp"
	"image/jpeg"
	"io"
	"log"
	"time"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		panic(err)
	}

	err = ch.Qos(
		3,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		panic(err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	forever := make(chan bool)

	go func() {
		println("Waiting for messages...")
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			d.Ack(false)

			var filename = string(d.Body)
			rc := store.Download(filename)

			println("Resizer - File " + filename + " downloaded")

			var buf = resizeImage(rc, 100)
			var reader = bytes.NewReader(buf.Bytes())
			store.UploadReader("thumb.jpg", reader, int64(reader.Len()))

			dot_count := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dot_count)
			time.Sleep(t * time.Second)
		}
	}()

	<-forever

}

func resizeImage(file io.Reader, width uint) *bytes.Buffer {
	img, err := jpeg.Decode(file)
	if err != nil {
		panic(err)
	}

	m := resize.Thumbnail(width, width, img, resize.Lanczos3)

	var buf = new(bytes.Buffer)
	err = jpeg.Encode(buf, m, nil)
	if err != nil {
		panic(err)
	}
	return buf
}
