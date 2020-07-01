package serializer

import (
	"fmt"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
)

// WriteProtobufToBinaryFile write protocol buffer message to binary file
func WriteProtobufToBinaryFile(message proto.Message, filename string) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("cannot marshal the message %v", err)
	}
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("cannot write file %v", err)
	}
	return nil
}

// ReadProtobufFromBinaryFile read protocol buffer message from binary file
func ReadProtobufFromBinaryFile(filename string, messsage proto.Message) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("cannot read binary data from file %v", err)
	}

	err = proto.Unmarshal(data, messsage)
	if err != nil {
		return fmt.Errorf("cannot unmarshal the data %v", err)
	}
	return nil

}

// WriteProtobufToJSON write protocol buffer message binary to JSON format file
func WriteProtobufToJSON(message proto.Message, filename string) error {
	data, err := ProtobufToJSON(message)
	if err != nil {
		return fmt.Errorf("cannot convert protobuf to JSON %v", err)
	}

	err = ioutil.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("cannot write data to file %v", err)
	}
	return nil

}
