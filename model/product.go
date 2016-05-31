package model

import (
	"../config"
	"../lib"
	"database/sql"
	"errors"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
)

const (
	ALL_PRODUCTS = iota
	DIAMOND_RING
	RING
	DIAMOND_PENDANT
	PENDANT
	DIAMOND_EARRING
	EARRING
	DIAMOND_BANGLE
	BANGLE
	DIAMOND_NOSEPIN
	NOSEPIN
	CHAIN
)
const DESIGN_ID_RANGE = 50

type ProductCategory struct {
	Id        int
	Name      string
	ShortName string
}

var productCategoryList []ProductCategory = []ProductCategory{
	ALL_PRODUCTS:    ProductCategory{Id: ALL_PRODUCTS, Name: "All products", ShortName: "AL"},
	DIAMOND_RING:    ProductCategory{Id: DIAMOND_RING, Name: "Ring", ShortName: "DR"},
	RING:            ProductCategory{Id: RING, Name: "Ring", ShortName: "GR"},
	DIAMOND_PENDANT: ProductCategory{Id: DIAMOND_PENDANT, Name: "Pendant", ShortName: "DP"},
	PENDANT:         ProductCategory{Id: PENDANT, Name: "Pendant", ShortName: "GP"},
	DIAMOND_EARRING: ProductCategory{Id: DIAMOND_EARRING, Name: "Earring", ShortName: "DT"},
	EARRING:         ProductCategory{Id: EARRING, Name: "Earring", ShortName: "GT"},
	DIAMOND_BANGLE:  ProductCategory{Id: DIAMOND_BANGLE, Name: "Bangle", ShortName: "DB"},
	BANGLE:          ProductCategory{Id: BANGLE, Name: "Bangle", ShortName: "GB"},
	DIAMOND_NOSEPIN: ProductCategory{Id: DIAMOND_NOSEPIN, Name: "Nosepin", ShortName: "DN"},
	NOSEPIN:         ProductCategory{Id: NOSEPIN, Name: "Nosepin", ShortName: "GN"},
	CHAIN:           ProductCategory{Id: CHAIN, Name: "Chain", ShortName: "SC"},
}

type Product struct {
	Id                     int
	EncId                  string
	Created                int32
	DirName                string
	Name                   string
	Category               int
	Description            string
	CenterStone            int
	CenterStoneCategory    int
	CenterStoneWt          string //In carats for diamond
	CenterStoneDim         string
	CenterStonePieces      string
	CenterStoneShape       int
	CenterStoneSetting     int
	AccentStoneOne         int
	AccentStoneOneCategory int
	AccentStoneOneWt       string //In carats for diamond
	AccentStoneOneDim      string
	AccentStoneOnePieces   string
	AccentStoneOneShape    int
	AccentStoneOneSetting  int
	AccentStoneTwo         int
	AccentStoneTwoCategory int
	AccentStoneTwoWt       string //In carats for diamond
	AccentStoneTwoDim      string
	AccentStoneTwoPieces   string
	AccentStoneTwoShape    int
	AccentStoneTwoSetting  int
	PrimaryMetal           int
	PrimaryMetalColor      int
	PrimaryMetalWt         float64 //In grams
	DeliveryMonth          string
	DeliveryDay            int
	CenterStoneClarity     int
	AccentStoneOneClarity  int
	AccentStoneTwoClarity  int
	Length                 float64
	Height                 float64
	Width                  float64
}

type ProductVariant struct {
	ImageIndex            string
	Price                 float64
	CenterStone           int
	AccentStoneOne        int
	AccentStoneTwo        int
	PrimaryMetal          int
	PrimaryMetalColor     int
	CenterStoneClarity    int
	AccentStoneOneClarity int
	AccentStoneTwoClarity int
}

type ProductInfo struct {
	ProductData     Product
	ProductVariants map[string]ProductVariant
}

type ByIncreasingPrice []ProductInfo
type ByDecreasingPrice []ProductInfo

func (info ByIncreasingPrice) Len() int {
	return len(info)
}
func (info ByIncreasingPrice) Swap(i, j int) {
	info[i], info[j] = info[j], info[i]
}
func (info ByIncreasingPrice) Less(i, j int) bool {
	var priceI, priceJ float64
	for _, variantData := range info[i].ProductVariants {
		priceI = variantData.Price
		break
	}
	for _, variantData := range info[j].ProductVariants {
		priceJ = variantData.Price
		break
	}
	return priceI < priceJ
}

