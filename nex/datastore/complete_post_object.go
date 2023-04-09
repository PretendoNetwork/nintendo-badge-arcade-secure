package nex_datastore

import (
	"github.com/PretendoNetwork/badge-arcade-secure/database"
	"github.com/PretendoNetwork/badge-arcade-secure/globals"

	nex "github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
)

func CompletePostObject(err error, client *nex.Client, callID uint32, param *datastore.DataStoreCompletePostParam) {
	// We update the version only if the post has been successful
	// This is done in order to prevent incomplete saves
	if param.IsSuccess {
		var initialVersion uint32 = 1
		database.UpdateUserPlayInfoVersion(param.DataID, initialVersion)
	}

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore.MethodCompletePostObject, nil)

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
