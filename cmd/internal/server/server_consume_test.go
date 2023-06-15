package server

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"testing"
)

func TestServerConsume(t *testing.T) {
	var tests []struct {
		name string
	} //here only 1 case, because it will wait for consuming
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ServerConsume()
		})
	}
} //integration

func Test_consumeMessages(t *testing.T) {
	type args struct {
		ch *amqp.Channel
	}
	var tests []struct {
		name string
		args args
	} // it can't consume and take 1 port multiple times, impossible
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consumeMessages(tt.args.ch)
		})
	}
} //integration

func Test_failOnError2(t *testing.T) {
	type args struct {
		err error
		msg string
	}
	var tests []struct {
		name string
		args args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			failOnError2(tt.args.err, tt.args.msg)
		})
	}
} //Unit Test
