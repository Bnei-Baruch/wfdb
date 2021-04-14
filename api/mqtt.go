package api

import (
	"encoding/json"
	"fmt"
	"github.com/Bnei-Baruch/wfdb/common"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strings"
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

func (a *App) SubMQTT(c mqtt.Client) {
	log.Info().Str("source", "MQTT").Msg("- Connected -")
	if token := a.Msg.Subscribe(common.ServiceTopic, byte(2), a.execMessage); token.Wait() && token.Error() != nil {
		log.Fatal().Str("source", "MQTT").Err(token.Error()).Msg("Subscription error")
	} else {
		log.Info().Str("source", "MQTT").Msg("Subscription - " + common.ServiceTopic)
	}

	if token := a.Msg.Subscribe(common.ExtPrefix+common.ServiceTopic, byte(2), a.execMessage); token.Wait() && token.Error() != nil {
		log.Fatal().Str("source", "MQTT").Err(token.Error()).Msg("Subscription error")
	} else {
		log.Info().Str("source", "MQTT").Msg("Subscription - " + common.ExtPrefix + common.ServiceTopic)
	}
}

func (a *App) LostMQTT(c mqtt.Client, err error) {
	log.Error().Str("source", "MQTT").Err(err).Msg("Lost Connection")
}

func (a *App) execMessage(c mqtt.Client, m mqtt.Message) {
	log.Debug().Str("source", "MQTT").Msgf("Received message: %s from topic: %s\n", m.Payload(), m.Topic())
	id := "false"
	s := strings.Split(m.Topic(), "/")
	p := string(m.Payload())

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
		topic = common.ServiceDataTopic + common.EP
	} else {
		topic = common.ServiceDataTopic + common.EP + "/" + id
	}
	message, err := json.Marshal(m)
	if err != nil {
		log.Error().Str("source", "MQTT").Err(err).Msg("Message parsing")
	}

	text := fmt.Sprintf(string(message))
	if token := a.Msg.Publish(topic, byte(2), false, text); token.Wait() && token.Error() != nil {
		log.Error().Str("source", "MQTT").Err(err).Msg("Send Respond")
	}
}

func (a *App) InitLogMQTT() {
	mqtt.DEBUG = NewPahoLogAdapter(zerolog.InfoLevel)
	mqtt.WARN = NewPahoLogAdapter(zerolog.WarnLevel)
	mqtt.CRITICAL = NewPahoLogAdapter(zerolog.ErrorLevel)
	mqtt.ERROR = NewPahoLogAdapter(zerolog.ErrorLevel)
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
