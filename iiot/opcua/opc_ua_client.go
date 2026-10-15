package opcua

import (
	"context"
	"errors"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
	"github.com/spf13/cast"
)

type OPCUADevice struct {
	Device  string      `json:"device"`
	TagList []*OPCUATag `json:"raw_tag_list"`
}

type OPCUATag struct {
	Tag          string    `json:"tag"`
	NodeAddr     string    `json:"node_addr"`
	Value        any       `json:"value"`
	Quality      string    `json:"quality"`
	LinerScaling float64   `json:"liner_scaling"` // 线性缩放
	ServerTime   time.Time `json:"server_time"`
	SourceTime   time.Time `json:"source_time"`
}

type OPCUANodeData struct {
	NodeAddr   string
	Value      any
	Quality    string
	ServerTime time.Time
	SourceTime time.Time
}

func NewOPCUAClientWithNoSignAndEncrypt(endpoint string, username string, password string) *OPCUAClientWithNoSignAndEncrypt {
	oc := new(OPCUAClientWithNoSignAndEncrypt)
	oc.Endpoint = endpoint
	oc.Username = username
	oc.Password = password
	if oc.Username == "" {
		oc.authMode = ua.UserTokenTypeAnonymous
	} else {
		oc.authMode = ua.UserTokenTypeUserName
	}
	oc.secMode = ua.MessageSecurityModeNone
	oc.secPolicy = ua.SecurityPolicyURINone
	return oc
}

// OPCUAClientWithNoSignAndEncrypt
type OPCUAClientWithNoSignAndEncrypt struct {
	Endpoint       string `json:"endpoint"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Client         *opcua.Client
	authMode       ua.UserTokenType
	secMode        ua.MessageSecurityMode
	secPolicy      string
	ServerEndpoint *ua.EndpointDescription
}

func (o *OPCUAClientWithNoSignAndEncrypt) GetServerEndpoint(ctx context.Context) error {
	endpoints, err := opcua.GetEndpoints(ctx, o.Endpoint)
	if err != nil {
		return err
	}
	var serverEndpoint *ua.EndpointDescription
	for _, e := range endpoints {
		if e.SecurityPolicyURI == o.secPolicy && e.SecurityMode == o.secMode && (serverEndpoint == nil || e.SecurityLevel >= serverEndpoint.SecurityLevel) {
			serverEndpoint = e
		}
	}
	if serverEndpoint == nil {
		o.ServerEndpoint = nil
		return errors.New("can not find matched EndpointDescription")
	}
	o.ServerEndpoint = serverEndpoint
	return nil
}

func (o *OPCUAClientWithNoSignAndEncrypt) Connect(ctx context.Context) error {
	if o.ServerEndpoint == nil {
		return errors.New("can not find matched EndpointDescription")
	}
	if o.Username == "" {
		c, err := opcua.NewClient(o.Endpoint,
			opcua.SecurityMode(o.secMode),
			opcua.SecurityFromEndpoint(o.ServerEndpoint, o.authMode))
		if err != nil {
			return err
		}
		if err = c.Connect(ctx); err != nil {
			return err
		}
		if c == nil {
			return errors.New("OPC UA Client Object is Null")
		}
		o.Client = c
	} else {
		c, err := opcua.NewClient(o.Endpoint,
			opcua.AuthUsername(o.Username, o.Password),
			opcua.SecurityMode(o.secMode),
			opcua.SecurityFromEndpoint(o.ServerEndpoint, o.authMode))
		if err != nil {
			return err
		}
		if err = c.Connect(ctx); err != nil {
			return err
		}
		if c == nil {
			return errors.New("OPC UA Client Object is Null")
		}
		o.Client = c
	}

	return nil
}

func (o *OPCUAClientWithNoSignAndEncrypt) Close(ctx context.Context) error {
	if o.Client != nil {
		return o.Client.Close(ctx)
	}
	return nil
}

func (o *OPCUAClientWithNoSignAndEncrypt) Read(ctx context.Context, deviceList []*OPCUADevice) error {
	if o.Client == nil {
		return errors.New("OPC UA Client is Null")
	}
	readIdList := make([]*ua.ReadValueID, 0)
	for _, dev := range deviceList {
		for _, tag := range dev.TagList {
			nodeId, err := ua.ParseNodeID(tag.NodeAddr)
			if err != nil {
				return err
			}
			if nodeId != nil {
				readId := new(ua.ReadValueID)
				readId.NodeID = nodeId
				readIdList = append(readIdList, readId)
			}
		}
	}
	nodeDataList := make(map[string]*OPCUANodeData)
	if len(readIdList) > 0 {
		req := new(ua.ReadRequest)
		req.MaxAge = 2000
		req.NodesToRead = readIdList
		req.TimestampsToReturn = ua.TimestampsToReturnBoth
		resp, err := o.Client.Read(ctx, req)
		if err != nil {
			return err
		}
		if resp == nil || resp.Results == nil {
			return errors.New("can not read any data")
		}
		if len(readIdList) != len(resp.Results) {
			return errors.New("ReadIdList Length != Results Length,Maybe NodeAddress Wong")
		}
		for i, result := range resp.Results {
			nodeData := new(OPCUANodeData)
			nodeData.NodeAddr = readIdList[i].NodeID.String()
			if result.Value != nil {
				nodeData.Value = result.Value.Value()
			} else {
				nodeData.Value = -9999
			}
			switch {
			case errors.Is(result.Status, ua.StatusOK):
				nodeData.Quality = "good"
			case errors.Is(result.Status, ua.StatusBad):
				nodeData.Quality = "bad"
			case errors.Is(result.Status, ua.StatusUncertain):
				nodeData.Quality = "uncertain"
			default:
				nodeData.Quality = cast.ToString(result.Status)
			}
			nodeData.ServerTime = result.ServerTimestamp
			nodeData.SourceTime = result.SourceTimestamp
			nodeDataList[nodeData.NodeAddr] = nodeData
		}
	}
	if len(nodeDataList) > 0 {
		for _, dev := range deviceList {
			for _, tag := range dev.TagList {
				if nodeData, ok := nodeDataList[tag.NodeAddr]; ok {
					tag.Value = nodeData.Value
					tag.Quality = nodeData.Quality
					tag.ServerTime = nodeData.ServerTime
					tag.SourceTime = nodeData.SourceTime
				}
			}
		}
	}
	return nil
}
