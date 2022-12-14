package main

// TODO: Verify currently known variables
// These variables can vary
type badgeArcadeMetaBinary struct {
	Unk1          uint32 // 02 00 00 00
	Unk2          uint32 // 02 00 00 00
	Unk3          uint32 // 00 00 00 00
	Unk4          uint32 // 00 00 00 00
	Unk5          uint32 // 00 00 00 00
	Unk6          uint32 // 01 00 00 00
	CollectionID  uint64
	PrizesIDs     []uint64
	UnkHash       []byte // 32 bytes long hash
}

