package main

import (
	"fmt"
	"os"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func prepareUpdateObject(err error, client *nex.Client, callID uint32, param *nexproto.DataStorePrepareUpdateParam) {
	dataID := param.DataID
	dataSize := param.Size
	dataVersion := getVersionByDataID(dataID)

	// TODO: This isn't a safe way for handling this. If the S3 server is down and the user quits,
	// it may lead to an incomplete save and cause errors!
	updateUserPlayInfoSize(dataID, dataSize)
	
	pReqUpdateInfo := nexproto.NewDataStoreReqUpdateInfo()

	key := fmt.Sprintf("%s/%011d-%05d", os.Getenv("DATASTORE_DATA_PATH"), dataID, dataVersion + 1)

	fieldKey := nexproto.NewDataStoreKeyValue()
	fieldKey.Key = "key"
	fieldKey.Value = key

	fieldACL := nexproto.NewDataStoreKeyValue()
	fieldACL.Key = "acl"
	fieldACL.Value = "private"

	fieldSignature := nexproto.NewDataStoreKeyValue()
	fieldSignature.Key = "signature"
	fieldSignature.Value = "signature" // TODO

	pReqUpdateInfo.Version = dataVersion + 1
	pReqUpdateInfo.Url = fmt.Sprintf("https://%s.%s/", os.Getenv("S3_BUCKET_NAME"), os.Getenv("DATASTORE_DATA_URL"))
	pReqUpdateInfo.RequestHeaders = []*nexproto.DataStoreKeyValue{}
	pReqUpdateInfo.FormFields = []*nexproto.DataStoreKeyValue{fieldKey, fieldACL, fieldSignature}
	pReqUpdateInfo.RootCaCert = []byte{}

	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteStructure(pReqUpdateInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreBadgeArcadeProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreMethodPrepareUpdateObject, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	nexServer.Send(responsePacket)
}
