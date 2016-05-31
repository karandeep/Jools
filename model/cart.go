package model

import (
	"../lib"
	"database/sql"
	"encoding/json"
	"log"
)

const MAX_PRODUCTS_IN_CART = 50

type CartProduct struct {
	Id    int
	Name  string
	Index string
	Price float64
	Image string
	Desc  string
	Qty   int
}

type Cart struct {
	NumProducts int
	Indices     map[string]int
	Products    []CartProduct
}

/*type ByIncreasingPrice []CartProduct

func (products ByIncreasingPrice) Len() int {
	return len(products)
}
func (products ByIncreasingPrice) Swap(i, j int) {
	products[i], products[j] = products[j], products[i]
}
func (products ByIncreasingPrice) Less(i, j int) bool {
	var priceI, priceJ float64
	priceI = products[i].Price
	priceJ = products[j].Price
	return priceI < priceJ
}*/

func GetCartForUser(userId int64) (string, error) {
	conn := lib.GetDBConnection()
	row := conn.QueryRow("SELECT products FROM Cart WHERE userId = ?", userId)
	var products string
	err := row.Scan(&products)
	return products, err
}

func createCartForUser(userId int64, products string) error {
	conn := lib.GetDBConnection()
	stmt, err := conn.Prepare("INSERT INTO Cart(userId, products) VALUES (?,?)")
	_, err = stmt.Exec(userId, products)
	return err
}

func updateCartForUser(userId int64, products string) error {
	conn := lib.GetDBConnection()
	stmt, err := conn.Prepare("UPDATE Cart SET products = ? WHERE userId = ? LIMIT 1")
	_, err = stmt.Exec(products, userId)
	return err
}

