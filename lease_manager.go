package klease

import "github.com/aws/aws-sdk-go/service/dynamodb"
import "time"
import "errors"

//KLeaseManager supports basic CRUD operations for Leases
type KLeaseManager struct {
	table           string
	client          *dynamodb.DynamoDB
	serialzer       *DynamoSerializer
	consistentReads bool //should never be true except in some tests
}

//NewKLeaseManager creates and returns a new manager
func NewKLeaseManager(table string, client *dynamodb.DynamoDB, serialzer *DynamoSerializer) KLeaseManager {

	if serialzer == nil {
		serialzer = NewDynamoSerializer()
	}

	manager := KLeaseManager{
		table:           table,
		client:          client,
		serialzer:       serialzer,
		consistentReads: false,
	}
	return manager
}

func newKLeaseManagerForTests(table string, client *dynamodb.DynamoDB, serialzer *DynamoSerializer, consistentReads bool) KLeaseManager {

	if serialzer == nil {
		serialzer = NewDynamoSerializer()
	}

	manager := KLeaseManager{
		table:           table,
		client:          client,
		serialzer:       serialzer,
		consistentReads: consistentReads,
	}
	return manager
}

func (k *KLeaseManager) CreateLeaseTableIfNotExists(readCapacity, writeCapacity int64) error {

	var request *dynamodb.CreateTableInput
	request.SetTableName(k.table)
	request.SetKeySchema(k.serialzer.GetKeySchema())
	request.SetAttributeDefinitions(k.serialzer.GetAttributeDefinitions())

	var throughput *dynamodb.ProvisionedThroughput
	throughput.SetReadCapacityUnits(readCapacity)
	throughput.SetWriteCapacityUnits(writeCapacity)
	request.SetProvisionedThroughput(throughput)

	_, err := k.client.CreateTable(request)
	if err != nil {
		if err.Error() != dynamodb.ErrCodeResourceInUseException {
			return err
		}
	}

	return nil
}

func (k *KLeaseManager) LeaseTableExists() (bool, error) {
	tableStatus, err := k.tableStatus()
	if err != nil {
		return false, err
	}
	return tableStatus == dynamodb.TableStatusActive, nil
}

func (k *KLeaseManager) tableStatus() (string, error) {
	var request *dynamodb.DescribeTableInput
	request = request.SetTableName(k.table)

	var result *dynamodb.DescribeTableOutput
	var err error
	result, err = k.client.DescribeTable(request)
	if err != nil {
		return "", err
	}

	return *result.Table.TableStatus, nil
}

func (k *KLeaseManager) WaitUntilLeaseTableExists(secondsBetweenPulls, timeoutSeconds int64) (bool, error) {
	secondsLeft := time.Duration(timeoutSeconds)
	durationBetweenPulls := time.Duration(secondsBetweenPulls)
	for {
		exists, err := k.LeaseTableExists()
		if err != nil {
			return false, err
		}

		if exists {
			return true, nil
		}

		if secondsLeft <= 0 {
			return false, nil
		}

		timeToSleep := min(durationBetweenPulls, secondsLeft)
		time.Sleep(timeToSleep * time.Second)
		secondsLeft -= timeToSleep
	}

}

func (k *KLeaseManager) ListLeases() ([]*KLease, error) {
	return k.list(-1)
}

func (k *KLeaseManager) IsLeaseTableEmpty() (bool, error) {
	leases, err := k.list(1)
	if err != nil {
		return true, err
	}
	return len(leases) < 1, nil
}

//internal lease lister
func (k *KLeaseManager) list(limit int64) ([]*KLease, error) {

	var request *dynamodb.ScanInput
	request.SetTableName(k.table)
	if limit >= 0 {
		request.SetLimit(limit)
	}

	result := []*KLease{}
	output, err := k.client.Scan(request)
	if err != nil {
		return result, err
	}

	for output != nil {
		//grab the actual results
		for _, item := range output.Items {
			result = append(result, k.serialzer.FromDynamoRecord(item))
		}

		//if this was the last page get out
		if len(output.LastEvaluatedKey) <= 0 {
			output = nil
		} else {
			//otherwise grab the next page of results
			request.SetExclusiveStartKey(output.LastEvaluatedKey)
			output, err = k.client.Scan(request)
			if err != nil {
				return result, err
			}
		}

	}

	return result, nil
}

