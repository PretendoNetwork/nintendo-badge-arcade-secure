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

	key := fmt.Sprintf("%s/%011d-%05d", os.Getenv("DATASTORE_DATA_PATH"), dataStorePrepareGetParam.DataID, dataVersion)
	dataSize, _ := s3ObjectSize(os.Getenv("S3_BUCKET_NAME"), key)

	pReqGetInfo.URL = fmt.Sprintf("http://%s.%s/%s", os.Getenv("S3_BUCKET_NAME"), os.Getenv("DATASTORE_DATA_URL"), key)
	pReqGetInfo.RequestHeaders = []*nexproto.DataStoreKeyValue{}
	pReqGetInfo.Size = uint32(dataSize)
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
