package controller

import (
	"calmwu.org/elbservice-operator/pkg/controller/elbservice"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, elbservice.Add)
}
