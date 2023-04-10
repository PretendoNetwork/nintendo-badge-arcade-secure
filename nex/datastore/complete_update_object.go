package nex_datastore

import (
	"github.com/PretendoNetwork/nintendo-badge-arcade-secure/database"
	"github.com/PretendoNetwork/nintendo-badge-arcade-secure/globals"

	nex "github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
)

func CompleteUpdateObject(err error, client *nex.Client, callID uint32, param *datastore.DataStoreCompleteUpdateParam) {
	// We update the version only if the update has been successful
	// This is done in order to prevent incomplete saves
	if param.IsSuccess {
		database.UpdateUserPlayInfoVersion(uint32(param.DataID), param.Version)
	}

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore.MethodCompleteUpdateObject, nil)

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
