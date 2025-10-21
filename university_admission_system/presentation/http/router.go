package presentationhttp

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"university_admission_system/application/services"
	"university_admission_system/docs"
	"university_admission_system/domain"
	appErrors "university_admission_system/pkg/errors"
	"university_admission_system/pkg/logger"
)

// RouterConfig aggregates dependencies required to expose HTTP endpoints.
type RouterConfig struct {
	SubmitService            *services.SubmitApplicationService
	IssueOfferService        *services.IssueOfferService
	AcceptOfferService       *services.AcceptOfferService
	ConfirmEnrollmentService *services.ConfirmEnrollmentService
	Logger                   logger.Logger
	EnableSwagger            bool
}

// NewRouter wires the HTTP routes with the provided services.
func NewRouter(cfg RouterConfig) http.Handler {
	h := &handler{
		submitService:            cfg.SubmitService,
		issueOfferService:        cfg.IssueOfferService,
		acceptOfferService:       cfg.AcceptOfferService,
		confirmEnrollmentService: cfg.ConfirmEnrollmentService,
		logger:                   cfg.Logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /applications", h.submitApplication)
	mux.HandleFunc("POST /applications/{applicationId}/issue-offer", h.issueOffer)
	mux.HandleFunc("POST /offers/{offerId}/accept", h.acceptOffer)
	mux.HandleFunc("POST /enrollments", h.confirmEnrollment)

	if cfg.EnableSwagger {
		mux.HandleFunc("GET /swagger.yaml", h.swaggerSpec)
		mux.HandleFunc("GET /swagger", h.swaggerUI)
	}

	return mux
}

type handler struct {
	submitService            *services.SubmitApplicationService
	issueOfferService        *services.IssueOfferService
	acceptOfferService       *services.AcceptOfferService
	confirmEnrollmentService *services.ConfirmEnrollmentService
	logger                   logger.Logger
}

type submitApplicationRequest struct {
	ApplicantID string `json:"applicantId"`
	ProgramID   string `json:"programId"`
}

type submitApplicationResponse struct {
	ApplicationID string    `json:"applicationId"`
	SubmittedAt   time.Time `json:"submittedAt"`
}

func (h *handler) submitApplication(w http.ResponseWriter, r *http.Request) {
	var req submitApplicationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	result, err := h.submitService.Submit(r.Context(), services.SubmitApplicationCommand{
		ApplicantID: req.ApplicantID,
		ProgramID:   req.ProgramID,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, submitApplicationResponse{
		ApplicationID: result.ApplicationID,
		SubmittedAt:   result.SubmittedAt,
	})
}

type issueOfferResponse struct {
	OfferID   string    `json:"offerId"`
	Score     float64   `json:"score"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func (h *handler) issueOffer(w http.ResponseWriter, r *http.Request) {
	applicationID := r.PathValue("applicationId")
	result, err := h.issueOfferService.Issue(r.Context(), services.IssueOfferCommand{
		ApplicationID: applicationID,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, issueOfferResponse{
		OfferID:   result.OfferID,
		Score:     result.Score,
		ExpiresAt: result.ExpiresAt,
	})
}

type acceptOfferResponse struct {
	AcceptedAt time.Time `json:"acceptedAt"`
}

func (h *handler) acceptOffer(w http.ResponseWriter, r *http.Request) {
	offerID := r.PathValue("offerId")
	result, err := h.acceptOfferService.Accept(r.Context(), services.AcceptOfferCommand{
		OfferID: offerID,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, acceptOfferResponse{
		AcceptedAt: result.AcceptedAt,
	})
}

type confirmEnrollmentRequest struct {
	ApplicationID string `json:"applicationId"`
	OfferID       string `json:"offerId"`
}

type confirmEnrollmentResponse struct {
	EnrollmentID string    `json:"enrollmentId"`
	ConfirmedAt  time.Time `json:"confirmedAt"`
}

func (h *handler) confirmEnrollment(w http.ResponseWriter, r *http.Request) {
	var req confirmEnrollmentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	result, err := h.confirmEnrollmentService.Confirm(r.Context(), services.ConfirmEnrollmentCommand{
		ApplicationID: req.ApplicationID,
		OfferID:       req.OfferID,
	})
	if err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, confirmEnrollmentResponse{
		EnrollmentID: result.EnrollmentID,
		ConfirmedAt:  result.ConfirmedAt,
	})
}

func (h *handler) swaggerSpec(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/yaml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(docs.OpenAPI)
}

func (h *handler) swaggerUI(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(swaggerUIHTML))
}

func (h *handler) writeError(w http.ResponseWriter, err error) {
	status, message := mapError(err)
	if status >= http.StatusInternalServerError && h.logger != nil {
		h.logger.Error("request failed", err, nil)
	}
	writeJSON(w, status, map[string]string{"error": message})
}

func decodeJSON(r *http.Request, target interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		return errors.New("invalid JSON payload")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, appErrors.ErrInvalidInput):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, appErrors.ErrNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, appErrors.ErrConflict),
		errors.Is(err, domain.ErrApplicationAlreadySubmitted),
		errors.Is(err, domain.ErrApplicationAlreadyScored),
		errors.Is(err, domain.ErrApplicationAlreadyOffered),
		errors.Is(err, domain.ErrOfferAlreadyAccepted),
		errors.Is(err, domain.ErrOfferExpired),
		errors.Is(err, domain.ErrEnrollmentAlreadyConfirmed):
		return http.StatusConflict, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
