package world

// Point ...
type Point struct {
	WorldID   uint
	ZoneID    uint
	X         float32
	Y         float32
	Z         float32
	otationX  int8
	RotationY int8
	RotationZ int8
	Relative  bool
}

// NewPoint ... Creates new Point object
func NewPoint(x, y, z float32) *Point {
	return &Point{X: x, Y: y, Z: z}
}

/*
        Point(float x, float y, float z, sbyte rotationX, sbyte rotationY, sbyte rotationZ)
        {
            X = x;
            Y = y;
            Z = z;
            RotationX = rotationX;
            RotationY = rotationY;
            RotationZ = rotationZ;
        }

        Point(uint zoneId, float x, float y, float z)
        {
            ZoneId = zoneId;
            X = x;
            Y = y;
            Z = z;
        }

        Point(uint worldId, uint zoneId, float x, float y, float z,
            sbyte rotationX, sbyte rotationY, sbyte rotationZ)
        {
            WorldId = worldId;
            ZoneId = zoneId;
            X = x;
            Y = y;
            Z = z;
            RotationX = rotationX;
            RotationY = rotationY;
            RotationZ = rotationZ;
        }

        Point Clone()
        {
            return new Point(WorldId, ZoneId, X, Y, Z, RotationX, RotationY, RotationZ);
        }
	}
*/