func (info ByDecreasingPrice) Len() int {
	return len(info)
}
func (info ByDecreasingPrice) Swap(i, j int) {
	info[i], info[j] = info[j], info[i]
}
func (info ByDecreasingPrice) Less(i, j int) bool {
	var priceI, priceJ float64
	for _, variantData := range info[i].ProductVariants {
		priceI = variantData.Price
		break
	}
	for _, variantData := range info[j].ProductVariants {
		priceJ = variantData.Price
		break
	}
	return priceI > priceJ
}

type ProductRefinement struct {
	Category          int
	Gemstone          int
	Setting           int
	PrimaryMetal      int
	PrimaryMetalColor int
	PriceMin          int
	PriceMax          int
	Sort              int
}

const (
	SORT_LATEST = iota
	SORT_PRICE_HIGH_LOW
	SORT_PRICE_LOW_HIGH
)
const PRODUCT_FETCH_FIELDS = "id,encId,created,dirName,name,category,centerStone,centerStoneWt,centerStoneDim,centerStonePieces,centerStoneCategory,centerStoneShape,centerStoneSetting,accentStoneOne,accentStoneOneWt,accentStoneOneDim,accentStoneOnePieces,accentStoneOneCategory,accentStoneOneShape,accentStoneOneSetting,accentStoneTwo,accentStoneTwoWt,accentStoneTwoDim,accentStoneTwoPieces,accentStoneTwoCategory,accentStoneTwoShape,accentStoneTwoSetting,primaryMetal,primaryMetalColor,primaryMetalWt,length,height,width"

func scanAllProductFields(rows *sql.Rows, products []Product) error {
	index := 0
	for rows.Next() {
		if err := rows.Scan(
			&products[index].Id,
			&products[index].EncId,
			&products[index].Created,
			&products[index].DirName,
			&products[index].Name,
			&products[index].Category,
			&products[index].CenterStone,
			&products[index].CenterStoneWt,
			&products[index].CenterStoneDim,
			&products[index].CenterStonePieces,
			&products[index].CenterStoneCategory,
			&products[index].CenterStoneShape,
			&products[index].CenterStoneSetting,
			&products[index].AccentStoneOne,
			&products[index].AccentStoneOneWt,
			&products[index].AccentStoneOneDim,
			&products[index].AccentStoneOnePieces,
			&products[index].AccentStoneOneCategory,
			&products[index].AccentStoneOneShape,
			&products[index].AccentStoneOneSetting,
			&products[index].AccentStoneTwo,
			&products[index].AccentStoneTwoWt,
			&products[index].AccentStoneTwoDim,
			&products[index].AccentStoneTwoPieces,
			&products[index].AccentStoneTwoCategory,
			&products[index].AccentStoneTwoShape,
			&products[index].AccentStoneTwoSetting,
			&products[index].PrimaryMetal,
			&products[index].PrimaryMetalColor,
			&products[index].PrimaryMetalWt,
			&products[index].Length,
			&products[index].Height,
			&products[index].Width,
		); err != nil {
			return err
		}
		index++
	}
	return nil
}
func scanAllProductFieldsFromRow(row *sql.Row, product *Product) error {
	if err := row.Scan(
		&product.Id,
		&product.EncId,
		&product.Created,
		&product.DirName,
		&product.Name,
		&product.Category,
		&product.CenterStone,
		&product.CenterStoneWt,
		&product.CenterStoneDim,
		&product.CenterStonePieces,
		&product.CenterStoneCategory,
		&product.CenterStoneShape,
		&product.CenterStoneSetting,
		&product.AccentStoneOne,
		&product.AccentStoneOneWt,
		&product.AccentStoneOneDim,
		&product.AccentStoneOnePieces,
		&product.AccentStoneOneCategory,
		&product.AccentStoneOneShape,
		&product.AccentStoneOneSetting,
		&product.AccentStoneTwo,
		&product.AccentStoneTwoWt,
		&product.AccentStoneTwoDim,
		&product.AccentStoneTwoPieces,
		&product.AccentStoneTwoCategory,
		&product.AccentStoneTwoShape,
		&product.AccentStoneTwoSetting,
		&product.PrimaryMetal,
		&product.PrimaryMetalColor,
		&product.PrimaryMetalWt,
		&product.Length,
		&product.Height,
		&product.Width,
	); err != nil {
		return err
	}
	return nil
}

