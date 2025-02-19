package device

type statusResponse struct {
	Status      string `json:"STATUS"`
	When        int    `json:"When"`
	Code        int    `json:"Code"`
	Msg         string `json:"Msg"`
	Description string `json:"Description"`
}

type VersionResponse struct {
	Status  []statusResponse `json:"STATUS"`
	Version []struct {
		BMMiner     string `json:"BMMiner"`
		API         string `json:"API"`
		Miner       string `json:"Miner"`
		CompileTime string `json:"CompileTime"`
		Type        string `json:"Type"`
	} `json:"VERSION"`
	Id int `json:"id"`
}

type SummaryResponse struct {
	Status  []statusResponse `json:"STATUS"`
	Summary []struct {
		Elapsed            int     `json:"Elapsed"`
		Ghs5s              float32 `json:"GHS 5s"`
		GhsAv              float32 `json:"GHS av"`
		Ghs30m             float32 `json:"GHS 30m"`
		FoundBlocks        int     `json:"Found Blocks"`
		Getwork            int     `json:"Getwork"`
		Accepted           int     `json:"Accepted"`
		Rejected           int     `json:"Rejected"`
		HardwareErrors     int     `json:"Hardware Errors"`
		Utility            float32 `json:"Utility"`
		Discarded          int     `json:"Discarded"`
		Stale              int     `json:"Stale"`
		GetFailures        int     `json:"Get Failures"`
		LocalWork          int     `json:"Local Work"`
		RemoteFailures     int     `json:"Remote Failures"`
		NetworkBlocks      int     `json:"Network Blocks"`
		TotalMH            float32 `json:"Total MH"`
		WorkUtility        float32 `json:"Work Utility"`
		DifficultyAccepted float32 `json:"Difficulty Accepted"`
		DifficultyRejected float32 `json:"Difficulty Rejected"`
		DifficultyStale    float32 `json:"Difficulty Stale"`
		BestShare          int     `json:"Best Share"`
		DeviceHardware     float32 `json:"Device Hardware"`
		DeviceRejected     float32 `json:"Device Rejected"`
		PoolRejected       float32 `json:"Pool Rejected"`
		PoolStale          float32 `json:"Pool Stale"`
		LastGetwork        int     `json:"Last Getwork"`
		RtHashrate         string  `json:"RT HASHRATE"`
		AvHashrate         string  `json:"AV HASHRATE"`
		TheoryHashrate     string  `json:"THEORY HASHRATE"`
	} `json:"SUMMARY"`
	Id int `json:"id"`
}

type StatsResponse struct {
	Status []statusResponse `json:"STATUS"`
	Stats  []struct {
		Stats          int     `json:"STATS"`
		Id             string  `json:"ID"`
		Elapsed        int     `json:"Elapsed"`
		Ghs5s          float32 `json:"GHS 5s"`
		GhsAv          float32 `json:"GHS av"`
		Rate30m        float32 `json:"rate_30m"`
		TotalRateIdeal float32 `json:"total_rateideal"`
		Mode           int     `json:"Mode"`
		MinerCount     int     `json:"miner_count"`
		Frequency      int     `json:"frequency"`
		FanNum         int     `json:"fan_num"`
		Fan1           int     `json:"fan1"`
		Fan2           int     `json:"fan2"`
		Fan3           int     `json:"fan3"`
		Fan4           int     `json:"fan4"`
		TempNum        int     `json:"temp_num"`
		Temp1          int     `json:"temp1"`
		Temp21         int     `json:"temp2_1"`
		Temp2          int     `json:"temp2"`
		Temp22         int     `json:"temp2_2"`
		Temp3          int     `json:"temp3"`
		Temp23         int     `json:"temp2_3"`
		TempInPcb1     string  `json:"temp_in_pcb_1"`
		TempInPcb2     string  `json:"temp_in_pcb_2"`
		TempInPcb3     string  `json:"temp_in_pcb_3"`
		TempOutPcb1    string  `json:"temp_out_pcb_1"`
		TempOutPcb2    string  `json:"temp_out_pcb_2"`
		TempOutPcb3    string  `json:"temp_out_pcb_3"`
		TempInChip1    string  `json:"temp_in_chip_1"`
		TempInChip2    string  `json:"temp_in_chip_2"`
		TempInChip3    string  `json:"temp_in_chip_3"`
		TempOutChip1   string  `json:"temp_out_chip_1"`
		TempOutChip2   string  `json:"temp_out_chip_2"`
		TempOutChip3   string  `json:"temp_out_chip_3"`
		TempMax        int     `json:"temp_max"`
		ChainAcn1      int     `json:"chain_acn1"`
		ChainAcn2      int     `json:"chain_acn2"`
		ChainAcn3      int     `json:"chain_acn3"`
		ChainAcs1      string  `json:"chain_acs1"`
		ChainAcs2      string  `json:"chain_acs2"`
		ChainAcs3      string  `json:"chain_acs3"`
		ChainHw1       int     `json:"chain_hw1"`
		ChainHw2       int     `json:"chain_hw2"`
		ChainHw3       int     `json:"chain_hw3"`
		NoMatchingWork int     `json:"no_matching_work"`
		//ChainRate1        string  `json:"chain_rate1"`
		//ChainRate2        string  `json:"chain_rate2"`
		//ChainRate3        string  `json:"chain_rate3"`
		ChainAvgHashrate1 string `json:"CHAIN AVG HASHRATE1"`
		ChainAvgHashrate2 string `json:"CHAIN AVG HASHRATE2"`
		ChainAvgHashrate3 string `json:"CHAIN AVG HASHRATE3"`
		Freq1             int    `json:"freq1"`
		Freq2             int    `json:"freq2"`
		Freq3             int    `json:"freq3"`
		MinerVersion      string `json:"miner_version"`
	} `json:"STATS"`
	Id int `json:"id"`
}

type PoolsResponse struct {
	Status []statusResponse `json:"STATUS"`
	Pools  []struct {
		Pool                int     `json:"POOL"`
		Url                 string  `json:"URL"`
		Status              string  `json:"Status"`
		Priority            int     `json:"Priority"`
		Quota               int     `json:"Quota"`
		LongPoll            string  `json:"LongPoll"`
		Gateworks           int     `json:"Gateworks"`
		Accepted            int     `json:"Accepted"`
		Rejected            int     `json:"Rejected"`
		Discarded           int     `json:"Discarded"`
		Stale               int     `json:"Stale"`
		GetFailures         int     `json:"Get Failures"`
		RemoteFailures      int     `json:"Remote Failures"`
		User                string  `json:"User"`
		LastShareTime       string  `json:"Last Share Time"`
		Diff                string  `json:"Diff"`
		Diff1Shares         int     `json:"Diff1 Shares"`
		ProxyType           string  `json:"Proxy Type"`
		Proxy               string  `json:"Proxy"`
		DifficultyAccepted  float64 `json:"Difficulty Accepted"`
		DifficultyRejected  float64 `json:"Difficulty Rejected"`
		DifficultyStale     float64 `json:"Difficulty Stale"`
		LastShareDifficulty float64 `json:"Last Share Difficulty"`
		HasStratum          bool    `json:"Has Stratum"`
		StratumActive       bool    `json:"Stratum Active"`
		StratumUrl          string  `json:"Stratum URL"`
		HasGBT              bool    `json:"Has GBT"`
		BestShare           float32 `json:"Best Share"`
		PoolRejected        float32 `json:"Pool Rejected%"`
		PoolStale           float32 `json:"Pool Stale%%"`
	} `json:"POOLS"`
	Id int `json:"id"`
}
