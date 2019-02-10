package serializer

import (
	"testing"
)

func Test_Serializer(t *testing.T) {
	client := NewClientSerializer()

	client.Deserializer(client.Serializer())
	client.Deserializer(client.Serializer())
}
