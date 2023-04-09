package database

import "log"

func UpdateUserPlayInfoVersion(dataID uint32, version uint32) {
	var err error
	_, err = postgres.Exec(`UPDATE pretendo_badge_arcade.user_play_info SET version=$1 WHERE data_id=$2`, version, dataID)
	if err != nil {
		log.Fatal(err)
	}
}
