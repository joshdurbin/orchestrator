package client

import (
	"database/sql"
	"time"
)

// APIResponseCode represents the response status code
type APIResponseCode int

const (
	// ERROR indicates an error response
	ERROR APIResponseCode = iota
	// OK indicates a successful response
	OK
)

// HttpStatus returns the HTTP status code for the response
func (code APIResponseCode) HttpStatus() int {
	if code == OK {
		return 200
	}
	return 500
}

// APIResponse is the standard response structure for all API calls
type APIResponse struct {
	Code    APIResponseCode `json:"Code"`
	Message string          `json:"Message"`
	Details interface{}     `json:"Details"`
}

// Instance represents a database instance, including its current configuration & status
type Instance struct {
	Key                          InstanceKey
	InstanceAlias                string
	Uptime                       uint
	ServerID                     uint
	ServerUUID                   string
	Version                      string
	VersionComment               string
	FlavorName                   string
	ReadOnly                     bool
	Binlog_format                string
	BinlogRowImage               string
	LogBinEnabled                bool
	LogSlaveUpdatesEnabled       bool // for API backwards compatibility
	LogReplicationUpdatesEnabled bool
	SelfBinlogCoordinates        BinlogCoordinates
	MasterKey                    InstanceKey
	MasterUUID                   string
	AncestryUUID                 string
	IsDetachedMaster             bool

	Slave_SQL_Running          bool // for API backwards compatibility
	ReplicationSQLThreadRuning bool
	Slave_IO_Running           bool // for API backwards compatibility
	ReplicationIOThreadRuning  bool
	ReplicationSQLThreadState  ReplicationThreadState
	ReplicationIOThreadState   ReplicationThreadState

	HasReplicationFilters bool
	GTIDMode              string
	SupportsOracleGTID    bool
	UsingOracleGTID       bool
	UsingMariaDBGTID      bool
	UsingPseudoGTID       bool
	ReadBinlogCoordinates BinlogCoordinates
	ExecBinlogCoordinates BinlogCoordinates
	IsDetached            bool
	RelaylogCoordinates   BinlogCoordinates
	LastSQLError          string
	LastIOError           string
	SecondsBehindMaster   sql.NullInt64
	SQLDelay              uint
	ExecutedGtidSet       string
	GtidPurged            string
	GtidErrant            string

	SlaveLagSeconds                   sql.NullInt64 // for API backwards compatibility
	ReplicationLagSeconds             sql.NullInt64
	SlaveHosts                        InstanceKeyMap // for API backwards compatibility
	Replicas                          InstanceKeyMap
	ClusterName                       string
	SuggestedClusterAlias             string
	DataCenter                        string
	Region                            string
	PhysicalEnvironment               string
	ReplicationDepth                  uint
	IsCoMaster                        bool
	HasReplicationCredentials         bool
	ReplicationCredentialsAvailable   bool
	SemiSyncAvailable                 bool
	SemiSyncPriority                  uint
	SemiSyncMasterPluginNewVersion    bool
	SemiSyncReplicaPluginNewVersion   bool
	SemiSyncMasterEnabled             bool
	SemiSyncReplicaEnabled            bool
	SemiSyncMasterTimeout             uint64
	SemiSyncMasterWaitForReplicaCount uint
	SemiSyncMasterStatus              bool
	SemiSyncMasterClients             uint
	SemiSyncReplicaStatus             bool

	LastSeenTimestamp    string
	IsLastCheckValid     bool
	IsUpToDate           bool
	IsRecentlyChecked    bool
	SecondsSinceLastSeen sql.NullInt64
	CountMySQLSnapshots  int

	IsCandidate          bool
	PromotionRule        CandidatePromotionRule
	IsDowntimed          bool
	DowntimeReason       string
	DowntimeOwner        string
	DowntimeEndTimestamp string
	ElapsedDowntime      time.Duration
	UnresolvedHostname   string
	AllowTLS             bool

	Problems []string

	LastDiscoveryLatency time.Duration

	// Group replication fields
	ReplicationGroupName                   string
	ReplicationGroupIsSinglePrimary        bool
	ReplicationGroupMemberState            string
	ReplicationGroupMemberRole             string
	ReplicationGroupMembers                InstanceKeyMap
	ReplicationGroupPrimaryInstanceKey     InstanceKey
}

