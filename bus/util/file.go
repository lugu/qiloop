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

// NewUnixAddr returns a unused UNIX socket address inluding the
// prefix unix://. The socket file will be collected when the program
// exit.
func NewUnixAddr() string {
	return "unix://" + makeTempFileName()
}
