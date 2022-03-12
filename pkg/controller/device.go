package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Todorov99/server/pkg/models"
	"github.com/Todorov99/server/pkg/repository"
	"github.com/Todorov99/server/pkg/service"
)

var device = models.Device{}
var deviceRepository = repository.CreateDeviceRepository()

type deviceController struct {
	deviceService service.DeviceService
}

func createDeviceController() IController {
	return &deviceController{
		deviceService: service.NewDeviceService(),
	}
}

func (d *deviceController) GetAll(w http.ResponseWriter, r *http.Request) {

}

func (d *deviceController) Get(w http.ResponseWriter, r *http.Request) {
	deviceID := getIDFromPathVariable(r)
	controllerLogger.Infof("Getting device with ID: %q", deviceID)

	//devices, err := deviceRepository.GetByID(deviceID)
	devices, err := d.deviceService.GetDeviceById(deviceID)
	respond(w, "", "Device GET query execution.", err, devices, http.StatusNotFound)
}

func (d *deviceController) Post(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&device)
	if err != nil {
		respond(w, "", "Device Post query", err, device, http.StatusInternalServerError)
		return
	}

	controllerLogger.Debugf("Post request with device: %q", device)
	//	err = deviceRepository.Add(device.Name, device.Description)
	err = d.deviceService.AddDevice(device)
	respond(w, "You successfully add your device: ", "Device POST query execution.", err, device, http.StatusConflict)
}

func (d *deviceController) Put(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&device)
	if err != nil {
		respond(w, "", "Device Post query", err, device, http.StatusInternalServerError)
		return
	}

	device.ID = getIDFromPathVariable(r)
	controllerLogger.Debugf("Updating device with ID: %q", device.ID)

	//device.ID = deviceID

	//err = deviceRepository.Update(device.Name, device.Description, deviceID)
	err = d.deviceService.UpdateDevice(device)
	respond(w, "You successfully update your device: ", "Device PUT query execution.", err, device, http.StatusConflict)
}

func (d *deviceController) Delete(w http.ResponseWriter, r *http.Request) {
	deviceID := getIDFromPathVariable(r)

	controllerLogger.Debugf("Deleting device with ID: %q", deviceID)
	//device, err := deviceRepository.Delete(deviceID)
	device, err := d.deviceService.DeleteDevice(deviceID)
	respond(w, "You successfully delete your device: ", "Device DELETE query execution.", err, device, http.StatusConflict)
}

func getAllDevices(w http.ResponseWriter, r *http.Request) {
	//devices, err := deviceRepository.GetAll()
	devices, err := deviceRepository.GetAll()
	respond(w, "", "Device GET query execution.", err, devices, http.StatusNotFound)
}
