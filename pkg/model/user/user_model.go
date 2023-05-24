package user

func (u *User) GetRolesStrings() []string {
	role := make([]string, len(u.Roles))
	for i, v := range u.Roles {
		role[i] = v.Name
	}
	return role
}
