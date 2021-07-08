package klease

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/ezoic/ezutil/amazon"
	"github.com/ezoic/log4go"
)

const (
	SNSArn = "arn:aws:sns:us-east-1:073320663057:data-matching"
)

func SendAlert(subject, message string) {

	snsSvc := sns.New(amazon.GetSession(amazon.DefaultRegion))

	input := &sns.PublishInput{
		Subject:  aws.String(subject),
		Message:  aws.String(message),
		TopicArn: aws.String(SNSArn),
	}
	_, err := snsSvc.Publish(input)
	if err != nil {
		log4go.Error("failed to post error message to sns: %s", err)
	}

}
