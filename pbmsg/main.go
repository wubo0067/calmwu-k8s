/*
 * @Author: calm.wu
 * @Date: 2019-08-16 19:59:21
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-08-16 20:06:39
 */

package main

import (
	"log"
	"pbmsg/ns_person"

	"github.com/gogo/protobuf/proto"
)

func main() {
	person_in := &ns_person.Person{
		Id:   0,
		Name: "test",
		Sex:  2,
	}

	log.Printf("person_in:%+v\n", person_in)

	data, _ := proto.Marshal(person_in)

	var person_out ns_person.Person
	proto.Unmarshal(data, &person_out)

	log.Printf("person_out:%+v\n", person_out)
}
