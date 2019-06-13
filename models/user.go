package models

type User struct {
	UserId    string `bson:"_id" json:"userId"`
	UserMoney int64  `bson:"user_money" json:"userMoney"`
	UserName  string `bson:"user_name" json:"userName"`
	UserPhone string `bson:"user_phone" json:"userPhone"`
}

type UserDTO struct {
	UserId    string `bson:"_id" json:"userId"`
	UserMoney int64  `bson:"user_money" json:"userMoney"`
	UserName  string `bson:"user_name" json:"userName"`
}

type TransferDTO struct {
	FromUserId string `json:"formUserId"`
	ToUserId   string `json:"toUserId"`
	Money      int64  `json:"money"`
}

type RegisteredDTO struct {
	UserPhone string `json:"userPhone"`
	SmsCode   string `json:"smsCode"`
}

const UserDb = "user"

func User2UserDto(model User) *UserDTO {

	userDto := new(UserDTO)
	userDto.UserId = model.UserId
	userDto.UserMoney = model.UserMoney
	userDto.UserName = model.UserName
	return userDto
}
