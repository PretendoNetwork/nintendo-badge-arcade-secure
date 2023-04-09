package database

import "github.com/PretendoNetwork/badge-arcade-secure/globals"

func UpdateUserPlayInfoVersion(dataID uint32, version uint32) {
	var err error
	_, err = postgres.Exec(`UPDATE pretendo_badge_arcade.user_play_info SET version=$1 WHERE data_id=$2`, version, dataID)
	if err != nil {
		globals.Logger.Error(err.Error())
	}
}
