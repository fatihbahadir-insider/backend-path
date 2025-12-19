package models

type Role uint

const (
	RoleUser  Role = iota + 1 
	RoleAdmin                            
	RoleMod                              
)

func (r Role) IsValid() bool {
	return r >= RoleUser && r <= RoleMod
}

func (r Role) String() string {
	names := map[Role]string{
		RoleUser:  "user",
		RoleAdmin: "admin",
		RoleMod:   "mod",
	}
	return names[r]
}

func ToRole(s string) Role {
	switch s {
	case "user":
		return RoleUser
	case "admin":
		return RoleAdmin
	case "mod":
		return RoleMod
	}
	return RoleUser
}