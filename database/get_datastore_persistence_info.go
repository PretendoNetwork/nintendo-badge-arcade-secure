package database

import (
	"database/sql"
	"github.com/PretendoNetwork/badge-arcade-secure/globals"
)

func GetDataStorePersistenceInfo(ownerID uint32, persistenceSlotID uint16) uint32 {
	var dataID uint32
	err := postgres.QueryRow(`SELECT data_id FROM pretendo_badge_arcade.user_play_info WHERE pid=$1 AND slot=$2`, ownerID, persistenceSlotID).Scan(&dataID)
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Error(err.Error())
	}

	return dataID
}
