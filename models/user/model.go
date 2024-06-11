package user

type User struct {
	Id                   string `bson:"_id,omitempty" json:"id"`
	Name                 string `validate:"required" bson:"name" json:"name"`
	Email                string `validate:"required,email" bson:"email" json:"email"`
	Password             string `validate:"required,min=10,max=15" bson:"password" json:"password"`
	ConfirmationPassword string `validate:"required,eqfield=Password" bson:"-" json:"confirmation_password"`
	CPF                  string `validate:"required" bson:"cpf" json:"cpf"`
	Phone                string `validate:"required" bson:"phone" json:"phone"`
	Code                 string `bson:"code,omitempty" json:"-"`
}
