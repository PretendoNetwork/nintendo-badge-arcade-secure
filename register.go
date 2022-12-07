package main

import (
	"strconv"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func register(err error, client *nex.Client, callID uint32, stationUrls []*nex.StationURL) {
	localStation := stationUrls[0]
	localStationURL := localStation.EncodeToString()
	pidConnectionID := uint32(nexServer.ConnectionIDCounter().Increment())
	client.SetConnectionID(pidConnectionID)
	client.SetLocalStationUrl(localStationURL)

	address := client.Address().IP.String()
	port := strconv.Itoa(client.Address().Port)
	natf := "0"
	natm := "0"
	type_ := "3"

	localStation.SetAddress(address)
	localStation.SetPort(port)
	localStation.SetNatf(natf)
	localStation.SetNatm(natm)
	localStation.SetType(type_)

	urlPublic := localStation.EncodeToString()

	pid := nexServer.FindClientFromConnectionID(pidConnectionID).PID()

	if !doesSessionExist(pid) {
		addPlayerSession(pid, []string{localStationURL, urlPublic}, address, port)
	} else {
		updatePlayerSessionAll(pid, []string{localStationURL, urlPublic}, address, port)
	}

	retval := nex.NewResultSuccess(nex.Errors.Core.Unknown)

	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteResult(retval) // Success
	rmcResponseStream.WriteUInt32LE(pidConnectionID)
	rmcResponseStream.WriteString(urlPublic)

	rmcResponseBody := rmcResponseStream.Bytes()

	// Build response packet
	rmcResponse := nex.NewRMCResponse(nexproto.SecureBadgeArcadeProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.SecureMethodRegister, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	var responsePacket nex.PacketInterface

	responsePacket, _ = nex.NewPacketV1(client, nil)
	responsePacket.SetVersion(1)

	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	nexServer.Send(responsePacket)
}
