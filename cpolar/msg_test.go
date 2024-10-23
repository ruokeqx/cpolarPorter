package cpolar

import (
	"encoding/json"
	"os"
	"testing"
)

func TestA(t *testing.T) {
	rc, err := os.ReadFile("./login.json")
	if err != nil {
		t.Fatal(err)
	}
	var d Response
	err = json.Unmarshal(rc, &d)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(d)
}