func (product *Product) GetMakingCharge() (float64, error) {
	var makingCharge float64
	if product.PrimaryMetal == GOLD_22KT {
		if product.PrimaryMetalWt < 2 {
			makingCharge = 600
		} else if product.PrimaryMetalWt <= 10 {
			makingCharge = 500
		} else {
			makingCharge = 450
		}
	} else if product.PrimaryMetal == GOLD_18KT ||
		product.PrimaryMetal == GOLD_14KT {
		if product.CenterStone == DIAMOND {
			if product.PrimaryMetalWt < 2 {
				makingCharge = 700
			} else if product.PrimaryMetalWt <= 10 {
				makingCharge = 650
			} else {
				makingCharge = 600
			}
		} else {
			if product.PrimaryMetalWt < 2 {
				makingCharge = 800
			} else if product.PrimaryMetalWt <= 10 {
				makingCharge = 700
			} else {
				makingCharge = 650
			}
		}
	} else if product.PrimaryMetal == SILVER {
		if product.CenterStone == DIAMOND {
			if product.PrimaryMetalWt < 2 {
				makingCharge = 700
			} else if product.PrimaryMetalWt <= 10 {
				makingCharge = 650
			} else {
				makingCharge = 600
			}
		} else {
			if product.PrimaryMetalWt < 2 {
				makingCharge = 1000
			} else if product.PrimaryMetalWt <= 10 {
				makingCharge = 750
			} else {
				makingCharge = 700
			}
		}
	}

	makingCharge *= product.PrimaryMetalWt

	if makingCharge <= 0.0 {
		return makingCharge, errors.New("Making charge coming out to be <= 0")
	}
	return makingCharge, nil
}

func (product *Product) GetPrice() (float64, error) {
	var price, centerStonePrice, accentStoneOnePrice, accentStoneTwoPrice float64

	makingCharge, err := product.GetMakingCharge()
	if err != nil {
		return price, err
	}

	metalPrice, err := GetMetalPrice(product.PrimaryMetal, product.PrimaryMetalWt)
	if err != nil {
		return price, err
	}

	if product.CenterStone != NO_STONE {
		weightList := strings.Split(product.CenterStoneWt, ";")
		numPrices := len(weightList)
		pieceCountList := strings.Split(product.CenterStonePieces, ";")

		for i := 0; i < numPrices; i++ {
			if weightList[i] == "" {
				break
			}
			curWeight, err := strconv.ParseFloat(weightList[i], 64)
			if err != nil {
				return price, err
			}
			pieceCount, err := strconv.ParseFloat(pieceCountList[i], 64)
			if err != nil {
				return price, err
			}

			tempPrice, err := GetGemstonePrice(product.CenterStone, curWeight/pieceCount, product.CenterStoneCategory, product.CenterStoneClarity)
			if err != nil {
				return price, err
			}
			centerStonePrice += (tempPrice * pieceCount)
		}
	}

	if product.AccentStoneOne != NO_STONE {
		weightList := strings.Split(product.AccentStoneOneWt, ";")
		numPrices := len(weightList)
		pieceCountList := strings.Split(product.AccentStoneOnePieces, ";")

		for i := 0; i < numPrices; i++ {
			if weightList[i] == "" {
				break
			}
			curWeight, err := strconv.ParseFloat(weightList[i], 64)
			if err != nil {
				return price, err
			}
			pieceCount, err := strconv.ParseFloat(pieceCountList[i], 64)
			if err != nil {
				return price, err
			}

			tempPrice, err := GetGemstonePrice(product.AccentStoneOne, curWeight/pieceCount, GEMSTONE_CAT_ZERO, product.AccentStoneOneClarity)
			if err != nil {
				return price, err
			}
			accentStoneOnePrice += (tempPrice * pieceCount)
		}

	}

	if product.AccentStoneTwo != NO_STONE {
		weightList := strings.Split(product.AccentStoneTwoWt, ";")
		numPrices := len(weightList)
		pieceCountList := strings.Split(product.AccentStoneTwoPieces, ";")

		for i := 0; i < numPrices; i++ {
			if weightList[i] == "" {
				break
			}
			curWeight, err := strconv.ParseFloat(weightList[i], 64)
			if err != nil {
				return price, err
			}
			pieceCount, err := strconv.ParseFloat(pieceCountList[i], 64)
			if err != nil {
				return price, err
			}

			tempPrice, err := GetGemstonePrice(product.AccentStoneTwo, curWeight/pieceCount, GEMSTONE_CAT_ZERO, product.AccentStoneTwoClarity)
			if err != nil {
				return price, err
			}
			accentStoneTwoPrice += (tempPrice * pieceCount)
		}
	}

	price = metalPrice + makingCharge + centerStonePrice + accentStoneOnePrice + accentStoneTwoPrice
	if price <= 0 {
		return price, errors.New("Product price coming out to be <= 0")
	}
	if price < 4000 {
		price = 3999
	}
	return price, nil
}

