package diff

import "github.com/apparentlymart/go-cty/cty"

// Change is an abstract type representing a single change operation as
// part of a Diff.
//
// Change is a closed interface, meaning that the only permitted
// implementations are those within this package.
type Change interface {
	changeSigil() changeImpl
}

// Embed changeImpl into a struct to make it a Change implementation
type changeImpl struct{}

func (c changeImpl) changeSigil() changeImpl {
	return c
}

// ReplaceChange is a Change implementation that represents replacing an
// existing value with an entirely new value.
//
// When adding a new element to a map value, this change type should be used
// with OldValue set to a null value of the appropriate type.
type ReplaceChange struct {
	changeImpl
	Path     Path
	OldValue cty.Value
	NewValue cty.Value
}

// DeleteChange is a Change implementation that represents removing an
// element from an indexable collection.
//
// For a list type, if the deleted element is not the final element in
// the list then the resulting "gap" is closed by renumbering all subsequent
// items. Therefore a Diff containing a sequence of DeleteChange operations
// on the same list must be careful to consider the new state of the element
// indices after each step, or present the deletions in reverse order to
// avoid such complexity.
type DeleteChange struct {
	changeImpl
	Path     Path
	OldValue cty.Value
}

// InsertChange is a Change implementation that represents inserting a new
// element into a list.
//
// When appending to a list, the Path should be to the not-yet-existing index
// and BeforeValue should be a null of the appropriate type.
type InsertChange struct {
	changeImpl
	Path        Path
	NewValue    cty.Value
	BeforeValue cty.Value
}

// AddChange is a Change implementation that represents adding a value to
// a set. The Path is to the set itself, and NewValue is the value to insert.
type AddChange struct {
	changeImpl
	Path     Path
	NewValue cty.Value
}

// RemoveChange is a Change implementation that represents removing a value
// from a set. The path is to the set itself, and OldValue is the value to
// remove.
//
// Note that it is not possible to remove an unknown value from a set
// because no two unknown values are equal, so a diff whose source value
// had sets with unknown members cannot be applied and is useful only
// for presentation to a user. Generally-speaking one should avoid including
// unknowns in the source value when creating a diff.
type RemoveChange struct {
	changeImpl
	Path     Path
	OldValue cty.Value
}

// Context is a funny sort of Change implementation that doesn't actually
// change anything but fails if the value at the given path doesn't match
// the given value.
//
// This can be used to add additional context to a diff so that merge
// conflicts can be detected.
type Context struct {
	changeImpl
	Path      Path
	WantValue cty.Value
}
