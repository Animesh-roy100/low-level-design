package main

/*
- Restaurant aggregates Category (one-to-many).
- Category aggregates Item (one-to-many).
- Item has VariantGroup (one-to-many) and AddOnGroup (one-to-many).
- VariantGroup has Variant (one-to-many), each with a price.
- AddOnGroup has AddOn (one-to-many).
- AddOn links to Variant via a many-to-many relationship for variant-specific pricing.
*/

type Address struct {
	Street  string
	City    string
	State   string
	ZipCode string
}

type Restaurant struct {
	RestaurantID string
	Name         string
	Address      Address
}

type Category struct {
	CategoryID   string
	Name         string
	Description  string
	RestaurantID string
}

type Item struct {
	ItemID       string
	Name         string
	Description  string
	CategoryID   string
	RestaurantID string
}

type VariantGroup struct {
	VariantGroupID string
	Name           string
	ItemID         string
	MinSelection   int
	MaxSelection   int
}

type Variant struct {
	VariantID      string
	Name           string
	Price          float64
	VariantGroupID string
}

type AddOnGroup struct {
	AddOnGroupID string
	Name         string
	MinSelection int
	MaxSelection int
	ItemID       string
}

type AddOn struct {
	AddOnID      string
	Name         string
	AddOnGroupID string
}

type AddOnPrice struct {
	AddOnPriceID string
	AddOnID      string
	VariantID    string
	Price        float64
}

// Repository Pattern
type RestaurantRepository interface {
	Save(restaurant Restaurant) error
	GetByID(restaurantID string) (Restaurant, error)
}

type CategoryRepository interface {
	Save(category Category) error
	GetByID(categoryID string) (Category, error)
	GetByRestaurantID(restaurantID string) ([]Category, error)
}

type ItemRepository interface {
	Save(item Item) error
	GetByID(itemID string) (Item, error)
	GetByCategoryID(categoryID string) ([]Item, error)
}

type VariantGroupRepository interface {
	Save(variantGroup VariantGroup) error
	GetByID(variantGroupID string) (VariantGroup, error)
	GetByItemID(itemID string) ([]VariantGroup, error)
}

type VariantRepository interface {
	Save(variant Variant) error
	GetByID(variantID string) (Variant, error)
	GetByVariantGroupID(variantGroupID string) ([]Variant, error)
}

type AddOnGroupRepository interface {
	Save(addOnGroup AddOnGroup) error
	GetByID(addOnGroupID string) (AddOnGroup, error)
	GetByItemID(itemID string) ([]AddOnGroup, error)
}

type AddOnRepository interface {
	Save(addOn AddOn) error
	GetByID(addOnID string) (AddOn, error)
	GetByAddOnGroupID(addOnGroupID string) ([]AddOn, error)
}

type AddOnPriceRepository interface {
	Save(addOnPrice AddOnPrice) error
	GetByVariantAndAddOn(variantID, addOnID string) (AddOnPrice, error)
}

// In-memory cache (for demonstration purposes)

// Item Builder Pattern
type ItemBuilder struct {
}

// Pricing Strategy Pattern
type PricingStrategy interface {
	CalculatePrice(variant Variant, addOns []AddOn) float64
}

//