func FetchCart(userId int64) (string, error) {
	curCart, err := GetCartForUser(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	var updatedCart Cart
	updatedCart.Indices = make(map[string]int, MAX_PRODUCTS_IN_CART)
	updatedCart.Products = make([]CartProduct, MAX_PRODUCTS_IN_CART)
	err = json.Unmarshal([]byte(curCart), &updatedCart)
	if err != nil {
		log.Println(err)
		return "", err
	}
	//Making sure that price info for products in cart is being refreshed, since prices may change
	for i := 0; i < updatedCart.NumProducts; i++ {
		product, err := GetProductFromIndex(updatedCart.Products[i].Index)
		if err != nil {
			return "", err
		}
		updatedCart.Products[i].Price, err = product.GetPrice()
		if err != nil {
			return "", err
		}
	}

	updatedCart.Products = updatedCart.Products[:updatedCart.NumProducts]
	updatedCartJson, err := json.Marshal(updatedCart)
	if err != nil {
		log.Println(err)
		return "", err
	}

	updatedCartStr := string(updatedCartJson)
	if curCart != updatedCartStr {
		updateCartForUser(userId, updatedCartStr)
	}
	return updatedCartStr, nil
}

func SyncCart(userId int64, products string, removedProducts string) (string, error) {
	curCart, err := GetCartForUser(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			err = createCartForUser(userId, products)
		}
		return products, err
	}
	var existingCart, localCart, removedCart, mergedCart Cart
	existingCart.Indices = make(map[string]int, MAX_PRODUCTS_IN_CART)
	existingCart.Products = make([]CartProduct, MAX_PRODUCTS_IN_CART)
	err = json.Unmarshal([]byte(curCart), &existingCart)
	if err != nil {
		log.Println("Existing Cart:", err)
		return "", err
	}

	localCart.Indices = make(map[string]int, MAX_PRODUCTS_IN_CART)
	localCart.Products = make([]CartProduct, MAX_PRODUCTS_IN_CART)
	if products != "" {
		err = json.Unmarshal([]byte(products), &localCart)
		if err != nil {
			log.Println("Local Cart:", err)
			return "", err
		}
	}

	removedCart.Indices = make(map[string]int, MAX_PRODUCTS_IN_CART)
	removedCart.Products = make([]CartProduct, MAX_PRODUCTS_IN_CART)
	if removedProducts != "" {
		err = json.Unmarshal([]byte(removedProducts), &removedCart)
		if err != nil {
			log.Println("Removed Cart:", err)
			return "", err
		}
	}

	mergedCart.Indices = make(map[string]int, MAX_PRODUCTS_IN_CART)
	mergedCart.Products = make([]CartProduct, MAX_PRODUCTS_IN_CART)
	var indexM, indexE, indexL int
	//If a product has been removed, it will not be present in localCart. It may be present in existingCart.
	for indexM = 0; indexM < MAX_PRODUCTS_IN_CART && indexE < existingCart.NumProducts && indexL < localCart.NumProducts; {
		if existingCart.Products[indexE].Index == localCart.Products[indexL].Index {
			mergedCart.Products[indexM] = localCart.Products[indexL]
			mergedCart.Indices[localCart.Products[indexL].Index] = indexM
			indexE++
			indexL++
			indexM++
		} else if localCart.Products[indexL].Price > existingCart.Products[indexE].Price {
			mergedCart.Products[indexM] = localCart.Products[indexL]
			mergedCart.Indices[localCart.Products[indexL].Index] = indexM
			indexL++
			indexM++
		} else {
			removed := false
			index := existingCart.Products[indexE].Index
			_, removed = removedCart.Indices[index]
			if !removed {
				mergedCart.Products[indexM] = existingCart.Products[indexE]
				mergedCart.Indices[existingCart.Products[indexL].Index] = indexM
				indexM++
			}
			indexE++
		}
	}
	if indexE < existingCart.NumProducts {
		for ; indexM < MAX_PRODUCTS_IN_CART && indexE < existingCart.NumProducts; indexE++ {
			removed := false
			index := existingCart.Products[indexE].Index
			_, removed = removedCart.Indices[index]
			if !removed {
				mergedCart.Products[indexM] = existingCart.Products[indexE]
				mergedCart.Indices[existingCart.Products[indexL].Index] = indexM
				indexM++
			}
		}
	}
	if indexL < localCart.NumProducts {
		for ; indexM < MAX_PRODUCTS_IN_CART && indexL < localCart.NumProducts; indexL++ {
			mergedCart.Products[indexM] = localCart.Products[indexL]
			mergedCart.Indices[localCart.Products[indexL].Index] = indexM
			indexM++
		}
	}
	mergedCart.NumProducts = indexM
	mergedCart.Products = mergedCart.Products[:indexM]

	//Making sure that price info from client is not being trusted
	for i := 0; i < indexM; i++ {
		product, err := GetProductFromIndex(mergedCart.Products[i].Index)
		if err != nil {
			return "", err
		}
		mergedCart.Products[i].Price, err = product.GetPrice()
		if err != nil {
			return "", err
		}
	}

	mergedProductsJson, err := json.Marshal(mergedCart)
	if err != nil {
		log.Println(err)
		return "", err
	}
	updateCartForUser(userId, string(mergedProductsJson))
	return string(mergedProductsJson), nil
}

func GetCostOfCartProducts(products string) (float64, error) {
	var totalCost float64
	var cart Cart
	cart.Indices = make(map[string]int, MAX_PRODUCTS_IN_CART)
	cart.Products = make([]CartProduct, MAX_PRODUCTS_IN_CART)
	err := json.Unmarshal([]byte(products), &cart)
	if err != nil {
		return totalCost, err
	}
	//Making sure that price info from client is not being trusted
	for i := 0; i < cart.NumProducts; i++ {
		product, err := GetProductFromIndex(cart.Products[i].Index)
		if err != nil {
			return totalCost, err
		}
		cart.Products[i].Price, err = product.GetPrice()
		if err != nil {
			return totalCost, err
		}
		totalCost += cart.Products[i].Price * float64(cart.Products[i].Qty)
	}
	return totalCost, nil
}
