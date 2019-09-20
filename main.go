package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"time"

	language "cloud.google.com/go/language/apiv1"
	"github.com/golang/protobuf/proto"
	"github.com/windy-tech/recommendation/pkg/text"
)

func main() {
	ctx, f := context.WithTimeout(context.Background(), time.Second*5)
	defer f()
	client, err := language.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	t, err := ioutil.ReadFile("./test.txt")
	if err != nil {
		log.Fatal(err)
	}
	resp, err := text.ClassifyText(ctx, client, string(t))
	printResp(resp, err)
}

func printResp(v proto.Message, err error) {
	if err != nil {
		log.Fatal(err)
	}
	proto.MarshalText(os.Stdout, v)
}
