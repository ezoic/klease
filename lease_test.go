package klease

import (
	"reflect"
	"testing"

	"github.com/nu7hatch/gouuid"
)

func TestNewKLeaseFromLease(t *testing.T) {
	type args struct {
		other KLease
	}
	tests := []struct {
		name string
		args args
		want *KLease
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewKLeaseFromLease(tt.args.other); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKLeaseFromLease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewKLease(t *testing.T) {
	type args struct {
		leaseKey   string
		leaseOwner string
		checkpoint *KCheckpoint
	}
	tests := []struct {
		name string
		args args
		want *KLease
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewKLease(tt.args.leaseKey, tt.args.leaseOwner, tt.args.checkpoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKLease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKLease_GetLeaseKey(t *testing.T) {
	tests := []struct {
		name string
		l    *KLease
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.GetLeaseKey(); got != tt.want {
				t.Errorf("KLease.GetLeaseKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKLease_GetLeaseCounter(t *testing.T) {
	tests := []struct {
		name string
		l    *KLease
		want int64
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.GetLeaseCounter(); got != tt.want {
				t.Errorf("KLease.GetLeaseCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKLease_GetLeaseOwner(t *testing.T) {
	tests := []struct {
		name string
		l    *KLease
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.GetLeaseOwner(); got != tt.want {
				t.Errorf("KLease.GetLeaseOwner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKLease_GetConcurrencyToken(t *testing.T) {
	tests := []struct {
		name string
		l    *KLease
		want *uuid.UUID
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.GetConcurrencyToken(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KLease.GetConcurrencyToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKLease_GetLastCounterIncrementNanos(t *testing.T) {
	tests := []struct {
		name string
		l    *KLease
		want int64
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.GetLastCounterIncrementNanos(); got != tt.want {
				t.Errorf("KLease.GetLastCounterIncrementNanos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKLease_GetCheckpoint(t *testing.T) {
	tests := []struct {
		name string
		l    *KLease
		want *KCheckpoint
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.GetCheckpoint(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KLease.GetCheckpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKLease_GetOwnerSwitchesSinceCheckpoint(t *testing.T) {
	tests := []struct {
		name string
		l    *KLease
		want int64
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.GetOwnerSwitchesSinceCheckpoint(); got != tt.want {
				t.Errorf("KLease.GetOwnerSwitchesSinceCheckpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKLease_GetParentShardIds(t *testing.T) {
	tests := []struct {
		name string
		l    *KLease
		want map[string]string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.GetParentShardIds(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KLease.GetParentShardIds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKLease_IsExpired(t *testing.T) {
	type args struct {
		leaseDurationNanos int64
		asOfNanos          int64
	}
	tests := []struct {
		name string
		l    *KLease
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.IsExpired(tt.args.leaseDurationNanos, tt.args.asOfNanos); got != tt.want {
				t.Errorf("KLease.IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKLease_SetLeaseKey(t *testing.T) {
	type args struct {
		leaseKey string
	}
	tests := []struct {
		name    string
		l       *KLease
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.l.SetLeaseKey(tt.args.leaseKey); (err != nil) != tt.wantErr {
				t.Errorf("KLease.SetLeaseKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKLease_SetLeaseCounter(t *testing.T) {
	type args struct {
		leaseCounter int64
	}
	tests := []struct {
		name    string
		l       *KLease
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.l.SetLeaseCounter(tt.args.leaseCounter); (err != nil) != tt.wantErr {
				t.Errorf("KLease.SetLeaseCounter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKLease_SetLeaseOwner(t *testing.T) {
	type args struct {
		leaseOwner string
	}
	tests := []struct {
		name    string
		l       *KLease
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.l.SetLeaseOwner(tt.args.leaseOwner); (err != nil) != tt.wantErr {
				t.Errorf("KLease.SetLeaseOwner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKLease_SetConcurrencyToken(t *testing.T) {
	type args struct {
		concurrencyToken *uuid.UUID
	}
	tests := []struct {
		name    string
		l       *KLease
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.l.SetConcurrencyToken(tt.args.concurrencyToken); (err != nil) != tt.wantErr {
				t.Errorf("KLease.SetConcurrencyToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKLease_SetLastCounterIncrementNanos(t *testing.T) {
	type args struct {
		lastCounterIncrementNanos int64
	}
	tests := []struct {
		name    string
		l       *KLease
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.l.SetLastCounterIncrementNanos(tt.args.lastCounterIncrementNanos); (err != nil) != tt.wantErr {
				t.Errorf("KLease.SetLastCounterIncrementNanos() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKLease_SetCheckpoint(t *testing.T) {
	type args struct {
		checkpoint *KCheckpoint
	}
	tests := []struct {
		name    string
		l       *KLease
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.l.SetCheckpoint(tt.args.checkpoint); (err != nil) != tt.wantErr {
				t.Errorf("KLease.SetCheckpoint() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKLease_SetOwnerSwitchesSinceCheckpoint(t *testing.T) {
	type args struct {
		ownerSwitchesSinceCheckpoint int64
	}
	tests := []struct {
		name    string
		l       *KLease
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.l.SetOwnerSwitchesSinceCheckpoint(tt.args.ownerSwitchesSinceCheckpoint); (err != nil) != tt.wantErr {
				t.Errorf("KLease.SetOwnerSwitchesSinceCheckpoint() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKLease_SetParentShardIds(t *testing.T) {
	type args struct {
		parentShardIds map[string]string
	}
	tests := []struct {
		name    string
		l       *KLease
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.l.SetParentShardIds(tt.args.parentShardIds); (err != nil) != tt.wantErr {
				t.Errorf("KLease.SetParentShardIds() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKLease_Update(t *testing.T) {
	type args struct {
		other *KLease
	}
	tests := []struct {
		name string
		l    *KLease
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Update(tt.args.other)
		})
	}
}

func TestKLease_HashCode(t *testing.T) {
	tests := []struct {
		name string
		l    *KLease
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.HashCode(); got != tt.want {
				t.Errorf("KLease.HashCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKLease_Equals(t *testing.T) {
	type args struct {
		other *KLease
	}
	tests := []struct {
		name string
		l    *KLease
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Equals(tt.args.other); got != tt.want {
				t.Errorf("KLease.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}
