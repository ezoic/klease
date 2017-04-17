package klease

import "github.com/aws/aws-sdk-go/service/dynamodb"

type DynamoSerializer struct {
	utils DynamoUtils
}

const OwnerSwitchesKey = "ownerSwitchesSinceCheckpoint"
const CheckpointSequenceNumberKey = "checkpoint"

//const CheckpointSubsequenceNumberKey = "checkpointSubSequenceNumber"
const ParentShardIdKey = "parentShardId"
const LeaseKeyKey = "leaseKey"
const LeaseOwnerKey = "leaseOwner"
const LeaseCounterKey = "leaseCounter"

//NewDynamoSerializer creates and returns a new manager
func NewDynamoSerializer() *DynamoSerializer {
	serializer := &DynamoSerializer{
		utils: DynamoUtils{},
	}
	return serializer
}

func (s *DynamoSerializer) ToDynamoRecord(lease *KLease) map[string]*dynamodb.AttributeValue {
	result := map[string]*dynamodb.AttributeValue{}
	result[LeaseKeyKey] = s.utils.CreateAttributeValueS(lease.GetLeaseKey())
	result[LeaseCounterKey] = s.utils.CreateAttributeValueN(lease.GetLeaseCounter())
	if lease.GetLeaseOwner() != "" {
		result[LeaseOwnerKey] = s.utils.CreateAttributeValueS(lease.GetLeaseOwner())
	}

	result[OwnerSwitchesKey] = s.utils.CreateAttributeValueN(lease.GetOwnerSwitchesSinceCheckpoint())
	if lease.GetCheckpoint() != nil {
		result[CheckpointSequenceNumberKey] = s.utils.CreateAttributeValueS(lease.GetCheckpoint().SequenceNumber)
	}

	//TODO if we need to add subsequence
	//result[CheckpointSubsequenceNumberKey] = s.utils.CreateAttributeValueN(lease.GetLeaseCounter())
	//if lease.GetParentShardIds() != nil && len(lease.GetParentShardIds()) > 0 {
	//TODO figure out parent shard id keys
	//result[ParentShardIdKey] = s.utils.CreateAttributeValueS(lease.GetParentShardIds())
	//}

	return result
}

func (s *DynamoSerializer) FromDynamoRecord(record map[string]*dynamodb.AttributeValue) *KLease {
	result := &KLease{}

	result.SetLeaseKey(s.utils.SafeGetS(record, LeaseKeyKey))
	result.SetLeaseOwner(s.utils.SafeGetS(record, LeaseOwnerKey))
	result.SetLeaseCounter(s.utils.SafeGetN(record, LeaseCounterKey))
	result.SetOwnerSwitchesSinceCheckpoint(s.utils.SafeGetN(record, OwnerSwitchesKey))
	//TODO if we need to add subsequence
	result.SetCheckpoint(&KCheckpoint{
		SequenceNumber: s.utils.SafeGetS(record, CheckpointSequenceNumberKey),
	})

	return result
}

func (s *DynamoSerializer) GetDynamoHashKey(leaseKey string) map[string]*dynamodb.AttributeValue {
	result := map[string]*dynamodb.AttributeValue{}
	result[LeaseKeyKey] = s.utils.CreateAttributeValueS(leaseKey)
	return result
}

func (s *DynamoSerializer) GetDynamoLeaseCounterExpectation(leaseCounter int64) map[string]*dynamodb.ExpectedAttributeValue {
	result := map[string]*dynamodb.ExpectedAttributeValue{}
	result[LeaseCounterKey] = &dynamodb.ExpectedAttributeValue{
		Value: s.utils.CreateAttributeValueN(leaseCounter),
	}
	return result
}

func (s *DynamoSerializer) GetDynamoLeaseOwnerExpectation(leaseOwner string) map[string]*dynamodb.ExpectedAttributeValue {
	result := map[string]*dynamodb.ExpectedAttributeValue{}

	eva := &dynamodb.ExpectedAttributeValue{}
	if leaseOwner == "" {
		eva.SetExists(false)
	} else {
		eva.SetValue(s.utils.CreateAttributeValueS(leaseOwner))
	}

	result[LeaseOwnerKey] = eva
	return result
}

