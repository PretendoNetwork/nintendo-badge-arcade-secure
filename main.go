package main

import (
	"fmt"
	"os"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var nexServer *nex.Server
var secureServer *nexproto.SecureBadgeArcadeProtocol
var s3Client *s3.S3

func main() {
	key := os.Getenv("S3_KEY")
	secret := os.Getenv("S3_SECRET")

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint: aws.String("http://" + os.Getenv("DATASTORE_DATA_URL")),
		Region: aws.String("us-east-1"),
	}

	newSession, _ := session.NewSession(s3Config)
	s3Client = s3.New(newSession)

	nexServer = nex.NewServer()
	nexServer.SetPrudpVersion(1)
	nexServer.SetPRUDPProtocolMinorVersion(3)
	nexServer.SetNexVersion(30716)
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

	nexServer.On("Connect", connect)

	secureServer := nexproto.NewSecureBadgeArcadeProtocol(nexServer)
	dataStoreBadgeArcadeProtocolServer := nexproto.NewDataStoreBadgeArcadeProtocol(nexServer)
	shopBadgeArcadePrococolServer := nexproto.NewShopBadgeArcadeProtocol(nexServer)

	secureServer.Register(register)
	secureServer.GetMaintenanceStatus(getMaintenanceStatus)

	dataStoreBadgeArcadeProtocolServer.GetPersistenceInfo(getPersistenceInfo)
	dataStoreBadgeArcadeProtocolServer.PostMetaBinary(postMetaBinary)
	dataStoreBadgeArcadeProtocolServer.PreparePostObject(preparePostObject)
	dataStoreBadgeArcadeProtocolServer.CompletePostObject(completePostObject)
	dataStoreBadgeArcadeProtocolServer.PrepareGetObject(prepareGetObject)
	dataStoreBadgeArcadeProtocolServer.GetMetaByOwnerId(getMetaByOwnerId)
	dataStoreBadgeArcadeProtocolServer.ChangeMeta(changeMeta)
	dataStoreBadgeArcadeProtocolServer.PrepareUpdateObject(prepareUpdateObject)
	dataStoreBadgeArcadeProtocolServer.CompleteUpdateObject(completeUpdateObject)

	shopBadgeArcadePrococolServer.PostPlayLog(postPlayLog)

	nexServer.Listen(":59401")
}
