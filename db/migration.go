package db

import "github.com/fabienbellanger/echo-boilerplate/entities"

// entitiesList lists all entities to automigrate.
var entitiesList = []interface{}{
	&entities.User{},
}
