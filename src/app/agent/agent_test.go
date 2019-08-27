package agent_test

import (
	"fmt"
	"github.com/onsi/gomega/types"
	"log"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/leadership-election/app/agent"
)

var run = 10000

var _ = Describe("Agent", func() {
	var agents map[string]*agent.Agent

	BeforeEach(func() {
		agents = make(map[string]*agent.Agent)

		var nodes []string

		// There are 3 intra network addresses and 7 fake addresses to simulate unresponsive agents
		for i := 3; i <= 12; i++ {
			nodes = append(nodes, fmt.Sprintf("127.0.0.1:%d", run+i))
		}

		for i := 0; i < 3; i++ {
			a := agent.New(
				i,
				nodes,

				// External address
				agent.WithPort(run+i),
				agent.WithLogger(log.New(GinkgoWriter, fmt.Sprintf("[AGENT %d]", i), log.LstdFlags)),
			)
			a.Start()
			agents[fmt.Sprintf("http://%s/v1/leader", a.Addr())] = a
		}
	})

	AfterEach(func() {
		run += 2 * len(agents)

		// We set up fake and not serviced addresses
		run += 7
	})

	It("returns a 200 if it is the leader", func() {
		Eventually(getLeaderStatusFunc(agents), 10).Should(haveSingleLeader(agents))
		Consistently(getLeaderStatusFunc(agents), 3).Should(haveSingleLeader(agents))
	})
})

func getLeaderStatusFunc(agents map[string]*agent.Agent) func() []int {
	return func() []int {
		var responses []int
		for addr := range agents {
			resp, err := http.Get(addr)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Or(Equal(http.StatusOK), Equal(http.StatusLocked)))

			responses = append(responses, resp.StatusCode)
		}

		return responses
	}
}

func haveSingleLeader(agents map[string]*agent.Agent) types.GomegaMatcher {
	var nonLeaders []int
	for i := 1; i < len(agents); i++ {
		nonLeaders = append(nonLeaders, http.StatusLocked)
	}

	return ConsistOf(append(nonLeaders, http.StatusOK))
}
