package klease

import (
	"reflect"
	"testing"

	"github.com/nu7hatch/gouuid"
)

func TestNewKLeaseRenewer(t *testing.T) {
	type args struct {
		leaseManager       *Manager
		workerId           string
		leaseDurationNanos int64
	}
	tests := []struct {
		name string
		args args
		want *Renewer
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewKLeaseRenewer(tt.args.leaseManager, tt.args.workerId, tt.args.leaseDurationNanos); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKLeaseRenewer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRenewer_RenewLeases(t *testing.T) {
	tests := []struct {
		name    string
		r       *Renewer
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.RenewLeases(); (err != nil) != tt.wantErr {
				t.Errorf("Renewer.RenewLeases() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRenewer_renewLease(t *testing.T) {
	type args struct {
		lease           *KLease
		renewLeaseTasks chan renewResult
		sem             chan bool
	}
	tests := []struct {
		name string
		r    *Renewer
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.renewLease(tt.args.lease, tt.args.renewLeaseTasks, tt.args.sem)
		})
	}
}

func TestRenewer_renewLeaseInner(t *testing.T) {
	type args struct {
		lease              *KLease
		renewEvenIfExpired bool
	}
	tests := []struct {
		name    string
		r       *Renewer
		args    args
		want    bool
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.renewLeaseInner(tt.args.lease, tt.args.renewEvenIfExpired)
			if (err != nil) != tt.wantErr {
				t.Errorf("Renewer.renewLeaseInner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Renewer.renewLeaseInner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRenewer_GetCurrentlyHeldLeases(t *testing.T) {
	tests := []struct {
		name string
		r    *Renewer
		want map[string]*KLease
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.GetCurrentlyHeldLeases(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Renewer.GetCurrentlyHeldLeases() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRenewer_GetCurrentlyHeldLease(t *testing.T) {
	type args struct {
		leaseKey string
	}
	tests := []struct {
		name string
		r    *Renewer
		args args
		want *KLease
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.GetCurrentlyHeldLease(tt.args.leaseKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Renewer.GetCurrentlyHeldLease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRenewer_getCopyOfHeldLease(t *testing.T) {
	type args struct {
		leaseKey string
		now      int64
	}
	tests := []struct {
		name string
		r    *Renewer
		args args
		want *KLease
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.getCopyOfHeldLease(tt.args.leaseKey, tt.args.now); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Renewer.getCopyOfHeldLease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRenewer_UpdateLease(t *testing.T) {
	type args struct {
		lease            *KLease
		concurrencyToken *uuid.UUID
	}
	tests := []struct {
		name    string
		r       *Renewer
		args    args
		want    bool
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.UpdateLease(tt.args.lease, tt.args.concurrencyToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("Renewer.UpdateLease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Renewer.UpdateLease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRenewer_AddLeasesToRenew(t *testing.T) {
	type args struct {
		newLeases []*KLease
	}
	tests := []struct {
		name    string
		r       *Renewer
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.AddLeasesToRenew(tt.args.newLeases); (err != nil) != tt.wantErr {
				t.Errorf("Renewer.AddLeasesToRenew() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRenewer_ClearCurentHeldLeases(t *testing.T) {
	tests := []struct {
		name string
		r    *Renewer
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.ClearCurentHeldLeases()
		})
	}
}

func TestRenewer_DropLease(t *testing.T) {
	type args struct {
		lease *KLease
	}
	tests := []struct {
		name string
		r    *Renewer
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.DropLease(tt.args.lease)
		})
	}
}

func TestRenewer_Init(t *testing.T) {
	tests := []struct {
		name    string
		r       *Renewer
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Init(); (err != nil) != tt.wantErr {
				t.Errorf("Renewer.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
