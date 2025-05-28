package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/VallfIK/bazaotdx/internal/models"
)

// CottageListItem represents a single cottage in the list
func NewCottageListItem(cottage models.Cottage, onEdit func(), onDelete func()) fyne.CanvasObject {
	return container.NewHBox(
		widget.NewLabel(fmt.Sprintf("%d. %s", cottage.ID, cottage.Name)),
		widget.NewButton("Edit", onEdit),
		widget.NewButton("Delete", onDelete),
	)
}
