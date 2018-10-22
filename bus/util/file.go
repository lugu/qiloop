package util

import (
	"io/ioutil"
	"log"
	"os"
)

// MakeTempFileName returns a non existing filename which will be
// deleted when the program exit.
func MakeTempFileName() string {
	f, err := ioutil.TempFile("", "qiloop-")
	if err != nil {
		log.Fatal(err)
	}
	name := f.Name()
	f.Close()
	os.Remove(name)
	return name
}
