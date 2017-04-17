package klease

import (
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func TestNewLeaseManager(t *testing.T) {
	type args struct {
		table     string
		client    *dynamodb.DynamoDB
		serialzer *DynamoSerializer
	}
	tests := []struct {
		name string
		args args
		want *Manager
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLeaseManager(tt.args.table, tt.args.client, tt.args.serialzer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLeaseManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newLeaseManagerForTests(t *testing.T) {
	type args struct {
		table           string
		client          *dynamodb.DynamoDB
		serialzer       *DynamoSerializer
		consistentReads bool
	}
	tests := []struct {
		name string
		args args
		want *Manager
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newLeaseManagerForTests(tt.args.table, tt.args.client, tt.args.serialzer, tt.args.consistentReads); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newLeaseManagerForTests() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_CreateLeaseTableIfNotExists(t *testing.T) {
	type args struct {
		readCapacity  int64
		writeCapacity int64
	}
	tests := []struct {
		name    string
		k       *Manager
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.k.CreateLeaseTableIfNotExists(tt.args.readCapacity, tt.args.writeCapacity); (err != nil) != tt.wantErr {
				t.Errorf("Manager.CreateLeaseTableIfNotExists() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_LeaseTableExists(t *testing.T) {
	tests := []struct {
		name    string
		k       *Manager
		want    bool
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.k.LeaseTableExists()
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.LeaseTableExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Manager.LeaseTableExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_tableStatus(t *testing.T) {
	tests := []struct {
		name    string
		k       *Manager
		want    string
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.k.tableStatus()
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.tableStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Manager.tableStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_WaitUntilLeaseTableExists(t *testing.T) {
	type args struct {
		secondsBetweenPulls int64
		timeoutSeconds      int64
	}
	tests := []struct {
		name    string
		k       *Manager
		args    args
		want    bool
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.k.WaitUntilLeaseTableExists(tt.args.secondsBetweenPulls, tt.args.timeoutSeconds)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.WaitUntilLeaseTableExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Manager.WaitUntilLeaseTableExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_ListLeases(t *testing.T) {
	tests := []struct {
		name    string
		k       *Manager
		want    []*KLease
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.k.ListLeases()
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.ListLeases() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.ListLeases() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_IsLeaseTableEmpty(t *testing.T) {
	tests := []struct {
		name    string
		k       *Manager
		want    bool
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.k.IsLeaseTableEmpty()
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.IsLeaseTableEmpty() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Manager.IsLeaseTableEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_list(t *testing.T) {
	type args struct {
		limit int64
	}
	tests := []struct {
		name    string
		k       *Manager
		args    args
		want    []*KLease
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.k.list(tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.list() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.list() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_CreateLeaseIfNotExists(t *testing.T) {
	type args struct {
		lease *KLease
	}
	tests := []struct {
		name    string
		k       *Manager
		args    args
		want    bool
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.k.CreateLeaseIfNotExists(tt.args.lease)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.CreateLeaseIfNotExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Manager.CreateLeaseIfNotExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetLease(t *testing.T) {
	type args struct {
		leaseKey string
	}
	tests := []struct {
		name    string
		k       *Manager
		args    args
		want    *KLease
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.k.GetLease(tt.args.leaseKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetLease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.GetLease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_RenewLease(t *testing.T) {
	type args struct {
		lease *KLease
	}
	tests := []struct {
		name    string
		k       *Manager
		args    args
		want    bool
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.k.RenewLease(tt.args.lease)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.RenewLease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Manager.RenewLease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_TakeLease(t *testing.T) {
	type args struct {
		lease *KLease
		owner string
	}
	tests := []struct {
		name    string
		k       *Manager
		args    args
		want    bool
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.k.TakeLease(tt.args.lease, tt.args.owner)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.TakeLease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Manager.TakeLease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_EvictLease(t *testing.T) {
	type args struct {
		lease *KLease
	}
	tests := []struct {
		name    string
		k       *Manager
		args    args
		want    bool
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.k.EvictLease(tt.args.lease)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.EvictLease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Manager.EvictLease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_DeleteAll(t *testing.T) {
	tests := []struct {
		name    string
		k       *Manager
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.k.DeleteAll(); (err != nil) != tt.wantErr {
				t.Errorf("Manager.DeleteAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_DeleteLease(t *testing.T) {
	type args struct {
		lease *KLease
	}
	tests := []struct {
		name    string
		k       *Manager
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.k.DeleteLease(tt.args.lease); (err != nil) != tt.wantErr {
				t.Errorf("Manager.DeleteLease() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_UpdateLease(t *testing.T) {
	type args struct {
		lease *KLease
	}
	tests := []struct {
		name    string
		k       *Manager
		args    args
		want    bool
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.k.UpdateLease(tt.args.lease)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.UpdateLease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Manager.UpdateLease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetCheckpoint(t *testing.T) {
	type args struct {
		shardId string
	}
	tests := []struct {
		name    string
		k       *Manager
		args    args
		want    *KCheckpoint
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.k.GetCheckpoint(tt.args.shardId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.GetCheckpoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.GetCheckpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_min(t *testing.T) {
	type args struct {
		a time.Duration
		b time.Duration
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := min(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("min() = %v, want %v", got, tt.want)
			}
		})
	}
}
