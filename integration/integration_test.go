package integration

import (

	"testing"
	
	"github.com/stretchr/testify/suite"

	"k8s-utils/conf"
	"k8s-utils/connection"

)

func TestIntegration(t *testing.T) {
	suite.Run(t, &intSpringbootSuite{})
}

func (s *intSpringbootSuite) TestSpringbootAppEnd2End() {
	s.podStatus()
	s.springbootStart()
	s.actuatorHealth()
	s.swagger()
}

func (s *intSpringbootSuite) SetupSuite() {
	s.config = conf.LoadConfig()
	s.connection = connection.KubeConnection(s.config.K8S.Cluster, s.config.K8S.InCluster)
}




