package node

import (
	"bytes"
	"os"
	"path"

	"github.com/algorand/go-algorand/config"
	"github.com/algorand/go-algorand/util/codecs"
)

func GetConfigs(n *AlgorandFullNode) string {
	configs, err := os.ReadFile(path.Join(n.rootDir, "config.json"))
	if err != nil {
		//return "config.json file not found"
		var buffer bytes.Buffer
		enc := codecs.NewFormattedJSONEncoder(&buffer)
		err = enc.Encode(defaultLocal)
		configs = buffer.Bytes()
	}
	return string(configs)
}

// from local_defaults.go
var defaultLocal = config.Local{
	Version:                                    21,
	AccountUpdatesStatsInterval:                5000000000,
	AccountsRebuildSynchronousMode:             1,
	AgreementIncomingBundlesQueueLength:        7,
	AgreementIncomingProposalsQueueLength:      25,
	AgreementIncomingVotesQueueLength:          10000,
	AnnounceParticipationKey:                   true,
	Archival:                                   false,
	BaseLoggerDebugLevel:                       4,
	BlockServiceCustomFallbackEndpoints:        "",
	BroadcastConnectionsLimit:                  -1,
	CadaverSizeTarget:                          1073741824,
	CatchpointFileHistoryLength:                365,
	CatchpointInterval:                         10000,
	CatchpointTracking:                         0,
	CatchupBlockDownloadRetryAttempts:          1000,
	CatchupBlockValidateMode:                   0,
	CatchupFailurePeerRefreshRate:              10,
	CatchupGossipBlockFetchTimeoutSec:          4,
	CatchupHTTPBlockFetchTimeoutSec:            4,
	CatchupLedgerDownloadRetryAttempts:         50,
	CatchupParallelBlocks:                      16,
	ConnectionsRateLimitingCount:               60,
	ConnectionsRateLimitingWindowSeconds:       1,
	DNSBootstrapID:                             "<network>.algorand.network",
	DNSSecurityFlags:                           1,
	DeadlockDetection:                          0,
	DeadlockDetectionThreshold:                 30,
	DisableLocalhostConnectionRateLimit:        true,
	DisableNetworking:                          false,
	DisableOutgoingConnectionThrottling:        false,
	EnableAccountUpdatesStats:                  false,
	EnableAgreementReporting:                   false,
	EnableAgreementTimeMetrics:                 false,
	EnableAssembleStats:                        false,
	EnableBlockService:                         false,
	EnableBlockServiceFallbackToArchiver:       true,
	EnableCatchupFromArchiveServers:            false,
	EnableDeveloperAPI:                         false,
	EnableGossipBlockService:                   true,
	EnableIncomingMessageFilter:                false,
	EnableLedgerService:                        false,
	EnableMetricReporting:                      false,
	EnableOutgoingNetworkMessageFiltering:      true,
	EnablePingHandler:                          true,
	EnableProcessBlockStats:                    false,
	EnableProfiler:                             false,
	EnableRequestLogger:                        false,
	EnableTopAccountsReporting:                 false,
	EnableVerbosedTransactionSyncLogging:       false,
	EndpointAddress:                            "127.0.0.1:0",
	FallbackDNSResolverAddress:                 "",
	ForceFetchTransactions:                     false,
	ForceRelayMessages:                         false,
	GossipFanout:                               4,
	IncomingConnectionsLimit:                   800,
	IncomingMessageFilterBucketCount:           5,
	IncomingMessageFilterBucketSize:            512,
	IsIndexerActive:                            false,
	LedgerSynchronousMode:                      2,
	LogArchiveMaxAge:                           "",
	LogArchiveName:                             "node.archive.log",
	LogSizeLimit:                               1073741824,
	MaxAPIResourcesPerAccount:                  100000,
	MaxCatchpointDownloadDuration:              7200000000000,
	MaxConnectionsPerIP:                        30,
	MinCatchpointFileDownloadBytesPerSecond:    20480,
	NetAddress:                                 "",
	NetworkMessageTraceServer:                  "",
	NetworkProtocolVersion:                     "",
	NodeExporterListenAddress:                  ":9100",
	NodeExporterPath:                           "./node_exporter",
	OptimizeAccountsDatabaseOnStartup:          false,
	OutgoingMessageFilterBucketCount:           3,
	OutgoingMessageFilterBucketSize:            128,
	ParticipationKeysRefreshInterval:           60000000000,
	PeerConnectionsUpdateInterval:              3600,
	PeerPingPeriodSeconds:                      0,
	PriorityPeers:                              map[string]bool{},
	ProposalAssemblyTime:                       250000000,
	PublicAddress:                              "",
	ReconnectTime:                              60000000000,
	ReservedFDs:                                256,
	RestConnectionsHardLimit:                   2048,
	RestConnectionsSoftLimit:                   1024,
	RestReadTimeoutSeconds:                     15,
	RestWriteTimeoutSeconds:                    120,
	RunHosted:                                  false,
	SuggestedFeeBlockHistory:                   3,
	SuggestedFeeSlidingWindowSize:              50,
	TLSCertFile:                                "",
	TLSKeyFile:                                 "",
	TelemetryToLog:                             true,
	TransactionSyncDataExchangeRate:            0,
	TransactionSyncSignificantMessageThreshold: 0,
	TxPoolExponentialIncreaseFactor:            2,
	TxPoolSize:                                 15000,
	TxSyncIntervalSeconds:                      60,
	TxSyncServeResponseSize:                    1000000,
	TxSyncTimeoutSeconds:                       30,
	UseXForwardedForAddressField:               "",
	VerifiedTranscationsCacheSize:              30000,
}
