package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Todorov99/sensorapi/pkg/dto"
	"github.com/Todorov99/sensorapi/pkg/server/config"
	"github.com/Todorov99/sensorapi/pkg/service"
	"github.com/dgrijalva/jwt-go"
)

type DeviceContoller interface {
	GenerateDeviceCfg(w http.ResponseWriter, r *http.Request, token *jwt.Token)
	IController
}

type deviceController struct {
	deviceService service.DeviceService
}

func NewDeviceController() DeviceContoller {
	return &deviceController{
		deviceService: service.NewDeviceService(),
	}
}

func (d *deviceController) GetAll(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	defer func() {
		r.Body.Close()
	}()

	devices, err := d.deviceService.GetAll(r.Context(), config.GetJWTUserIDClaim(token))
	response(w, "Device GET query execution.", err, devices, http.StatusNotFound)
}

func (d *deviceController) GetByID(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	defer func() {
		r.Body.Close()
	}()

	deviceID := getIDFromPathVariable(r)
	controllerLogger.Infof("Getting device with ID: %q", deviceID)

	devices, err := d.deviceService.GetById(r.Context(), deviceID, config.GetJWTUserIDClaim(token))
	response(w, "Device GET query execution.", err, devices, http.StatusNotFound)
}

func (d *deviceController) Post(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	defer func() {
		r.Body.Close()
	}()

	addDeviceDto := dto.AddUpdateDeviceDto{}
	err := json.NewDecoder(r.Body).Decode(&addDeviceDto)
	if err != nil {
		response(w, "Device Post query", err, addDeviceDto, http.StatusInternalServerError)
		return
	}

	controllerLogger.Debugf("Post request with device: %q", addDeviceDto)

	err = d.deviceService.Add(r.Context(), addDeviceDto, config.GetJWTUserIDClaim(token))
	response(w, "Device POST query execution...", err, addDeviceDto, http.StatusConflict)
}

func (d *deviceController) Put(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	defer func() {
		r.Body.Close()
	}()

	updateDeviceDto := dto.AddUpdateDeviceDto{}
	err := json.NewDecoder(r.Body).Decode(&updateDeviceDto)
	if err != nil {
		response(w, "Device Post query", err, updateDeviceDto, http.StatusInternalServerError)
		return
	}

	updateDeviceDto.ID = getIDFromPathVariable(r)
	controllerLogger.Debugf("Updating device with ID: %q", updateDeviceDto.ID)

	err = d.deviceService.Update(r.Context(), updateDeviceDto, config.GetJWTUserIDClaim(token))
	response(w, "Device PUT query execution.", err, updateDeviceDto, http.StatusConflict)
}

func (d *deviceController) Delete(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	defer func() {
		r.Body.Close()
	}()

	deviceID := getIDFromPathVariable(r)
	controllerLogger.Debugf("Deleting device with ID: %q", deviceID)

	device, err := d.deviceService.Delete(r.Context(), deviceID, config.GetJWTUserIDClaim(token))
	response(w, "Device DELETE query execution.", err, device, http.StatusConflict)
}

func (d *deviceController) GenerateDeviceCfg(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	defer func() {
		r.Body.Close()
	}()

	deviceID := getIDFromPathVariable(r)

	filename, err := d.deviceService.GenerateDeviceCfg(r.Context(), deviceID, config.GetJWTUserIDClaim(token))
	if err != nil {
		response(w, "Failed generating device cfg", err, nil, http.StatusBadRequest)
		return
	}

	serverFile(w, r, filename)
}
