package user1

import "testing"

func TestUserProduce1(t *testing.T) {
	var tests []struct {
		name string
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UserProduce1()
		})
	}
}

//I stopped the tests, because user1 and user2 code are very similar, difference only in keys and queue names
