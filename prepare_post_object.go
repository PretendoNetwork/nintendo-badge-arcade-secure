package main

import (
	"fmt"
	"os"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func preparePostObject(err error, client *nex.Client, callID uint32, param *nexproto.DataStorePreparePostParam) {
	pid := client.PID()
	var slot uint16 = 0

	dataID := getDataStorePersistenceInfo(pid, slot)
	dataSize := param.Size
	var dataVersion uint32 = 1

	// TODO: This isn't a safe way for handling this. If the S3 server is down and the user quits,
	// it may lead to an incomplete save and cause errors!
	postUserPlayInfo(dataID, dataVersion, dataSize)
	
	pReqPostInfo := nexproto.NewDataStoreReqPostInfo()

	key := fmt.Sprintf("%s/%011d-%05d", os.Getenv("DATASTORE_DATA_PATH"), dataID, dataVersion)

	fieldKey := nexproto.NewDataStoreKeyValue()
	fieldKey.Key = "key"
	fieldKey.Value = key

	fieldACL := nexproto.NewDataStoreKeyValue()
	fieldACL.Key = "acl"
	fieldACL.Value = "private"

	fieldSignature := nexproto.NewDataStoreKeyValue()
	fieldSignature.Key = "signature"
	fieldSignature.Value = "signature" // TODO

	pReqPostInfo.URL = fmt.Sprintf("https://%s.%s/", os.Getenv("S3_BUCKET_NAME"), os.Getenv("DATASTORE_DATA_URL"))
	pReqPostInfo.RequestHeaders = []*nexproto.DataStoreKeyValue{}
	pReqPostInfo.FormFields = []*nexproto.DataStoreKeyValue{fieldKey, fieldACL, fieldSignature}
	pReqPostInfo.RootCACert = []byte{}

	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteStructure(pReqPostInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreBadgeArcadeProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreMethodPreparePostObject, rmcResponseBody)

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
