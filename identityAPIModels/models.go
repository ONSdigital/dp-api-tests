package identityAPIModels

type API struct {
	ID                string `bson:"id" json:"id"`
	Name              string `bson:"name" json:"name"`
	Email             string `bson:"email" json:"email"`
	Password          string `bson:"password" json:"password"`
	UserType          string `bson:"user_type" json:"user_type"`
	TemporaryPassword bool   `bson:"temporary_password" json:"temporary_password"`
	Migrated          bool   `bson:"migrated" json:"migrated"`
	Deleted           bool   `bson:"deleted" json:"deleted"`
}

type Mongo struct {
	ID                string `bson:"id" json:"id"`
	Name              string `bson:"name" json:"name"`
	Email             string `bson:"email" json:"email"`
	Password          string `bson:"password" json:"password"`
	UserType          string `bson:"user_type" json:"user_type"`
	TemporaryPassword bool   `bson:"temporary_password" json:"temporary_password"`
	Migrated          bool   `bson:"migrated" json:"migrated"`
	Deleted           bool   `bson:"deleted" json:"deleted"`
}

type NewTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (m Mongo) HashedPassword() []byte {
	return []byte(m.Password)
}
