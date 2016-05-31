package model

import (
	"errors"
	"strconv"
)

const NO_STONE = -1
const NUM_GEMSTONES = 7
const GEMSTONE_SHORTNAME_LEN = 2
const GEMSTONE_CLARITY_SHORTNAME_LEN = 2
const (
	DIAMOND = iota
	BLACK_DIAMOND
	EMERALD
	RUBY
	SAPPHIRE_PINK
	SAPPHIRE_YELLOW
	SAPPHIRE_BLUE
)

const (
	CENTER_STONE = iota
	ACCENT_STONE_ONE
	ACCENT_STONE_TWO
)

const (
	GEMSTONE_CAT_ZERO = iota //Round stones
	GEMSTONE_CAT_ONE         //Fancy Shaped stones
	GEMSTONE_CAT_TWO         //Non faceted stones
)
const (
	GEM_ROUND = iota
	GEM_SQUARE
	GEM_MARQUISE
	GEM_CUSHION
	GEM_TRILLION
	GEM_BAGUTTE
)
const (
	GEM_CLARITY_ZERO = iota
	GEM_CLARITY_ONE
	GEM_CLARITY_TWO
)
const (
	GEM_SET_PRONG = iota
	GEM_SET_BEZEL
	GEM_SET_PAVE
	GEM_SET_CHANNEL
)

type Gemstone struct {
	Id         int
	Name       string
	ShortName  string
	Birthstone string
}

type GemstoneClarity struct {
	Id        int
	Name      string
	ShortName string
}

type GemstoneShape struct {
	Id   int
	Name string
}

type GemstoneSetting struct {
	Id   int
	Name string
}

var GemstoneList []Gemstone = []Gemstone{
	DIAMOND:         Gemstone{Id: DIAMOND, Name: "Diamond", ShortName: "DW", Birthstone: "April"},
	BLACK_DIAMOND:   Gemstone{Id: BLACK_DIAMOND, Name: "Black Diamond", ShortName: "DB", Birthstone: "April"},
	EMERALD:         Gemstone{Id: EMERALD, Name: "Emerald", ShortName: "EM", Birthstone: "May"},
	RUBY:            Gemstone{Id: RUBY, Name: "Ruby", ShortName: "RY", Birthstone: "July"},
	SAPPHIRE_PINK:   Gemstone{Id: SAPPHIRE_PINK, Name: "Pink Sapphire", ShortName: "SP", Birthstone: "September"},
	SAPPHIRE_YELLOW: Gemstone{Id: SAPPHIRE_YELLOW, Name: "Yellow Sapphire", ShortName: "SY", Birthstone: "September"},
	SAPPHIRE_BLUE:   Gemstone{Id: SAPPHIRE_BLUE, Name: "Blue Sapphire", ShortName: "SB", Birthstone: "September"},
}
var GemstoneReverseMap map[string]int = map[string]int{
	"DW": DIAMOND,
	"DB": BLACK_DIAMOND,
	"EM": EMERALD,
	"RY": RUBY,
	"SP": SAPPHIRE_PINK,
	"SY": SAPPHIRE_YELLOW,
	"SB": SAPPHIRE_BLUE,
}

var GemstoneClarityList []GemstoneClarity = []GemstoneClarity{
	GEM_CLARITY_ZERO: GemstoneClarity{Id: GEM_CLARITY_ZERO, Name: "SI IJ", ShortName: "C0"},
	GEM_CLARITY_ONE:  GemstoneClarity{Id: GEM_CLARITY_ONE, Name: "SI GH", ShortName: "C1"},
	GEM_CLARITY_TWO:  GemstoneClarity{Id: GEM_CLARITY_TWO, Name: "VVS EF", ShortName: "C2"},
}
var GemstoneClarityReverseMap map[string]int = map[string]int{
	"C0": GEM_CLARITY_ZERO,
	"C1": GEM_CLARITY_ONE,
	"C2": GEM_CLARITY_TWO,
}

var GemstoneShapeList []GemstoneShape = []GemstoneShape{
	GEM_ROUND:    GemstoneShape{Id: GEM_ROUND, Name: "Round"},
	GEM_SQUARE:   GemstoneShape{Id: GEM_SQUARE, Name: "Square"},
	GEM_MARQUISE: GemstoneShape{Id: GEM_MARQUISE, Name: "Marquise"},
	GEM_CUSHION:  GemstoneShape{Id: GEM_CUSHION, Name: "Cushion"},
	GEM_TRILLION: GemstoneShape{Id: GEM_TRILLION, Name: "Trillion"},
	GEM_BAGUTTE:  GemstoneShape{Id: GEM_BAGUTTE, Name: "Bagutte"},
}

