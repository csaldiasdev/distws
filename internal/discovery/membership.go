package discovery

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"

	"github.com/google/uuid"
	"github.com/hashicorp/serf/serf"
)

const nodePropsKey = "node_props"

type config struct {
	NodePort       uint
	Tags           map[string]string
	StartJoinAddrs []string
}

type NodeProps struct {
	NodeIp            string    `json:"node_ip"`
	NodeId            uuid.UUID `json:"node_id"`
	RepositoryRpcPort uint      `json:"repository_rpc_port"`
	HubRpcPort        uint      `json:"hub_rpc_port"`
	RaftPort          uint      `json:"raft_port"`
}

type Membership struct {
	config
	serf            *serf.Serf
	events          chan serf.Event
	handleJoinFunc  func(*NodeProps)
	handleLeaveFunc func(*NodeProps)
}

func (m *Membership) isLocal(member serf.Member) bool {
	return m.serf.LocalMember().Name == member.Name
}

func (m *Membership) Members() []serf.Member {
	return m.serf.Members()
}

func (m *Membership) AliveMembers() []serf.Member {
	activeMembers := make([]serf.Member, 0)

	for _, v := range m.serf.Members() {
		if v.Status == serf.StatusAlive {
			activeMembers = append(activeMembers, v)
		}
	}

	return activeMembers
}

func (m *Membership) Leave() error {
	return m.serf.Leave()
}

func (m *Membership) handleJoin(member serf.Member) {
	props, _ := decodeProps(member)
	m.handleJoinFunc(props)
}

func (m *Membership) handleLeave(member serf.Member) {
	props, _ := decodeProps(member)
	m.handleLeaveFunc(props)
}

func (m *Membership) eventHandler() {
	for e := range m.events {
		switch e.EventType() {
		case serf.EventMemberJoin:
			for _, member := range e.(serf.MemberEvent).Members {
				if m.isLocal(member) {
					continue
				}
				m.handleJoin(member)
			}
		case serf.EventMemberLeave, serf.EventMemberFailed:
			for _, member := range e.(serf.MemberEvent).Members {
				if m.isLocal(member) {
					continue
				}
				m.handleLeave(member)
			}
		}
	}
}

func (m *Membership) setupSerf() (err error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", m.NodePort))

	if err != nil {
		return err
	}

	config := serf.DefaultConfig()
	config.Init()
	config.MemberlistConfig.BindAddr = addr.IP.String()
	config.MemberlistConfig.BindPort = addr.Port

	m.events = make(chan serf.Event)
	config.EventCh = m.events

	config.Tags = m.Tags
	config.NodeName = uuid.NewString()

	m.serf, err = serf.Create(config)

	if err != nil {
		return err
	}

	go m.eventHandler()

	if m.StartJoinAddrs != nil {
		_, err := m.serf.Join(m.StartJoinAddrs, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewMembership(nodePort uint, props NodeProps, handleJoinFunc func(*NodeProps), handleLeaveFunc func(*NodeProps), StartJoinAddrs ...string) (*Membership, error) {

	bytesProps, _ := json.Marshal(props)

	sEnc := base64.StdEncoding.EncodeToString(bytesProps)

	c := &Membership{
		config: config{
			NodePort:       nodePort,
			Tags:           map[string]string{nodePropsKey: sEnc},
			StartJoinAddrs: StartJoinAddrs,
		},
		handleJoinFunc:  handleJoinFunc,
		handleLeaveFunc: handleLeaveFunc,
	}

	err := c.setupSerf()

	if err != nil {
		return nil, err
	}

	return c, nil
}

func decodeProps(member serf.Member) (*NodeProps, error) {
	bytesValue, err := base64.StdEncoding.DecodeString(member.Tags[nodePropsKey])

	if err != nil {
		return nil, err
	}

	var props NodeProps

	if json.Unmarshal(bytesValue, &props) != nil {
		return nil, err
	}

	return &props, nil
}
