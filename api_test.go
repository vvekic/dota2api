package dota2api

import (
	"flag"
	"testing"

	"os"

	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
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

func (s *Suite) fixtureTest(v interface{}, path string) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return errors.Wrap(err, "error creating fixture")
	}
	if writeFixtures {
		if err := os.MkdirAll("fixtures", os.ModeDir); err != nil {
			return errors.Wrap(err, "error creating dir")
		}
		if err := ioutil.WriteFile(path, b, 0666); err != nil {
			return errors.Wrap(err, "error writing fixture file")
		}
	} else {
		f, err := ioutil.ReadFile(path)
		if err != nil {
			return errors.Wrap(err, "error reading fixture file")
		}
		s.Equal(f, b)
	}
	return nil
}

func (s *Suite) TestGetMatchHistoryBySequenceNum() {
	res, err := s.c.GetMatchHistoryBySequenceNum(2081736343, 10)
	if !s.NoError(err) {
		s.T().FailNow()
	}
	if !s.NotNil(res) {
		s.T().FailNow()
	}
	if !s.NoError(s.fixtureTest(res, "fixtures/GetMatchHistoryBySequenceNum.json")) {
		s.T().FailNow()
	}
}

func (s *Suite) TestGetMatchHistory() {
	res, err := s.c.GetMatchHistory()
	if !s.NoError(err) {
		s.T().FailNow()
	}
	if !s.NotNil(res) {
		s.T().FailNow()
	}
	if !s.NoError(s.fixtureTest(res, "fixtures/GetMatchHistory.json")) {
		s.T().FailNow()
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
