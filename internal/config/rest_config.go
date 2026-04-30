package config

import (
	"fmt"

	"log"
	"os"

	"github.com/FPRPL26/rpl-be/db"
	"github.com/FPRPL26/rpl-be/internal/api/controller"
	"github.com/FPRPL26/rpl-be/internal/api/repository"
	"github.com/FPRPL26/rpl-be/internal/api/routes"
	"github.com/FPRPL26/rpl-be/internal/api/service"
	"github.com/FPRPL26/rpl-be/internal/middleware"
	"github.com/FPRPL26/rpl-be/internal/pkg/cron"
	mailer "github.com/FPRPL26/rpl-be/internal/pkg/email"
	"github.com/FPRPL26/rpl-be/internal/pkg/midtrans"
	"github.com/gin-gonic/gin"
)

type RestConfig struct {
	server *gin.Engine
}

func NewRest() RestConfig {
	db := db.New()
	app := gin.Default()
	server := NewRouter(app, db)
	middleware := middleware.New(db)

	cron.StartMediaCron(db)

	var (
		//=========== (PACKAGE) ===========//
		mailerService   mailer.Mailer            = mailer.New()
		midtransService midtrans.MidtransService = midtrans.NewMidtransService()
		// oauthService  oauth.Oauth   = oauth.New()
		// awsS3Service  storage.AwsS3 = storage.NewAwsS3()

		//=========== (REPOSITORY) ===========//
		userRepository              repository.UserRepository                    = repository.NewUser(db)
		refreshTokenRepository      repository.RefreshTokenRepository            = repository.NewRefreshTokenRepository(db)
		taskRepository              repository.TaskRepository                    = repository.NewTask(db)
		classRepository             repository.ClassRepository                   = repository.NewClassRepository(db)
		scheduleRepository          repository.ScheduleRepository                = repository.NewScheduleRepository(db)
		tutorProfileRepository      repository.TutorProfileRepository            = repository.NewTutorProfileRepository(db)
		userClassRequestRepository  repository.ClassRequestRepository            = repository.NewClassRequestRepository(db)
		transactionRepository       repository.ClassTransactionRepository        = repository.NewClassTransactionRepository(db)
		barterTransactionRepository repository.BarterSkillTransactionRepository  = repository.NewBarterSkillTransactionRepository(db)
		classRequestTransRepository repository.ClassRequestTransactionRepository = repository.NewClassRequestTransactionRepository(db)
		reviewRepository            repository.ReviewRepository                  = repository.NewReviewRepository(db)
		//=========== (SERVICE) ===========//
		authService         service.AuthService             = service.NewAuth(userRepository, refreshTokenRepository, mailerService, db)
		taskService         service.TaskService             = service.NewTask(taskRepository)
		classService        service.ClassService            = service.NewClass(classRepository, scheduleRepository, transactionRepository, reviewRepository, db)
		tutorProfileService service.TutorService            = service.NewTutorService(tutorProfileRepository)
		classRequestService service.ClassRequestService     = service.NewClassRequestService(userClassRequestRepository)
		transactionService  service.ClassTransactionService = service.NewClassTransactionService(db, transactionRepository, scheduleRepository, classRepository, userRepository, midtransService)
		reviewService       service.ReviewService           = service.NewReviewService(reviewRepository, transactionRepository, barterTransactionRepository, classRequestTransRepository)
		// userService                   service.UserService                   = service.NewUser(userRepository, userDisciplineRepository, disciplineGroupConsolidatorRepository, disciplineListDocumentConsolidatorRepository, packageRepository, db)

		//=========== (CONTROLLER) ===========//
		authController         controller.AuthController             = controller.NewAuth(authService)
		taskController         controller.TaskController             = controller.NewTask(taskService)
		classController        controller.ClassController            = controller.NewClass(classService)
		tutorProfileController controller.TutorController            = controller.NewTutorController(tutorProfileService)
		classRequestController controller.ClassRequestController     = controller.NewClassRequest(classRequestService)
		transactionController  controller.ClassTransactionController = controller.NewClassTransactionController(transactionService)
		reviewController       controller.ReviewController           = controller.NewReviewController(reviewService)
	// userController                   controller.UserController                   = controller.NewUser(userService)
	)

	// Register all routes
	routes.Auth(server, authController, middleware)
	routes.Task(server, taskController, middleware)
	routes.Class(server, classController, middleware)
	routes.Tutor(server, tutorProfileController, middleware)
	routes.ClassRequest(server, classRequestController, middleware)
	routes.ClassTransaction(server, transactionController, reviewController, middleware)
	routes.Review(server, reviewController, middleware)

	return RestConfig{
		server: server,
	}
}

func (ap *RestConfig) Start() {
	port := os.Getenv("APP_PORT")
	host := os.Getenv("APP_HOST")
	if port == "" {
		port = "8080"
	}

	serve := fmt.Sprintf("%s:%s", host, port)
	if err := ap.server.Run(serve); err != nil {
		log.Panicf("failed to start server: %s", err)
	}
	log.Println("server start on port ", serve)
}
