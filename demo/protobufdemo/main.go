package main

import (
	"MMO/demo/protobufdemo/pb"
	"fmt"
	"google.golang.org/protobuf/proto"
)

func main() {
	// 定义一个Person对象
	person := &pb.Person{
		Name: "XXY",
		Age:  23,
		Emails: []string{
			"MF21330093@smail.nju.edu.cn",
			"1152024415@qq.com",
		},
		Phones: []*pb.PhoneNumber{
			&pb.PhoneNumber{
				Number: "123456789",
				Type:   pb.PhoneType_MOBILE,
			},
			&pb.PhoneNumber{
				Number: "987654321",
				Type:   pb.PhoneType_HOME,
			},
		},
	}

	// encode
	// 序列化
	data, err := proto.Marshal(person)
	if err != nil {
		fmt.Println("marshal failed, err:", err)
		return
	}
	// data是一个[]byte类型的数据，可以将其存储到文件中或者发送到网络中

	// decode
	newdata := &pb.Person{}
	err = proto.Unmarshal(data, newdata)
	if err != nil {
		fmt.Println("unmarshal failed, err:", err)
		return
	}
	fmt.Println("源数据：", person)
	fmt.Println("反序列化后的数据：", newdata)
}
