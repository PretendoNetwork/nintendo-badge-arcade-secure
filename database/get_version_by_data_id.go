package database

import "github.com/PretendoNetwork/badge-arcade-secure/globals"

func GetVersionByDataID(dataID uint32) uint32 {
	var version uint32
	err := postgres.QueryRow(`SELECT version FROM pretendo_badge_arcade.user_play_info WHERE data_id=$1`, dataID).Scan(&version)
	if err != nil {
		globals.Logger.Error(err.Error())
	}

	return version
}
