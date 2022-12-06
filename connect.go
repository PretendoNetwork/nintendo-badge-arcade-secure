package main

import (
	"github.com/PretendoNetwork/nex-go"
)

func connect(packet nex.PacketInterface) {
	payload := packet.Payload()

	stream := nex.NewStreamIn(payload, nexServer)

	// TODO: Error check!!
	ticketData, _ := stream.ReadBuffer()
	requestData, _ := stream.ReadBuffer()

	serverKey := nex.DeriveKerberosKey(2, []byte(nexServer.KerberosPassword()))

	ticket := nex.NewKerberosTicketInternalData()
	ticket.Decrypt(nex.NewStreamIn(ticketData, nexServer), serverKey)

	// TODO: Check timestamp here

	sessionKey := ticket.SessionKey()
	kerberos := nex.NewKerberosEncryption(sessionKey)

	decryptedRequestData := kerberos.Decrypt(requestData)
	checkDataStream := nex.NewStreamIn(decryptedRequestData, nexServer)

	userPID := checkDataStream.ReadUInt32LE()
	_ = checkDataStream.ReadUInt32LE() //CID of secure server station url
	responseCheck := checkDataStream.ReadUInt32LE()

	responseValueStream := nex.NewStreamOut(nexServer)
	responseValueStream.WriteUInt32LE(responseCheck + 1)

	responseValueBufferStream := nex.NewStreamOut(nexServer)
	responseValueBufferStream.WriteBuffer(responseValueStream.Bytes())

	nexServer.AcknowledgePacket(packet, responseValueBufferStream.Bytes())

	packet.Sender().UpdateRC4Key(sessionKey)
	packet.Sender().SetSessionKey(sessionKey)

	packet.Sender().SetPID(userPID)
}
