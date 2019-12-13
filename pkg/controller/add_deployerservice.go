package controller

import (
	"github.com/dotmesh-io/dotscience-operator/pkg/controller/deployerservice"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, deployerservice.Add)
}
