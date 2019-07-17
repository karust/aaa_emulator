package charmodel

import (
	"fmt"

	"../../../common/packet"
)

// FixedDecalAsset ...
type FixedDecalAsset struct {
	AssetID     uint32
	AssetWeight float32
}

// Parse ... parse decal asset from packet
func (char *FixedDecalAsset) Parse(reader *packet.Reader) {

}

// FaceModel ...
type FaceModel struct {
	MovableDecalAssetID uint32
	MovableDecalWeight  float32
	MovableDecalScale   float32
	MovableDecalRotate  float32
	MovableDecalMoveX   int16
	MovableDecalMoveY   int16
	FixedDecalAsset     []FixedDecalAsset
	DiffuseMapID        uint32
	NormalMapID         uint32
	EyelashMapID        uint32
	NormalMapWeight     float32
	LipColor            uint32
	LeftPupilColor      uint32
	RightPupilColor     uint32
	EyebrowColor        uint32
	DecoColor           uint32
	Modifier            []byte
}

// Parse ... parse face model from packet
func (face *FaceModel) Parse(reader *packet.Reader) {
	face.MovableDecalAssetID = reader.UInt()       // type
	face.MovableDecalWeight = reader.Float()       // weight
	face.MovableDecalScale = reader.Float()        // scale
	face.MovableDecalRotate = reader.Float()       // rotate
	face.MovableDecalMoveX = int16(reader.Short()) // moveX
	face.MovableDecalMoveY = int16(reader.Short()) // moveY

	mAssets := reader.Pisc(4)
	fmt.Println("mAssets", mAssets)
	face.FixedDecalAsset = make([]FixedDecalAsset, 6)
	face.FixedDecalAsset[0].AssetID = uint32(mAssets[0])
	face.FixedDecalAsset[1].AssetID = uint32(mAssets[1])
	face.FixedDecalAsset[2].AssetID = uint32(mAssets[2])
	face.FixedDecalAsset[3].AssetID = uint32(mAssets[3])

	mAssets = reader.Pisc(2)
	face.FixedDecalAsset[4].AssetID = uint32(mAssets[0])
	face.FixedDecalAsset[5].AssetID = uint32(mAssets[1])

	// for 3.0.3.0
	mMap := reader.Pisc(3)
	for i := 0; i < 6; i++ {
		face.FixedDecalAsset[i].AssetWeight = reader.Float() // weight
	}

	// почему-то нет такого в 3+
	//DiffuseMapId = stream.ReadUInt32();
	//NormalMapId = stream.ReadUInt32();
	//EyelashMapId = stream.ReadUInt32();
	face.DiffuseMapID = uint32(mMap[0])
	face.NormalMapID = uint32(mMap[1])
	face.EyelashMapID = uint32(mMap[2])

	face.NormalMapWeight = reader.Float()
	face.LipColor = reader.UInt()        // lip
	face.LeftPupilColor = reader.UInt()  // leftPupil
	face.RightPupilColor = reader.UInt() // rightPupil
	face.EyebrowColor = reader.UInt()    // eyebrow
	face.DecoColor = reader.UInt()       // deco
	face.Modifier = reader.Bytes()
}

// CharacterModel ... character model
type CharacterModel struct {
	HairColorID uint32
	SkinColorID uint32
	UnkID       float32
	Face        FaceModel
}

// New ... returns object of CharacterModel
func New() *CharacterModel {
	return &CharacterModel{}
}

// Parse ... parse character model from packet
func (char *CharacterModel) Parse(reader *packet.Reader) {
	ext := reader.Byte()
	// None
	if ext == 0 {
		return
	}
	char.HairColorID = reader.UInt() // HairColorId
	reader.UInt()                    // type for 3.0.3.0
	reader.UInt()                    // defaultHairColor for 3.0.3.0
	reader.UInt()                    // twoToneHair for 3.0.3.0
	reader.Float()                   // twoToneFirstWidth for 3.0.3.0
	reader.Float()                   // twoToneSecondWidth for 3.0.3.0

	// Hair
	if ext == 1 {
		return
	}

	char.SkinColorID = reader.UInt()
	reader.UInt()               // type for 3.0.3.0
	reader.UInt()               // type for 3.0.3.0
	char.UnkID = reader.Float() // weight

	// Skin
	if ext == 2 {
		return
	}

	// Face
	char.Face.Parse(reader)
}
