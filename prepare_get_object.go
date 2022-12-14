package main

import (
	"fmt"
	"os"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func prepareGetObject(err error, client *nex.Client, callID uint32, dataStorePrepareGetParam *nexproto.DataStorePrepareGetParam) {
	pReqGetInfo := nexproto.NewDataStoreReqGetInfo()
	
	dataVersion := getVersionByDataID(dataStorePrepareGetParam.DataID)
	dataSize := getSizeByDataID(dataStorePrepareGetParam.DataID)

	pReqGetInfo.URL = fmt.Sprintf("https://%s.%s/%s/%011d-%05d", os.Getenv("S3_BUCKET_NAME"), os.Getenv("DATASTORE_DATA_URL"), os.Getenv("DATASTORE_DATA_PATH"), dataStorePrepareGetParam.DataID, dataVersion)
	pReqGetInfo.RequestHeaders = []*nexproto.DataStoreKeyValue{}
	pReqGetInfo.Size = dataSize
	pReqGetInfo.RootCA = []byte{}
	pReqGetInfo.DataID = dataStorePrepareGetParam.DataID

	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteStructure(pReqGetInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreBadgeArcadeProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreMethodPrepareGetObject, rmcResponseBody)

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
