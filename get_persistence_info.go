package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func getPersistenceInfo(err error, client *nex.Client, callID uint32, ownerID uint32, persistenceSlotID uint16) {
	dataID := getDataStorePersistenceInfo(ownerID, persistenceSlotID)
	
	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreBadgeArcadeProtocolID, callID)

	if dataID != 0 {
		pPersistenceInfo := nexproto.NewDataStorePersistenceInfo()
		pPersistenceInfo.OwnerID = ownerID
		pPersistenceInfo.PersistenceSlotID = persistenceSlotID
		pPersistenceInfo.DataID = dataID

		rmcResponseStream := nex.NewStreamOut(nexServer)

		rmcResponseStream.WriteStructure(pPersistenceInfo)

		rmcResponseBody := rmcResponseStream.Bytes()

		rmcResponse.SetSuccess(nexproto.DataStoreMethodGetPersistenceInfo, rmcResponseBody)
	} else {
		rmcResponse.SetError(nex.Errors.DataStore.NotFound)
	}

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
