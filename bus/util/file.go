package util

import (
	"io/ioutil"
	"log"
	"os"
)

// makeTempFileName returns a non existing filename which will be
// deleted when the program exit.
func makeTempFileName() string {
	f, err := ioutil.TempFile("", "qiloop-")
	if err != nil {
		log.Fatal(err)
	}
	name := f.Name()
	f.Close()
	os.Remove(name)
	return name
}

// NewUnixAddr returns an unused UNIX socket address inluding the
// prefix unix://. The socket file will be collected when the program
// exits.
func NewUnixAddr() string {
	return "unix://" + makeTempFileName()
}

// NewPipeAddr returns an unused UNIX socket address including the
// prefix pipe://. The socket file will be collected when the program
// extis.
func NewPipeAddr() string {
	return "pipe://" + makeTempFileName()
}
