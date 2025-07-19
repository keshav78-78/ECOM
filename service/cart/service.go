package cart

import (
	"fmt"

	"github.com/keshav78-78/ECOM/types"
)

// getCartItemsIDs extracts product IDs from a list of cart items.
// This is a helper function for our service logic.
func getCartItemsIDs(items []types.CartItem) ([]int, error) {
	productIDs := make([]int, len(items))
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for product %d", item.ProductID)
		}
		productIDs[i] = item.ProductID
	}
	return productIDs, nil
}

// createOrder handles the core business logic of creating an order.
// It's a method on the Handler struct, so it can access the necessary stores.
func (h *Handler) createOrder(products []types.Product, items []types.CartItem, userID int) (int, float64, error) {
	// Create a map for quick product lookups by ID.
	productMap := make(map[int]types.Product)
	for _, product := range products {
		productMap[product.ID] = product
	}

	// Check if all products are in stock with the requested quantity.
	if err := checkIfCartIsInStock(items, productMap); err != nil {
		return 0, 0, err
	}

	// Calculate the total price based on the items and their quantities.
	totalPrice := calculateTotalPrice(items, productMap)

	// Reduce the stock for each product in the database.
	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity
		if err := h.productStore.UpdateProduct(product); err != nil {
			// If updating stock fails for any product, stop the entire process.
			return 0, 0, err
		}
	}

	// Create the main order entry in the database.
	orderID, err := h.store.CreateOrder(types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "some default address", // This should probably come from the user's profile.
	})
	if err != nil {
		return 0, 0, err
	}

	// After the main order is created, create an entry for each item in that order.
	for _, item := range items {
		if err := h.store.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
		}); err != nil {
			// This is a critical failure. Ideally, this whole process should be in a
			// database transaction that can be rolled back if this step fails.
			return 0, 0, err
		}
	}

	return orderID, totalPrice, nil
}

// checkIfCartIsInStock is a helper function to validate product availability.
func checkIfCartIsInStock(cartItems []types.CartItem, products map[int]types.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range cartItems {
		product, ok := products[item.ProductID]
		if !ok {
			return fmt.Errorf("product %d is not available in the store", item.ProductID)
		}
		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %s is not available in the requested quantity", product.Name)
		}
	}

	return nil
}

// calculateTotalPrice is a helper function to calculate the final price of all items.
func calculateTotalPrice(items []types.CartItem, products map[int]types.Product) float64 {
	var total float64
	for _, item := range items {
		product := products[item.ProductID]
		total += product.Price * float64(item.Quantity)
	}
	return total
}
