package dota2api

import (
	"flag"
	"testing"

	"os"

	"encoding/json"
	"io/ioutil"

	"github.com/stretchr/testify/suite"
)

var writeFixtures bool

func init() {
	flag.BoolVar(&writeFixtures, "write-fixtures", false, "write the fixtures instead of testing against them")
	flag.Parse()
}

type Suite struct {
	suite.Suite
	c *Client
}

func (s *Suite) SetupSuite() {
	key, ok := os.LookupEnv("DOTA2_API_KEY")
	if !ok {
		s.FailNow("DOTA2_API_KEY environment variable not set")
	}
	s.c = NewClient(key)
}

func (s *Suite) TestGetMatchHistoryBySequenceNum() {
	res, err := s.c.GetMatchHistoryBySequenceNum(2081736343, 10)
	if !s.NoError(err) {
		s.T().FailNow()
	}
	if !s.NotNil(res) {
		s.T().FailNow()
	}
	b, err := json.MarshalIndent(res, "", "  ")
	if !s.NoError(err) {
		s.T().FailNow()
	}
	if writeFixtures {
		if !s.NoError(os.MkdirAll("fixtures", os.ModeDir)) {
			s.T().FailNow()
		}
		if !s.NoError(ioutil.WriteFile("fixtures/GetMatchHistoryBySequenceNum.json", b, 0666)) {
			s.T().FailNow()
		}
		s.T().SkipNow()
	}
	f, err := ioutil.ReadFile("fixtures/GetMatchHistoryBySequenceNum.json")
	if !s.NoError(err) {
		s.T().FailNow()
	}
	s.Equal(f, b)
}

func TestSteamSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
