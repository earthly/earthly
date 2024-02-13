package main

import (
	"log"

	"github.com/earthly/earthly/examples/readme/proto/pb"
	"github.com/golang/protobuf/proto"
)

func main() {
	dt, err := proto.Marshal(&pb.Person{
		Name:  "John Doe",
		Id:    1,
		Email: "john@example.com",
	})
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	person := new(pb.Person)
	err = proto.Unmarshal(dt, person)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	log.Printf("Name: %s\n", person.GetName())
	log.Printf("Email: %s\n", person.GetEmail())
}
