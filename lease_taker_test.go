package klease

import (
	"reflect"
	"testing"
)

func TestNewKLeaseTaker(t *testing.T) {
	type args struct {
		manager            *Manager
		workerId           string
		leaseDurationNanos int64
	}
	tests := []struct {
		name string
		args args
		want *Taker
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewKLeaseTaker(tt.args.manager, tt.args.workerId, tt.args.leaseDurationNanos); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKLeaseTaker() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaker_WithMaxLeasesForWorker(t *testing.T) {
	type args struct {
		maxLeasesForWorker int64
	}
	tests := []struct {
		name string
		t    *Taker
		args args
		want *Taker
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.WithMaxLeasesForWorker(tt.args.maxLeasesForWorker); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Taker.WithMaxLeasesForWorker() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaker_WithMaxLeasesToStealAtOneTime(t *testing.T) {
	type args struct {
		maxLeasesToStealAtOnce int64
	}
	tests := []struct {
		name string
		t    *Taker
		args args
		want *Taker
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.WithMaxLeasesToStealAtOneTime(tt.args.maxLeasesToStealAtOnce); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Taker.WithMaxLeasesToStealAtOneTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaker_TakeLeases(t *testing.T) {
	tests := []struct {
		name string
		t    *Taker
		want map[string]*KLease
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.TakeLeases(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Taker.TakeLeases() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaker_updateAllLeases(t *testing.T) {
	tests := []struct {
		name    string
		t       *Taker
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.updateAllLeases(); (err != nil) != tt.wantErr {
				t.Errorf("Taker.updateAllLeases() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTaker_getExpiredLeases(t *testing.T) {
	tests := []struct {
		name string
		t    *Taker
		want []*KLease
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.getExpiredLeases(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Taker.getExpiredLeases() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaker_computeLeasesToTake(t *testing.T) {
	type args struct {
		expiredLeases []*KLease
	}
	tests := []struct {
		name string
		t    *Taker
		args args
		want []*KLease
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.computeLeasesToTake(tt.args.expiredLeases); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Taker.computeLeasesToTake() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaker_chooseLeasesToSteal(t *testing.T) {
	type args struct {
		leaseCounts map[string]int64
		needed      int64
		target      int64
	}
	tests := []struct {
		name string
		t    *Taker
		args args
		want []*KLease
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.chooseLeasesToSteal(tt.args.leaseCounts, tt.args.needed, tt.args.target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Taker.chooseLeasesToSteal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaker_computeLeaseCounts(t *testing.T) {
	type args struct {
		expiredLeases []*KLease
	}
	tests := []struct {
		name string
		t    *Taker
		args args
		want map[string]int64
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.computeLeaseCounts(tt.args.expiredLeases); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Taker.computeLeaseCounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaker_GetWorkerId(t *testing.T) {
	tests := []struct {
		name string
		t    *Taker
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.GetWorkerId(); got != tt.want {
				t.Errorf("Taker.GetWorkerId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaker_contains(t *testing.T) {
	type args struct {
		a *KLease
		b []*KLease
	}
	tests := []struct {
		name string
		t    *Taker
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.contains(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Taker.contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaker_shuffle(t *testing.T) {
	type args struct {
		a []*KLease
	}
	tests := []struct {
		name string
		t    *Taker
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.shuffle(tt.args.a)
		})
	}
}
