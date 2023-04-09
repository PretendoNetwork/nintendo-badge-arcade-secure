package database

import "log"

func PostFreePlayDataMetaInfo(dataID uint64, ownerID uint32, metaBinary []byte, createdTime uint64, period uint16, flag uint32) {
	var err error
	_, err = postgres.Exec(`INSERT INTO pretendo_badge_arcade.free_play_data(
		data_id,
		owner_id,
		meta_binary,
		created_time,
		updated_time,
		period,
		flag,
		referred_time
	)
	VALUES (
		$1,
		$2,
		$3,
		$4,
		$4,
		$5,
		$6,
		$4
	) ON CONFLICT DO NOTHING`, dataID, ownerID, metaBinary, createdTime, period, flag)
	if err != nil {
		log.Fatal(err)
	}
}