// BinlogCoordinates holds binlog file and position information
type BinlogCoordinates struct {
	LogFile string
	LogPos  int64
	Type    BinlogType
}

// BinlogType represents the type of binary log
type BinlogType int

const (
	BinaryLog BinlogType = iota
	RelayLog
)

// InstanceKeyMap is a map of instance keys
type InstanceKeyMap map[InstanceKey]bool

// ReplicationThreadState represents the state of a replication thread
type ReplicationThreadState int

const (
	ReplicationThreadStateNoThread ReplicationThreadState = iota
	ReplicationThreadStateStopped
	ReplicationThreadStateRunning
	ReplicationThreadStateOther
)

// CandidatePromotionRule describes the promotion preference/rule for an instance
type CandidatePromotionRule string

const (
	MustPromoteRule       CandidatePromotionRule = "must"
	PreferPromoteRule     CandidatePromotionRule = "prefer"
	NeutralPromoteRule    CandidatePromotionRule = "neutral"
	PreferNotPromoteRule  CandidatePromotionRule = "prefer_not"
	MustNotPromoteRule    CandidatePromotionRule = "must_not"
)

// ClusterInfo contains information about a MySQL cluster
type ClusterInfo struct {
	ClusterName                            string
	ClusterAlias                           string
	ClusterDomain                          string
	CountInstances                         uint
	HeuristicLag                           int64
	HasAutomatedMasterRecovery             bool
	HasAutomatedIntermediateMasterRecovery bool
}

// PoolInstancesMap maps instance keys by pool name
type PoolInstancesMap map[string][]InstanceKey

// ReplicationAnalysis provides analysis of replication issues
type ReplicationAnalysis struct {
	AnalyzedInstanceKey              InstanceKey
	AnalyzedInstanceMasterKey        InstanceKey
	ClusterDetails                   ClusterInfo
	Analysis                         string
	Description                      string
	StructureAnalysis                []string
	IsMaster                         bool
	IsCoMaster                       bool
	LastCheckValid                   bool
	LastCheckPartialSuccess          bool
	CountReplicas                    uint
	CountValidReplicas               uint
	CountValidReplicatingReplicas    uint
	CountReplicasFailingToConnectToMaster uint
	CountDowntimedReplicas           uint
	ReplicationDepth                 uint
	IsDowntimed                      bool
	IsReplicasDowntimed              bool
	DowntimeEndTimestamp             string
	DowntimeRemainingSeconds         int
	IsBinlogServer                   bool
	PseudoGTIDImmediateTopology      bool
	OracleGTIDImmediateTopology      bool
	MariaDBGTIDImmediateTopology     bool
	BinlogServerImmediateTopology    bool
	SemiSyncMasterEnabled            bool
	SemiSyncMasterStatus             bool
	SemiSyncMasterWaitForReplicaCount uint
	SemiSyncMasterClients            uint
	CountSemiSyncReplicasEnabled     uint
	IsActionableRecovery             bool
	ProcessingNodeHostname           string
	ProcessingNodeToken              string
	CountAdditionalAgreeingNodes     int
	StartActivePeriod                string
	SkippableDueToDowntime           bool
	GTIDMode                         string
	MinReplicaGTIDMode               string
	MaxReplicaGTIDMode               string
	MaxReplicaGTIDErrant             string
	CommandHint                      string
	IsReadOnly                       bool
}

// TopologyRecovery represents a recovery process
type TopologyRecovery struct {
	Id                         int64
	UID                        string
	AnalysisEntry              ReplicationAnalysis
	SuccessorKey               *InstanceKey
	SuccessorAlias             string
	SuccessorBinlogCoordinates *BinlogCoordinates
	IsActive                   bool
	IsSuccessful               bool
	LostReplicas               InstanceKeyMap
	ParticipatingInstanceKeys  InstanceKeyMap
	AllErrors                  []string
	RecoveryStartTimestamp     string
	RecoveryEndTimestamp       string
	ProcessingNodeHostname     string
	ProcessingNodeToken        string
	Acknowledged               bool
	AcknowledgedAt             string
	AcknowledgedBy             string
	AcknowledgedComment        string
	LastDetectionId            int64
	RelatedRecoveryId          int64
	Type                       string
	RecoveryType               string
}

