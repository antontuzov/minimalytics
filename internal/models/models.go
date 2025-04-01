package models

type (
	DailyStat struct {
		Day   string `json:"day"`
		Count int    `json:"count"`
	}

	UniqueVisitStat struct {
		Day   string `json:"day"`
		Count int    `json:"count"`
	}

	PageStat struct {
		Path  string `json:"path"`
		Count int    `json:"count"`
	}

	ReferrerStat struct {
		Referrer string `json:"referrer"`
		Count    int    `json:"count"`
	}

	DeviceStat struct {
		Device string `json:"device"`
		Count  int    `json:"count"`
	}

	BrowserStat struct {
		Browser string `json:"browser"`
		Count   int    `json:"count"`
	}

	PageView struct {
		Timestamp string `json:"timestamp"`
		Path      string `json:"path"`
		Referrer  string `json:"referrer"`
		UserAgent string `json:"user_agent"`
		IP        string `json:"ip"`
	}
)
