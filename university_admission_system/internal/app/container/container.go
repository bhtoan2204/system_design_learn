package container

import (
	"context"
	"sync"

	"university_admission_system/application/services"
	"university_admission_system/domain"
	"university_admission_system/infrastructure/id"
	"university_admission_system/infrastructure/repository/memory"
	"university_admission_system/infrastructure/seed"
	"university_admission_system/pkg/clock"
	"university_admission_system/pkg/config"
	"university_admission_system/pkg/logger"
	"university_admission_system/pkg/validator"
)

// Container wires dependencies following a simple service registry pattern.
type Container struct {
	cfg       *config.Config
	logger    logger.Logger
	validator validator.Validator
	clock     clock.Clock
	idGen     domain.IDGenerator

	applicantRepo   *memory.ApplicantRepository
	applicationRepo *memory.ApplicationRepository
	offerRepo       *memory.OfferRepository
	enrollmentRepo  *memory.EnrollmentRepository

	submitOnce sync.Once
	submitSvc  *services.SubmitApplicationService
	issueOnce  sync.Once
	issueSvc   *services.IssueOfferService
	acceptOnce sync.Once
	acceptSvc  *services.AcceptOfferService
	enrollOnce sync.Once
	enrollSvc  *services.ConfirmEnrollmentService
}

// New constructs a new Container instance.
func New(cfg *config.Config, log logger.Logger) *Container {
	return &Container{
		cfg:       cfg,
		logger:    log,
		validator: validator.New(),
		clock:     clock.SystemClock{},
		idGen:     id.RandomGenerator{},
	}
}

// Config exposes the loaded configuration.
func (c *Container) Config() *config.Config {
	return c.cfg
}

// Logger exposes the configured logger.
func (c *Container) Logger() logger.Logger {
	return c.logger
}

// Validator exposes the shared validator.
func (c *Container) Validator() validator.Validator {
	return c.validator
}

// Clock returns the clock implementation.
func (c *Container) Clock() clock.Clock {
	return c.clock
}

// ApplicantRepository provides the applicant repository.
func (c *Container) ApplicantRepository() domain.ApplicantRepository {
	return c.ensureApplicantRepo()
}

// ApplicationRepository provides the application repository.
func (c *Container) ApplicationRepository() domain.ApplicationRepository {
	return c.ensureApplicationRepo()
}

// OfferRepository provides the offer repository.
func (c *Container) OfferRepository() domain.OfferRepository {
	return c.ensureOfferRepo()
}

// EnrollmentRepository provides the enrollment repository.
func (c *Container) EnrollmentRepository() domain.EnrollmentRepository {
	return c.ensureEnrollmentRepo()
}

// SubmitApplicationService returns the orchestrator for submission flow.
func (c *Container) SubmitApplicationService() *services.SubmitApplicationService {
	c.submitOnce.Do(func() {
		c.submitSvc = services.NewSubmitApplicationService(
			c.ensureApplicantRepo(),
			c.ensureApplicationRepo(),
			c.idGen,
			c.clock,
			c.validator,
		)
	})
	return c.submitSvc
}

// IssueOfferService returns the orchestrator for issuing offers.
func (c *Container) IssueOfferService() *services.IssueOfferService {
	c.issueOnce.Do(func() {
		c.issueSvc = services.NewIssueOfferService(
			c.ensureApplicationRepo(),
			c.ensureApplicantRepo(),
			c.ensureOfferRepo(),
			c.idGen,
			domain.DefaultScoreCalculator{},
			c.clock,
			c.validator,
			c.cfg.MinimumScore,
		)
	})
	return c.issueSvc
}

// AcceptOfferService returns the orchestrator for offer acceptance.
func (c *Container) AcceptOfferService() *services.AcceptOfferService {
	c.acceptOnce.Do(func() {
		c.acceptSvc = services.NewAcceptOfferService(
			c.ensureOfferRepo(),
			c.clock,
			c.validator,
		)
	})
	return c.acceptSvc
}

// ConfirmEnrollmentService returns the orchestrator for enrollment confirmation.
func (c *Container) ConfirmEnrollmentService() *services.ConfirmEnrollmentService {
	c.enrollOnce.Do(func() {
		c.enrollSvc = services.NewConfirmEnrollmentService(
			c.ensureApplicationRepo(),
			c.ensureOfferRepo(),
			c.ensureEnrollmentRepo(),
			c.idGen,
			c.clock,
			c.validator,
		)
	})
	return c.enrollSvc
}

// SeedDemoData populates in-memory repositories for demo usage.
func (c *Container) SeedDemoData(ctx context.Context) error {
	if !c.cfg.SeedDemoData {
		return nil
	}
	if _, err := seed.SeedData(ctx, c.idGen, c.ensureApplicantRepo(), c.ensureApplicationRepo(), c.ensureOfferRepo(), c.ensureEnrollmentRepo()); err != nil {
		return err
	}
	if c.logger != nil {
		c.logger.Info("seed data loaded", nil)
	}
	return nil
}

func (c *Container) ensureApplicantRepo() *memory.ApplicantRepository {
	if c.applicantRepo == nil {
		c.applicantRepo = memory.NewApplicantRepository()
	}
	return c.applicantRepo
}

func (c *Container) ensureApplicationRepo() *memory.ApplicationRepository {
	if c.applicationRepo == nil {
		c.applicationRepo = memory.NewApplicationRepository()
	}
	return c.applicationRepo
}

func (c *Container) ensureOfferRepo() *memory.OfferRepository {
	if c.offerRepo == nil {
		c.offerRepo = memory.NewOfferRepository()
	}
	return c.offerRepo
}

func (c *Container) ensureEnrollmentRepo() *memory.EnrollmentRepository {
	if c.enrollmentRepo == nil {
		c.enrollmentRepo = memory.NewEnrollmentRepository()
	}
	return c.enrollmentRepo
}
