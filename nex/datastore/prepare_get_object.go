package nex_datastore

import (
	"fmt"
	"os"

	"github.com/PretendoNetwork/nintendo-badge-arcade-secure/database"
	"github.com/PretendoNetwork/nintendo-badge-arcade-secure/globals"

	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
)

func PrepareGetObject(err error, client *nex.Client, callID uint32, dataStorePrepareGetParam *datastore.DataStorePrepareGetParam) {
	pReqGetInfo := datastore.NewDataStoreReqGetInfo()

	dataVersion := database.GetVersionByDataID(uint32(dataStorePrepareGetParam.DataID))

	key := fmt.Sprintf("%s/%011d-%05d", os.Getenv("PN_NBA_CONFIG_S3_PATH"), dataStorePrepareGetParam.DataID, dataVersion)
	dataSize, err := globals.S3ObjectSize(os.Getenv("PN_NBA_CONFIG_S3_BUCKET"), key)
	if err != nil {
		globals.Logger.Error(err.Error())
	}

	pReqGetInfo.URL = fmt.Sprintf("%s/%s", os.Getenv("PN_NBA_CONFIG_S3_ENDPOINT"), key)
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