func GetProduct(id int) (Product, error) {
	conn := lib.GetDBConnection()
	var product Product
	row := conn.QueryRow("SELECT "+PRODUCT_FETCH_FIELDS+" FROM Product WHERE id = ?", id)

	err := scanAllProductFieldsFromRow(row, &product)
	if err == nil {
		//Using this by default, since there's no value for these in DB currently
		product.CenterStoneClarity = GEM_CLARITY_ONE
		product.AccentStoneOneClarity = GEM_CLARITY_ONE
		product.AccentStoneTwoClarity = GEM_CLARITY_ONE
	}
	_, product.DeliveryMonth, product.DeliveryDay = GetExpectedDelivery(-1)
	product.DeliveryMonth = lib.MonthAbbr(product.DeliveryMonth)
	return product, err
}

func (product *Product) updateBaseFromVariant(variant ProductVariant) {
	product.CenterStone = variant.CenterStone
	product.AccentStoneOne = variant.AccentStoneOne
	product.AccentStoneTwo = variant.AccentStoneTwo
	product.PrimaryMetal = variant.PrimaryMetal
	product.PrimaryMetalColor = variant.PrimaryMetalColor
	product.CenterStoneClarity = variant.CenterStoneClarity
	product.AccentStoneOneClarity = variant.AccentStoneOneClarity
	product.AccentStoneTwoClarity = variant.AccentStoneTwoClarity
}

func (product *Product) getVariantInfo(price float64, imageIndex string) ProductVariant {
	var variant ProductVariant
	variant.ImageIndex = imageIndex
	variant.Price = price
	variant.CenterStone = product.CenterStone
	variant.AccentStoneOne = product.AccentStoneOne
	variant.AccentStoneTwo = product.AccentStoneTwo
	variant.PrimaryMetal = product.PrimaryMetal
	variant.PrimaryMetalColor = product.PrimaryMetalColor
	variant.CenterStoneClarity = product.CenterStoneClarity
	variant.AccentStoneOneClarity = product.AccentStoneOneClarity
	variant.AccentStoneTwoClarity = product.AccentStoneTwoClarity
	return variant
}

func GetProductFromIndex(uniqueVariantIndex string) (Product, error) {
	var product Product
	productInfo := strings.Split(uniqueVariantIndex, "-")
	productId, err := strconv.Atoi(productInfo[0])
	if err != nil {
		return product, err
	}
	product, err = GetProduct(productId)
	if err != nil {
		return product, err
	}
	product.PrimaryMetalColor = MetalColorReverseMap[productInfo[1][0:0]]
	var stoneIndex = 1
	var clarityIndex = 0
	if product.CenterStone != NO_STONE {
		product.CenterStone = GemstoneReverseMap[productInfo[1][stoneIndex:stoneIndex+GEMSTONE_SHORTNAME_LEN]]
		stoneIndex += GEMSTONE_SHORTNAME_LEN
		//Only for Diamond, the clarity index shud be anything other than 0
		product.CenterStoneClarity = 0
		if product.CenterStone == DIAMOND {
			product.CenterStoneClarity = GemstoneClarityReverseMap[productInfo[2][clarityIndex:clarityIndex+GEMSTONE_CLARITY_SHORTNAME_LEN]]
		}
		clarityIndex += GEMSTONE_CLARITY_SHORTNAME_LEN
	}
	if product.AccentStoneOne != NO_STONE {
		product.AccentStoneOne = GemstoneReverseMap[productInfo[1][stoneIndex:stoneIndex+GEMSTONE_SHORTNAME_LEN]]
		stoneIndex += GEMSTONE_SHORTNAME_LEN
		product.AccentStoneOneClarity = 0
		if product.AccentStoneOne == DIAMOND {
			product.AccentStoneOneClarity = GemstoneClarityReverseMap[productInfo[2][clarityIndex:clarityIndex+GEMSTONE_CLARITY_SHORTNAME_LEN]]
		}
		clarityIndex += GEMSTONE_CLARITY_SHORTNAME_LEN
	}
	if product.AccentStoneTwo != NO_STONE {
		product.AccentStoneTwo = GemstoneReverseMap[productInfo[1][stoneIndex:stoneIndex+GEMSTONE_SHORTNAME_LEN]]
		stoneIndex += GEMSTONE_SHORTNAME_LEN
		product.AccentStoneTwoClarity = 0
		if product.AccentStoneTwo == DIAMOND {
			product.AccentStoneTwoClarity = GemstoneClarityReverseMap[productInfo[2][clarityIndex:clarityIndex+GEMSTONE_CLARITY_SHORTNAME_LEN]]
		}
		clarityIndex += GEMSTONE_CLARITY_SHORTNAME_LEN
	}
	middlePartLen := len(productInfo[1])
	product.PrimaryMetal = MetalColorReverseMap[productInfo[1][middlePartLen-METAL_SHORTNAME_LEN:]]
	return product, err
}

