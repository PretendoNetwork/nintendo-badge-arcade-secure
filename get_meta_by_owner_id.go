package main

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func getMetaByOwnerId(err error, client *nex.Client, callID uint32, param *nexproto.DataStoreGetMetaByOwnerIdParam) {
	pMetaInfo := make([]*nexproto.DataStoreMetaInfo, 0)
	var pHasNext bool = false // Unknown

	if len(param.OwnerIDs) != len(param.DataTypes) {
		// Not sure if this is possible in the first place
		fmt.Println("WARNING: Owner ID and DataType length mismatch")
	}
	
	for i, _ := range param.OwnerIDs {
		switch param.DataTypes[i] {
		case 100: // Free Play Data
			pMetaInfo = append(pMetaInfo, freePlayDataToDataStoreMetaInfo(param.OwnerIDs[i], param.DataTypes[i]))
		default:
			fmt.Println("WARNING: Unknown DataType: %d", param.DataTypes[i])
		}
	}

	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteListStructure(pMetaInfo)
	rmcResponseStream.WriteBool(pHasNext)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreBadgeArcadeProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreBadgeArcadeMethodGetMetaByOwnerId, rmcResponseBody)

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