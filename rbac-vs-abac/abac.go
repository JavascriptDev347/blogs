package main

type Policy func(user User, product Product, action Action, env Environment) bool

func OwnershipPolicy(user User, product Product, action Action, env Environment) bool {
	if action != ActionUpdateProduct && action != ActionDeleteProduct {
		return true // bu policy faqat update/delete ga tegishli, boshqasiga aralashmaydi
	}
	if user.Role == RoleAdmin {
		return true // admin har doim o'tadi
	}
	return product.OwnerID == user.ID
}

type Engine struct {
	Policies []Policy
}

func (e Engine) CanABAC(user User, product Product, action Action, env Environment) bool {
	for _, policy := range e.Policies {
		if !policy(user, product, action, env) {
			return false
		}
	}
	return true
}
