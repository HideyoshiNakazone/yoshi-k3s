package resources

import (
	"github.com/HideyoshiNakazone/yoshi-k3s/pkg/ssh_handler"
	"reflect"
	"testing"
)

func TestNewNodeConfig(t *testing.T) {
	type args struct {
		name             string
		connectionConfig *ssh_handler.SshConfig
	}
	tests := []struct {
		name string
		args args
		want *NodeConfig
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNodeConfig(tt.args.name, tt.args.connectionConfig); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNodeConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeConfig_GetConnectionConfig(t *testing.T) {
	type fields struct {
		name             string
		connectionConfig *ssh_handler.SshConfig
	}
	tests := []struct {
		name   string
		fields fields
		want   *ssh_handler.SshConfig
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NodeConfig{
				name:             tt.fields.name,
				connectionConfig: tt.fields.connectionConfig,
			}
			if got := n.GetConnectionConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConnectionConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeConfig_GetName(t *testing.T) {
	type fields struct {
		name             string
		connectionConfig *ssh_handler.SshConfig
	}
	tests := []struct {
		name   string
		fields fields
		want   *string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NodeConfig{
				name:             tt.fields.name,
				connectionConfig: tt.fields.connectionConfig,
			}
			if got := n.GetName(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeConfig_IsValid(t *testing.T) {
	type fields struct {
		name             string
		connectionConfig *ssh_handler.SshConfig
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NodeConfig{
				name:             tt.fields.name,
				connectionConfig: tt.fields.connectionConfig,
			}
			if err := n.IsValid(); (err != nil) != tt.wantErr {
				t.Errorf("IsValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
