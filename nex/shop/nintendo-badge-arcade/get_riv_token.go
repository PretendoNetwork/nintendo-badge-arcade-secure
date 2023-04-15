package nex_shop_nintendo_badge_arcade

import (
	"github.com/PretendoNetwork/nintendo-badge-arcade-secure/globals"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go/shop/nintendo-badge-arcade"
)

func GetRivToken(err error, client *nex.Client, callID uint32, itemCode string, referenceID []byte) {
	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	// Stubbed
	rmcResponseStream.WriteString("")

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.ProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.MethodGetRivToken, rmcResponseBody)
	rmcResponse.SetCustomID(nexproto.CustomProtocolID)

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
