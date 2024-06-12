package main

import (
	"github.com/julienschmidt/httprouter"

	user_profile_api "github.com/envmission/template-api/delivery/user-profile-api"
)

type basicProjectRepository struct {
	// repository taro disini
}

func enableBasicModule(
	router *httprouter.Router,
) {
	// Buat Delivery handler disini.

	deliveryUserProfile := user_profile_api.New()

	// Project

	router.OPTIONS("/user", corsOptional(Empty))
	router.GET("/user", corsOptional(deliveryUserProfile.Get))
	router.PUT("/user", corsOptional(deliveryUserProfile.Put))
	router.DELETE("/user", corsOptional(deliveryUserProfile.Delete))

	// Sebenernya: selain corsOptional, ada proteksi Authorization pakai withAuth

	// Eg: orang dengan role "admin" "viewer" boleh nge liat 'GET'
	// tapi cuma "admin" yang bisa nambah user profile "PUT", dan delete "DELETE"

	// router.GET("/user", corsOptional(withAuth(authUC, []string{"admin", "viewer"}, deliveryUserProfile.Get)))
	// router.PUT("/user", corsOptional(withAuth(authUC, []string{"admin"}, deliveryUserProfile.Put)))
	// router.DELETE("/user", corsOptional(withAuth(authUC, []string{"admin"}, deliveryUserProfile.Delete)))

}
