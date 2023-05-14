package utility

import (
	"github.com/PretendoNetwork/nintendo-badge-arcade-secure/database"
	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
)

func FreePlayDataToDataStoreMetaInfo(ownerID uint32, dataType uint16) *datastore.DataStoreMetaInfo {
	dataID, metaBinary, createdTime, updatedTime, period, flag, referredTime := database.GetFreePlayDataMetaInfoByOwnerID(ownerID)

	if dataID == 0 {
		return nil
	}

	metaInfo := datastore.NewDataStoreMetaInfo()

	metaInfo.DataID = uint64(dataID)
	metaInfo.OwnerID = ownerID
	metaInfo.Size = 0 // What?
	metaInfo.Name = "FreePlayData"
	metaInfo.DataType = dataType
	metaInfo.MetaBinary = metaBinary
	metaInfo.Permission = datastore.NewDataStorePermission()
	metaInfo.Permission.Permission = 0 // Unknown
	metaInfo.Permission.RecipientIds = []uint32{}
	metaInfo.DelPermission = datastore.NewDataStorePermission()
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
	metaInfo.Ratings = []*datastore.DataStoreRatingInfoWithSlot{}

	return metaInfo
}
