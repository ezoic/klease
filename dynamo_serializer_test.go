package klease

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	uuid "github.com/nu7hatch/gouuid"
)

func TestDynamoSerializer_ToDynamoRecord(t *testing.T) {
	type args struct {
		lease *KLease
	}
	token, _ := uuid.NewV4()
	SequenceNumber := "testSequenceNumber"
	leaseCounter := int64(2)
	leaseCounterString := "2"
	leaseKey := "testShard"
	leaseOwner := "testOwner"
	OwnerSwitches := int64(1)
	OwnerSwitchesString := "1"

	tests := []struct {
		name string
		s    *DynamoSerializer
		args args
		want map[string]*dynamodb.AttributeValue
	}{
		{
			name: "withAllLeaseValues",
			s:    NewDynamoSerializer(),
			args: args{lease: &KLease{
				checkpoint: &KCheckpoint{
					AppName:        "testApp",
					StreamName:     "testStream",
					SequenceNumber: SequenceNumber,
				},
				ownerSwitchesSinceCheckpoint: OwnerSwitches,
				leaseKey:                     leaseKey,
				leaseOwner:                   leaseOwner,
				leaseCounter:                 leaseCounter,
				concurrencyToken:             token,
				lastCounterIncrementNanos:    123456789,
			},
			},
			want: map[string]*dynamodb.AttributeValue{
				LeaseKeyKey: &dynamodb.AttributeValue{
					S: &leaseKey,
				},
				LeaseCounterKey: &dynamodb.AttributeValue{
					N: &leaseCounterString,
				},
				LeaseOwnerKey: &dynamodb.AttributeValue{
					S: &leaseOwner,
				},
				OwnerSwitchesKey: &dynamodb.AttributeValue{
					N: &OwnerSwitchesString,
				},
				CheckpointSequenceNumberKey: &dynamodb.AttributeValue{
					S: &SequenceNumber,
				},
			},
		},
		{
			name: "withoutOwner",
			s:    NewDynamoSerializer(),
			args: args{lease: &KLease{
				checkpoint: &KCheckpoint{
					AppName:        "testApp",
					StreamName:     "testStream",
					SequenceNumber: SequenceNumber,
				},
				ownerSwitchesSinceCheckpoint: OwnerSwitches,
				leaseKey:                     leaseKey,
				leaseCounter:                 leaseCounter,
				concurrencyToken:             token,
				lastCounterIncrementNanos:    123456789,
			},
			},
			want: map[string]*dynamodb.AttributeValue{
				LeaseKeyKey: &dynamodb.AttributeValue{
					S: &leaseKey,
				},
				LeaseCounterKey: &dynamodb.AttributeValue{
					N: &leaseCounterString,
				},
				OwnerSwitchesKey: &dynamodb.AttributeValue{
					N: &OwnerSwitchesString,
				},
				CheckpointSequenceNumberKey: &dynamodb.AttributeValue{
					S: &SequenceNumber,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.ToDynamoRecord(tt.args.lease); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DynamoSerializer.ToDynamoRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDynamoSerializer_FromDynamoRecord(t *testing.T) {
	type args struct {
		record map[string]*dynamodb.AttributeValue
	}

	SequenceNumber := "testSequenceNumber"
	leaseCounter := int64(2)
	leaseCounterString := "2"
	leaseKey := "testShard"
	leaseOwner := "testOwner"
	OwnerSwitches := int64(1)
	OwnerSwitchesString := "1"

	tests := []struct {
		name string
		s    *DynamoSerializer
		args args
		want *KLease
	}{
		{
			name: "withAllLeaseValues",
			s:    NewDynamoSerializer(),
			args: args{record: map[string]*dynamodb.AttributeValue{
				LeaseKeyKey: &dynamodb.AttributeValue{
					S: &leaseKey,
				},
				LeaseCounterKey: &dynamodb.AttributeValue{
					N: &leaseCounterString,
				},
				LeaseOwnerKey: &dynamodb.AttributeValue{
					S: &leaseOwner,
				},
				OwnerSwitchesKey: &dynamodb.AttributeValue{
					N: &OwnerSwitchesString,
				},
				CheckpointSequenceNumberKey: &dynamodb.AttributeValue{
					S: &SequenceNumber,
				},
			},
			},
			want: &KLease{
				checkpoint: &KCheckpoint{
					SequenceNumber: SequenceNumber,
				},
				ownerSwitchesSinceCheckpoint: OwnerSwitches,
				leaseKey:                     leaseKey,
				leaseOwner:                   leaseOwner,
				leaseCounter:                 leaseCounter,
			},
		},
		{
			name: "withoutOwner",
			s:    NewDynamoSerializer(),
			args: args{record: map[string]*dynamodb.AttributeValue{
				LeaseKeyKey: &dynamodb.AttributeValue{
					S: &leaseKey,
				},
				LeaseCounterKey: &dynamodb.AttributeValue{
					N: &leaseCounterString,
				},
				OwnerSwitchesKey: &dynamodb.AttributeValue{
					N: &OwnerSwitchesString,
				},
				CheckpointSequenceNumberKey: &dynamodb.AttributeValue{
					S: &SequenceNumber,
				},
			},
			},
			want: &KLease{
				checkpoint: &KCheckpoint{
					SequenceNumber: SequenceNumber,
				},
				ownerSwitchesSinceCheckpoint: OwnerSwitches,
				leaseKey:                     leaseKey,
				leaseCounter:                 leaseCounter,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.FromDynamoRecord(tt.args.record); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DynamoSerializer.FromDynamoRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDynamoSerializer_GetDynamoHashKey(t *testing.T) {
	type args struct {
		leaseKey string
	}
	tests := []struct {
		name string
		s    *DynamoSerializer
		args args
		want map[string]*dynamodb.AttributeValue
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetDynamoHashKey(tt.args.leaseKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DynamoSerializer.GetDynamoHashKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDynamoSerializer_GetDynamoLeaseCounterExpectation(t *testing.T) {
	type args struct {
		leaseCounter int64
	}
	tests := []struct {
		name string
		s    *DynamoSerializer
		args args
		want map[string]*dynamodb.ExpectedAttributeValue
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetDynamoLeaseCounterExpectation(tt.args.leaseCounter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DynamoSerializer.GetDynamoLeaseCounterExpectation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDynamoSerializer_GetDynamoLeaseOwnerExpectation(t *testing.T) {
	type args struct {
		leaseOwner string
	}
	tests := []struct {
		name string
		s    *DynamoSerializer
		args args
		want map[string]*dynamodb.ExpectedAttributeValue
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetDynamoLeaseOwnerExpectation(tt.args.leaseOwner); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DynamoSerializer.GetDynamoLeaseOwnerExpectation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDynamoSerializer_GetDynamoNonexistantExpectation(t *testing.T) {
	tests := []struct {
		name string
		s    *DynamoSerializer
		want map[string]*dynamodb.ExpectedAttributeValue
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetDynamoNonexistantExpectation(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DynamoSerializer.GetDynamoNonexistantExpectation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDynamoSerializer_GetDynamoLeaseCounterUpdate(t *testing.T) {
	type args struct {
		leaseCounter int64
	}
	tests := []struct {
		name string
		s    *DynamoSerializer
		args args
		want map[string]*dynamodb.AttributeValueUpdate
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetDynamoLeaseCounterUpdate(tt.args.leaseCounter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DynamoSerializer.GetDynamoLeaseCounterUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDynamoSerializer_GetDynamoTakeLeaseUpdate(t *testing.T) {
	type args struct {
		oldOwner   string
		leaseOwner string
	}
	tests := []struct {
		name string
		s    *DynamoSerializer
		args args
		want map[string]*dynamodb.AttributeValueUpdate
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetDynamoTakeLeaseUpdate(tt.args.oldOwner, tt.args.leaseOwner); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DynamoSerializer.GetDynamoTakeLeaseUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDynamoSerializer_GetDynamoEvictLeaseUpdate(t *testing.T) {
	tests := []struct {
		name string
		s    *DynamoSerializer
		want map[string]*dynamodb.AttributeValueUpdate
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetDynamoEvictLeaseUpdate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DynamoSerializer.GetDynamoEvictLeaseUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDynamoSerializer_GetDynamoUpdateLeaseUpdate(t *testing.T) {
	type args struct {
		lease *KLease
	}
	tests := []struct {
		name string
		s    *DynamoSerializer
		args args
		want map[string]*dynamodb.AttributeValueUpdate
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetDynamoUpdateLeaseUpdate(tt.args.lease); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DynamoSerializer.GetDynamoUpdateLeaseUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDynamoSerializer_GetKeySchema(t *testing.T) {
	tests := []struct {
		name string
		s    *DynamoSerializer
		want []*dynamodb.KeySchemaElement
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetKeySchema(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DynamoSerializer.GetKeySchema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDynamoSerializer_GetAttributeDefinitions(t *testing.T) {
	tests := []struct {
		name string
		s    *DynamoSerializer
		want []*dynamodb.AttributeDefinition
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetAttributeDefinitions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DynamoSerializer.GetAttributeDefinitions() = %v, want %v", got, tt.want)
			}
		})
	}
}
