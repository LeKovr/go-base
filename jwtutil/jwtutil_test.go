package jwtutil

import (
	"log"
	"os"
	"reflect"
	"testing"
)

const (
	testBlockKey = "asdfghjklqwertyu"
)

// -----------------------------------------------------------------------------
// kvstore definition

type PhoneData struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

// -----------------------------------------------------------------------------
func TestJWT(t *testing.T) {

	logger := log.New(os.Stderr, "", 0)

	jwt, _ := New(logger, &Flags{BlockKey: testBlockKey})

	value := &PhoneData{Phone: "1234", Code: "abc"}
	encoded, err := jwt.Cryptor.Encode("appKey", value)
	if err != nil {
		t.Error("Unexpected encode error:", err)
	}

	ret := new(PhoneData)
	if err = jwt.Cryptor.Decode("appKey", encoded, ret); err != nil {
		t.Error("Unexpected decode error:", err)
	}

	if !reflect.DeepEqual(value, ret) {
		t.Error("Expected equal data but got:", value, ret)
	}

}
