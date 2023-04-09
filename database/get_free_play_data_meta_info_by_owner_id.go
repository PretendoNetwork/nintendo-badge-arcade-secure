package database

import (
	"database/sql"
	"log"
)

func GetFreePlayDataMetaInfoByOwnerID(ownerID uint32) (uint64, []byte, uint64, uint64, uint16, uint32, uint64) {
	var dataID uint64
	var metaBinary []byte
	var createdTime uint64
	var updatedTime uint64
	var period uint16
	var flag uint32
	var referredTime uint64
	err := postgres.QueryRow(`SELECT data_id, meta_binary, created_time, updated_time, period, flag, referred_time FROM pretendo_badge_arcade.free_play_data WHERE owner_id=$1`, ownerID).Scan(&dataID, &metaBinary, &createdTime, &updatedTime, &period, &flag, &referredTime)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}

	return dataID, metaBinary, createdTime, updatedTime, period, flag, referredTime
}
