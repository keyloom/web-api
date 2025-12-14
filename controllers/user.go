package controllers

import "github.com/keyloom/web-api/core"

type UserController struct{}

var _ core.Controller = (*UserController)(nil)

func (uc *UserController) RegisterRoutes() {
	// Route registration logic goes here
}
