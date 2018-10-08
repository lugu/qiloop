package main

import (
	"bufio"
	"flag"
	"github.com/lugu/qiloop/bus/client"
	"github.com/lugu/qiloop/bus/net"
	"log"
	"os"
	"time"
)

func Fuzz(data []byte) int {
	return 0
}

func test(endpoint net.EndPoint, user, token string, i int) {
	log.Printf("Authentication attempt: %d", i)
	err := client.AuthenticateUser(endpoint, user, token)
	if err == nil {
		log.Printf("Authentication succedded: %s", token)
		exit(0)
	} else {
		log.Printf("Authentication failed: %s", token)
	}
}

func exit(status int) {
	duration, _ := time.ParseDuration("1s")
	time.Sleep(duration)
	log.Fatalf("exiting...")
}

func filter(hdr *net.Header) (matched bool, keep bool) {
	log.Printf("received response (%d)", hdr.ID)
	return false, true
}

func consumer(msg *net.Message) error {
	return nil
}

func closer(err error) {
	exit(1)
}

func main() {
	var serverURL = flag.String("qi-url", "tcp://127.0.0.1:9559",
		"server address")
	var dictionnary = flag.String("dictionary", "", "dictionary file")
	var user = flag.String("user", "", "auth user")

	flag.Parse()

	endpoint, err := net.DialEndPoint(*serverURL)
	if err != nil {
		log.Fatalf("failed to contact %s: %s", *serverURL, err)
	}

	endpoint.AddHandler(filter, consumer, closer)

	file, err := os.Open(*dictionnary)
	if err != nil {
		log.Fatalf("failed to open %s: %s", *dictionnary, err)
	}

	defer file.Close()
	r := bufio.NewReader(file)

	for i := 0; i < 10; i++ {
		password, err := r.ReadString('\n')
		if err != nil {
			log.Fatalf("failed to read password: %s", err)
		}
		go test(endpoint, *user, password, i)
	}
	time.Sleep(time.Hour * 10)
}
