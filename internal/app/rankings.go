package app

import (
	"net/http"
	"sort"
	"time"
)

type modelRanking struct {
	Rank         int     `json:"rank"`
	PreviousRank int     `json:"previous_rank,omitempty"`
	Model        string  `json:"model_name"`
	Vendor       string  `json:"vendor"`
	Tokens       int64   `json:"total_tokens"`
	Share        float64 `json:"share"`
	Growth       float64 `json:"growth_pct"`
}
type vendorRanking struct {
	Rank        int     `json:"rank"`
	Vendor      string  `json:"vendor"`
	Tokens      int64   `json:"total_tokens"`
	Share       float64 `json:"share"`
	Growth      float64 `json:"growth_pct"`
	ModelsCount int     `json:"models_count"`
	TopModel    string  `json:"top_model"`
}
type rankingMover struct {
	Model       string  `json:"model_name"`
	Vendor      string  `json:"vendor"`
	RankDelta   int     `json:"rank_delta"`
	CurrentRank int     `json:"current_rank"`
	Growth      float64 `json:"growth_pct"`
}
type rankingTotals struct {
	model             string
	current, previous int64
}

// modelVendor preserves the default labels used by older callers and tests.
func modelVendor(model string) string {
	return providerForModel(model, []modelProvider{
		{Name: "OpenAI", Prefixes: []string{"gpt-", "o1", "o3", "o4"}, Priority: 10},
		{Name: "Anthropic", Prefixes: []string{"claude"}, Priority: 20},
		{Name: "Google", Prefixes: []string{"gemini"}, Priority: 30},
		{Name: "DeepSeek", Prefixes: []string{"deepseek"}, Priority: 40},
		{Name: "Alibaba", Prefixes: []string{"qwen", "qwq"}, Priority: 50},
	}).Name
}

func rankingDuration(period string) (time.Duration, bool) {
	switch period {
	case "today":
		return 24 * time.Hour, true
	case "week":
		return 7 * 24 * time.Hour, true
	case "month":
		return 30 * 24 * time.Hour, true
	case "year":
		return 365 * 24 * time.Hour, true
	default:
		return 0, false
	}
}

func growthPercent(current, previous int64) float64 {
	if previous == 0 {
		if current > 0 {
			return 100
		}
		return 0
	}
	return float64(current-previous) / float64(previous) * 100
}

func (s *Service) rankings(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		period = "week"
	}
	duration, ok := rankingDuration(period)
	if !ok {
		writeError(w, 400, "invalid_request", "period must be today, week, month, or year")
		return
	}
	now := time.Now().UTC()
	providers := s.providers(r)
	start := now.Add(-duration)
	previousStart := start.Add(-duration)
	rows, err := s.db.Query(r.Context(), `select model,coalesce(sum(prompt_tokens+completion_tokens) filter(where created_at >= $1),0),coalesce(sum(prompt_tokens+completion_tokens) filter(where created_at < $1),0) from usage_records where created_at >= $2 and created_at < $3 group by model`, start, previousStart, now)
	if err != nil {
		writeError(w, 500, "internal_error", "could not load rankings")
		return
	}
	defer rows.Close()
	totals := []rankingTotals{}
	for rows.Next() {
		var item rankingTotals
		if rows.Scan(&item.model, &item.current, &item.previous) == nil {
			totals = append(totals, item)
		}
	}
	if rows.Err() != nil {
		writeError(w, 500, "internal_error", "could not load rankings")
		return
	}
	previous := append([]rankingTotals(nil), totals...)
	sort.Slice(previous, func(i, j int) bool { return previous[i].previous > previous[j].previous })
	previousRanks := map[string]int{}
	rank := 0
	for _, item := range previous {
		if item.previous > 0 {
			rank++
			previousRanks[item.model] = rank
		}
	}
	sort.Slice(totals, func(i, j int) bool { return totals[i].current > totals[j].current })
	var allTokens int64
	for _, item := range totals {
		allTokens += item.current
	}
	models := []modelRanking{}
	for _, item := range totals {
		if item.current <= 0 || len(models) == 20 {
			continue
		}
		share := 0.0
		if allTokens > 0 {
			share = float64(item.current) / float64(allTokens)
		}
		models = append(models, modelRanking{len(models) + 1, previousRanks[item.model], item.model, providerForModel(item.model, providers).Name, item.current, share, growthPercent(item.current, item.previous)})
	}
	type vendorTotals struct {
		current, previous int64
		models            map[string]int64
	}
	byVendor := map[string]*vendorTotals{}
	for _, item := range totals {
		vendor := providerForModel(item.model, providers).Name
		if byVendor[vendor] == nil {
			byVendor[vendor] = &vendorTotals{models: map[string]int64{}}
		}
		byVendor[vendor].current += item.current
		byVendor[vendor].previous += item.previous
		if item.current > 0 {
			byVendor[vendor].models[item.model] = item.current
		}
	}
	vendors := []vendorRanking{}
	for name, item := range byVendor {
		if item.current <= 0 {
			continue
		}
		top := ""
		var topTokens int64
		for model, tokens := range item.models {
			if tokens > topTokens {
				top, topTokens = model, tokens
			}
		}
		vendors = append(vendors, vendorRanking{Vendor: name, Tokens: item.current, Share: float64(item.current) / float64(allTokens), Growth: growthPercent(item.current, item.previous), ModelsCount: len(item.models), TopModel: top})
	}
	sort.Slice(vendors, func(i, j int) bool { return vendors[i].Tokens > vendors[j].Tokens })
	for i := range vendors {
		vendors[i].Rank = i + 1
	}
	movers, droppers := []rankingMover{}, []rankingMover{}
	for _, item := range models {
		if item.PreviousRank == 0 {
			continue
		}
		delta := item.PreviousRank - item.Rank
		mover := rankingMover{item.Model, item.Vendor, delta, item.Rank, item.Growth}
		if delta > 0 {
			movers = append(movers, mover)
		} else if delta < 0 {
			droppers = append(droppers, mover)
		}
	}
	sort.Slice(movers, func(i, j int) bool { return movers[i].RankDelta > movers[j].RankDelta })
	sort.Slice(droppers, func(i, j int) bool { return droppers[i].RankDelta < droppers[j].RankDelta })
	if len(movers) > 6 {
		movers = movers[:6]
	}
	if len(droppers) > 6 {
		droppers = droppers[:6]
	}
	writeJSON(w, 200, map[string]any{"period": period, "models": models, "vendors": vendors, "top_movers": movers, "top_droppers": droppers, "total_tokens": allTokens, "updated_at": now})
}
