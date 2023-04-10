package database

import "github.com/PretendoNetwork/nintendo-badge-arcade-secure/globals"

func PostUserPlayInfo(dataID uint32, ownerID uint32, slot uint16) {
	var err error
	_, err = postgres.Exec(`INSERT INTO pretendo_badge_arcade.user_play_info(
		data_id,
		pid,
		slot,
		version
	)
	VALUES (
		$1,
		$2,
		$3,
		0
	) ON CONFLICT DO NOTHING`, dataID, ownerID, slot)
	if err != nil {
		globals.Logger.Error(err.Error())
	}
}
