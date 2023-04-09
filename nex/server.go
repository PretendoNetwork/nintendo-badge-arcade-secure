package nex

import (
	"fmt"
	"os"

	"github.com/PretendoNetwork/badge-arcade-secure/database"
	"github.com/PretendoNetwork/badge-arcade-secure/globals"
	"github.com/PretendoNetwork/badge-arcade-secure/prudp"

	nex "github.com/PretendoNetwork/nex-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func StartNEXServer() {
	key := os.Getenv("S3_KEY")
	secret := os.Getenv("S3_SECRET")

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String("http://" + os.Getenv("DATASTORE_DATA_URL")),
		Region:      aws.String("us-east-1"),
	}

	newSession, _ := session.NewSession(s3Config)
	globals.S3Client = s3.New(newSession)

	globals.NEXServer = nex.NewServer()
	globals.NEXServer.SetPRUDPVersion(1)
	globals.NEXServer.SetPRUDPProtocolMinorVersion(3)
	globals.NEXServer.SetDefaultNEXVersion(&nex.NEXVersion{
		Major: 3,
		Minor: 7,
		Patch: 16,
	})
	globals.NEXServer.SetKerberosPassword(os.Getenv("KERBEROS_PASSWORD"))
	globals.NEXServer.SetAccessKey("82d5962d")

	globals.NEXServer.On("Data", func(packet *nex.PacketV1) {
		request := packet.RMCRequest()

		fmt.Println("==Badge Arcade - Secure==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID())
		fmt.Printf("Method ID: %#v\n", request.MethodID())
		fmt.Println("====================")
	})

	globals.NEXServer.On("Kick", func(packet *nex.PacketV1) {
		fmt.Println("Leaving")
		database.DeletePlayerSession(packet.Sender().PID())
	})

	globals.NEXServer.On("Connect", prudp.Connect)

	registerNEXProtocols()

	globals.NEXServer.Listen(":59401")
}
