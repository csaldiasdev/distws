package agent

import (
	"fmt"
	"net"
	"net/http"

	"github.com/csaldiasdev/distws/internal/discovery"
	"github.com/csaldiasdev/distws/internal/httpserver"
	"github.com/csaldiasdev/distws/internal/repository"
	"github.com/csaldiasdev/distws/internal/util"
	"github.com/csaldiasdev/distws/internal/wshub"

	"github.com/google/uuid"
)

type Agent struct {
	serfMembership *discovery.Membership
	repository     repository.Repository
	wsHub          *wshub.Hub
	httpServer     *http.Server
	configuration  AgentConfiguration
}

type AgentConfiguration struct {
	RepositoryRpcPort uint     `json:"repository_rpc_port"`
	HubRpcPort        uint     `json:"hub_rpc_port"`
	RaftPort          uint     `json:"raft_port"`
	HttpPort          uint     `json:"http_port"`
	SerfPort          uint     `json:"serf_port"`
	SerfMembers       []string `json:"serf_members"`
}

func NewAgent(config AgentConfiguration) (*Agent, error) {
	localIp, err := util.GetLocalIp()

	if err != nil {
		return nil, err
	}

	nodeId := uuid.New()
	repo, err := repository.NewRepository(nodeId.String(), localIp, config.RaftPort, config.RepositoryRpcPort)

	if err != nil {
		return nil, err
	}

	wshub := wshub.NewHub(nodeId.String(), localIp, config.HubRpcPort, repo)

	joinFunc := func(np *discovery.NodeProps) {
		repo.AddNode(np.NodeId.String(), np.NodeIp, np.RaftPort, np.RepositoryRpcPort)
		wshub.AddHubNode(np.NodeId.String(), np.NodeIp, np.HubRpcPort)
	}

	leaveFunc := func(np *discovery.NodeProps) {}

	nodeMember, err := discovery.NewMembership(config.SerfPort, discovery.NodeProps{
		NodeIp:            localIp,
		NodeId:            nodeId,
		RepositoryRpcPort: config.RepositoryRpcPort,
		HubRpcPort:        config.HubRpcPort,
		RaftPort:          config.RaftPort,
	}, joinFunc, leaveFunc, config.SerfMembers...)

	if err != nil {
		return nil, err
	}

	httpsvr := httpserver.NewHTTPServer(wshub)

	return &Agent{
		serfMembership: nodeMember,
		repository:     repo,
		wsHub:          wshub,
		httpServer:     httpsvr,
		configuration:  config,
	}, nil
}

func (a *Agent) Run() error {
	httpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.configuration.HttpPort))

	if err != nil {
		return err
	}

	return a.httpServer.Serve(httpListener)
}
