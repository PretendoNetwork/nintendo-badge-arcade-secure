package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func getMaintenanceStatus(err error, client *nex.Client, callID uint32) {
	// TODO: Don't use hardcoded variables
	var maintenanceStatus uint16 = 0xFFFF
	var maintenanceTime uint32 = 0 // Time of maintenance in Unix time (I think)
	var isSuccess bool = true

	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteUInt16LE(maintenanceStatus)
	rmcResponseStream.WriteUInt32LE(maintenanceTime)
	rmcResponseStream.WriteBool(isSuccess)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.SecureBadgeArcadeProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.SecureBadgeArcadeMethodGetMaintenanceStatus, rmcResponseBody)

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
