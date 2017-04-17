package klease

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func TestDynamoUtils_CreateAttributeValueS(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		u    *DynamoUtils
		args args
		want *dynamodb.AttributeValue
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.CreateAttributeValueS(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DynamoUtils.CreateAttributeValueS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDynamoUtils_CreateAttributeValueSS(t *testing.T) {
	type args struct {
		s []*string
	}
	tests := []struct {
		name string
		u    *DynamoUtils
		args args
		want *dynamodb.AttributeValue
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.CreateAttributeValueSS(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DynamoUtils.CreateAttributeValueSS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDynamoUtils_CreateAttributeValueN(t *testing.T) {
	type args struct {
		n int64
	}
	tests := []struct {
		name string
		u    *DynamoUtils
		args args
		want *dynamodb.AttributeValue
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.CreateAttributeValueN(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DynamoUtils.CreateAttributeValueN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDynamoUtils_SafeGetN(t *testing.T) {
	type args struct {
		record map[string]*dynamodb.AttributeValue
		key    string
	}
	tests := []struct {
		name string
		u    *DynamoUtils
		args args
		want int64
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.SafeGetN(tt.args.record, tt.args.key); got != tt.want {
				t.Errorf("DynamoUtils.SafeGetN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDynamoUtils_SafeGetS(t *testing.T) {
	type args struct {
		record map[string]*dynamodb.AttributeValue
		key    string
	}
	tests := []struct {
		name string
		u    *DynamoUtils
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.SafeGetS(tt.args.record, tt.args.key); got != tt.want {
				t.Errorf("DynamoUtils.SafeGetS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDynamoUtils_SafeGetSS(t *testing.T) {
	type args struct {
		record map[string]*dynamodb.AttributeValue
		key    string
	}
	tests := []struct {
		name string
		u    *DynamoUtils
		args args
		want []*string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.SafeGetSS(tt.args.record, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DynamoUtils.SafeGetSS() = %v, want %v", got, tt.want)
			}
		})
	}
}