// RecoveryStep represents a single step in a recovery process
type RecoveryStep struct {
	RecoveryUID string
	AuditAt     string
	Message     string
}

// Maintenance represents a maintenance window
type Maintenance struct {
	MaintenanceId  uint
	Key            InstanceKey
	BeginTimestamp string
	SecondsElapsed uint
	IsActive       bool
	Owner          string
	Reason         string
}

// MaintenanceInstanceKey is an alias for Maintenance
type MaintenanceInstanceKey Maintenance

// Agent represents an orchestrator agent
type Agent struct {
	Hostname           string
	Port               int
	LastSubmitted      string
	AvailableLocalSnapshots []string
	AvailableSnapshotHosts  []string
	TotalSecondsUnavailable int64
	AvailableDiskSpaceRatio float64
	LogicalVolume          *AgentLogicalVolume
}

// AgentLogicalVolume represents an LVM logical volume on an agent
type AgentLogicalVolume struct {
	Name          string
	IsActive      bool
	HasSnapshot   bool
	SnapshotName  string
	DataPath      string
	SnapshotPath  string
	MySQLPort     int
	MySQLDataPath string
	MySQLDiskPath string
	FileSystem    string
}

// AgentSeed represents a seed operation between agents
type AgentSeed struct {
	SeedId            int64
	TargetHostname    string
	SourceHostname    string
	StartTimestamp    string
	EndTimestamp      string
	IsComplete        bool
	IsSuccessful      bool
}

// AgentSeedState represents the state of a seed operation
type AgentSeedState struct {
	SeedId        int64
	StateTimestamp string
	State         string
	ErrorMessage  string
}

// AuditEntry represents an audit log entry
type AuditEntry struct {
	AuditId           int64
	AuditTimestamp    string
	AuditType         string
	AuditInstanceKey  InstanceKey
	Message           string
}

// HostnameResolveCache represents a hostname resolution cache entry
type HostnameResolveCache struct {
	Hostname         string
	ResolvedHostname string
}

// Tag represents an instance tag
type Tag struct {
	TagName  string
	TagValue string
}

// DiscoveryMetric represents discovery metrics
type DiscoveryMetric struct {
	Timestamp      string
	InstanceKey    InstanceKey
	DurationMillis int64
}

// DiscoveryQueueMetric represents discovery queue metrics
type DiscoveryQueueMetric struct {
	Timestamp      string
	QueueName      string
	QueueLength    int
}

// BackendQueryMetric represents backend query metrics
type BackendQueryMetric struct {
	Timestamp      string
	Query          string
	DurationMillis int64
}

// WriteBufferMetric represents write buffer metrics
type WriteBufferMetric struct {
	Timestamp      string
	BufferSize     int
}

// RaftMembershipHealth represents Raft cluster health
type RaftMembershipHealth struct {
	Healthy    bool
	Reason     string
}

// RaftFollowerHealthReport represents a follower's health report
type RaftFollowerHealthReport struct {
	Hostname           string
	Token              string
	RaftBind           string
	RaftAdvertise      string
	IsAvailable        bool
	AvailabilityReason string
}

// RaftState represents the state of Raft cluster
type RaftState struct {
	Leader     string
	Peer       string
	Peers      []string
	IsLeader   bool
	IsFollower bool
	State      string
}

// AutomatedRecoveryFilter represents a recovery filter
type AutomatedRecoveryFilter struct {
	Pattern      string
	IsPromotion  bool
}

// BlockedTopologyRecovery represents a blocked recovery
type BlockedTopologyRecovery struct {
	FailedInstanceKey InstanceKey
	ClusterName       string
	Analysis          string
	BlockingRecoveryId int64
	BlockingRecovery  *TopologyRecovery
}
