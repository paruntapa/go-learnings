package cart

import (
	"backend-apis/types"
	"fmt"
)

func GetCartItemsIDs(items []types.CartItem) ([]int, error) {
	productsIds := make([]int, len(items))

	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("Invalid Quantity for product %d", item.ProductID)
		}

		productsIds[i] = item.ProductID
	}

	return productsIds, nil
}

func (h *Handler) CreateOrder(ps []types.Product, items []types.CartItem, userID int) (int, float64, error) {
	productMap := make(map[int]types.Product)

	for _, product := range ps {
		productMap[product.ID] = product
	}

	if err := CheckIfCartIsInStock(items, productMap); err != nil {
		return 0, 0, nil
	}

	totalPrice := CalculateTotalPrice(items, productMap)

	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity

		h.productStore.UpdateProduct(product)
	}

	orderID, err := h.store.CreateOrder(types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "Sailok Colony, Karbari Grant, ShimlaByPass Road",
	})

	if err != nil {
		return 0, 0, err
	}

	for _, item := range items {
		h.store.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
		})
	}

	return orderID, totalPrice, nil

}

func CalculateTotalPrice(cartItems []types.CartItem, products map[int]types.Product) float64 {
	var total float64

	for _, item := range cartItems {
		product := products[item.ProductID]
		total += product.Price * float64(item.Quantity)
	}

	return total
}

func CheckIfCartIsInStock(cartItems []types.CartItem, products map[int]types.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("Your Cart is EMPTY")
	}

	for _, item := range cartItems {
		product, ok := products[item.ProductID]

		if !ok {
			return fmt.Errorf("Product %d is not available in the store, please refresh your cart", item.ProductID)
		}

		if product.Quantity < item.Quantity {
			fmt.Errorf("Product %s is not avaiable in the quantity requested", product.Name)
		}
	}

	return nil
}
