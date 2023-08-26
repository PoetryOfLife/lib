package hs256

import (
	"testing"
	"time"
)

func TestHS256(t *testing.T) {

	token, err := GenerateTokenUsingHs256(1, "eden")
	if err != nil {
		t.Errorf("generate token failed:%s", err.Error())
	}
	t.Logf("token:%s", token)
	time.Sleep(time.Second * 2)
	claims, err := ParseTokenHS256(token)
	if err != nil {
		t.Errorf("parse token failed:%s", err.Error())
	}
	t.Logf("claims:%++v", claims)
}
