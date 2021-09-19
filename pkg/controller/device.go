package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Todorov99/server/pkg/models"
	"github.com/Todorov99/server/pkg/repository"
)

var device = models.Device{}
var deviceRepository = repository.CreateDeviceRepository()

type deviceController struct{}

func createDeviceController() IController {
	return &deviceController{}
}

func (d *deviceController) Get(w http.ResponseWriter, r *http.Request) {
	deviceID := getIDFromPathVariable(r)

	controllerLogger.Infof("Getting device with ID: %q", deviceID)

	devices, err := deviceRepository.GetByID(deviceID)
	respond(w, "", "Device GET query execution.", err, devices, http.StatusNotFound)
}

func (d *deviceController) Post(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&device)
	if err != nil {
		respond(w, "", "Device Post query", err, device, http.StatusInternalServerError)
		return
	}

	controllerLogger.Debugf("Post request with device: %q", device)
	err = deviceRepository.Add(device.Name, device.Description)
	respond(w, "You successfully add your device: ", "Device POST query execution.", err, device, http.StatusConflict)
}

func (d *deviceController) Put(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&device)
	if err != nil {
		respond(w, "", "Device Post query", err, device, http.StatusInternalServerError)
		return
	}

	deviceID := getIDFromPathVariable(r)
	controllerLogger.Debugf("Updating device with ID: %q", deviceID)

	err = deviceRepository.Update(device.Name, device.Description, deviceID)
	respond(w, "You successfully update your device: ", "Device PUT query execution.", err, device, http.StatusConflict)
}

func (d *deviceController) Delete(w http.ResponseWriter, r *http.Request) {
	deviceID := getIDFromPathVariable(r)

	controllerLogger.Debugf("Deleting device with ID: %q", deviceID)
	device, err := deviceRepository.Delete(deviceID)
	respond(w, "You successfully delete your device: ", "Device DELETE query execution.", err, device, http.StatusConflict)
}

func getAllDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := deviceRepository.GetAll()
	respond(w, "", "Device GET query execution.", err, devices, http.StatusNotFound)
}
