package models

type PlantIcon string

const (
	BananaPlant  PlantIcon = "bananaPlant"
	BigCactus    PlantIcon = "bigCactus"
	BigPlant     PlantIcon = "bigPlant"
	BigRose      PlantIcon = "bigRose"
	ChilliPlant  PlantIcon = "chilliPlant"
	Daisy        PlantIcon = "daisy"
	FlowerBed    PlantIcon = "flowerBed"
	Flower       PlantIcon = "flower"
	LeafyPlant   PlantIcon = "leafyPlant"
	MediumPlant  PlantIcon = "mediumPlant"
	RedTulip     PlantIcon = "redTulip"
	SeaweedPlant PlantIcon = "seaweedPlant"
	ShortPlant   PlantIcon = "shortPlant"
	SkinnyPlant  PlantIcon = "skinnyPlant"
	SmallCactus  PlantIcon = "smallCactus"
	SmallPlant   PlantIcon = "smallPlant"
	SmallRose    PlantIcon = "smallRose"
	SpikyPlant   PlantIcon = "spikyPlant"
	TallPlant    PlantIcon = "tallPlant"
	ThreeFlowers PlantIcon = "threeFlowers"
	TwoFlowers   PlantIcon = "twoFlowers"
	TwoPlants    PlantIcon = "twoPlants"
	WhiteFlower  PlantIcon = "whiteFlower"
	YellowTulip  PlantIcon = "yellowTulip"
)

var validPlantIcons = []PlantIcon{
	BananaPlant, BigCactus, BigPlant, BigRose, ChilliPlant, Daisy,
	FlowerBed, Flower, LeafyPlant, MediumPlant, RedTulip, SeaweedPlant,
	ShortPlant, SkinnyPlant, SmallCactus, SmallPlant, SmallRose, SpikyPlant,
	TallPlant, ThreeFlowers, TwoFlowers, TwoPlants, WhiteFlower, YellowTulip,
}

func (pi PlantIcon) IsValid() bool {
	for _, valid := range validPlantIcons {
		if pi == valid {
			return true
		}
	}
	return false
}

type Plant struct {
	ID        int64 `gorm:"primaryKey"`
	Name      string
	Note      string
	TagColor  string
	UserID    int64
	User      User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Reminders []Reminder `gorm:"foreignKey:PlantID;constraint:OnDelete:CASCADE"`
	PlantIcon PlantIcon
}
