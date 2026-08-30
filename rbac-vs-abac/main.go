package main

import "fmt"

func main() {
	customer := User{
		ID:   "1",
		Name: "User",
		Role: RoleCustomer,
	}

	admin := User{
		ID:   "2",
		Name: "Admin",
		Role: RoleAdmin,
	}

	seller := User{
		ID:   "3",
		Name: "Seller",
		Role: RoleSeller,
	}

	seller2 := User{
		ID:   "4",
		Name: "Seller 2",
		Role: RoleSeller,
	}

	product := Product{
		ID:       "1",
		OwnerID:  seller.ID,
		Name:     "Product 1",
		Category: "Category 1",
	}
	engine := Engine{
		Policies: []Policy{OwnershipPolicy},
	}

	fmt.Println("=== Stsenariy 1: seller2 BOSHQA sellerning mahsulotini tahrirlamoqchi ===")
	fmt.Printf("RBAC natija: %v\n", CanRBAC(seller2, ActionUpdateProduct))
	fmt.Printf("ABAC natija: %v\n\n", engine.CanABAC(seller2, product, ActionUpdateProduct, Environment{}))

	fmt.Println("=== Stsenariy 2: seller1 O'Z mahsulotini tahrirlamoqchi ===")
	fmt.Printf("RBAC natija: %v\n", CanRBAC(seller, ActionUpdateProduct))
	fmt.Printf("ABAC natija: %v\n\n", engine.CanABAC(seller, product, ActionUpdateProduct, Environment{}))

	fmt.Println("=== Stsenariy 3: admin istalgan mahsulotni tahrirlamoqchi ===")
	fmt.Printf("RBAC natija: %v\n", CanRBAC(admin, ActionUpdateProduct))
	fmt.Printf("ABAC natija: %v\n\n", engine.CanABAC(admin, product, ActionUpdateProduct, Environment{}))

	fmt.Println("=== Stsenariy 4: customer mahsulotni tahrirlamoqchi ===")
	fmt.Printf("RBAC natija: %v\n", CanRBAC(customer, ActionUpdateProduct))
	fmt.Printf("ABAC natija: %v\n\n", engine.CanABAC(customer, product, ActionUpdateProduct, Environment{}))
}