func (product *Product) GetIndices() (string, string) {
	var imageIndex, uniqueIndex, temp string
	imageIndex = product.DirName
	uniqueIndex = strconv.Itoa(product.Id) + "-"
	imageIndex += MetalColorList[product.PrimaryMetalColor].ShortName
	uniqueIndex += MetalColorList[product.PrimaryMetalColor].ShortName

	if product.CenterStone != NO_STONE {
		imageIndex += GemstoneList[product.CenterStone].ShortName
		uniqueIndex += GemstoneList[product.CenterStone].ShortName
		temp += GemstoneClarityList[product.CenterStoneClarity].ShortName
	}
	if product.AccentStoneOne != NO_STONE {
		imageIndex += GemstoneList[product.AccentStoneOne].ShortName
		uniqueIndex += GemstoneList[product.AccentStoneOne].ShortName
		temp += GemstoneClarityList[product.AccentStoneOneClarity].ShortName
	}
	if product.AccentStoneTwo != NO_STONE {
		imageIndex += GemstoneList[product.AccentStoneTwo].ShortName
		uniqueIndex += GemstoneList[product.AccentStoneTwo].ShortName
		temp += GemstoneClarityList[product.AccentStoneTwoClarity].ShortName
	}

	uniqueIndex += MetalList[product.PrimaryMetal].ShortName
	if temp != "" {
		uniqueIndex += "-" + temp
	}
	return imageIndex, uniqueIndex
}

func (originalProduct Product) GetBaseVariant() (map[string]ProductVariant, error) {
	product := originalProduct
	var variant map[string]ProductVariant
	variant = make(map[string]ProductVariant, 1)
	imageIndex, uniqueIndex := product.GetIndices()
	price, err := product.GetPrice()
	if err != nil {
		return variant, err
	}
	variant[uniqueIndex] = product.getVariantInfo(price, imageIndex)
	return variant, nil
}

