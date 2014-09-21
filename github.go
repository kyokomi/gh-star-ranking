package ghstar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"appengine"
	"appengine/urlfetch"
)

const (
	gitHubStarURL = "https://api.github.com/search/repositories?q=language:%s&sort=star&order=desc"
	gitHubUserURL = "https://github.com/%s"

	snapshotKey = "lang:%s_data:%s"
)

// Ranking struct GitHub Star Ranking.
type Ranking struct {
	Item
	Rank     int64
	LastRank int64
	//	FullName string
	//	StargazersCount int64
	LastStargazersCount int64
	UserURL             string
	//	OwnerAvatarURL string
	//	HtmlURL string
	//	UpdatedAt       string
	//	CreatedAt       string
}

// parse time default.
func (r Ranking) ParseTime(at string) time.Time {
	t, _ := time.Parse(time.RFC3339, at)
	return t
}

// format time default.
func (r Ranking) FormatTime(at string) string {
	return r.ParseTime(at).Format("2006-01-02")
}

// ResponseRanking struct GitHub star Ranking for each Language.
type ResponseRanking struct {
	Language string
	Rankings []Ranking
}

func createSnapshotKey(lang string, t time.Time) string {
	return fmt.Sprintf(snapshotKey, lang, t.Format("2006-01-02"))
}

func readGitHubStarRanking(c appengine.Context, lang string) ([]Ranking, error) {

	client := urlfetch.Client(c)
	res, err := client.Get(fmt.Sprintf(gitHubStarURL, lang))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// parse
	var reps Repositories
	if err := json.Unmarshal(data, &reps); err != nil {
		return nil, err
	}

	return newRankings(&reps), nil
}

func newRankings(reps *Repositories) []Ranking {

	rankings := make([]Ranking, len(reps.Items))
	for idx, item := range reps.Items {
		var lastRank int64
		var lastStargazersCount int64

		rankings[idx] = Ranking{
			Item:                item,
			Rank:                int64(idx + 1),
			LastRank:            lastRank,
			LastStargazersCount: lastStargazersCount,
			UserURL:             fmt.Sprintf(gitHubUserURL, item.Owner.Login),
		}
	}
	return rankings
}
