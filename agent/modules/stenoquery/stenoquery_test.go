// Copyright 2019 Jason Ertel (jertel). All rights reserved.
//
// This program is distributed under the terms of version 2 of the
// GNU General Public License.  See LICENSE for further details.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

package stenoquery

import (
	"strconv"
	"testing"
	"time"

	"github.com/khulnasoft/khulnasoft-soc/model"
	"github.com/stretchr/testify/assert"
)

func TestInitStenoQuery(tester *testing.T) {
	cfg := make(map[string]interface{})
	sq := NewStenoQuery(nil)
	err := sq.Init(cfg)
	assert.Error(tester, err)
	assert.Equal(tester, DEFAULT_EXECUTABLE_PATH, sq.executablePath)
	assert.Equal(tester, DEFAULT_PCAP_OUTPUT_PATH, sq.pcapOutputPath)
	assert.Equal(tester, DEFAULT_PCAP_INPUT_PATH, sq.pcapInputPath)
	assert.Equal(tester, DEFAULT_TIMEOUT_MS, sq.timeoutMs)
	assert.Equal(tester, DEFAULT_EPOCH_REFRESH_MS, sq.epochRefreshMs)
	assert.Equal(tester, DEFAULT_DATA_LAG_MS, sq.dataLagMs)
}

func TestDataLag(tester *testing.T) {
	cfg := make(map[string]interface{})
	sq := NewStenoQuery(nil)
	sq.Init(cfg)
	lagDate := sq.getDataLagDate()
	assert.False(tester, lagDate.After(time.Now()), "expected data lag datetime to be before current datetime")
}

func TestCreateQuery(tester *testing.T) {
	sq := NewStenoQuery(nil)

	job := model.NewJob()
	job.Filter.BeginTime, _ = time.Parse(time.RFC3339, "2006-01-02T15:05:05Z")
	job.Filter.EndTime, _ = time.Parse(time.RFC3339, "2006-01-02T15:06:05Z")
	expectedQuery := "before 2006-01-02T15:06:05Z and after 2006-01-02T15:05:05Z"
	query := sq.CreateQuery(job)
	assert.Equal(tester, expectedQuery, query)

	job.Filter.SrcIp = "1.2.3.4"
	query = sq.CreateQuery(job)
	expectedQuery = expectedQuery + " and host " + job.Filter.SrcIp
	assert.Equal(tester, expectedQuery, query)

	job.Filter.DstIp = "1.2.1.2"
	query = sq.CreateQuery(job)
	expectedQuery = expectedQuery + " and host " + job.Filter.DstIp
	assert.Equal(tester, expectedQuery, query)

	job.Filter.SrcPort = 123
	query = sq.CreateQuery(job)
	expectedQuery = expectedQuery + " and port " + strconv.Itoa(job.Filter.SrcPort)
	assert.Equal(tester, expectedQuery, query)

	job.Filter.DstPort = 123
	query = sq.CreateQuery(job)
	expectedQuery = expectedQuery + " and port " + strconv.Itoa(job.Filter.DstPort)
	assert.Equal(tester, expectedQuery, query)
}
