package auth

import (
	"log"
	"time"

	"github.com/TerrexTech/go-authserver-query/model"
	"github.com/pkg/errors"
)

func Login(db DBI, ts TokenStoreI, user *model.User) (*model.AuthResponse, error) {
	user, err := db.Login(user)
	if err != nil {
		err = errors.Wrap(err, "Login error")
		return nil, err
	}

	// Access Token
	accessExp := 15 * time.Minute
	claims := &model.Claims{
		Role: user.Role,
		Sub:  user.UUID,
	}
	accessToken, err := model.NewAccessToken(accessExp, claims)
	if err != nil {
		err = errors.Wrap(err, "Login Error: Error generating Access-Token")
		log.Println(err)
		return nil, err
	}

	// Refresh Token
	refreshExp := (24 * 7) * time.Hour
	refreshToken, err := model.NewRefreshToken(refreshExp, user.UUID)
	if err != nil {
		err = errors.Wrap(err, "Error generating Refresh-Token")
		log.Println(err)
		return nil, err
	}
	err = ts.Set(refreshToken)
	// We continue executing the code even if storing refresh-token fails since other parts
	// of application might still be accessible.
	if err != nil {
		err = errors.Wrapf(
			err,
			"Error storing RefreshToken in TokenStorage for UserID: %s", user.UUID,
		)
		log.Println(err)
	}

	userData := map[string]interface{}{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	}
	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserData:     userData,
	}, nil
}
