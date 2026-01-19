package main

import (
	"net/http"
	"time"

	"github.com/mysecodgit/go_accounting/internal/service"
	"github.com/mysecodgit/go_accounting/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

type application struct {
	config  config
	store   store.Storage
	service service.Service
	logger  *zap.SugaredLogger
}

type config struct {
	addr        string
	db          dbConfig
	env         string
	apiURL      string
	frontendURL string
	auth        authConfig
}

type authConfig struct {
	basic basicConfig
}

type basicConfig struct {
	user string
	pass string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	//cors
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:7174"}, // frontend URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Max age in seconds
	}))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.checkHealthHandler)

		r.Route("/users", func(r chi.Router) {
			r.Get("/", app.getUsersHandler)
			r.Post("/", app.createUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				r.Get("/", app.getUserHandler)
				r.Put("/", app.updateUserHandler)
				r.Delete("/", app.deleteUserHandler)
			})
		})

		r.Route("/buildings", func(r chi.Router) {
			r.Get("/", app.getBuildingsHandler)
			r.Post("/", app.createBuildingHandler)
			r.Route("/{buildingID}", func(r chi.Router) {
				r.Get("/", app.getBuildingHandler)
				r.Put("/", app.updateBuildingHandler)
				r.Delete("/", app.deleteBuildingHandler)
				r.Get("/available-units", app.getAvailableUnitsByBuildingIDHandler)

				r.Route("/units", func(r chi.Router) {
					r.Get("/", app.getUnitsHandler)
					r.Post("/", app.createUnitHandler)
					r.Route("/{unitID}", func(r chi.Router) {
						r.Get("/", app.getUnitHandler)
						r.Get("/active_lease", app.getActiveLeaseByUnitIDHandler)
						r.Get("/readings", app.getReadingsByUnitHandler)
						r.Put("/", app.updateUnitHandler)
						r.Delete("/", app.deleteUnitHandler)
					})
				})

				// people
				r.Route("/people", func(r chi.Router) {
					r.Get("/", app.getPeopleHandler)
					r.Post("/", app.createPersonHandler)
					r.Route("/{peopleID}", func(r chi.Router) {
						r.Get("/", app.getPersonHandler)
						r.Get("/units", app.getUnitsByPeopleHandler)
						r.Get("/available-credits", app.getAvailableCreditsHandler)
						r.Put("/", app.updatePersonHandler)
						r.Delete("/", app.deletePersonHandler)
					})
				})

				// accounts
				r.Route("/accounts", func(r chi.Router) {
					r.Get("/", app.getAccountsHandler)
					r.Post("/", app.createAccountHandler)
					r.Route("/{accountID}", func(r chi.Router) {
						r.Get("/", app.getAccountHandler)
						r.Put("/", app.updateAccountHandler)
						r.Delete("/", app.deleteAccountHandler)
					})
				})

				// items
				r.Route("/items", func(r chi.Router) {
					r.Get("/", app.getItemsHandler)
					r.Post("/", app.createItemHandler)
					r.Route("/{itemID}", func(r chi.Router) {
						r.Get("/", app.getItemHandler)
						r.Put("/", app.updateItemHandler)
						r.Delete("/", app.deleteItemHandler)
					})
				})

				// invoices
				r.Route("/invoices", func(r chi.Router) {
					r.Get("/", app.getInvoicesHandler)
					r.Post("/preview", app.previewInvoiceSplitsHandler)
					r.Post("/", app.createInvoiceHandler)
					r.Route("/{invoiceID}", func(r chi.Router) {
						r.Get("/", app.getInvoiceHandler)
						r.Put("/", app.updateInvoiceHandler)
						r.Get("/payments", app.getPaymentsHandler)
						r.Get("/applied-discounts", app.getInvoiceDiscountsHandler)
						r.Post("/apply-discount", app.applyInvoiceDiscountHandler)
						r.Get("/available-credits", app.getInvoiceAvailableCreditsHandler)
						r.Get("/applied-credits", app.getInvoiceAppliedCreditsHandler)
						r.Post("/apply-credit", app.applyInvoiceCreditHandler)
					})
				})

				r.Route("/invoice-payments", func(r chi.Router) {
					r.Post("/", app.createInvoicePaymentHandler)
					r.Get("/", app.getInvoicePaymentsHandler)
					r.Route("/{invoicePaymentID}", func(r chi.Router) {
						r.Get("/", app.getInvoicePaymentHandler)
						r.Put("/", app.updateInvoicePaymentHandler)
					})

				})

				r.Route("/sales-receipts", func(r chi.Router) {
					r.Get("/", app.getSalesReceiptsHandler)
					r.Post("/", app.createSalesReceiptHandler)
					r.Route("/{salesReceiptID}", func(r chi.Router) {
						r.Get("/", app.getSalesReceiptHandler)
						r.Put("/", app.updateSalesReceiptHandler)
					})
				})

				r.Route("/leases", func(r chi.Router) {
					r.Get("/", app.getLeasesHandler)
					r.Post("/", app.createLeaseHandler)
					r.Route("/{leaseID}", func(r chi.Router) {
						r.Get("/", app.getLeaseHandler)
						r.Put("/", app.updateLeaseHandler)
					})
				})

				r.Route("/credit-memos", func(r chi.Router) {
					r.Get("/", app.getAllCreditMemoHandler)
					r.Post("/", app.createCreditMemoHandler)
					r.Route("/{creditMemoID}", func(r chi.Router) {
						r.Get("/", app.getCreditMemoHandler)
						r.Put("/", app.updateCreditMemoHandler)
					})
				})

				r.Route("/checks", func(r chi.Router) {
					r.Get("/", app.getChecksHandler)
					r.Post("/", app.createCheckHandler)
					r.Route("/{checkID}", func(r chi.Router) {
						r.Get("/", app.getCheckHandler)
						r.Put("/", app.updateCheckHandler)
						// r.Delete("/", app.deleteCheckHandler)
					})
				})

				r.Route("/journals", func(r chi.Router) {
					r.Get("/", app.getJournalsHandler)
					r.Post("/", app.createJournalHandler)
					r.Route("/{journalID}", func(r chi.Router) {
						r.Get("/", app.getJournalHandler)
						r.Put("/", app.updateJournalHandler)
					})
				})

				r.Route("/readings", func(r chi.Router) {
					r.Get("/", app.getReadingsHandler)
					r.Get("/latest", app.getLatestReadingHandler)
					r.Post("/", app.createReadingHandler)
					r.Route("/{readingID}", func(r chi.Router) {
						r.Get("/", app.getReadingHandler)
						r.Put("/", app.updateReadingHandler)
					})
				})

				r.Route("/reports", func(r chi.Router) {
					r.Get("/balance-sheet", app.getBalanceSheetHandler)
					r.Get("/trial-balance", app.getTrialBalanceHandler)
					r.Get("/customer-balance-summary", app.getCustomerBalanceSummaryHandler)
				})

			})
		})

		r.Route("/people_types", func(r chi.Router) {
			r.Get("/", app.getPeopleTypesHandler)
			r.Post("/", app.createPeopleTypeHandler)
			r.Route("/{peopleTypeID}", func(r chi.Router) {
				r.Get("/", app.getPeopleTypeHandler)
				r.Put("/", app.updatePeopleTypeHandler)
				r.Delete("/", app.deletePeopleTypeHandler)
			})
		})

		r.Route("/account_types", func(r chi.Router) {
			r.Get("/", app.getAccountTypesHandler)
			r.Post("/", app.createAccountTypeHandler)
			r.Route("/{accountTypeID}", func(r chi.Router) {
				r.Get("/", app.getAccountTypeHandler)
				r.Put("/", app.updateAccountTypeHandler)
				r.Delete("/", app.deleteAccountTypeHandler)
			})
		})

		// r.Route("/posts", func(r chi.Router) {
		// 	r.Post("/", app.createPostHandler)

		// 	r.Route("/{postID}", func(r chi.Router) {
		// 		r.Use(app.postsContextMiddleware)

		// 		r.Get("/", app.getPostHandler)
		// 		r.Delete("/", app.deletePostHandler)
		// 		r.Patch("/", app.updatePostHandler)

		// 	})
		// })

		// r.Route("/users", func(r chi.Router) {
		// 	r.Put("/activate/{token}", app.activateUserHandler)

		// 	r.Route("/{userID}", func(r chi.Router) {
		// 		r.Use(app.userContextMiddleware)

		// 		r.Get("/", app.getUserHandler)
		// 		r.Put("/follow", app.followUserHandler)
		// 		r.Put("/unfollow", app.unfollowUserHandler)
		// 	})

		// 	r.Group(func(r chi.Router) {
		// 		r.Get("/feed", app.getUserFeedHandler)
		// 	})
		// })

		// r.Route("/authentication", func(r chi.Router) {
		// 	r.Post("/user", app.registerUserHandler)
		// })
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("server has started", "addr", app.config.addr, "env", app.config.env)

	return srv.ListenAndServe()
}
