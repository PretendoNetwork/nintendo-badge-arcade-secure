package prudp

import (
	"github.com/PretendoNetwork/badge-arcade-secure/globals"

	"github.com/PretendoNetwork/nex-go"
)

func Connect(packet nex.PacketInterface) {
	payload := packet.Payload()

	stream := nex.NewStreamIn(payload, globals.NEXServer)

	// TODO: Error check!!
	ticketData, _ := stream.ReadBuffer()
	requestData, _ := stream.ReadBuffer()

	serverKey := nex.DeriveKerberosKey(2, []byte(globals.NEXServer.KerberosPassword()))

	ticket := nex.NewKerberosTicketInternalData()
	ticket.Decrypt(nex.NewStreamIn(ticketData, globals.NEXServer), serverKey)

	// TODO: Check timestamp here

	sessionKey := ticket.SessionKey()
	kerberos := nex.NewKerberosEncryption(sessionKey)

	decryptedRequestData := kerberos.Decrypt(requestData)
	checkDataStream := nex.NewStreamIn(decryptedRequestData, globals.NEXServer)

	userPID := checkDataStream.ReadUInt32LE()
	_ = checkDataStream.ReadUInt32LE() //CID of secure server station url
	responseCheck := checkDataStream.ReadUInt32LE()

	responseValueStream := nex.NewStreamOut(globals.NEXServer)
	responseValueStream.WriteUInt32LE(responseCheck + 1)

	responseValueBufferStream := nex.NewStreamOut(globals.NEXServer)
	responseValueBufferStream.WriteBuffer(responseValueStream.Bytes())

	globals.NEXServer.AcknowledgePacket(packet, responseValueBufferStream.Bytes())

	packet.Sender().UpdateRC4Key(sessionKey)
	packet.Sender().SetSessionKey(sessionKey)

	packet.Sender().SetPID(userPID)
}
