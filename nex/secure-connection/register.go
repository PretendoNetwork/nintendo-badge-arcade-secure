package nex_secure_connection

import (
	"strconv"

	"github.com/PretendoNetwork/badge-arcade-secure/database"
	"github.com/PretendoNetwork/badge-arcade-secure/globals"
	secure_connection "github.com/PretendoNetwork/nex-protocols-go/secure-connection"

	nex "github.com/PretendoNetwork/nex-go"
)

func Register(err error, client *nex.Client, callID uint32, stationUrls []*nex.StationURL) {
	localStation := stationUrls[0]
	localStationURL := localStation.EncodeToString()
	pidConnectionID := uint32(globals.NEXServer.ConnectionIDCounter().Increment())
	client.SetConnectionID(pidConnectionID)
	client.SetLocalStationURL(localStationURL)

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

	pid := globals.NEXServer.FindClientFromConnectionID(pidConnectionID).PID()

	if !database.DoesSessionExist(pid) {
		database.AddPlayerSession(pid, []string{localStationURL, urlPublic}, address, port)
	} else {
		database.UpdatePlayerSessionAll(pid, []string{localStationURL, urlPublic}, address, port)
	}

	retval := nex.NewResultSuccess(nex.Errors.Core.Unknown)

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteResult(retval) // Success
	rmcResponseStream.WriteUInt32LE(pidConnectionID)
	rmcResponseStream.WriteString(urlPublic)

	rmcResponseBody := rmcResponseStream.Bytes()

	// Build response packet
	rmcResponse := nex.NewRMCResponse(secure_connection.ProtocolID, callID)
	rmcResponse.SetSuccess(secure_connection.MethodRegister, rmcResponseBody)

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

	globals.NEXServer.Send(responsePacket)
}
