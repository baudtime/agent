package outputs

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/baudtime/agent/plugin"
	"github.com/baudtime/agent/vars"
	"github.com/baudtime/baudtime/msg"
	"github.com/baudtime/baudtime/msg/gateway"
	"github.com/baudtime/baudtime/tcp/client"
	ts "github.com/baudtime/baudtime/util/time"
	"github.com/go-kit/kit/log/level"
	"gopkg.in/alecthomas/kingpin.v2"
)

var endpoints = kingpin.Flag("output.baudtime.endpoints", "baudtime endpoints").Default("debug.jd.local:8088").String()

func init() {
	register("baudtime", plugin.DefaultEnabled, newBaudtime)
}

type baudtime struct {
	client *client.Client
}

func newBaudtime() (plugin.Output, error) {
	domain, port, err := parseEndpoint(*endpoints)
	if err != nil {
		level.Error(vars.Logger).Log("err", err)
		return nil, err
	}

	addrProvider := client.NewDnsAddrProvider(domain, port)
	cli := client.NewGatewayClient("name", addrProvider)

	return &baudtime{client: cli}, nil
}

func (bt *baudtime) Connect() error {
	return nil
}

func (bt *baudtime) Write(metrics []plugin.Metric) error {
	t := ts.FromTime(time.Now())
	request := &gateway.AddRequest{}

	for _, metric := range metrics {
		labels := append(metric.Labels, msg.Label{"__name__", metric.Name})
		sort.Slice(labels, func(i, j int) bool {
			return labels[i].Name < labels[i].Name
		})

		request.Series = append(request.Series, &msg.Series{
			Labels: labels,
			Points: []msg.Point{
				{T: t, V: metric.Value},
			},
		})
	}

	if len(request.Series) > 0 {
		ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
		resp, err := bt.client.SyncRequest(ctx, request)
		if err != nil {
			return err
		}
		if resp, ok := resp.(*msg.GeneralResponse); ok && resp.Status != msg.StatusCode_Succeed {
			return errors.New(resp.Message)
		}
	}

	return nil
}

func (bt *baudtime) Close() error {
	return bt.client.Close()
}

func parseEndpoint(endpoint string) (domain string, port int, err error) {
	domainAndPort := strings.Split(endpoint, ":")
	domain = strings.TrimSpace(domainAndPort[0])
	port, err = strconv.Atoi(strings.TrimSpace(domainAndPort[1]))
	return
}
