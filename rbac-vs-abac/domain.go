package main

import "time"

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleCustomer Role = "customer"
	RoleSeller   Role = "seller"
)

type Action string

const (
	ActionCreateProduct Action = "create_product"
	ActionUpdateProduct Action = "update_product"
	ActionDeleteProduct Action = "delete_product"
	ActionViewProduct   Action = "view_product"
)

type User struct {
	ID   string
	Role Role
	Name string
}

type Product struct {
	ID       string
	OwnerID  string
	Name     string
	Category string
}

// Environment holds the current environment state for abac
type Environment struct {
	CurrentTime time.Time
	IPAddress   string
}
