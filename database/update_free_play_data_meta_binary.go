package database

import "log"

func UpdateFreePlayDataMetaBinary(dataID uint64, metaBinary []byte, updatedTime uint64) {
	var err error
	_, err = postgres.Exec(`UPDATE pretendo_badge_arcade.free_play_data SET meta_binary=$1, updated_time=$2 WHERE data_id=$3`, metaBinary, updatedTime, dataID)
	if err != nil {
		log.Fatal(err)
	}
}