func (s *DynamoSerializer) GetDynamoNonexistantExpectation() map[string]*dynamodb.ExpectedAttributeValue {
	result := map[string]*dynamodb.ExpectedAttributeValue{}

	eva := &dynamodb.ExpectedAttributeValue{}
	eva.SetExists(false)
	result[LeaseKeyKey] = eva

	return result
}

func (s *DynamoSerializer) GetDynamoLeaseCounterUpdate(leaseCounter int64) map[string]*dynamodb.AttributeValueUpdate {
	result := map[string]*dynamodb.AttributeValueUpdate{}
	put := dynamodb.AttributeActionPut
	result[LeaseCounterKey] = &dynamodb.AttributeValueUpdate{
		Value:  s.utils.CreateAttributeValueN(leaseCounter + 1),
		Action: &put,
	}
	return result
}

func (s *DynamoSerializer) GetDynamoTakeLeaseUpdate(oldOwner, leaseOwner string) map[string]*dynamodb.AttributeValueUpdate {
	result := map[string]*dynamodb.AttributeValueUpdate{}
	put := dynamodb.AttributeActionPut
	result[LeaseOwnerKey] = &dynamodb.AttributeValueUpdate{
		Value:  s.utils.CreateAttributeValueS(leaseOwner),
		Action: &put,
	}

	if oldOwner != "" {
		add := dynamodb.AttributeActionAdd
		result[OwnerSwitchesKey] = &dynamodb.AttributeValueUpdate{
			Value:  s.utils.CreateAttributeValueN(1),
			Action: &add,
		}
	}

	return result
}

func (s *DynamoSerializer) GetDynamoEvictLeaseUpdate() map[string]*dynamodb.AttributeValueUpdate {
	result := map[string]*dynamodb.AttributeValueUpdate{}
	del := dynamodb.AttributeActionDelete
	result[LeaseOwnerKey] = &dynamodb.AttributeValueUpdate{
		Value:  nil,
		Action: &del,
	}
	return result
}

func (s *DynamoSerializer) GetDynamoUpdateLeaseUpdate(lease *KLease) map[string]*dynamodb.AttributeValueUpdate {
	result := map[string]*dynamodb.AttributeValueUpdate{}
	put := dynamodb.AttributeActionPut
	result[LeaseOwnerKey] = &dynamodb.AttributeValueUpdate{
		Value:  s.utils.CreateAttributeValueS(lease.GetLeaseOwner()),
		Action: &put,
	}
	//TODO if we need subsequence number
	result[CheckpointSequenceNumberKey] = &dynamodb.AttributeValueUpdate{
		Value:  s.utils.CreateAttributeValueS(lease.GetCheckpoint().SequenceNumber),
		Action: &put,
	}

	result[OwnerSwitchesKey] = &dynamodb.AttributeValueUpdate{
		Value:  s.utils.CreateAttributeValueN(lease.GetOwnerSwitchesSinceCheckpoint()),
		Action: &put,
	}

	return result
}

func (s *DynamoSerializer) GetKeySchema() []*dynamodb.KeySchemaElement {
	result := []*dynamodb.KeySchemaElement{}
	lkk := LeaseKeyKey
	typeHash := dynamodb.KeyTypeHash
	result = append(result, &dynamodb.KeySchemaElement{
		AttributeName: &lkk,
		KeyType:       &typeHash,
	})
	return result
}

func (s *DynamoSerializer) GetAttributeDefinitions() []*dynamodb.AttributeDefinition {
	result := []*dynamodb.AttributeDefinition{}
	lkk := LeaseKeyKey
	typeS := dynamodb.ScalarAttributeTypeS
	result = append(result, &dynamodb.AttributeDefinition{
		AttributeName: &lkk,
		AttributeType: &typeS,
	})
	return result
}
