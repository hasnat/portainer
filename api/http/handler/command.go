package handler

import (
	"github.com/portainer/portainer"
	httperror "github.com/portainer/portainer/http/error"
	"github.com/portainer/portainer/http/proxy"
	"github.com/portainer/portainer/http/security"

	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

// CommandHandler represents an HTTP API handler for managing Docker commands.
type CommandHandler struct {
	*mux.Router
	Logger                      *log.Logger
	authorizeCommandManagement bool
	CommandService             portainer.CommandService
	FileService                 portainer.FileService
	ProxyManager                *proxy.Manager
}

const (
	// ErrCommandManagementDisabled is an error raised when trying to access the commands management commands
	// when the server has been started with the --external-commands flag
	ErrCommandManagementDisabled = portainer.Error("Command management is disabled")
)

// NewCommandHandler returns a new instance of CommandHandler.
func NewCommandHandler(bouncer *security.RequestBouncer, authorizeCommandManagement bool) *CommandHandler {
	h := &CommandHandler{
		Router: mux.NewRouter(),
		Logger: log.New(os.Stderr, "", log.LstdFlags),
		authorizeCommandManagement: authorizeCommandManagement,
	}
	h.Handle("/commands",
		bouncer.AdministratorAccess(http.HandlerFunc(h.handlePostCommands))).Methods(http.MethodPost)
	h.Handle("/commands",
		bouncer.RestrictedAccess(http.HandlerFunc(h.handleGetCommands))).Methods(http.MethodGet)
	h.Handle("/commands/{id}",
		bouncer.AdministratorAccess(http.HandlerFunc(h.handleGetCommand))).Methods(http.MethodGet)
	h.Handle("/commands/{id}",
		bouncer.AdministratorAccess(http.HandlerFunc(h.handlePutCommand))).Methods(http.MethodPut)
	h.Handle("/commands/{id}/access",
		bouncer.AdministratorAccess(http.HandlerFunc(h.handlePutCommandAccess))).Methods(http.MethodPut)
	h.Handle("/commands/{id}",
		bouncer.AdministratorAccess(http.HandlerFunc(h.handleDeleteCommand))).Methods(http.MethodDelete)

	return h
}

type (
	postCommandsRequest struct {
		Name                string `valid:"required"`
		Image                 string `valid:"required"`
		Command           string `valid:"-"`
	}

	postCommandsResponse struct {
		ID int `json:"Id"`
	}

	putCommandAccessRequest struct {
		AuthorizedUsers []int `valid:"-"`
		AuthorizedTeams []int `valid:"-"`
	}

	putCommandsRequest struct {
		Name                string `valid:"-"`
		Image                 string `valid:"-"`
		Command           string `valid:"-"`
	}
)

// handleGetCommands handles GET requests on /commands
func (handler *CommandHandler) handleGetCommands(w http.ResponseWriter, r *http.Request) {

	commands, err := handler.CommandService.Commands()
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	filteredCommands := commands

	encodeJSON(w, filteredCommands, handler.Logger)
}

// handlePostCommands handles POST requests on /commands
func (handler *CommandHandler) handlePostCommands(w http.ResponseWriter, r *http.Request) {
	if !handler.authorizeCommandManagement {
		httperror.WriteErrorResponse(w, ErrCommandManagementDisabled, http.StatusServiceUnavailable, handler.Logger)
		return
	}

	var req postCommandsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.WriteErrorResponse(w, ErrInvalidJSON, http.StatusBadRequest, handler.Logger)
		return
	}

	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		httperror.WriteErrorResponse(w, ErrInvalidRequestFormat, http.StatusBadRequest, handler.Logger)
		return
	}

	command := &portainer.Command{
		Name:      req.Name,
		Image:       req.Image,
		Command: req.Command,
	}

	err = handler.CommandService.CreateCommand(command)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}


	encodeJSON(w, &postCommandsResponse{ID: int(command.ID)}, handler.Logger)
}

// handleGetCommand handles GET requests on /commands/:id
func (handler *CommandHandler) handleGetCommand(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	commandID, err := strconv.Atoi(id)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusBadRequest, handler.Logger)
		return
	}

	command, err := handler.CommandService.Command(portainer.CommandID(commandID))
	if err == portainer.ErrCommandNotFound {
		httperror.WriteErrorResponse(w, err, http.StatusNotFound, handler.Logger)
		return
	} else if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	encodeJSON(w, command, handler.Logger)
}

// handlePutCommandAccess handles PUT requests on /commands/:id/access
func (handler *CommandHandler) handlePutCommandAccess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	commandID, err := strconv.Atoi(id)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusBadRequest, handler.Logger)
		return
	}

	var req putCommandAccessRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.WriteErrorResponse(w, ErrInvalidJSON, http.StatusBadRequest, handler.Logger)
		return
	}

	_, err = govalidator.ValidateStruct(req)
	if err != nil {
		httperror.WriteErrorResponse(w, ErrInvalidRequestFormat, http.StatusBadRequest, handler.Logger)
		return
	}

	command, err := handler.CommandService.Command(portainer.CommandID(commandID))
	if err == portainer.ErrCommandNotFound {
		httperror.WriteErrorResponse(w, err, http.StatusNotFound, handler.Logger)
		return
	} else if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}


	err = handler.CommandService.UpdateCommand(command.ID, command)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}
}

// handlePutCommand handles PUT requests on /commands/:id
func (handler *CommandHandler) handlePutCommand(w http.ResponseWriter, r *http.Request) {
	if !handler.authorizeCommandManagement {
		httperror.WriteErrorResponse(w, ErrCommandManagementDisabled, http.StatusServiceUnavailable, handler.Logger)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	commandID, err := strconv.Atoi(id)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusBadRequest, handler.Logger)
		return
	}

	var req putCommandsRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.WriteErrorResponse(w, ErrInvalidJSON, http.StatusBadRequest, handler.Logger)
		return
	}

	_, err = govalidator.ValidateStruct(req)
	if err != nil {
		httperror.WriteErrorResponse(w, ErrInvalidRequestFormat, http.StatusBadRequest, handler.Logger)
		return
	}

	command, err := handler.CommandService.Command(portainer.CommandID(commandID))
	if err == portainer.ErrCommandNotFound {
		httperror.WriteErrorResponse(w, err, http.StatusNotFound, handler.Logger)
		return
	} else if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	if req.Name != "" {
		command.Name = req.Name
	}

	if req.Image != "" {
		command.Image = req.Image
	}

	if req.Command != "" {
		command.Command = req.Command
	}




	err = handler.CommandService.UpdateCommand(command.ID, command)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}
}

// handleDeleteCommand handles DELETE requests on /commands/:id
func (handler *CommandHandler) handleDeleteCommand(w http.ResponseWriter, r *http.Request) {
	if !handler.authorizeCommandManagement {
		httperror.WriteErrorResponse(w, ErrCommandManagementDisabled, http.StatusServiceUnavailable, handler.Logger)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	commandID, err := strconv.Atoi(id)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusBadRequest, handler.Logger)
		return
	}



	if err == portainer.ErrCommandNotFound {
		httperror.WriteErrorResponse(w, err, http.StatusNotFound, handler.Logger)
		return
	} else if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	handler.ProxyManager.DeleteProxy(string(commandID))
	handler.ProxyManager.DeleteExtensionProxies(string(commandID))

	err = handler.CommandService.DeleteCommand(portainer.CommandID(commandID))
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}


}
