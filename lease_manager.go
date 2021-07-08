package klease

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

//Manager supports basic CRUD operations for Leases
type Manager struct {
	table           string
	client          *dynamodb.DynamoDB
	serialzer       *DynamoSerializer
	consistentReads bool //should never be true except in some tests
}

//NewLeaseManager creates and returns a new manager
func NewLeaseManager(table string, client *dynamodb.DynamoDB, serialzer *DynamoSerializer) *Manager {

	if serialzer == nil {
		serialzer = NewDynamoSerializer()
	}

	manager := &Manager{
		table:           table,
		client:          client,
		serialzer:       serialzer,
		consistentReads: false,
	}
	return manager
}

func newLeaseManagerForTests(table string, client *dynamodb.DynamoDB, serialzer *DynamoSerializer, consistentReads bool) *Manager {

	if serialzer == nil {
		serialzer = NewDynamoSerializer()
	}

	manager := &Manager{
		table:           table,
		client:          client,
		serialzer:       serialzer,
		consistentReads: consistentReads,
	}
	return manager
}

func (k *Manager) CreateLeaseTableIfNotExists() error {
	//l4g.Debug("makin table: %s", k.table)
	request := &dynamodb.CreateTableInput{}
	request.SetTableName(k.table)
	request.SetKeySchema(k.serialzer.GetKeySchema())
	request.SetAttributeDefinitions(k.serialzer.GetAttributeDefinitions())
	request.BillingMode = aws.String("PAY_PER_REQUEST")

	_, err := k.client.CreateTable(request)
	if err != nil {
		SendAlert(
			fmt.Sprintf("Failed To Create Lease Table For %s", k.table),
			fmt.Sprintf("Error creating lease table for %s: %s", k.table, err),
		)
		if strings.Contains(err.Error(), dynamodb.ErrCodeResourceInUseException) == false {
			return err
		}
	}

	return nil
}

func (k *Manager) LeaseTableExists() (bool, error) {
	tableStatus, err := k.tableStatus()
	if err != nil {
		return false, err
	}
	return tableStatus == dynamodb.TableStatusActive, nil
}

func (k *Manager) tableStatus() (string, error) {
	request := &dynamodb.DescribeTableInput{}
	request = request.SetTableName(k.table)

	result := &dynamodb.DescribeTableOutput{}
	var err error
	result, err = k.client.DescribeTable(request)
	if err != nil {
		return "", err
	}

	return *result.Table.TableStatus, nil
}

func (k *Manager) WaitUntilLeaseTableExists(secondsBetweenPulls, timeoutSeconds int64) (bool, error) {
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

func (k *Manager) ListLeases() ([]*KLease, error) {
	return k.list(-1)
}

func (k *Manager) ListShards() ([]string, error) {
	leases, err := k.list(-1)
	if err != nil {
		return []string{}, err
	}
	leasesSlice := make([]string, len(leases))

	idx := 0
	for _, value := range leases {
		leasesSlice[idx] = value.GetLeaseKey()
		idx++
	}
	return leasesSlice, nil
}

func (k *Manager) IsLeaseTableEmpty() (bool, error) {
	leases, err := k.list(1)
	if err != nil {
		return true, err
	}
	return len(leases) < 1, nil
}

//internal lease lister
func (k *Manager) list(limit int64) ([]*KLease, error) {

	request := &dynamodb.ScanInput{}
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

func (k *Manager) CreateLeaseIfNotExists(lease *KLease) (bool, error) {
	if lease == nil {
		return false, errors.New("lease cannot be nil")
	}

	request := &dynamodb.PutItemInput{}
	request.SetTableName(k.table)
	request.SetItem(k.serialzer.ToDynamoRecord(lease))
	request.SetExpected(k.serialzer.GetDynamoNonexistantExpectation())

	_, err := k.client.PutItem(request)
	if err != nil {
		if isErrCodeConditionalCheckFailedException(err) {
			//failed because it already existed
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (k *Manager) GetLease(leaseKey string) (*KLease, error) {

	if leaseKey == "" {
		return nil, errors.New("leaseKey cannot be empty")
	}

	request := &dynamodb.GetItemInput{}
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

func (k *Manager) RenewLease(lease *KLease) (bool, error) {

	if lease == nil {
		return false, errors.New("lease cannot be nil")
	}

	request := &dynamodb.UpdateItemInput{}
	request.SetTableName(k.table)
	request.SetKey(k.serialzer.GetDynamoHashKey(lease.GetLeaseKey()))
	request.SetExpected(k.serialzer.GetDynamoLeaseCounterExpectation(lease.GetLeaseCounter()))
	request.SetAttributeUpdates(k.serialzer.GetDynamoLeaseCounterUpdate(lease.GetLeaseCounter()))

	_, err := k.client.UpdateItem(request)
	if err != nil {
		if isErrCodeConditionalCheckFailedException(err) {
			return false, nil
		}
		return false, err
	}

	lease.SetLeaseCounter(lease.GetLeaseCounter() + 1)
	return true, nil
}

func (k *Manager) TakeLease(lease *KLease, owner string) (bool, error) {
	if lease == nil {
		return false, errors.New("lease cannot be nil")
	}
	if owner == "" {
		return false, errors.New("owner cannot be empty")
	}

	request := &dynamodb.UpdateItemInput{}
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
		if isErrCodeConditionalCheckFailedException(err) {
			return false, nil
		}
		return false, err
	}

	lease.SetLeaseCounter(lease.GetLeaseCounter() + 1)
	lease.SetLeaseOwner(owner)
	return true, nil

}

func (k *Manager) EvictLease(lease *KLease) (bool, error) {
	if lease == nil {
		return false, errors.New("lease cannot be nil")
	}

	request := &dynamodb.UpdateItemInput{}
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
		if isErrCodeConditionalCheckFailedException(err) {
			return false, nil
		}
		return false, err
	}

	lease.SetLeaseCounter(lease.GetLeaseCounter() + 1)
	lease.SetLeaseOwner("")
	return true, nil
}

func (k *Manager) DeleteAll() error {
	allLeases, err := k.ListLeases()
	if err != nil {
		return err
	}
	for _, lease := range allLeases {
		request := &dynamodb.DeleteItemInput{}
		request.SetTableName(k.table)
		request.SetKey(k.serialzer.GetDynamoHashKey(lease.GetLeaseKey()))

		_, err := k.client.DeleteItem(request)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *Manager) DeleteLease(lease *KLease) error {
	if lease == nil {
		return errors.New("lease cannot be nil")
	}

	request := &dynamodb.DeleteItemInput{}
	request.SetTableName(k.table)
	request.SetKey(k.serialzer.GetDynamoHashKey(lease.GetLeaseKey()))

	_, err := k.client.DeleteItem(request)

	return err

}

func (k *Manager) UpdateLease(lease *KLease) (bool, error) {
	if lease == nil {
		return false, errors.New("lease cannot be nil")
	}

	request := &dynamodb.UpdateItemInput{}
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
		if isErrCodeConditionalCheckFailedException(err) {
			return false, nil
		}
		return false, err
	}

	lease.SetLeaseCounter(lease.GetLeaseCounter() + 1)

	return true, nil

}

func (k *Manager) GetCheckpoint(shardId string) (*KCheckpoint, error) {
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

func isErrCodeConditionalCheckFailedException(err error) bool {
	return strings.HasPrefix(err.Error(), dynamodb.ErrCodeConditionalCheckFailedException)
}
