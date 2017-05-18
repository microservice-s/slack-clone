package messages

import (
	"testing"

	"fmt"

	"github.com/aethanol/challenges-aethanol/apiserver/models/users"
)

func createNewChannel() *NewChannel {
	return &NewChannel{
		Name: "test",
	}
}

func TestNewChannelToChannel(t *testing.T) {
	nc := createNewChannel()
	creator := &users.User{
		ID: 1234,
	}
	c, err := nc.ToChannel(creator)
	fmt.Println(c.Members)
	if err != nil {
		t.Errorf("erro onverting NewChannel to Channel: %s\n", err.Error())
	}
	if nil == c {
		t.Fatalf("ToChannel() returned nil\n")
	}

}
