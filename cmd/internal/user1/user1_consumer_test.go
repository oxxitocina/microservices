package user1

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"testing"
)

func Test_failOnError(t *testing.T) {
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
			failOnError(tt.args.err, tt.args.msg)
		})
	}
}

func Test_userConsume(t *testing.T) {
	type args struct {
		ch *amqp.Channel
	}
	var tests []struct {
		name string
		args args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userConsume(tt.args.ch)
		})
	}
}
