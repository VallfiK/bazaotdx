package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/VallfIK/bazaotdx/internal/models"
	"github.com/VallfIK/bazaotdx/internal/service"
)

// CottagesTab represents the tab for managing cottages
func NewCottagesTab(window fyne.Window, cottageService *service.CottageService, bookingService *service.BookingService) fyne.CanvasObject {
	var cottagesList *widget.List
	var refreshBtn *widget.Button

	// Update cottages list
	updateCottagesList := func() {
		cottages, err := cottageService.GetAllCottages()
		if err != nil {
			dialog.ShowError(fmt.Errorf("Failed to load cottages: %v", err), window)
			return
		}

		// Create a new list with the correct number of items
		cottagesList = widget.NewList(
			func() int { return len(cottages) },
			func() fyne.CanvasObject {
				return container.NewHBox(
					widget.NewLabel("Loading..."),
					widget.NewButton("Edit", nil),
					widget.NewButton("Delete", nil),
				)
			},
			func(id widget.ListItemID, item fyne.CanvasObject) {
				cottage := cottages[id]
				hbox := item.(*fyne.Container)
				label := hbox.Objects[0].(*widget.Label)
				editBtn := hbox.Objects[1].(*widget.Button)
				deleteBtn := hbox.Objects[2].(*widget.Button)

				label.SetText(fmt.Sprintf("%d. %s", cottage.ID, cottage.Name))
				editBtn.OnTapped = func() {
					editCottageDialog(window, cottageService, cottage, cottagesList)
				}
				deleteBtn.OnTapped = func() {
					deleteCottageDialog(window, cottageService, cottage, cottagesList)
				}
			},
		)
	}

	// Add new cottage
	addCottageBtn := widget.NewButton("Add New Cottage", func() {
		editCottageDialog(window, cottageService, models.Cottage{}, cottagesList)
	})

	// Refresh button
	refreshBtn = widget.NewButton("Refresh", updateCottagesList)

	// Create initial empty list
	cottagesList = widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject { return container.NewHBox() },
		func(_ widget.ListItemID, _ fyne.CanvasObject) {},
	)

	// Initial update
	updateCottagesList()

	// Create main container
	mainContainer := container.NewVBox(
		container.NewHBox(addCottageBtn, refreshBtn),
		cottagesList,
	)

	return mainContainer
}

// editCottageDialog shows dialog for editing cottage
func editCottageDialog(window fyne.Window, cottageService *service.CottageService, cottage models.Cottage, cottagesList *widget.List) {
	nameEntry := widget.NewEntry()
	if cottage.Name != "" {
		nameEntry.SetText(cottage.Name)
	}

	saveBtn := widget.NewButton("Save", func() {
		name := nameEntry.Text
		if name == "" {
			dialog.ShowError(fmt.Errorf("Name cannot be empty"), window)
			return
		}

		if cottage.ID > 0 {
			// Update existing cottage
			if err := cottageService.UpdateCottageName(cottage.ID, name); err != nil {
				dialog.ShowError(fmt.Errorf("Failed to update cottage: %v", err), window)
				return
			}
		} else {
			// Add new cottage
			if err := cottageService.AddCottage(name); err != nil {
				dialog.ShowError(fmt.Errorf("Failed to add cottage: %v", err), window)
				return
			}
		}

		// Update list
		cottagesList.Refresh()
	})

	cancelBtn := widget.NewButton("Cancel", func() {
		window.Close()
	})

	content := container.NewVBox(
		widget.NewLabel("Cottage Name:"),
		nameEntry,
		container.NewHBox(saveBtn, cancelBtn),
	)

	if cottage.ID > 0 {
		dialog.ShowCustom("Edit Cottage", "Save", content, window)
	} else {
		dialog.ShowCustom("Add New Cottage", "Save", content, window)
	}
}

// deleteCottageDialog shows dialog for deleting cottage
func deleteCottageDialog(window fyne.Window, cottageService *service.CottageService, cottage models.Cottage, cottagesList *widget.List) {
	confirmDialog := dialog.NewConfirm(
		"Delete Cottage",
		"Are you sure you want to delete this cottage?",
		func(confirmed bool) {
			if confirmed {
				if err := cottageService.DeleteCottage(cottage.ID); err != nil {
					dialog.ShowError(fmt.Errorf("Failed to delete cottage: %v", err), window)
					return
				}
				cottagesList.Refresh()
			}
		},
		window,
	)
	confirmDialog.Show()
}
