package mlxutils

import (
	"testing"

	sriovnetworkv1 "github.com/k8snetworkplumbingwg/sriov-network-operator/api/v1"
)

func TestHandleTotalVfs(t *testing.T) {
	tests := []struct {
		name                    string
		current                 int
		next                    int
		requested               int
		wantAttribute           int
		wantReboot              bool
		wantChangeWithoutReboot bool
	}{
		{
			name:          "preserves spare capacity",
			current:       8,
			next:          8,
			requested:     4,
			wantAttribute: -1,
		},
		{
			name:                    "cancels staged reduction to policy value",
			current:                 8,
			next:                    4,
			requested:               4,
			wantAttribute:           8,
			wantChangeWithoutReboot: true,
		},
		{
			name:                    "cancels staged reduction below policy value",
			current:                 8,
			next:                    0,
			requested:               4,
			wantAttribute:           8,
			wantChangeWithoutReboot: true,
		},
		{
			name:          "increases insufficient capacity",
			current:       4,
			next:          4,
			requested:     8,
			wantAttribute: 8,
			wantReboot:    true,
		},
		{
			name:          "disables sriov explicitly",
			current:       8,
			next:          8,
			requested:     0,
			wantAttribute: 0,
			wantReboot:    true,
		},
		{
			name:                    "cancels staged increase when sriov is disabled",
			current:                 0,
			next:                    8,
			requested:               0,
			wantAttribute:           0,
			wantChangeWithoutReboot: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			current := &MlxNic{TotalVfs: tt.current}
			next := &MlxNic{TotalVfs: tt.next}
			attributes := &MlxNic{TotalVfs: -1}
			iface := sriovnetworkv1.Interface{NumVfs: tt.requested}

			totalVfs, needReboot, changeWithoutReboot := HandleTotalVfs(
				current,
				next,
				attributes,
				iface,
				false,
				nil,
			)

			if totalVfs != tt.requested {
				t.Fatalf("totalVfs = %d, want %d", totalVfs, tt.requested)
			}
			if attributes.TotalVfs != tt.wantAttribute {
				t.Errorf("attributes.TotalVfs = %d, want %d", attributes.TotalVfs, tt.wantAttribute)
			}
			if needReboot != tt.wantReboot {
				t.Errorf("needReboot = %t, want %t", needReboot, tt.wantReboot)
			}
			if changeWithoutReboot != tt.wantChangeWithoutReboot {
				t.Errorf("changeWithoutReboot = %t, want %t", changeWithoutReboot, tt.wantChangeWithoutReboot)
			}
		})
	}
}
