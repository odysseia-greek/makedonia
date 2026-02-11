package atomos

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/odysseia-greek/agora/aristoteles"
	"github.com/odysseia-greek/agora/aristoteles/models"
	eupalinos "github.com/odysseia-greek/agora/eupalinos/stomion"
	"github.com/odysseia-greek/agora/plato/config"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/plato/service"
	"github.com/odysseia-greek/delphi/aristides/diplomat"
	pb "github.com/odysseia-greek/delphi/aristides/proto"
	"google.golang.org/grpc/metadata"
)

const (
	defaultIndex    string = "dictionary"
	defaultMinNGram string = "3"
	defaultMaxNGram string = "5"
	envMaxNGram     string = "MAX_NGRAM"
	envMinNGram     string = "MIN_NGRAM"
)

func CreateNewConfig() (*DemokritosHandler, error) {
	tls := config.BoolFromEnv(config.EnvTlSKey)

	var cfg models.Config
	ambassador, err := diplomat.NewClientAmbassador(diplomat.DEFAULTADDRESS)

	healthy := ambassador.WaitForHealthyState()
	if !healthy {
		logging.Info("ambassador service not ready - restarting seems the only option")
		os.Exit(1)
	}

	traceId := uuid.New().String()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	md := metadata.New(map[string]string{service.HeaderKey: traceId})
	ctx = metadata.NewOutgoingContext(context.Background(), md)
	vaultConfig, err := ambassador.GetSecret(ctx, &pb.VaultRequest{})
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}

	elasticService := aristoteles.ElasticService(tls)

	cfg = models.Config{
		Service:     elasticService,
		Username:    vaultConfig.ElasticUsername,
		Password:    vaultConfig.ElasticPassword,
		ElasticCERT: vaultConfig.ElasticCERT,
	}

	elastic, err := aristoteles.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	err = aristoteles.HealthCheck(elastic)
	if err != nil {
		return nil, err
	}

	min := config.StringFromEnv(envMinNGram, defaultMinNGram)
	max := config.StringFromEnv(envMaxNGram, defaultMaxNGram)

	minNGram, err := strconv.Atoi(min)
	if err != nil {
		return nil, err
	}
	maxNGram, err := strconv.Atoi(max)
	if err != nil {
		return nil, err
	}

	index := config.StringFromEnv(config.EnvIndex, defaultIndex)
	searchWord := config.StringFromEnv(config.EnvSearchWord, config.DefaultSearchWord)
	policyName := config.StringFromEnv("HOT_POLICY_NAME", "hot_plain")

	eupalinosAddress := config.StringFromEnv(config.EnvEupalinosService, config.DefaultEupalinosService)
	logging.Debug(fmt.Sprintf("creating new eupalinos client: %s", eupalinosAddress))
	queue, err := eupalinos.NewEupalinosClient(eupalinosAddress)
	if err != nil {
		logging.Error(err.Error())
	}

	logging.Debug("waiting for queue to be ready")
	queueHealthy := queue.WaitForHealthyState()
	if !queueHealthy {
		logging.Debug("no queue that is healthy")
	}

	channel := config.StringFromEnv(config.EnvJobName, config.DefaultJobName)

	var buf bytes.Buffer

	return &DemokritosHandler{
		Index:      index,
		Created:    0,
		SearchWord: searchWord,
		Elastic:    elastic,
		MaxNGram:   maxNGram,
		MinNGram:   minNGram,
		PolicyName: policyName,
		Buf:        buf,
		Ambassador: ambassador,
		Eupalinos:  queue,
		Channel:    channel,
	}, nil
}
