package database

import "log"

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
		log.Fatal(err)
	}
}
