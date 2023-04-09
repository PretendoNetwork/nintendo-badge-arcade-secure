package nex

import (
	"github.com/PretendoNetwork/badge-arcade-secure/globals"
	nex_datastore "github.com/PretendoNetwork/badge-arcade-secure/nex/datastore"
	nex_datastore_nintendo_badge_arcade "github.com/PretendoNetwork/badge-arcade-secure/nex/datastore/nintendo-badge-arcade"
	nex_secure_connection "github.com/PretendoNetwork/badge-arcade-secure/nex/secure-connection"
	nex_secure_connection_nintendo_badge_arcade "github.com/PretendoNetwork/badge-arcade-secure/nex/secure-connection/nintendo-badge-arcade"
	nex_shop_nintendo_badge_arcade "github.com/PretendoNetwork/badge-arcade-secure/nex/shop/nintendo-badge-arcade"
	datastore_nintendo_badge_arcade "github.com/PretendoNetwork/nex-protocols-go/datastore/nintendo-badge-arcade"
	secure_connection_nintendo_badge_arcade "github.com/PretendoNetwork/nex-protocols-go/secure-connection/nintendo-badge-arcade"
	shop_nintendo_badge_arcade "github.com/PretendoNetwork/nex-protocols-go/shop/nintendo-badge-arcade"
)

func registerNEXProtocols() {

	secureConnectionProtocol := secure_connection_nintendo_badge_arcade.NewSecureConnectionNintendoBadgeArcadeProtocol(globals.NEXServer)

	secureConnectionProtocol.Register(nex_secure_connection.Register)
	secureConnectionProtocol.GetMaintenanceStatus(nex_secure_connection_nintendo_badge_arcade.GetMaintenanceStatus)

	dataStoreNintendoBadgeArcadeProtocol := datastore_nintendo_badge_arcade.NewDataStoreNintendoBadgeArcadeProtocol(globals.NEXServer)

	dataStoreNintendoBadgeArcadeProtocol.GetPersistenceInfo(nex_datastore.GetPersistenceInfo)
	dataStoreNintendoBadgeArcadeProtocol.PostMetaBinary(nex_datastore.PostMetaBinary)
	dataStoreNintendoBadgeArcadeProtocol.PreparePostObject(nex_datastore.PreparePostObject)
	dataStoreNintendoBadgeArcadeProtocol.CompletePostObject(nex_datastore.CompletePostObject)
	dataStoreNintendoBadgeArcadeProtocol.PrepareGetObject(nex_datastore.PrepareGetObject)
	dataStoreNintendoBadgeArcadeProtocol.GetMetaByOwnerID(nex_datastore_nintendo_badge_arcade.GetMetaByOwnerID)
	dataStoreNintendoBadgeArcadeProtocol.ChangeMeta(nex_datastore.ChangeMeta)
	dataStoreNintendoBadgeArcadeProtocol.PrepareUpdateObject(nex_datastore.PrepareUpdateObject)
	dataStoreNintendoBadgeArcadeProtocol.CompleteUpdateObject(nex_datastore.CompleteUpdateObject)

	shopNintendoBadgeArcadePrococol := shop_nintendo_badge_arcade.NewShopNintendoBadgeArcadeProtocol(globals.NEXServer)

	shopNintendoBadgeArcadePrococol.PostPlayLog(nex_shop_nintendo_badge_arcade.PostPlayLog)
}
