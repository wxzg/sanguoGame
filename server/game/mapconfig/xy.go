package mapconfig

var Mapwidth = 200
var MapHeight = 200

func ToPosition(x, y int) int{
	return x + y*MapHeight
}