func (originalProduct Product) GetAllCombosForProduct() (map[string]ProductVariant, error) {
	product := originalProduct
	var productVariants map[string]ProductVariant
	configData := config.GetConfig()
	var centerStoneCounter, accentStoneOneCounter, accentStoneTwoCounter, primaryMetalCounter int
	var centerStoneClarityCounter, accentStoneOneClarityCounter, accentStoneTwoClarityCounter int
	centerStoneCounterLimit, accentStoneOneCounterLimit, accentStoneTwoCounterLimit, primaryMetalCounterLimit := 1, 1, 1, NUM_METALS
	centerStoneClarityCounterLimit, accentStoneOneClarityCounterLimit, accentStoneTwoClarityCounterLimit := 1, 1, 1

	if originalProduct.CenterStone != NO_STONE {
		centerStoneCounterLimit = NUM_GEMSTONES
	}
	if originalProduct.AccentStoneOne != NO_STONE {
		accentStoneOneCounterLimit = NUM_GEMSTONES
	}
	if originalProduct.AccentStoneTwo != NO_STONE {
		accentStoneTwoCounterLimit = NUM_GEMSTONES
	}

	productVariants = make(map[string]ProductVariant, centerStoneCounterLimit*accentStoneOneCounterLimit*accentStoneTwoCounterLimit*NUM_METAL_VARIANTS)
	for centerStoneCounter = 0; centerStoneCounter < centerStoneCounterLimit; centerStoneCounter++ {
		centerStoneClarityCounterLimit = 1
		if product.CenterStone != NO_STONE {
			product.CenterStone = centerStoneCounter
			if product.CenterStone == DIAMOND {
				centerStoneClarityCounterLimit = 3
			}
		}

		for centerStoneClarityCounter = 0; centerStoneClarityCounter < centerStoneClarityCounterLimit; centerStoneClarityCounter++ {
			product.CenterStoneClarity = centerStoneClarityCounter

			for accentStoneOneCounter = 0; accentStoneOneCounter < accentStoneOneCounterLimit; accentStoneOneCounter++ {
				accentStoneOneClarityCounterLimit = 1
				if product.AccentStoneOne != NO_STONE {
					product.AccentStoneOne = accentStoneOneCounter
					if product.AccentStoneOne == DIAMOND {
						accentStoneOneClarityCounterLimit = 3
					}
				}

				for accentStoneOneClarityCounter = 0; accentStoneOneClarityCounter < accentStoneOneClarityCounterLimit; accentStoneOneClarityCounter++ {
					product.AccentStoneOneClarity = accentStoneOneClarityCounter

					for accentStoneTwoCounter = 0; accentStoneTwoCounter < accentStoneTwoCounterLimit; accentStoneTwoCounter++ {
						accentStoneTwoClarityCounterLimit = 1
						if product.AccentStoneTwo != NO_STONE {
							product.AccentStoneTwo = accentStoneTwoCounter
							if product.AccentStoneTwo == DIAMOND {
								accentStoneTwoClarityCounterLimit = 3
							}
						}

						for accentStoneTwoClarityCounter = 0; accentStoneTwoClarityCounter < accentStoneTwoClarityCounterLimit; accentStoneTwoClarityCounter++ {
							product.AccentStoneTwoClarity = accentStoneTwoClarityCounter

							for primaryMetalCounter = 0; primaryMetalCounter < primaryMetalCounterLimit; primaryMetalCounter++ {
								product.PrimaryMetal = primaryMetalCounter
								price, err := product.GetPrice()
								if err != nil {
									return productVariants, err
								}
								product.PrimaryMetalColor = METAL_WHITE
								imageIndex, uniqueIndex := product.GetIndices()
								if configData.ENV != config.PRODUCTION {
									if _, found := productVariants[uniqueIndex]; found {
										log.Println("Overwrite happening - Original variant:", productVariants[uniqueIndex], "New info:", product)
										return productVariants, errors.New("About to overwrite a variant")
									}
								}
								productVariants[uniqueIndex] = product.getVariantInfo(price, imageIndex)
								if product.PrimaryMetal == GOLD_22KT ||
									product.PrimaryMetal == GOLD_18KT ||
									product.PrimaryMetal == GOLD_14KT {
									product.PrimaryMetalColor = METAL_YELLOW
									imageIndex, uniqueIndex = product.GetIndices()
									if configData.ENV != config.PRODUCTION {
										if _, found := productVariants[uniqueIndex]; found {
											log.Println("Overwrite happening - Original variant:", productVariants[uniqueIndex], "New info:", product)
											return productVariants, errors.New("About to overwrite a variant")
										}
									}
									productVariants[uniqueIndex] = product.getVariantInfo(price, imageIndex)

									product.PrimaryMetalColor = METAL_ROSE
									imageIndex, uniqueIndex = product.GetIndices()
									if configData.ENV != config.PRODUCTION {
										if _, found := productVariants[uniqueIndex]; found {
											log.Println("Overwrite happening - Original variant:", productVariants[uniqueIndex], "New info:", product)
											return productVariants, errors.New("About to overwrite a variant")
										}
									}
									productVariants[uniqueIndex] = product.getVariantInfo(price, imageIndex)
								}
							}
						}
					}
				}
			}
		}
	}

	return productVariants, nil
}

func GetInitialProductList(refinement ProductRefinement) ([]ProductInfo, error) {
	conn := lib.GetDBConnection()

	var products []Product
	var productList []ProductInfo
	var productCount int
	var err error
	var row *sql.Row
	var rows *sql.Rows
	if refinement.Category == ALL_PRODUCTS || refinement.Category == -1 {
		row = conn.QueryRow("SELECT COUNT(*) FROM Product")
		rows, err = conn.Query("SELECT " + PRODUCT_FETCH_FIELDS + " FROM Product ORDER by id DESC")
	} else {
		row = conn.QueryRow("SELECT COUNT(*) FROM Product WHERE category = ?", refinement.Category)
		rows, err = conn.Query("SELECT "+PRODUCT_FETCH_FIELDS+" FROM Product WHERE category = ? ORDER by id DESC", refinement.Category)
	}

	if err != nil {
		return productList, err
	}
	err = row.Scan(&productCount)
	if err != nil {
		return productList, err
	}
	products = make([]Product, productCount)
	productList = make([]ProductInfo, productCount)
	err = scanAllProductFields(rows, products)
	if err != nil {
		return productList, err
	}

	return getProductsToShow(products, productList, refinement, productCount)
}

