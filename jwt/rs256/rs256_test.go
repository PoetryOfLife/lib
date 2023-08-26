package rs256

import (
	"testing"
	"time"
)

func TestRS256(t *testing.T) {

	token, err := GenerateTokenUsingRS256(1, "eden")
	if err != nil {
		t.Errorf("generate token failed:%s", err.Error())
	}
	t.Logf("token:%s", token)
	time.Sleep(time.Second * 1)
	claims, err := ParseTokenRs256(token)
	if err != nil {
		t.Errorf("parse token failed:%s", err.Error())
	}
	t.Logf("claims:%++v", claims)
}
