package main

import (
	"crypto/rand"
	"io"
	"io/ioutil"
	"time"

	"github.com/cheggaaa/pb/v3"
)

func main() {
	simple()
	ioOperation()
}

func simple() {
	count := 100
	// create and start new bar
	bar := pb.StartNew(count)

	// start bar from 'default' template
	// bar := pb.Default.Start(count)

	// start bar from 'simple' template
	// bar := pb.Simple.Start(count)

	// start bar from 'full' template
	//bar := pb.Full.Start(count)

	for i := 0; i < count; i++ {
		bar.Increment()
		time.Sleep(time.Millisecond * 1000)
	}
	bar.Finish()
}

func ioOperation() {
	var limit int64 = 1024 * 1024 * 500
	// we will copy 200 Mb from /dev/rand to /dev/null
	reader := io.LimitReader(rand.Reader, limit)
	writer := ioutil.Discard

	// start new bar
	bar := pb.Full.Start64(limit)

	// create proxy reader
	barReader := bar.NewProxyReader(reader)
	// copy from proxy reader
	_, _ = io.Copy(writer, barReader)
	// finish bar
	bar.Finish()
}
