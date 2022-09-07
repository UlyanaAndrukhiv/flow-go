package main

import (
	"context"
	"flag"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	flowsdk "github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access"
	client "github.com/onflow/flow-go-sdk/access/grpc"

	"github.com/onflow/flow-go/integration/benchmark"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module/metrics"
	"github.com/onflow/flow-go/utils/unittest"
)

type LoadCase struct {
	tps      int
	duration time.Duration
}

func main() {
	sleep := flag.Duration("sleep", 0, "duration to sleep before benchmarking starts")
	loadTypeFlag := flag.String("load-type", "token-transfer", "type of loads (\"token-transfer\", \"add-keys\", \"computation-heavy\", \"event-heavy\", \"ledger-heavy\", \"const-exec\")")
	tpsFlag := flag.String("tps", "1", "transactions per second (TPS) to send, accepts a comma separated list of values if used in conjunction with `tps-durations`")
	tpsDurationsFlag := flag.String("tps-durations", "0", "duration that each load test will run, accepts a comma separted list that will be applied to multiple values of the `tps` flag (defaults to infinite if not provided, meaning only the first tps case will be tested; additional values will be ignored)")
	chainIDStr := flag.String("chain", string(flowsdk.Emulator), "chain ID")
	accessNodes := flag.String("access", net.JoinHostPort("127.0.0.1", "3569"), "access node address")
	serviceAccountPrivateKeyHex := flag.String("servPrivHex", unittest.ServiceAccountPrivateKeyHex, "service account private key hex")
	logLvl := flag.String("log-level", "info", "set log level")
	metricport := flag.Uint("metricport", 8080, "port for /metrics endpoint")
	pushgateway := flag.String("pushgateway", "127.0.0.1:9091", "host:port for pushgateway")
	profilerEnabled := flag.Bool("profiler-enabled", false, "whether to enable the auto-profiler")
	_ = flag.Bool("track-txs", false, "deprecated")
	accountMultiplierFlag := flag.Int("account-multiplier", 100, "number of accounts to create per load tps")
	feedbackEnabled := flag.Bool("feedback-enabled", true, "wait for trannsaction execution before submitting new transaction")
	maxConstExecTxSizeInBytes := flag.Uint("const-exec-max-tx-size", flow.DefaultMaxTransactionByteSize/10, "max byte size of constant exec transaction size to generate")
	authAccNumInConstExecTx := flag.Uint("const-exec-num-authorizer", 1, "num of authorizer for each constant exec transaction to generate")
	argSizeInByteInConstExecTx := flag.Uint("const-exec-arg-size", 100, "byte size of tx argument for each constant exec transaction to generate")
	payerKeyCountInConstExecTx := flag.Uint("const-exec-payer-key-count", 2, "num of payer keys for each constant exec transaction to generate")
	flag.Parse()

	chainID := flowsdk.ChainID([]byte(*chainIDStr))

	// parse log level and apply to logger
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
	lvl, err := zerolog.ParseLevel(strings.ToLower(*logLvl))
	if err != nil {
		log.Fatal().Err(err).Msg("invalid log level")
	}
	log = log.Level(lvl)

	server := metrics.NewServer(log, *metricport, *profilerEnabled)
	<-server.Ready()
	loaderMetrics := metrics.NewLoaderCollector()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sp := benchmark.NewStatsPusher(ctx, log, *pushgateway, "loader", prometheus.DefaultGatherer)
	defer sp.Stop()

	addressGen := flowsdk.NewAddressGenerator(chainID)
	serviceAccountAddress := addressGen.NextAddress()
	log.Info().Msgf("Service Address: %v", serviceAccountAddress)
	fungibleTokenAddress := addressGen.NextAddress()
	log.Info().Msgf("Fungible Token Address: %v", fungibleTokenAddress)
	flowTokenAddress := addressGen.NextAddress()
	log.Info().Msgf("Flow Token Address: %v", flowTokenAddress)

	// sleep in order to ensure the testnet is up and running
	if *sleep > 0 {
		log.Info().Msgf("Sleeping for %v before starting benchmark", sleep)
		time.Sleep(*sleep)
	}

	accessNodeAddrs := strings.Split(*accessNodes, ",")
	clients := make([]access.Client, 0, len(accessNodeAddrs))
	for _, addr := range accessNodeAddrs {
		client, err := client.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal().Str("addr", addr).Err(err).Msgf("unable to initialize flow client")
		}
		clients = append(clients, client)
	}

	// run load cases
	for i, c := range parseLoadCases(log, tpsFlag, tpsDurationsFlag) {
		log.Info().Str("load_type", *loadTypeFlag).Int("number", i).Int("tps", c.tps).Dur("duration", c.duration).Msgf("Running load case...")

		loaderMetrics.SetTPSConfigured(c.tps)

		var lg *benchmark.ContLoadGenerator
		if c.tps > 0 {
			var err error
			lg, err = benchmark.New(
				ctx,
				log,
				loaderMetrics,
				clients,
				benchmark.NetworkParams{
					ServAccPrivKeyHex:     *serviceAccountPrivateKeyHex,
					ServiceAccountAddress: &serviceAccountAddress,
					FungibleTokenAddress:  &fungibleTokenAddress,
					FlowTokenAddress:      &flowTokenAddress,
				},
				benchmark.LoadParams{
					TPS:              c.tps,
					NumberOfAccounts: c.tps * *accountMultiplierFlag,
					LoadType:         benchmark.LoadType(*loadTypeFlag),
					FeedbackEnabled:  *feedbackEnabled,
				},
				benchmark.ConstExecParams{
					MaxTxSizeInByte: *maxConstExecTxSizeInBytes,
					AuthAccountNum:  *authAccNumInConstExecTx,
					ArgSizeInByte:   *argSizeInByteInConstExecTx,
					PayerKeyCount:   *payerKeyCountInConstExecTx,
				},
			)
			if err != nil {
				log.Fatal().Err(err).Msgf("unable to create new cont load generator")
			}

			err = lg.Init()
			if err != nil {
				log.Fatal().Err(err).Msgf("unable to init loader")
			}
			lg.Start()
		}

		// if the duration is 0, we run this case forever
		if c.duration.Nanoseconds() == 0 {
			for {
				time.Sleep(time.Minute)
			}
		}

		time.Sleep(c.duration)

		if lg != nil {
			lg.Stop()
		}
	}
}

func parseLoadCases(log zerolog.Logger, tpsFlag, tpsDurationsFlag *string) []LoadCase {
	tpsStrings := strings.Split(*tpsFlag, ",")
	var cases []LoadCase
	for _, s := range tpsStrings {
		t, err := strconv.ParseInt(s, 0, 32)
		if err != nil {
			log.Fatal().Err(err).Str("value", s).
				Msg("could not parse tps flag, expected comma separated list of integers")
		}
		cases = append(cases, LoadCase{tps: int(t)})
	}

	tpsDurationsStrings := strings.Split(*tpsDurationsFlag, ",")
	for i := range cases {
		if i >= len(tpsDurationsStrings) {
			break
		}

		// ignore empty entries (implying that case will run indefinitely)
		if tpsDurationsStrings[i] == "" {
			continue
		}

		d, err := time.ParseDuration(tpsDurationsStrings[i])
		if err != nil {
			log.Fatal().Err(err).Str("value", tpsDurationsStrings[i]).
				Msg("could not parse tps-durations flag, expected comma separated list of durations")
		}
		cases[i].duration = d
	}

	return cases
}