func (k *KLeaseManager) CreateLeaseIfNotExists(lease *KLease) (bool, error) {
	if lease == nil {
		return false, errors.New("lease cannot be nil")
	}

	var request *dynamodb.PutItemInput
	request.SetTableName(k.table)
	request.SetItem(k.serialzer.ToDynamoRecord(*lease))
	request.SetExpected(k.serialzer.GetDynamoNonexistantExpectation())

	_, err := k.client.PutItem(request)
	if err != nil {
		if err.Error() == dynamodb.ErrCodeConditionalCheckFailedException {
			//failed because it already existed
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (k *KLeaseManager) GetLease(leaseKey string) (*KLease, error) {

	if leaseKey == "" {
		return nil, errors.New("leaseKey cannot be empty")
	}

	var request *dynamodb.GetItemInput
	request.SetTableName(k.table)
	request.SetKey(k.serialzer.GetDynamoHashKey(leaseKey))
	request.SetConsistentRead(k.consistentReads)

	result, err := k.client.GetItem(request)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	return k.serialzer.FromDynamoRecord(result.Item), nil
}

func (k *KLeaseManager) RenewLease(lease *KLease) (bool, error) {

	if lease == nil {
		return false, errors.New("lease cannot be nil")
	}

	var request *dynamodb.UpdateItemInput
	request.SetTableName(k.table)
	request.SetKey(k.serialzer.GetDynamoHashKey(lease.GetLeaseKey()))
	request.SetExpected(k.serialzer.GetDynamoLeaseCounterExpectation(lease.GetLeaseCounter()))
	request.SetAttributeUpdates(k.serialzer.GetDynamoLeaseCounterUpdate(lease.GetLeaseCounter()))

	_, err := k.client.UpdateItem(request)
	if err != nil {
		if err.Error() == dynamodb.ErrCodeConditionalCheckFailedException {
			return false, nil
		}
		return false, err
	}

	lease.SetLeaseCounter(lease.GetLeaseCounter() + 1)
	return true, nil
}

func (k *KLeaseManager) TakeLease(lease *KLease, owner string) (bool, error) {
	if lease == nil {
		return false, errors.New("lease cannot be nil")
	}
	if owner == "" {
		return false, errors.New("owner cannot be empty")
	}

	var request *dynamodb.UpdateItemInput
	request.SetTableName(k.table)
	request.SetKey(k.serialzer.GetDynamoHashKey(lease.GetLeaseKey()))
	request.SetExpected(k.serialzer.GetDynamoLeaseCounterExpectation(lease.GetLeaseCounter()))

	updates := k.serialzer.GetDynamoLeaseCounterUpdate(lease.GetLeaseCounter())
	ownerUpdate := k.serialzer.GetDynamoTakeLeaseUpdate(lease.GetLeaseOwner(), owner)
	//merging the two update maps
	for k, v := range ownerUpdate {
		updates[k] = v
	}

	request.SetAttributeUpdates(updates)

	if lease.GetLeaseOwner() != "" && lease.GetLeaseOwner() != owner {
		lease.SetOwnerSwitchesSinceCheckpoint(lease.GetOwnerSwitchesSinceCheckpoint() + 1)
	}

	_, err := k.client.UpdateItem(request)
	if err != nil {
		if err.Error() == dynamodb.ErrCodeConditionalCheckFailedException {
			return false, nil
		}
		return false, err
	}

	lease.SetLeaseCounter(lease.GetLeaseCounter() + 1)
	lease.SetLeaseOwner(owner)
	return true, nil

}

func (k *KLeaseManager) EvictLease(lease *KLease) (bool, error) {
	if lease == nil {
		return false, errors.New("lease cannot be nil")
	}

	var request *dynamodb.UpdateItemInput
	request.SetTableName(k.table)
	request.SetKey(k.serialzer.GetDynamoHashKey(lease.GetLeaseKey()))
	request.SetExpected(k.serialzer.GetDynamoLeaseOwnerExpectation(lease.GetLeaseOwner()))

	updates := k.serialzer.GetDynamoLeaseCounterUpdate(lease.GetLeaseCounter())
	ownerUpdate := k.serialzer.GetDynamoEvictLeaseUpdate()
	//merging the two update maps
	for k, v := range ownerUpdate {
		updates[k] = v
	}

	request.SetAttributeUpdates(updates)

	_, err := k.client.UpdateItem(request)
	if err != nil {
		if err.Error() == dynamodb.ErrCodeConditionalCheckFailedException {
			return false, nil
		}
		return false, err
	}

	lease.SetLeaseCounter(lease.GetLeaseCounter() + 1)
	lease.SetLeaseOwner("")
	return true, nil
}

func (k *KLeaseManager) DeleteAll() error {
	allLeases, err := k.ListLeases()
	if err != nil {
		return err
	}
	for _, lease := range allLeases {
		var request *dynamodb.DeleteItemInput
		request.SetTableName(k.table)
		request.SetKey(k.serialzer.GetDynamoHashKey(lease.GetLeaseKey()))

		_, err := k.client.DeleteItem(request)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *KLeaseManager) DeleteLease(lease *KLease) error {
	if lease == nil {
		return errors.New("lease cannot be nil")
	}

	var request *dynamodb.DeleteItemInput
	request.SetTableName(k.table)
	request.SetKey(k.serialzer.GetDynamoHashKey(lease.GetLeaseKey()))

	_, err := k.client.DeleteItem(request)

	return err

}

func (k *KLeaseManager) UpdateLease(lease *KLease) (bool, error) {
	if lease == nil {
		return false, errors.New("lease cannot be nil")
	}

	var request *dynamodb.UpdateItemInput
	request.SetTableName(k.table)
	request.SetKey(k.serialzer.GetDynamoHashKey(lease.GetLeaseKey()))
	request.SetExpected(k.serialzer.GetDynamoLeaseCounterExpectation(lease.GetLeaseCounter()))

	updates := k.serialzer.GetDynamoLeaseCounterUpdate(lease.GetLeaseCounter())
	leaseUpdate := k.serialzer.GetDynamoUpdateLeaseUpdate(lease)
	//merging the two update maps
	for k, v := range leaseUpdate {
		updates[k] = v
	}

	request.SetAttributeUpdates(updates)

	_, err := k.client.UpdateItem(request)
	if err != nil {
		if err.Error() == dynamodb.ErrCodeConditionalCheckFailedException {
			return false, nil
		}
		return false, err
	}

	lease.SetLeaseCounter(lease.GetLeaseCounter() + 1)

	return true, nil

}

func (k *KLeaseManager) GetCheckpoint(shardId string) (*KCheckpoint, error) {
	lease, err := k.GetLease(shardId)
	if err != nil {
		return nil, err
	}
	if lease != nil {
		return lease.GetCheckpoint(), nil
	}
	return nil, nil
}

func min(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
