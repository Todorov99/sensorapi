package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Todorov99/server/pkg/dto"
	"github.com/Todorov99/server/pkg/service"
)

var device = dto.Device{}

type deviceController struct {
	deviceService service.IService
}

func NewDeviceController() IController {
	return &deviceController{
		deviceService: service.NewDeviceService(),
	}
}

func (d *deviceController) GetAll(w http.ResponseWriter, r *http.Request) {
	devices, err := d.deviceService.GetAll()
	response(w, "", "Device GET query execution.", err, devices, http.StatusNotFound)
}

func (d *deviceController) Get(w http.ResponseWriter, r *http.Request) {
	deviceID := getIDFromPathVariable(r)
	controllerLogger.Infof("Getting device with ID: %q", deviceID)

	devices, err := d.deviceService.GetById(deviceID)
	response(w, "", "Device GET query execution.", err, devices, http.StatusNotFound)
}

func (d *deviceController) Post(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&device)
	if err != nil {
		response(w, "", "Device Post query", err, device, http.StatusInternalServerError)
		return
	}

	controllerLogger.Debugf("Post request with device: %q", device)
	err = d.deviceService.Add(device)
	response(w, "You successfully add your device: ", "Device POST query execution.", err, device, http.StatusConflict)
}

func (d *deviceController) Put(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&device)
	if err != nil {
		response(w, "", "Device Post query", err, device, http.StatusInternalServerError)
		return
	}

	device.ID = getIDFromPathVariable(r)
	controllerLogger.Debugf("Updating device with ID: %q", device.ID)

	err = d.deviceService.Update(device)
	response(w, "You successfully update your device: ", "Device PUT query execution.", err, device, http.StatusConflict)
}

func (d *deviceController) Delete(w http.ResponseWriter, r *http.Request) {
	deviceID := getIDFromPathVariable(r)
	controllerLogger.Debugf("Deleting device with ID: %q", deviceID)

	device, err := d.deviceService.Delete(deviceID)
	response(w, "You successfully delete your device: ", "Device DELETE query execution.", err, device, http.StatusConflict)
}
