package fluentdsyslogs

/**
 * Panther is a scalable, powerful, cloud-native SIEM written in Golang/React.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import (
	"net"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers"
	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers/timestamp"
)

var RFC5424Desc = `Fluentd syslog parser for the RFC3164 format (ie. BSD-syslog messages)
Reference: https://docs.fluentd.org/parser/syslog#rfc3164-log`

// nolint:lll
type RFC5424 struct {
	Priority  *uint8                      `json:"pri" validate:"required" description:"Priority is calculated by (Facility * 8 + Severity). The lower this value, the higher importance of the log message."`
	Hostname  *string                     `json:"host,omitempty" description:"Hostname identifies the machine that originally sent the syslog message."`
	Ident     *string                     `json:"ident,omitempty" description:"Appname identifies the device or application that originated the syslog message."`
	ProcID    *string                     `json:"pid,omitempty" description:"ProcID is often the process ID, but can be any value used to enable log analyzers to detect discontinuities in syslog reporting."`
	MsgID     *string                     `json:"msgid,omitempty" description:"MsgID identifies the type of message. For example, a firewall might use the MsgID 'TCPIN' for incoming TCP traffic."`
	ExtraData *string                     `json:"extradata,omitempty" description:"ExtraData contains syslog strucured data as string"`
	Message   *string                     `json:"message,omitempty" description:"Message contains free-form text that provides information about the event."`
	Timestamp *timestamp.FluentdTimestamp `json:"time,omitempty" description:"Timestamp of the syslog message in UTC."`
	Tag       *string                     `json:"tag,omitempty" description:"Tag of the syslog message"`
	// NOTE: added to end of struct to allow expansion later
	parsers.PantherLog
}

// RFC5424Parser parses fluentd syslog logs in the RFC5424 format
type RFC5424Parser struct{}

func (p *RFC5424Parser) New() parsers.LogParser {
	return &RFC5424Parser{}
}

// Parse returns the parsed events or nil if parsing failed
func (p *RFC5424Parser) Parse(log string) []*parsers.PantherLog {
	rfc5424 := &RFC5424{}

	err := jsoniter.UnmarshalFromString(log, rfc5424)
	if err != nil {
		zap.L().Debug("failed to parse log", zap.Error(err))
		return nil
	}

	rfc5424.updatePantherFields(p)

	if err := parsers.Validator.Struct(rfc5424); err != nil {
		zap.L().Debug("failed to validate log", zap.Error(err))
		return nil
	}

	return rfc5424.Logs()
}

// LogType returns the log type supported by this parser
func (p *RFC5424Parser) LogType() string {
	return "Fluentd.Syslog5424"
}

func (event *RFC5424) updatePantherFields(p *RFC5424Parser) {
	event.SetCoreFields(p.LogType(), (*timestamp.RFC3339)(event.Timestamp), event)
	if event.Hostname != nil {
		// The hostname should be a FQDN, but may also be an IP address. Check for IP, otherwise
		// add as a domain name. https://tools.ietf.org/html/rfc3164#section-6.2.4
		hostname := *event.Hostname
		if net.ParseIP(hostname) != nil {
			event.AppendAnyIPAddresses(hostname)
		} else {
			event.AppendAnyDomainNames(hostname)
		}
	}
}
