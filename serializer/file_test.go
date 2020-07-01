package serializer_test

import (
	"testing"

	"github.com/naraycitra/pcbook/pb"
	"github.com/naraycitra/pcbook/sample"
	"github.com/naraycitra/pcbook/serializer"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestWriteProtobufToBinaryFile(t *testing.T) {
	t.Parallel()

	binaryFile := "../tmp/laptop.bin"
	laptop1 := sample.NewLaptop()
	err := serializer.WriteProtobufToBinaryFile(laptop1, binaryFile)
	require.NoError(t, err)

	laptop2 := &pb.Laptop{}
	err = serializer.ReadProtobufFromBinaryFile(binaryFile, laptop2)
	require.NoError(t, err)
	require.True(t, proto.Equal(laptop1, laptop2))

	jsonFile := "../tmp/laptop.json"
	err = serializer.WriteProtobufToJSON(laptop1, jsonFile)
	require.NoError(t, err)

}

func TestReadProtobufFromBinaryFile(t *testing.T) {
	t.Parallel()

	binaryFile := "../tmp/laptop.bin"
	laptop1 := &pb.Laptop{}
	err := serializer.ReadProtobufFromBinaryFile(binaryFile, laptop1)
	require.NoError(t, err)
}

func TestWriteProtobufToJSON(t *testing.T) {
	t.Parallel()

	binaryFile := "../tmp/laptop.bin"
	jsonFile := "../tmp/laptop.json"

	laptop1 := &pb.Laptop{}
	err := serializer.ReadProtobufFromBinaryFile(binaryFile, laptop1)
	require.NoError(t, err)

	err = serializer.WriteProtobufToJSON(laptop1, jsonFile)
	require.NoError(t, err)

}