var GemstoneSettingList []GemstoneSetting = []GemstoneSetting{
	GEM_SET_PRONG:   GemstoneSetting{Id: GEM_SET_PRONG, Name: "Prong"},
	GEM_SET_BEZEL:   GemstoneSetting{Id: GEM_SET_BEZEL, Name: "Bezel"},
	GEM_SET_PAVE:    GemstoneSetting{Id: GEM_SET_PAVE, Name: "Pave"},
	GEM_SET_CHANNEL: GemstoneSetting{Id: GEM_SET_CHANNEL, Name: "Channel"},
}

func GetGemstoneList() []Gemstone {
	return GemstoneList
}

func GetGemstone(gemstoneId int) Gemstone {
	return GemstoneList[gemstoneId]
}

func GetGemstonePrice(gemstoneId int, weight float64, category int, clarity int) (float64, error) {
	var price float64
	var priceMultiplier float64
	if weight <= 0 {
		return price, errors.New("Invalid weight value passed for gemstone: Weight - " + strconv.FormatFloat(weight, 'f', -1, 64))
	}

	if gemstoneId == DIAMOND || gemstoneId == BLACK_DIAMOND {
	} else if gemstoneId == RUBY || gemstoneId == SAPPHIRE_PINK ||
		gemstoneId == SAPPHIRE_YELLOW || gemstoneId == SAPPHIRE_BLUE {
		weight *= 1.14
	} else if gemstoneId == EMERALD {
		weight *= 0.79
	}
	if gemstoneId == DIAMOND {
		if category == GEMSTONE_CAT_ZERO {
			if weight < 0.01 {
				if clarity == GEM_CLARITY_ZERO {
					priceMultiplier = 42000
				} else if clarity == GEM_CLARITY_ONE {
					priceMultiplier = 55000
				} else if clarity == GEM_CLARITY_TWO {
					priceMultiplier = 98000
				}
			} else if weight <= 0.02 {
				if clarity == GEM_CLARITY_ZERO {
					priceMultiplier = 40000
				} else if clarity == GEM_CLARITY_ONE {
					priceMultiplier = 50000
				} else if clarity == GEM_CLARITY_TWO {
					priceMultiplier = 95000
				}
			} else if weight <= 0.05 {
				if clarity == GEM_CLARITY_ZERO {
					priceMultiplier = 42000
				} else if clarity == GEM_CLARITY_ONE {
					priceMultiplier = 55000
				} else if clarity == GEM_CLARITY_TWO {
					priceMultiplier = 98000
				}
			} else if weight <= 0.11 {
				if clarity == GEM_CLARITY_ZERO {
					priceMultiplier = 55000
				} else if clarity == GEM_CLARITY_ONE {
					priceMultiplier = 65000
				} else if clarity == GEM_CLARITY_TWO {
					priceMultiplier = 110000
				}
			} else if weight <= 0.19 {
				if clarity == GEM_CLARITY_ZERO {
					priceMultiplier = 65000
				} else if clarity == GEM_CLARITY_ONE {
					priceMultiplier = 80000
				} else if clarity == GEM_CLARITY_TWO {
					priceMultiplier = 120000
				}
			} else if weight <= 0.25 {
				if clarity == GEM_CLARITY_ZERO {
					priceMultiplier = 80000
				} else if clarity == GEM_CLARITY_ONE {
					priceMultiplier = 100000
				} else if clarity == GEM_CLARITY_TWO {
					priceMultiplier = 135000
				}
			}
		} else if category == GEMSTONE_CAT_ONE {
			if weight <= 0.05 {
				priceMultiplier = 50000
			} else if weight <= 0.10 {
				priceMultiplier = 60000
			} else if weight <= 0.18 {
				priceMultiplier = 80000
			} else if weight <= 0.25 {
				priceMultiplier = 95000
			}
		}
	} else if gemstoneId == RUBY {
		if category != GEMSTONE_CAT_TWO {
			if weight <= 0.05 {
				priceMultiplier = 10000
			} else if weight <= 0.14 {
				priceMultiplier = 12000
			} else if weight < 0.25 {
				priceMultiplier = 18000
			} else if weight < 0.50 {
				priceMultiplier = 22000
			} else if weight <= 1.00 {
				priceMultiplier = 30000
			}
		} else {
			if weight <= 0.05 {
				priceMultiplier = 8500
			} else if weight <= 0.14 {
				priceMultiplier = 10200
			} else if weight < 0.25 {
				priceMultiplier = 15300
			} else if weight < 0.50 {
				priceMultiplier = 18700
			} else if weight <= 1.00 {
				priceMultiplier = 25500
			}
		}
	} else if gemstoneId == SAPPHIRE_BLUE {
		if category != GEMSTONE_CAT_TWO {
			if weight <= 0.05 {
				priceMultiplier = 10000
			} else if weight <= 0.14 {
				priceMultiplier = 12000
			} else if weight < 0.25 {
				priceMultiplier = 18000
			} else if weight < 0.50 {
				priceMultiplier = 25000
			} else if weight <= 1.00 {
				priceMultiplier = 33000
			}
		} else {
			if weight <= 0.05 {
				priceMultiplier = 7650
			} else if weight <= 0.14 {
				priceMultiplier = 10200
			} else if weight < 0.25 {
				priceMultiplier = 15300
			} else if weight < 0.50 {
				priceMultiplier = 21250
			} else if weight <= 1.00 {
				priceMultiplier = 28050
			}
		}
	} else if gemstoneId == EMERALD {
		if category != GEMSTONE_CAT_TWO {
			if weight <= 0.05 {
				priceMultiplier = 9000
			} else if weight <= 0.14 {
				priceMultiplier = 12000
			} else if weight < 0.25 {
				priceMultiplier = 18000
			} else if weight < 0.50 {
				priceMultiplier = 22000
			} else if weight <= 1.00 {
				priceMultiplier = 30000
			}
		} else {
			if weight <= 0.05 {
				priceMultiplier = 7650
			} else if weight <= 0.14 {
				priceMultiplier = 10200
			} else if weight < 0.25 {
				priceMultiplier = 15300
			} else if weight < 0.50 {
				priceMultiplier = 18700
			} else if weight <= 1.00 {
				priceMultiplier = 25500
			}
		}
	} else if gemstoneId == BLACK_DIAMOND {
		if weight <= 0.05 {
			priceMultiplier = 24000
		} else if weight <= 0.20 {
			priceMultiplier = 28000
		} else if weight <= 1.00 {
			priceMultiplier = 30000
		}
	} else if gemstoneId == SAPPHIRE_YELLOW {
		if category != GEMSTONE_CAT_TWO {
			if weight <= 0.05 {
				priceMultiplier = 9000
			} else if weight <= 0.14 {
				priceMultiplier = 12000
			} else if weight < 0.25 {
				priceMultiplier = 18000
			} else if weight < 0.50 {
				priceMultiplier = 22000
			} else if weight <= 1.00 {
				priceMultiplier = 30000
			}
		} else {
			if weight <= 0.05 {
				priceMultiplier = 7650
			} else if weight <= 0.14 {
				priceMultiplier = 10200
			} else if weight < 0.25 {
				priceMultiplier = 15300
			} else if weight < 0.50 {
				priceMultiplier = 18700
			} else if weight <= 1.00 {
				priceMultiplier = 25500
			}
		}
	} else if gemstoneId == SAPPHIRE_PINK {
		if category != GEMSTONE_CAT_TWO {
			if weight <= 0.05 {
				priceMultiplier = 8000
			} else if weight <= 0.14 {
				priceMultiplier = 11000
			} else if weight < 0.25 {
				priceMultiplier = 15000
			} else if weight < 0.50 {
				priceMultiplier = 20000
			} else if weight <= 1.00 {
				priceMultiplier = 28000
			}
		} else {
			if weight <= 0.05 {
				priceMultiplier = 6800
			} else if weight <= 0.14 {
				priceMultiplier = 9350
			} else if weight < 0.25 {
				priceMultiplier = 12750
			} else if weight < 0.50 {
				priceMultiplier = 17000
			} else if weight <= 1.00 {
				priceMultiplier = 23800
			}
		}
	} else {
		return price, errors.New("Invalid gemstoneId id passed")
	}

	if priceMultiplier <= 0 {
		return price, errors.New("Pricing data for gemstone unavailable:Gem Id - " + strconv.Itoa(gemstoneId) + "; Weight - " + strconv.FormatFloat(weight, 'f', -1, 64))
	}

	price = weight * priceMultiplier
	if price <= 0 {
		return price, errors.New("gemstoneId price coming out to be <= 0")
	}
	return price, nil
}
