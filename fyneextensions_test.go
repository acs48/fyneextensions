package fyneextensions

import (
	"fmt"
	"fyne.io/fyne/v2/app"
	"testing"
)

// Sample test function (replace "FunctionName" with an actual function you want to test)
func TestFyneextensions(t *testing.T) {
	// Test code here
}

// Example function demonstrating a basic Fyne window creation
func ExampleFyneextensions() {
	// Instantiate the Fyne application
	a := app.New()

	// Create a new window
	w := a.NewWindow("Fyne Window")

	// Set the main window content
	// TODO: Add more elements to the window as needed for the demonstration

	// Show and run the application
	w.ShowAndRun()

	// The application run in a different goroutine and could not be represented in console output.
	// Let's print something instead
	fmt.Println("A new Fyne window has been created and shown.")

	// Output: A new Fyne window has been created and shown.
}
