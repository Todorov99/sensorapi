package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Todorov99/server/pkg/dto"
	"github.com/Todorov99/server/pkg/service"
)

type deviceController struct {
	deviceService service.IService
}

func NewDeviceController() IController {
	return &deviceController{
		deviceService: service.NewDeviceService(),
	}
}

func (d *deviceController) GetAll(w http.ResponseWriter, r *http.Request) {
	devices, err := d.deviceService.GetAll(r.Context())
	response(w, "Device GET query execution.", err, devices, http.StatusNotFound)
}

func (d *deviceController) Get(w http.ResponseWriter, r *http.Request) {
	deviceID := getIDFromPathVariable(r)
	controllerLogger.Infof("Getting device with ID: %q", deviceID)

	devices, err := d.deviceService.GetById(r.Context(), deviceID)
	response(w, "Device GET query execution.", err, devices, http.StatusNotFound)
}

func (d *deviceController) Post(w http.ResponseWriter, r *http.Request) {
	addDeviceDto := dto.AddUpdateDeviceDto{}
	err := json.NewDecoder(r.Body).Decode(&addDeviceDto)
	if err != nil {
		response(w, "Device Post query", err, addDeviceDto, http.StatusInternalServerError)
		return
	}

	controllerLogger.Debugf("Post request with device: %q", addDeviceDto)
	err = d.deviceService.Add(r.Context(), addDeviceDto)
	response(w, "Device POST query execution...", err, addDeviceDto, http.StatusConflict)
}

func (d *deviceController) Put(w http.ResponseWriter, r *http.Request) {
	updateDeviceDto := dto.AddUpdateDeviceDto{}
	err := json.NewDecoder(r.Body).Decode(&updateDeviceDto)
	if err != nil {
		response(w, "Device Post query", err, updateDeviceDto, http.StatusInternalServerError)
		return
	}

	updateDeviceDto.ID = getIDFromPathVariable(r)
	controllerLogger.Debugf("Updating device with ID: %q", updateDeviceDto.ID)

	err = d.deviceService.Update(r.Context(), updateDeviceDto)
	response(w, "Device PUT query execution.", err, updateDeviceDto, http.StatusConflict)
}

func (d *deviceController) Delete(w http.ResponseWriter, r *http.Request) {
	deviceID := getIDFromPathVariable(r)
	controllerLogger.Debugf("Deleting device with ID: %q", deviceID)

	device, err := d.deviceService.Delete(r.Context(), deviceID)
	response(w, "Device DELETE query execution.", err, device, http.StatusConflict)
}
