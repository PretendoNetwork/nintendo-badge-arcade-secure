package main

import (
	"fmt"
	"os"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
	nexsecure "github.com/PretendoNetwork/nex-protocols-common-go/secure-connection"
)

var nexServer *nex.Server
var secureServer *nexproto.SecureProtocol

func main() {
	nexServer = nex.NewServer()
	nexServer.SetPrudpVersion(1)
	nexServer.SetNexVersion(30500)
	nexServer.SetKerberosKeySize(32)
	nexServer.SetKerberosPassword(os.Getenv("KERBEROS_PASSWORD"))
	nexServer.SetAccessKey("82d5962d")

	nexServer.On("Data", func(packet *nex.PacketV1) {
		request := packet.RMCRequest()

		fmt.Println("==Badge Arcade - Secure==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID())
		fmt.Printf("Method ID: %#v\n", request.MethodID())
		fmt.Println("====================")
	})

	nexServer.On("Kick", func(packet *nex.PacketV1) {
		fmt.Println("Leaving")
		deletePlayerSession(packet.Sender().PID())
	})

	secureServer := nexsecure.NewCommonSecureConnectionProtocol(nexServer)
	dataStoreBadgeArcadePrococolServer := nexproto.NewDataStoreBadgeArcadeProtocol(nexServer)
	shopBadgeArcadePrococolServer := nexproto.NewShopBadgeArcadeProtocol(nexServer)
	_ = shopBadgeArcadePrococolServer

	secureServer.AddConnection(func(rvcid uint32, urls []string, ip, port string) {
		pid := nexServer.FindClientFromConnectionID(rvcid).PID()
		addPlayerSession(pid, urls, ip, port)
	})
	secureServer.UpdateConnection(func(rvcid uint32, urls []string, ip, port string) {
		pid := nexServer.FindClientFromConnectionID(rvcid).PID()
		updatePlayerSessionAll(pid, urls, ip, port)
	})
	secureServer.DoesConnectionExist(func(rvcid uint32) bool {
		pid := nexServer.FindClientFromConnectionID(rvcid).PID()
		return doesSessionExist(pid)
	})

	secureServer.GetMaintenanceStatus(getMaintenanceStatus)

	dataStoreBadgeArcadePrococolServer.GetPersistenceInfo(getPersistenceInfo)
	dataStoreBadgeArcadePrococolServer.PrepareGetObject(prepareGetObject)
	dataStoreBadgeArcadePrococolServer.GetMetaByOwnerId(getMetaByOwnerId)
	dataStoreBadgeArcadePrococolServer.ChangeMeta(changeMeta)

	shopBadgeArcadePrococolServer.PostPlayLog(postPlayLog)

	nexServer.Listen(":59401")
}