package main

struct Agent []str {
name: str
networkAddress: str
port: int
sendRequest: func
}

struct DomainCharacter []str {
	domain: str
	agent: Agent
	response: []str
}



func main() {
	// respond to command from human admin panel
	// domain input
	domain := 'unfiltered.co.uk'

	// call first agent to get character, categories and content for this domain
	// get agent name and network
	agentA = new Agent {'DomainCharacter', '127.0.0.1', 8082}

	// call agent and get record of communication - it might be a pipeline or parallel mesh
	// agentA may call its own team to plan and work out the best way to solve the problem and to solve it
	domainCharacter := agentA->SendRequest(domain)



}

func (Agent a) sendRequest(domain) : DomainCharacter {
	response []str = {'subject1', 'subject2', 'subject3'}
	return new DomainCharacter(domain, a, response)
}


