package listener

import (
	"net/http"

	"github.com/cloudfoundry/bosh-agent/bootstrapper/installer"
	"github.com/cloudfoundry/bosh-agent/internal/github.com/cloudfoundry/bosh-utils/logger"
)

const StatusUnprocessableEntity = 422

type SelfUpdateHandler struct {
	Logger    logger.Logger
	installer installer.Installer
}

func (h *SelfUpdateHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "PUT" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var err installer.Error

	err = h.installer.Install(req.Body)
	if err != nil {
		if err.SystemError() {
			rw.WriteHeader(http.StatusInternalServerError)
		} else {
			rw.WriteHeader(StatusUnprocessableEntity)
		}
		if _, wErr := rw.Write([]byte(err.Error())); wErr != nil {
			h.Logger.Warn("SelfUpdateHandler", "Failed to write error to buffer: %s", wErr.Error())
		}
		h.Logger.Error("SelfUpdateHandler", "failed to install package: %s", err.Error())
		return
	}

	h.Logger.Info("SelfUpdateHandler", "successfully installed package")
}
