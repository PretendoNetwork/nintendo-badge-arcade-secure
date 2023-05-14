package nex_shop_nintendo_badge_arcade

import (
	"github.com/PretendoNetwork/nintendo-badge-arcade-secure/globals"

	"github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go/shop/nintendo-badge-arcade"
)

func PostPlayLog(err error, client *nex.Client, callID uint32, param *nexproto.ShopPostPlayLogParam) {
	// TODO: Do something with the data
	rmcResponse := nex.NewRMCResponse(nexproto.ProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.MethodPostPlayLog, nil)
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
