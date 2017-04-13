package klease

import "github.com/aws/aws-sdk-go/service/dynamodb"
import "strconv"

type DynamoUtils struct{}

func (u *DynamoUtils) CreateAttributeValueS(s string) *dynamodb.AttributeValue {
	t := dynamodb.AttributeValue{}
	return t.SetS(s)
}

func (u *DynamoUtils) CreateAttributeValueSS(s []*string) *dynamodb.AttributeValue {
	t := dynamodb.AttributeValue{}
	return t.SetSS(s)
}

func (u *DynamoUtils) CreateAttributeValueN(n int64) *dynamodb.AttributeValue {
	t := dynamodb.AttributeValue{}
	return t.SetN(strconv.FormatInt(n, 10))
}

func (u *DynamoUtils) SafeGetN(record map[string]*dynamodb.AttributeValue, key string) int64 {
	if av, ok := record[key]; ok {
		n, err := strconv.ParseInt(*av.N, 10, 64)
		if err != nil {
			//todo
			n = 0
		}
		return n
	}
	return 0
}

func (u *DynamoUtils) SafeGetS(record map[string]*dynamodb.AttributeValue, key string) string {
	if av, ok := record[key]; ok {
		return *av.S
	}
	return ""
}

func (u *DynamoUtils) SafeGetSS(record map[string]*dynamodb.AttributeValue, key string) []*string {
	if av, ok := record[key]; ok {
		return av.SS
	}
	return nil
}
