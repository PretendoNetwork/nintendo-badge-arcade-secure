package main

// TODO: Verify currently known variables
// These variables can vary
type badgeArcadeMetaBinary struct {
	Unknown1                uint32 // 02 00 00 00
	SetupStatus             uint32 // 0, 1 or 2 depending on initialization stage
	Unknown2                uint32 // 00 00 00 00
	PracticePlaysPlayed     uint32 // Number of practice machines played? (not sure)
	NotWinnerPracticeBadges uint32 // Number of practice badges which don't have free play prizes
	PracticePlaysLeft       uint32 // Number of practice machines plays left? (not sure)
	Unknown3                uint64
	Unknown4                []uint64
	Unknown5                []byte // 32 bytes long hash
}

