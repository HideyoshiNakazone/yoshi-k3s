package resources

import "testing"

func TestK3sNode_Init_Checks_If_Valid(t *testing.T) {
	n := &K3sNode{}

	err := n.Init("127.0.0.1", "token", []string{"worker"},
		"latest", map[string]string{})

	if err == nil {
		t.Errorf("Error: expected to failt")
	}

	nValid := &K3sNode{nodeType: MASTER}

	err = nValid.Init("127.0.0.1", "token", []string{"worker"},
		"latest", map[string]string{})

	if err != nil {
		t.Errorf("Error: expected to pass")
	}
}
