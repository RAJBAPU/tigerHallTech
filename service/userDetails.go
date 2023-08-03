package service

import (
	models "simpl_pr/model"
)

func (tg *User) GetUserDetails(currentUser *models.TgUser) (userResponse *UserResponse, err error) {

	useDetails, err := tg.User.GetYpUserById(currentUser.Id)
	if err != nil {
		return
	}

	userResponse = &UserResponse{
		Id:    useDetails.Id,
		Name:  useDetails.Name,
		Email: useDetails.Email,
	}
	return

}
