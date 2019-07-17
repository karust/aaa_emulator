package objects

import (
	"fmt"

	"../../common/packet"
)

// CharacterModel ... character model
type CharacterModel struct {
}

func (char *CharacterModel) Parse(reader *packet.Reader) {
	ext := reader.Bool()
	fmt.Println("EXT:", ext)
}
