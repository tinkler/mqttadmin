package parser

import "testing"

func TestParseProtoFile(t *testing.T) {
	p, err := ParseProtoFile("../../api/proto/role/v1/role.proto")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(p.sequenceMap)
}
