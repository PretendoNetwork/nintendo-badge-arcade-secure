package nex_datastore

import (
	"fmt"
	"os"
	"time"

	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/datastore"

	"github.com/PretendoNetwork/nintendo-badge-arcade-secure/database"
	"github.com/PretendoNetwork/nintendo-badge-arcade-secure/globals"
)

func PreparePostObject(err error, client *nex.Client, callID uint32, param *datastore.DataStorePreparePostParam) {
	pid := client.PID()
	var slot uint16 = 0

	dataID := database.GetDataStorePersistenceInfo(pid, slot)

	bucket := os.Getenv("PN_NBA_CONFIG_S3_BUCKET")
	key := fmt.Sprintf("%s/%011d-%05d", os.Getenv("PN_NBA_CONFIG_S3_PATH"), dataID, 1)

	input := &globals.PostObjectInput{
		Bucket:    bucket,
		Key:       key,
		ExpiresIn: time.Minute * 15,
	}

	res, _ := globals.S3PresignClient.PresignPostObject(input)

	fieldKey := datastore.NewDataStoreKeyValue()
	fieldKey.Key = "key"
	fieldKey.Value = key

	fieldCredential := datastore.NewDataStoreKeyValue()
	fieldCredential.Key = "X-Amz-Credential"
	fieldCredential.Value = res.Credential

	fieldSecurityToken := datastore.NewDataStoreKeyValue()
	fieldSecurityToken.Key = "X-Amz-Security-Token"
	fieldSecurityToken.Value = ""

	fieldAlgorithm := datastore.NewDataStoreKeyValue()
	fieldAlgorithm.Key = "X-Amz-Algorithm"
	fieldAlgorithm.Value = "AWS4-HMAC-SHA256"

	fieldDate := datastore.NewDataStoreKeyValue()
	fieldDate.Key = "X-Amz-Date"
	fieldDate.Value = res.Date

	fieldPolicy := datastore.NewDataStoreKeyValue()
	fieldPolicy.Key = "policy"
	fieldPolicy.Value = res.Policy

	fieldSignature := datastore.NewDataStoreKeyValue()
	fieldSignature.Key = "X-Amz-Signature"
	fieldSignature.Value = res.Signature

	pReqPostInfo := datastore.NewDataStoreReqPostInfo()

	pReqPostInfo.DataID = uint64(dataID)
	pReqPostInfo.URL = res.URL
	pReqPostInfo.RequestHeaders = []*datastore.DataStoreKeyValue{}
	pReqPostInfo.FormFields = []*datastore.DataStoreKeyValue{
		fieldKey,
		fieldCredential,
		fieldSecurityToken,
		fieldAlgorithm,
		fieldDate,
		fieldPolicy,
		fieldSignature,
	}
	pReqPostInfo.RootCACert = []byte{}

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteStructure(pReqPostInfo)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore.MethodPreparePostObject, rmcResponseBody)

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