func (product *Product) CreateKeywords() string {
	keywords := "buy %PRODUCT_NAME% %PRODUCT_CATEGORY%,%PRODUCT_NAME% %PRODUCT_CATEGORY% price in india,%PRODUCT_NAME% %PRODUCT_CATEGORY% price,%PRODUCT_NAME% %PRODUCT_CATEGORY%,price of %PRODUCT_NAME% %PRODUCT_CATEGORY%,%PRODUCT_NAME% %PRODUCT_CATEGORY% india,%PRODUCT_NAME% %PRODUCT_CATEGORY% specifications,%PRODUCT_NAME% %PRODUCT_CATEGORY% review,%PRODUCT_NAME% %PRODUCT_CATEGORY% price,%PRODUCT_NAME% %PRODUCT_CATEGORY% review,%PRODUCT_NAME% %PRODUCT_CATEGORY% customization,%PRODUCT_NAME% %PRODUCT_CATEGORY% design, %PRODUCT_NAME% %PRODUCT_CATEGORY% setting"
	keywords = strings.Replace(keywords, "%PRODUCT_NAME%", product.Name, -1)
	keywords = strings.Replace(keywords, "%PRODUCT_CATEGORY%", productCategoryList[product.Category].Name, -1)
	return keywords
}

func (product *Product) CreateDescription() string {
	var description string
	if product.Category == DIAMOND_RING && product.CenterStone == DIAMOND {
		description += GemstoneShapeList[product.CenterStoneShape].Name + " " + GemstoneList[product.CenterStone].Name + " " + productCategoryList[product.Category].Name + " set in " + MetalList[product.PrimaryMetal].Karat + " " + MetalColorList[product.PrimaryMetalColor].Name + " " + MetalList[product.PrimaryMetal].Name
	} else if product.Category == DIAMOND_RING || product.Category == RING {
		description += GemstoneShapeList[product.CenterStoneShape].Name + " " + GemstoneList[product.CenterStone].Name + " " + MetalList[product.PrimaryMetal].Karat + " " + MetalColorList[product.PrimaryMetalColor].Name + " " + MetalList[product.PrimaryMetal].Name + " " + productCategoryList[product.Category].Name
	} else {
		description += GemstoneShapeList[product.CenterStoneShape].Name + " " + GemstoneList[product.CenterStone].Name + " " + MetalList[product.PrimaryMetal].Karat + " " + MetalColorList[product.PrimaryMetalColor].Name + " " + MetalList[product.PrimaryMetal].Name + " " + productCategoryList[product.Category].Name
	}

	if product.AccentStoneOne != NO_STONE {
		description += " with " + GemstoneList[product.AccentStoneOne].Name
		if product.AccentStoneTwo != NO_STONE {
			description += " and " + GemstoneList[product.AccentStoneTwo].Name
		}
	}
	return description
}

func (product *Product) CreateTitle() string {
	title := product.Name + " " + productCategoryList[product.Category].Name
	return title
}

func (product *Product) CreateImageUrl(imageIndex string) string {
	configData := config.GetConfig()
	imageUrl := configData.STATIC_URL + "/images/designs/" + product.DirName + "/" + imageIndex + "-2.JPG"
	return imageUrl
}

