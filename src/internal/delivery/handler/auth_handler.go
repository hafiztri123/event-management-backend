package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/hafiztri123/src/internal/model"
	errs "github.com/hafiztri123/src/internal/pkg/error"
	"github.com/hafiztri123/src/internal/pkg/response"
	"github.com/hafiztri123/src/internal/service"
)

// AuthHandler defines the interface for authentication-related HTTP handlers.
type AuthHandler interface {
    Register(w http.ResponseWriter, r *http.Request)
    Login(w http.ResponseWriter, r *http.Request)
}

// authHandler implements the AuthHandler interface.
type authHandler struct {
    authService service.AuthService
    validator   *validator.Validate
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewAuthHandler(authService service.AuthService) AuthHandler {
    return &authHandler{
        authService: authService,
        validator:   validator.New(),
    }
}

// Register godoc
// @Summary      Register new user
// @Description  Register a new user with email, password, and full name
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body service.RegisterInput true "Registration Details"
// @Success      201  {object}  response.Response
// @Failure      400  {object}  response.Response{message=string} "Invalid input"
// @Failure      409  {object}  response.Response{message=string} "User already exists"
// @Failure      500  {object}  response.Response{message=string} "Server error"
// @Router       /auth/register [post]
func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
    var input model.RegisterInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        HandleErrorResponse(w, err)
        return
    }

    if err := h.validator.Struct(input); err != nil {
        HandleErrorResponse(w, errs.NewValidationError("Request is not valid"))
        return
    }

    err := h.authService.Register(&input)
    if err != nil {
        HandleErrorResponse(w, err)
        return
    }

    respondWithJSON(w, http.StatusCreated, response.Response{
        Timestamp: time.Now(),
    })
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body service.LoginInput true "Login Credentials"
// @Success      200  {object}  response.Response{data=service.LoginResponse}
// @Failure      400  {object}  response.Response{message=string} "Invalid input"
// @Failure      401  {object}  response.Response{message=string} "Invalid credentials"
// @Failure      500  {object}  response.Response{message=string} "Server error"
// @Router       /auth/login [post]
func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
    var input model.LoginInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        HandleErrorResponse(w, err)
        return
    }

    loginResponse, err := h.authService.Login(&input)
    if err != nil {
        HandleErrorResponse(w, err)
        return
    }

    respondWithJSON(w, http.StatusOK, response.Response{
        Timestamp: time.Now(),
        Data:      loginResponse,
    })
}

