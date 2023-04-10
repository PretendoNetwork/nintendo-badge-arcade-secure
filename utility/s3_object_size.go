package utility

import (
	"github.com/PretendoNetwork/nintendo-badge-arcade-secure/globals"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func S3ObjectSize(bucket, key string) (uint64, error) {
	headObj := s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	result, err := globals.S3Client.HeadObject(&headObj)
	if err != nil {
		return 0, err
	}

	return uint64(aws.Int64Value(result.ContentLength)), nil
}
