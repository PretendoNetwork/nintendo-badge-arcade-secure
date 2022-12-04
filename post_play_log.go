package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func postPlayLog(err error, client *nex.Client, callID uint32, param *nexproto.ShopPostPlayLogParam) {
	// TODO: Implement Custom ID to RMCResponse and do something with the data
	rmcResponse := nex.NewRMCResponse(nexproto.ShopBadgeArcadeProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.ShopBadgeArcadeMethodPostPlayLog, nil)

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