package nex_datastore

import (
	"fmt"

	"github.com/PretendoNetwork/badge-arcade-secure/globals"
	"github.com/PretendoNetwork/badge-arcade-secure/utility"

	nex "github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"
)

func ChangeMeta(err error, client *nex.Client, callID uint32, param *datastore.DataStoreChangeMetaParam) {
	// TODO: This doesn't seem right
	switch param.DataType {
	case 0: // Free Play Data?
		utility.ChangeFreePlayDataMeta(param.DataID, param.MetaBinary)
	default:
		fmt.Println("WARNING: Unknown DataType: %d", param.DataType)
	}

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore.MethodChangeMeta, nil)

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
