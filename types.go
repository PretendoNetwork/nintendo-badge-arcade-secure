package main

// TODO: Verify currently known variables
type badgeArcadeMetaBinary struct {
	Unk1          uint32 // 02 00 00 00
	Unk2          uint32 // 02 00 00 00
	Unk3          uint32 // 00 00 00 00
	Unk4          uint32 // 00 00 00 00
	Unk5          uint32 // 00 00 00 00
	CollectionIDs []uint64
	PrizesIDs     []uint64
	UnkHash       []byte // 32 bytes long hash
}

