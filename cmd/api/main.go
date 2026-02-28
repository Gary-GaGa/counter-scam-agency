package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	domainpersonnel "counter-scam-agency/internal/domain/personnel"
	mongoinfra "counter-scam-agency/internal/infrastructure/persistence/mongo"
	httphandler "counter-scam-agency/internal/interface/in/http"
	defenserepo "counter-scam-agency/internal/interface/out/persistence/mongo/defense"
	invrepo "counter-scam-agency/internal/interface/out/persistence/mongo/investigation"
	"counter-scam-agency/internal/interface/out/persistence/mongo/mission"
	"counter-scam-agency/internal/interface/out/persistence/mongo/player"
	usecasedefense "counter-scam-agency/internal/usecase/defense"
	usecase "counter-scam-agency/internal/usecase/investigation"
	usecasepersonnel "counter-scam-agency/internal/usecase/personnel"
)

func main() {
	ctx := context.Background()

	cfg := mongoinfra.Config{
		URI:      getenv("MONGO_URI", "mongodb://localhost:27017"),
		Database: getenv("MONGO_DB", "counter_scam_agency"),
	}

	client, err := mongoinfra.NewClient(ctx, cfg)
	if err != nil {
		log.Fatalf("mongo connect: %v", err)
	}
	defer func() { _ = mongoinfra.Disconnect(ctx, client) }()

	db, err := mongoinfra.NewDatabase(client, cfg)
	if err != nil {
		log.Fatalf("mongo database: %v", err)
	}

	// Repositories.
	missionsRepo := mission.NewMongoRepository(db)
	investigationsRepo := invrepo.NewMongoRepository(db)
	playersRepo := player.NewMongoRepository(db)
	basesRepo := defenserepo.NewMongoRepository(db)

	// Usecases.
	investigationSvc := usecase.NewService(missionsRepo, investigationsRepo, playersRepo)
	personnelSvc := usecasepersonnel.NewService(playersRepo, buildSkillCatalog())
	defenseSvc := usecasedefense.NewService(basesRepo)

	// HTTP Server.
	srv := httphandler.NewServer(investigationSvc, personnelSvc, defenseSvc)
	corsOrigin := getenv("CORS_ORIGIN", "*")
	handler := httphandler.SecurityHeadersMiddleware(httphandler.CORSMiddleware(corsOrigin)(srv.Handler()))

	addr := getenv("ADDR", ":8080")
	httpServer := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	// Graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("API server listening on %s", addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}

	fmt.Println("server stopped")
}

func buildSkillCatalog() []domainpersonnel.Skill {
	return []domainpersonnel.Skill{
		{
			ID:                 "skill-logic-1",
			Type:               domainpersonnel.SkillTypeAnalysis,
			Name:               "矛盾掃描",
			Description:        "短時間內分析對話矛盾，提高破案效率。",
			CooldownSeconds:    10,
			ReputationRequired: 0,
		},
		{
			ID:                 "skill-nego-1",
			Type:               domainpersonnel.SkillTypeNegotiation,
			Name:               "談判增幅",
			Description:        "提升談判成功率與獲得資訊品質。",
			CooldownSeconds:    12,
			ReputationRequired: 5,
		},
		{
			ID:                 "skill-defense-1",
			Type:               domainpersonnel.SkillTypeDefense,
			Name:               "防禦強化",
			Description:        "強化資安防線，降低詐騙滲透風險。",
			CooldownSeconds:    15,
			ReputationRequired: 10,
		},
		{
			ID:                 "skill-forensics-1",
			Type:               domainpersonnel.SkillTypeForensics,
			Name:               "數位鑑識",
			Description:        "深度取證分析，揭露隱藏的數位足跡。",
			CooldownSeconds:    18,
			ReputationRequired: 15,
		},
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
