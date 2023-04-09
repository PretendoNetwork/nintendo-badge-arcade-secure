package nex_datastore

import (
	"fmt"
	"os"

	"github.com/PretendoNetwork/badge-arcade-secure/database"
	"github.com/PretendoNetwork/badge-arcade-secure/globals"
	"github.com/PretendoNetwork/badge-arcade-secure/utility"

	nex "github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
)

func PrepareGetObject(err error, client *nex.Client, callID uint32, dataStorePrepareGetParam *datastore.DataStorePrepareGetParam) {
	pReqGetInfo := datastore.NewDataStoreReqGetInfo()

	dataVersion := database.GetVersionByDataID(dataStorePrepareGetParam.DataID)

	key := fmt.Sprintf("%s/%011d-%05d", os.Getenv("DATASTORE_DATA_PATH"), dataStorePrepareGetParam.DataID, dataVersion)
	dataSize, _ := utility.S3ObjectSize(os.Getenv("S3_BUCKET_NAME"), key)

	pReqGetInfo.URL = fmt.Sprintf("http://%s.%s/%s", os.Getenv("S3_BUCKET_NAME"), os.Getenv("DATASTORE_DATA_URL"), key)
	pReqGetInfo.RequestHeaders = []*datastore.DataStoreKeyValue{}
	pReqGetInfo.Size = uint32(dataSize)
	pReqGetInfo.RootCA = []byte{}
	pReqGetInfo.DataID = dataStorePrepareGetParam.DataID

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteStructure(pReqGetInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore.MethodPrepareGetObject, rmcResponseBody)

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
