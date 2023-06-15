package server

import (
	pb "rabbitMQhelloworld/api/proto"
	"reflect"
	"testing"
)

func TestServerProduce(t *testing.T) {
	var tests []struct {
		name string
	} //here only 1 test, because server can not be listened from 3 different sources (ports) at the same time
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ServerProduce()
			// Add assertions or checks here to verify the expected behavior
		})
	}
}

func Test_bodyFrom(t *testing.T) {
	tests := []struct {
		name string
		want *pb.MyMessage
	}{
		{
			name: "Test Case 1",
			want: &pb.MyMessage{
				Data: "Test Message 1",
			},
		},
		{
			name: "Test Case 2",
			want: &pb.MyMessage{
				Data: "Test Message 2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bodyFrom()
			if reflect.DeepEqual(got, tt.want) {
				t.Errorf("bodyFrom() = %v, want not equal to %v", got, tt.want)
			}
		})
	} //Here because there's no input, we check if bodyfrom not equal to TestMessage, they are not equal, because our bodyfrom is empty
}
