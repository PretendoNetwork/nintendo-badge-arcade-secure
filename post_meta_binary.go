package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func postMetaBinary(err error, client *nex.Client, callID uint32, param *nexproto.DataStorePreparePostParam) {
	pid := client.PID()
	var slot uint16 = 0

	switch param.DataType {
		case 100: // Free Play Data
		dataStorePostParamToFreePlayData(pid, param)
	}

	postPersistenceInfo(uint64(pid), pid, slot)

	rmcResponseStream := nex.NewStreamOut(nexServer)

	// We are using the PID as the Data ID as it is easier to handle
	rmcResponseStream.WriteUInt64LE(uint64(pid))

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreBadgeArcadeProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreMethodPostMetaBinary, rmcResponseBody)

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
