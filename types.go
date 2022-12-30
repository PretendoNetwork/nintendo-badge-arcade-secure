package main

// TODO: Verify currently known variables
// These variables can vary
type badgeArcadeMetaBinary struct {
	Unknown1                uint32 // 02 00 00 00
	SetupStatus             uint32 // 0, 1 or 2 depending on initialization stage
	Unknown2                uint32 // 00 00 00 00
	WinnerPracticeBadges    uint32 // Number of practice badges which have free plays prizes
	NotWinnerPracticeBadges uint32 // Number of practice badges which don't have free play prizes
	PracticeCraneAvailable  uint32 // 1 if the practice crane hasn't been played yet, 0 if it has
	Unknown3                uint64
	Unknown4                []uint64
	Unknown5                []byte // 32 bytes long hash
}

