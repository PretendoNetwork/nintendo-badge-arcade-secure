package nex_datastore

import (
	"fmt"
	"os"

	"github.com/PretendoNetwork/badge-arcade-secure/database"
	"github.com/PretendoNetwork/badge-arcade-secure/globals"

	nex "github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
)

func PrepareUpdateObject(err error, client *nex.Client, callID uint32, param *datastore.DataStorePrepareUpdateParam) {
	dataID := param.DataID
	dataVersion := database.GetVersionByDataID(uint32(dataID))

	pReqUpdateInfo := datastore.NewDataStoreReqUpdateInfo()

	key := fmt.Sprintf("%s/%011d-%05d", os.Getenv("DATASTORE_DATA_PATH"), dataID, dataVersion+1)

	fieldKey := datastore.NewDataStoreKeyValue()
	fieldKey.Key = "key"
	fieldKey.Value = key

	fieldACL := datastore.NewDataStoreKeyValue()
	fieldACL.Key = "acl"
	fieldACL.Value = "private"

	fieldSignature := datastore.NewDataStoreKeyValue()
	fieldSignature.Key = "signature"
	fieldSignature.Value = "signature" // TODO

	pReqUpdateInfo.Version = dataVersion + 1
	pReqUpdateInfo.Url = fmt.Sprintf("http://%s.%s/", os.Getenv("S3_BUCKET_NAME"), os.Getenv("DATASTORE_DATA_URL"))
	pReqUpdateInfo.RequestHeaders = []*datastore.DataStoreKeyValue{}
	pReqUpdateInfo.FormFields = []*datastore.DataStoreKeyValue{fieldKey, fieldACL, fieldSignature}
	pReqUpdateInfo.RootCaCert = []byte{}

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteStructure(pReqUpdateInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore.MethodPrepareUpdateObject, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.NEXServer.Send(responsePacket)
}
