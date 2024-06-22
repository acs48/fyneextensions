package fyneextensions

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

/*
The Actionable interface has two methods:

- GetActions() *ActionItem
- GetCanvas() Fyne.Canvas

These are implemented by types that have actions and operate on a Fyne canvas.
*/
type Actionable interface {
	GetActions() *ActionItem
	GetCanvas() fyne.Canvas
}

/*
ActionItem is a struct that represents an action. An action typically has a name, a state
resources associated with it (icons), an action function to trigger and some binding states.
Instead of a function to trigger, it could have sub-actions to create complex "menu" systems

The fields are as follows:
- Name is a binding.String which provides a way to display a name for an action and observe changes to this name.
- CriticalName is a bool that determines if the name is critical and should always be rendered.
- AlwaysShowAsContainer is a bool that defines if the action should always be represented as a container even if there are no sub-actions.
- Resources is a slice of fyne.Resource which can be used for representing the action in UI, like using an icon. The state of the item will force the related resource to be shown
- Triggered is a function that will be invoked when the action is triggered.
- SubActions are nested actions.
- HasDynamicStates is a bool that defines if the action has dynamic states that can change.
- Disabler, Hider, and Stater are binding variables which provide a way to control the disabled, hidden, and state properties of an action and observe changes to these properties.

The `Actionable` interface and its related items aim to provide a structured way to represent and manipulate actions in a Fyne application.
*/
type ActionItem struct {
	Name                  binding.String
	CriticalName          bool
	AlwaysShowAsContainer bool
	Resources             []fyne.Resource
	Triggered             func(int)
	SubActions            []*ActionItem
	HasDynamicStates      bool

	Disabler binding.Bool
	Hider    binding.Bool
	Stater   binding.Int
}

// NewActionItem function is a factory function for creating new action items.
func NewActionItem(name string, isNameCritical, alwaysShowAsContainer bool, resources []fyne.Resource, disabled bool, hidden bool, dynamicStates bool, state int, triggered func(int), subActions []*ActionItem) (action *ActionItem) {
	action = &ActionItem{
		Name:                  binding.NewString(),
		CriticalName:          isNameCritical,
		AlwaysShowAsContainer: alwaysShowAsContainer,
		Resources:             resources,
		Triggered:             triggered,
		SubActions:            subActions,
		Disabler:              binding.NewBool(),
		Hider:                 binding.NewBool(),
		HasDynamicStates:      dynamicStates,
	}
	action.Name.Set(name)
	action.Disabler.Set(disabled)
	action.Hider.Set(hidden)
	if dynamicStates {
		action.Stater = binding.NewInt()
		action.Stater.Set(state)
	}
	return
}

// AppendActions method is to append sub actions to an existing ActionItem.
func (ai *ActionItem) AppendActions(subActions ...*ActionItem) {
	ai.SubActions = append(ai.SubActions, subActions...)
}
