package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Bnei-Baruch/wfdb/common"
	"github.com/Bnei-Baruch/wfdb/models"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/url"
	"strings"
	"time"
)

type MqttPayload struct {
	Action  string      `json:"action,omitempty"`
	ID      string      `json:"id,omitempty"`
	Name    string      `json:"name,omitempty"`
	Source  string      `json:"src,omitempty"`
	Error   error       `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Result  string      `json:"result,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (a *App) Init() error {
	log.Info().Str("source", "APP").Msgf("Init MQTT")

	serverURL, err := url.Parse(common.SERVER)
	if err != nil {
		log.Error().Str("source", "MQTT").Err(err).Msg("environmental variable must be a valid URL")
	}

	cliCfg := autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{serverURL},
		KeepAlive:         10,
		ConnectRetryDelay: 3 * time.Second,
		OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
			log.Info().Str("source", "APP").Msgf("MQTT connection up")
			if _, err := cm.Subscribe(context.Background(), &paho.Subscribe{
				Subscriptions: map[string]paho.SubscribeOptions{
					common.ServiceTopic: {QoS: byte(1)},
				},
			}); err != nil {
				log.Error().Str("source", "MQTT").Err(err).Msg("client.Subscribe")
				return
			}
			log.Info().Str("source", "APP").Msgf("MQTT subscription made")
		},
		OnConnectError: func(err error) { log.Error().Str("source", "MQTT").Err(err).Msg("error whilst attempting connection") },
		ClientConfig: paho.ClientConfig{
			ClientID: "wfdb_mqtt_client",
			//Router: paho.RegisterHandler(common.WorkflowExec, m.execMessage),
			Router:        paho.NewStandardRouter(),
			OnClientError: func(err error) { log.Error().Str("source", "MQTT").Err(err).Msg("server requested disconnect:") },
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					log.Error().Str("source", "MQTT").Err(err).Msgf("MQTT server requested disconnect: %d", d.Properties.ReasonString)
				} else {
					log.Error().Str("source", "MQTT").Err(err).Msgf("MQTT server requested disconnect: %d", d.ReasonCode)
				}
			},
		},
	}

	cliCfg.SetUsernamePassword(common.USERNAME, []byte(common.PASSWORD))

	debugLog := NewPahoLogAdapter(zerolog.DebugLevel)

	cliCfg.Debug = debugLog
	cliCfg.PahoDebug = debugLog

	a.Msg, err = autopaho.NewConnection(context.Background(), cliCfg)
	if err != nil {
		return err
	}

	cliCfg.Router.RegisterHandler(common.ServiceTopic, a.execMessage)

	return nil
}

func (a *App) execMessage(m *paho.Publish) {
	log.Debug().Str("source", "MQTT").Msgf("Received message: %s from topic: %s\n", string(m.Payload), m.Topic)
	id := "false"
	s := strings.Split(m.Topic, "/")
	p := string(m.Payload)

	if s[0] == "kli" && len(s) == 5 {
		id = s[4]
	} else if s[0] == "exec" && len(s) == 4 {
		id = s[3]
	}

	if id == "false" {
		switch p {
		case "start":
			//	go a.startExecMqtt(p)
			//case "stop":
			//	go a.stopExecMqtt(p)
			//case "status":
			//	go a.execStatusMqtt(p)
		}
	}

	if id != "false" {
		switch p {
		case "start":
			//	go a.startExecMqttByID(p, id)
			//case "stop":
			//	go a.stopExecMqttByID(p, id)
			//case "status":
			//	go a.execStatusMqttByID(p, id)
			//case "cmdstat":
			//	go a.cmdStatMqtt(p, id)
			//case "progress":
			//	go a.getProgressMqtt(p, id)
			//case "report":
			//	go a.getReportMqtt(p, id)
			//case "alive":
			//	go a.isAliveMqtt(p, id)
		}
	}
}

func (a *App) SendRespond(id string, m *MqttPayload) {
	var topic string

	if id == "false" {
		topic = common.ServiceDataTopic
	} else {
		topic = common.ServiceDataTopic + "/" + id
	}
	message, err := json.Marshal(m)
	if err != nil {
		log.Error().Str("source", "MQTT").Err(err).Msg("Message parsing")
	}

	pa, err := a.Msg.Publish(context.Background(), &paho.Publish{
		QoS:     byte(1),
		Retain:  false,
		Topic:   topic,
		Payload: message,
	})
	if err != nil {
		log.Error().Str("source", "MQTT").Err(err).Msgf("MQTT Publish error: %d ", pa.Properties.ReasonString)
	}
}

func (a *App) SendMessage(id string) {
	var topic string
	var m interface{}
	date := time.Now().Format("2006-01-02")

	if id == "ingest" {
		topic = common.MonitorIngestTopic
		m, _ = models.FindIngest(a.DB, "date", date)
	}

	if id == "trimmer" {
		topic = common.MonitorTrimmerTopic
		m, _ = models.FindTrimmer(a.DB, "date", date)
	}

	if id == "archive" {
		topic = common.MonitorArchiveTopic
		m, _ = models.FindKmFiles(a.DB, "date", date)
	}

	if id == "trim" {
		topic = common.StateTrimmerTopic
		m, _ = models.GetFilesToTrim(a.DB)
	}

	if id == "drim" {
		topic = common.StateDgimaTopic
		m, _ = models.GetFilesToDgima(a.DB)
	}

	if id == "bdika" {
		topic = common.StateArichaTopic
		m, _ = models.GetBdika(a.DB)
	}

	if id == "jobs" {
		topic = common.StateJobsTopic
		m, _ = models.GetActiveJobs(a.DB)
	}

	if id == "langcheck" {
		topic = common.StateLangcheckTopic
		var s models.State
		s.StateID = "langcheck"
		_ = s.GetState(a.DB)
		m = s.Data
	}

	message, err := json.Marshal(m)
	if err != nil {
		log.Error().Str("monitor", "MQTT").Err(err).Msg("Message parsing")
	}

	pa, err := a.Msg.Publish(context.Background(), &paho.Publish{
		QoS:     byte(1),
		Retain:  true,
		Topic:   topic,
		Payload: message,
	})
	if err != nil {
		log.Error().Str("source", "MQTT").Err(err).Msg("Publish: Topic - " + topic + " " + pa.Properties.ReasonString)
	}

	log.Debug().Str("source", "MQTT").Str("json", string(message)).Msg("Publish: Topic - " + topic)
}

type PahoLogAdapter struct {
	level zerolog.Level
}

func NewPahoLogAdapter(level zerolog.Level) *PahoLogAdapter {
	return &PahoLogAdapter{level: level}
}

func (a *PahoLogAdapter) Println(v ...interface{}) {
	log.Debug().Str("source", "MQTT").Msgf("%s", fmt.Sprint(v...))
}

func (a *PahoLogAdapter) Printf(format string, v ...interface{}) {
	log.Debug().Str("source", "MQTT").Msgf("%s", fmt.Sprintf(format, v...))
}
