package klease

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/ezoic/sol/amazon"

	uuid "github.com/nu7hatch/gouuid"
)

func TestCoordinator_RunTaker(t *testing.T) {

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(amazon.Auth.AccessKey, amazon.Auth.SecretKey, ""),
	}))
	dDB := dynamodb.New(sess)
	manager := NewLeaseManager("testLease", dDB, NewDynamoSerializer())
	err := manager.CreateLeaseTableIfNotExists()
	if err != nil {
		t.Fatal(err)
	}
	err = manager.DeleteAll()
	if err != nil {
		t.Fatal(err)
	}
	_, err = manager.CreateLeaseIfNotExists(NewKLease("shard1", "", nil))
	if err != nil {
		t.Fatal(err)
	}
	_, err = manager.CreateLeaseIfNotExists(NewKLease("shard2", "", nil))
	if err != nil {
		t.Fatal(err)
	}
	_, err = manager.CreateLeaseIfNotExists(NewKLease("shard3", "", nil))
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		c    *Coordinator
	}{
		{
			name: "coldStart",
			c:    newLeaseCoordinatorForTest(manager, "testWorker", 1000, 100),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.RunTaker()
			count := len(tt.c.GetAssignments())
			if count != 3 {
				t.Errorf("Coordinator.RunTaker() should have taken 3 leases, got %d", count)
			}

		})
	}
}

func TestCoordinator_RunRenewer(t *testing.T) {
	tests := []struct {
		name    string
		c       *Coordinator
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.RunRenewer(); (err != nil) != tt.wantErr {
				t.Errorf("Coordinator.RunRenewer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCoordinator_GetAssignments(t *testing.T) {
	tests := []struct {
		name string
		c    *Coordinator
		want []*KLease
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetAssignments(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Coordinator.GetAssignments() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoordinator_GetCurrentlyHeldLease(t *testing.T) {
	type args struct {
		leaseKey string
	}
	tests := []struct {
		name string
		c    *Coordinator
		args args
		want *KLease
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetCurrentlyHeldLease(tt.args.leaseKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Coordinator.GetCurrentlyHeldLease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoordinator_GetWorkerId(t *testing.T) {
	tests := []struct {
		name string
		c    *Coordinator
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetWorkerId(); got != tt.want {
				t.Errorf("Coordinator.GetWorkerId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoordinator_Stop(t *testing.T) {
	tests := []struct {
		name string
		c    *Coordinator
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Stop()
		})
	}
}

func TestCoordinator_StopLeaseTaker(t *testing.T) {
	tests := []struct {
		name string
		c    *Coordinator
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.StopLeaseTaker()
		})
	}
}

func TestCoordinator_DropLease(t *testing.T) {
	type args struct {
		lease *KLease
	}
	tests := []struct {
		name string
		c    *Coordinator
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.DropLease(tt.args.lease)
		})
	}
}

func TestCoordinator_IsRunning(t *testing.T) {
	tests := []struct {
		name string
		c    *Coordinator
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsRunning(); got != tt.want {
				t.Errorf("Coordinator.IsRunning() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoordinator_UpdateLease(t *testing.T) {
	type args struct {
		lease            *KLease
		concurrencyToken *uuid.UUID
	}
	tests := []struct {
		name    string
		c       *Coordinator
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.UpdateLease(tt.args.lease, tt.args.concurrencyToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("Coordinator.UpdateLease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Coordinator.UpdateLease() = %v, want %v", got, tt.want)
			}
		})
	}
}
