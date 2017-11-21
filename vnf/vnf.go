package vnf

import (
	"time"
	"syscall"

	"github.com/sofianinho/vnf-api-golang/vnf/types"
	"github.com/sofianinho/vnf-api-golang/scheduler"
	"github.com/sofianinho/vnf-api-golang/storage"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	//DefaultVNFName is the name of the lte-softmodem exe folder inside bin folder
	DefaultVNFName = "oai"
	//DefaultLTEMName is the name of the LTE-M lte-softmodem exe folder inside bin folder
	DefaultLTEMName = "lte-m"
)

var (
	sessionsGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "sessions",
		Help: "Sessions",
	})
	configsGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "configurations",
		Help: "Configuration files",
	})
	instancesGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "instances",
		Help: "Instances of VNF",
	})
	latencyHistogramVec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "vnf_action_duration_ms",
		Help:    "How long it took to process a specific action, in a specific host",
		Buckets: []float64{300, 1200, 5000},
	}, []string{"action"})
)

func observeAction(action string, start time.Time) {
	latencyHistogramVec.WithLabelValues(action).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
}

func init() {
	prometheus.MustRegister(sessionsGauge)
	prometheus.MustRegister(configsGauge)
	prometheus.MustRegister(instancesGauge)
	prometheus.MustRegister(latencyHistogramVec)
}

type vnf struct {
	instances	scheduler.API
	storage		storage.API
}

//API is the interface to use to create, delete, and intract with configurations and VNFs (instances)
type API interface{
	SessionNew()(*types.Session, error)
	SessionGet(id string)(*types.Session)
	SessionDelete(id string)(error)
	SessionCount()(int,error)
	SessionList()([]*types.Session, error)
	SessionStatus()(*types.Status, error)

	ConfigNew(sID, ver string, conf *types.VNFParams, alias string, tags []string)(*types.Config, error)
	ConfigGet(sID,cID string)(*types.Config, error)
	ConfigGetByAlias(sID,alias string)(*types.Config, error)
	ConfigUpdate(sID, cID, ver string, conf *types.VNFParams, alias string, tags []string)(error)
	ConfigDelete(sID,cID string)(error)
	ConfigCount(sID string)(int,error)
	ConfigList(sID string)([]*types.Config, error)

	InstanceNew(sID string, conf *types.Config, alias string, tags []string, artefact string)(*types.Instance, error)
	InstanceGet(sID, iID string)(*types.Instance, error)
	InstanceGetByAlias(sID,alias string)(*types.Instance, error)
	InstanceUpdate(sID,iID string, conf *types.Config, alias string, tags []string, artefact string)(error)
	InstanceDelete(sID, iID string)(error)
	InstanceCount(sID string)(int, error)
	InstanceList(sID string)([]*types.Instance, error)
	InstanceRun(*types.Instance)(error)
	InstanceKill(syscall.Signal, string)(error)
	InstanceStatus(sID,iID string)(string, error)


}

//New returns a new VNF interface
func New(i scheduler.API, s storage.API) (API){
	return &vnf{instances: i, storage: s}
}

func (v *vnf) setGauges(){
	s, _ := v.storage.SessionCount()
	sessionsGauge.Set(float64(s))

	c, _ := v.storage.ConfigsCount()
	configsGauge.Set(float64(c))

	i, _ := v.storage.InstancesCount()
	instancesGauge.Set(float64(i))
}