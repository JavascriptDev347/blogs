package main

var rolePermissions = map[Role][]Action{
	RoleAdmin:    {ActionCreateProduct, ActionDeleteProduct, ActionUpdateProduct, ActionViewProduct},
	RoleCustomer: {ActionViewProduct},
	RoleSeller:   {ActionCreateProduct, ActionUpdateProduct, ActionViewProduct},
}

func CanRBAC(user User, action Action) bool {
	permissions, ok := rolePermissions[user.Role]
	if !ok {
		return false
	}
	for _, perm := range permissions {
		if perm == action {
			return true
		}
	}
	return false
}
