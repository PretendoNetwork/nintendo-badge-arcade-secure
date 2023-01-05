package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func freePlayDataToDataStoreMetaInfo(ownerID uint32, dataType uint16) *nexproto.DataStoreMetaInfo {
	dataID, metaBinary, createdTime, updatedTime, period, flag, referredTime := getFreePlayDataMetaInfoByOwnerID(ownerID)

	if dataID == 0 {
		return nil
	}

	metaInfo := nexproto.NewDataStoreMetaInfo()

	metaInfo.DataID = dataID
	metaInfo.OwnerID = ownerID
	metaInfo.Size = 0 // What?
	metaInfo.Name = "FreePlayData"
	metaInfo.DataType = dataType
	metaInfo.MetaBinary = metaBinary
	metaInfo.Permission = nexproto.NewDataStorePermission()
	metaInfo.Permission.Permission = 0 // Unknown
	metaInfo.Permission.RecipientIds = []uint32{}
	metaInfo.DelPermission = nexproto.NewDataStorePermission()
	metaInfo.DelPermission.Permission = 3 // Unknown
	metaInfo.DelPermission.RecipientIds = []uint32{}
	metaInfo.CreatedTime = nex.NewDateTime(createdTime)
	metaInfo.UpdatedTime = nex.NewDateTime(updatedTime)
	metaInfo.Period = period
	metaInfo.Status = 0      // Unknown
	metaInfo.ReferredCnt = 0 // Unknown
	metaInfo.ReferDataID = 0 // Unknown
	metaInfo.Flag = flag
	metaInfo.ReferredTime = nex.NewDateTime(referredTime)
	metaInfo.ExpireTime = nex.NewDateTime(671075926016) // December 31st, year 9999
	metaInfo.Tags = []string{}
	metaInfo.Ratings = []*nexproto.DataStoreRatingInfoWithSlot{}

	return metaInfo
}

func dataStorePostParamToFreePlayData(ownerID uint32, metaInfo *nexproto.DataStorePreparePostParam) {
	dateTime := nex.NewDateTime(0)
	createdTime := dateTime.Now()

	// We are setting the PID as the Data ID, as it is easier to handle
	postFreePlayDataMetaInfo(uint64(ownerID), ownerID, metaInfo.MetaBinary, createdTime, metaInfo.Period, metaInfo.Flag)
}

func changeFreePlayDataMeta(dataID uint64, metaBinary []byte) {
	dateTime := nex.NewDateTime(0)
	updatedTime := dateTime.Now()

	updateFreePlayDataMetaBinary(dataID, metaBinary, updatedTime)
}

func s3ObjectSize(bucket, key string) (uint64, error) {
	headObj := s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key: aws.String(key),
	}

	result, err := s3Client.HeadObject(&headObj)
	if err != nil {
		return 0, err
	}

	return uint64(aws.Int64Value(result.ContentLength)), nil
}
