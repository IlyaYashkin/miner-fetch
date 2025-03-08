package device

type VersionCommand struct {
	Response VersionResponse
}

func (v *VersionCommand) GetResponse() any {
	return &v.Response
}

func (v *VersionCommand) GetCommand() string {
	return "version"
}

type SummaryCommand struct {
	Response SummaryResponse
}

func (s *SummaryCommand) GetResponse() any {
	return &s.Response
}

func (s *SummaryCommand) GetCommand() string {
	return "summary"
}

type StatsCommand struct {
	Response StatsResponse
}

func (s *StatsCommand) GetResponse() any {
	return &s.Response
}

func (s *StatsCommand) GetCommand() string {
	return "stats"
}

func (s *StatsCommand) GetStats() Stats {
	if len(s.Response.Stats) >= 2 {
		return s.Response.Stats[1]
	}

	return Stats{}
}

type PoolsCommand struct {
	Response PoolsResponse
}

func (p *PoolsCommand) GetResponse() any {
	return &p.Response
}

func (p *PoolsCommand) GetCommand() string {
	return "pools"
}