func getProductsToShow(products []Product, productList []ProductInfo, refinement ProductRefinement, productCount int) ([]ProductInfo, error) {
	var err error
	gemstone := 0
	accentStone := 0
	metal := 0
	metalColor := 0
	for i := 0; i < productCount; i++ {
		if products[i].Id != 0 {
			if refinement.Gemstone == -1 {
				gemstone = i % NUM_GEMSTONES
			} else {
				gemstone = refinement.Gemstone
			}
			if refinement.PrimaryMetal == -1 {
				metal = i % NUM_METALS
			} else {
				metal = refinement.PrimaryMetal
			}
			if refinement.PrimaryMetalColor == -1 {
				if metal == GOLD_22KT || metal == GOLD_18KT || metal == GOLD_14KT {
					metalColor = i % NUM_METAL_COLORS
				} else {
					metalColor = METAL_WHITE
				}
			} else {
				metalColor = refinement.PrimaryMetalColor
			}
			if products[i].AccentStoneOne != -1 {
				accentStone = (accentStone + 3) % NUM_GEMSTONES
				products[i].AccentStoneOne = accentStone
			}
			products[i].CenterStone = gemstone
			products[i].PrimaryMetal = metal
			products[i].PrimaryMetalColor = metalColor
			productList[i].ProductData = products[i]
			productList[i].ProductVariants, err = products[i].GetBaseVariant()
			if err != nil {
				log.Println(err)
			}
			if refinement.PriceMin != -1 || refinement.PriceMax != -1 {
				variantIndex := ""
				productList[i].ProductVariants, err = products[i].GetAllCombosForProduct()
				productList[i].ProductVariants, variantIndex, err = getProductVariantInPriceRange(productList[i].ProductVariants, refinement)
				if err != nil {
					productList[i].ProductData = Product{}
				} else {
					productList[i].ProductData.updateBaseFromVariant(productList[i].ProductVariants[variantIndex])
				}
			}
		}
	}

	if refinement.Sort == SORT_LATEST {
	} else if refinement.Sort == SORT_PRICE_HIGH_LOW {
		sort.Sort(ByDecreasingPrice(productList))
	} else if refinement.Sort == SORT_PRICE_LOW_HIGH {
		sort.Sort(ByIncreasingPrice(productList))
	}
	return productList, nil
}

func GetNextProductList(lastId int, refinement ProductRefinement) ([]ProductInfo, error) {
	PRODUCT_LOT_SIZE := 9
	limitId := lastId - DESIGN_ID_RANGE //Since ids of designs uploaded may be missing values
	conn := lib.GetDBConnection()
	products := make([]Product, PRODUCT_LOT_SIZE)
	productList := make([]ProductInfo, PRODUCT_LOT_SIZE)

	var err error
	var rows *sql.Rows
	if refinement.Category == ALL_PRODUCTS {
		rows, err = conn.Query("SELECT "+PRODUCT_FETCH_FIELDS+" FROM Product WHERE id < ? AND id > ? ORDER BY id DESC LIMIT ?", lastId, limitId, PRODUCT_LOT_SIZE)
	} else {
		rows, err = conn.Query("SELECT "+PRODUCT_FETCH_FIELDS+" FROM Product WHERE category = ? AND id < ? AND id > ? ORDER BY id DESC LIMIT ?", refinement.Category, lastId, limitId, PRODUCT_LOT_SIZE)
	}
	if err != nil {
		return productList, err
	}

	err = scanAllProductFields(rows, products)
	if err != nil || products[0].Id == 0 {
		//No matches were found. Try once more without limiting - in case there are x consecutive ids missing
		if refinement.Category == ALL_PRODUCTS {
			rows, err = conn.Query("SELECT "+PRODUCT_FETCH_FIELDS+" FROM Product WHERE id < ? ORDER BY id DESC LIMIT ?", lastId, PRODUCT_LOT_SIZE)
		} else {
			rows, err = conn.Query("SELECT "+PRODUCT_FETCH_FIELDS+" FROM Product WHERE category = ? AND id < ? ORDER BY id DESC LIMIT ?", refinement.Category, lastId, PRODUCT_LOT_SIZE)
		}
		if err != nil {
			return productList, err
		}
		err = scanAllProductFields(rows, products)
		if err != nil {
			return productList, err
		}
	}

	return getProductsToShow(products, productList, refinement, PRODUCT_LOT_SIZE)
}

func getProductVariantInPriceRange(variants map[string]ProductVariant, refinement ProductRefinement) (map[string]ProductVariant, string, error) {
	var variantInRange map[string]ProductVariant
	variantInRange = make(map[string]ProductVariant, 1)
	if refinement.PriceMax == -1 {
		refinement.PriceMax = math.MaxInt32
	}
	found := false
	variantIndex := ""
	priceMinFloat := float64(refinement.PriceMin)
	priceMaxFloat := float64(refinement.PriceMax)
	for index, variantData := range variants {
		if variantData.Price >= priceMinFloat && variantData.Price <= priceMaxFloat {
			if (refinement.Gemstone == -1 || variantData.CenterStone == refinement.Gemstone) && (refinement.PrimaryMetal == -1 || refinement.PrimaryMetal == variantData.PrimaryMetal) {
				variantInRange[index] = variantData
				variantIndex = index
				found = true
				break
			}
		}
	}
	if found {
		return variantInRange, variantIndex, nil
	} else {
		return variantInRange, variantIndex, errors.New("No variant in price range")
	}
}
