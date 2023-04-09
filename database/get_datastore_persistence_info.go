package database

import (
	"database/sql"
	"log"
)

func GetDataStorePersistenceInfo(ownerID uint32, persistenceSlotID uint16) uint32 {
	var dataID uint32
	err := postgres.QueryRow(`SELECT data_id FROM pretendo_badge_arcade.user_play_info WHERE pid=$1 AND slot=$2`, ownerID, persistenceSlotID).Scan(&dataID)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}

	return dataID
}
