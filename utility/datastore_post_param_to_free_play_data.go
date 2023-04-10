package utility

import (
	"github.com/PretendoNetwork/nintendo-badge-arcade-secure/database"
	nex "github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
)

func DataStorePostParamToFreePlayData(ownerID uint32, metaInfo *datastore.DataStorePreparePostParam) {
	dateTime := nex.NewDateTime(0)
	createdTime := dateTime.Now()

	// We are setting the PID as the Data ID, as it is easier to handle
	database.PostFreePlayDataMetaInfo(ownerID, ownerID, metaInfo.MetaBinary, createdTime, metaInfo.Period, metaInfo.Flag)
}
