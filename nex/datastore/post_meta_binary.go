package nex_datastore

import (
	"github.com/PretendoNetwork/badge-arcade-secure/database"
	"github.com/PretendoNetwork/badge-arcade-secure/globals"
	"github.com/PretendoNetwork/badge-arcade-secure/utility"

	nex "github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
)

func PostMetaBinary(err error, client *nex.Client, callID uint32, param *datastore.DataStorePreparePostParam) {
	pid := client.PID()
	var slot uint16 = 0

	switch param.DataType {
	case 100: // Free Play Data
		utility.DataStorePostParamToFreePlayData(pid, param)
	}

	database.PostUserPlayInfo(uint64(pid), pid, slot)

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	// We are using the PID as the Data ID as it is easier to handle
	rmcResponseStream.WriteUInt64LE(uint64(pid))

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore.MethodPostMetaBinary, rmcResponseBody)

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
