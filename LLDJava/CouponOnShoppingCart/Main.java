package LLDJava.CouponOnShoppingCart;

import java.math.BigDecimal;
import java.math.RoundingMode;
import java.util.*;

/* ---------- Domain ---------- */

enum ProductType {
    DECORATIVE_GOODS,
    ELECTRONIC_GOODS,
    FURNITURE_GOODS,
}

abstract class Product {
    private final String name;
    private final BigDecimal price;
    private final ProductType productType;

    protected Product(String name, BigDecimal price, ProductType productType) {
        this.name = name;
        this.price = price;
        this.productType = productType;
    }

    public String getName() { return name; }
    public BigDecimal getPrice() { return price; }
    public ProductType getProductType() { return productType; }
}

/* Use a single concrete product class unless special behavior is needed */
class ConcreteProduct extends Product {
    public ConcreteProduct(String name, BigDecimal price, ProductType productType) {
        super(name, price, productType);
    }
}

/* Line item (supports quantity) */
class LineItem {
    private final Product product;
    private final int quantity;

    public LineItem(Product product, int quantity) {
        if (quantity <= 0) throw new IllegalArgumentException("quantity must be >= 1");
        this.product = product;
        this.quantity = quantity;
    }

    public Product getProduct() { return product; }
    public int getQuantity() { return quantity; }
}

/* ---------- Cart context (per-calculation state) ---------- */

class CartContext {
    private final Map<ProductType, Integer> seenCounts = new EnumMap<>(ProductType.class);

    public int getSeenCount(ProductType type) {
        return seenCounts.getOrDefault(type, 0);
    }

    public void incrementSeen(ProductType type) {
        seenCounts.put(type, getSeenCount(type) + 1);
    }
}

/* ---------- Coupon abstraction ---------- */

interface Coupon {
    /**
     * Apply coupon to currentPrice for the given product using cart context.
     * Should return a new BigDecimal (do not mutate the incoming).
     */
    BigDecimal apply(BigDecimal currentPrice, Product product, CartContext ctx);
}

/* Flat N% off on every individual item */
class FlatDiscountPercentageCoupon implements Coupon {
    private final BigDecimal percent; // e.g., 10 for 10%

    public FlatDiscountPercentageCoupon(BigDecimal percent) {
        this.percent = percent;
    }

    public BigDecimal apply(BigDecimal currentPrice, Product product, CartContext ctx) {
        BigDecimal multiplier = BigDecimal.ONE.subtract(percent.divide(BigDecimal.valueOf(100)));
        return currentPrice.multiply(multiplier);
    }
}

/* D% off on items of eligible product types */
class ProductTypeDiscountPercentageCoupon implements Coupon {
    private final BigDecimal percent;
    private final Set<ProductType> eligibleTypes;

    public ProductTypeDiscountPercentageCoupon(BigDecimal percent, Set<ProductType> eligibleTypes) {
        this.percent = percent;
        this.eligibleTypes = new HashSet<>(eligibleTypes);
    }

    public BigDecimal apply(BigDecimal currentPrice, Product product, CartContext ctx) {
        if (eligibleTypes.contains(product.getProductType())) {
            BigDecimal multiplier = BigDecimal.ONE.subtract(percent.divide(BigDecimal.valueOf(100)));
            return currentPrice.multiply(multiplier);
        }
        return currentPrice;
    }
}

/* P% off on next item(s) of same type: applies when ctx.getSeenCount(type) > 0 */
class NextItemDiscountPercentageCoupon implements Coupon {
    private final BigDecimal percent;

    public NextItemDiscountPercentageCoupon(BigDecimal percent) {
        this.percent = percent;
    }

    public BigDecimal apply(BigDecimal currentPrice, Product product, CartContext ctx) {
        int seen = ctx.getSeenCount(product.getProductType());
        if (seen > 0) {
            BigDecimal multiplier = BigDecimal.ONE.subtract(percent.divide(BigDecimal.valueOf(100)));
            return currentPrice.multiply(multiplier);
        }
        return currentPrice;
    }
}

/* ---------- Shopping cart ---------- */

class ShoppingCart {
    private final List<LineItem> items = new ArrayList<>();
    private final List<Coupon> coupons = new ArrayList<>();

    public void addItem(LineItem item) {
        items.add(item);
    }

    public void addCoupon(Coupon coupon) {
        coupons.add(coupon);
    }

    /**
     * Calculate total applying coupons sequentially to each item (unit-by-unit).
     * Returns total rounded to 2 decimals (HALF_UP).
     */
    public BigDecimal calculateTotal() {
        CartContext ctx = new CartContext();
        BigDecimal total = BigDecimal.ZERO;

        for (LineItem li : items) {
            Product product = li.getProduct();
            BigDecimal baseUnitPrice = product.getPrice();

            // Apply per-unit to handle NextItem semantics correctly for qty>1
            for (int i = 0; i < li.getQuantity(); i++) {
                BigDecimal unitPrice = baseUnitPrice;

                // apply coupons in the configured order
                for (Coupon c : coupons) {
                    unitPrice = c.apply(unitPrice, product, ctx);
                }

                // normalize/round each unit price to 2 decimal places
                BigDecimal unitRounded = unitPrice.setScale(2, RoundingMode.HALF_UP);
                total = total.add(unitRounded);

                // update context after processing this unit
                ctx.incrementSeen(product.getProductType());
            }
        }

        return total.setScale(2, RoundingMode.HALF_UP);
    }
}

/* ---------- Demo main ---------- */

public class Main {
    public static void main(String[] args) {
        Product item1 = new ConcreteProduct("lights", BigDecimal.valueOf(1000.00), ProductType.DECORATIVE_GOODS);
        Product item2 = new ConcreteProduct("sofa", BigDecimal.valueOf(2000.00), ProductType.FURNITURE_GOODS);
        Product item3 = new ConcreteProduct("fancy-lights", BigDecimal.valueOf(1200.00), ProductType.DECORATIVE_GOODS);
        Product item4 = new ConcreteProduct("basic-android-phone", BigDecimal.valueOf(10000.00), ProductType.ELECTRONIC_GOODS);

        ShoppingCart cart = new ShoppingCart();

        // Add items (quantity is supported)
        cart.addItem(new LineItem(item1, 1));
        cart.addItem(new LineItem(item2, 1));
        cart.addItem(new LineItem(item3, 1));
        cart.addItem(new LineItem(item4, 1));

        // Configure coupons and order of application
        cart.addCoupon(new FlatDiscountPercentageCoupon(BigDecimal.valueOf(10))); // 10% off all
        cart.addCoupon(new ProductTypeDiscountPercentageCoupon(BigDecimal.valueOf(3),
                Set.of(ProductType.DECORATIVE_GOODS, ProductType.FURNITURE_GOODS))); // extra 3% for some types
        cart.addCoupon(new NextItemDiscountPercentageCoupon(BigDecimal.valueOf(5))); // 5% off for next item of same type

        BigDecimal total = cart.calculateTotal();

        System.out.println("Final total: " + total); // printed with 2 decimals
    }
}
