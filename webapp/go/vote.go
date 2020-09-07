package main

import (
	"fmt"
	"strings"
)

// Vote Model
type Vote struct {
	ID          int
	UserID      int
	CandidateID int
	Keyword     string
}

func getVoteCountByCandidateID(candidateID int) (count int) {
	row := db.QueryRow("SELECT COUNT(*) AS count FROM votes WHERE candidate_id = ?", candidateID)
	row.Scan(&count)
	return
}

func getVoteCountByPartyName(partyName string) (count int) {
	row := db.QueryRow("SELECT COUNT(*) AS count FROM votes WHERE political_party = ?", partyName)
	row.Scan(&count)
	return
}

func getUserVotedCount(userID int) (count int) {
	row := db.QueryRow("SELECT COUNT(*) AS count FROM votes WHERE user_id =  ?", userID)
	row.Scan(&count)
	return
}

func createVote(userID int, candidateID int, keyword string, politicalParty string, count int) {
	if count == 1 {
		db.Exec("INSERT INTO votes (user_id, candidate_id, keyword, political_party) VALUES (?, ?, ?, ?)",
			userID, candidateID, keyword, politicalParty)
	} else {
		valueArgs := []interface{}{}
		valueStrings := []string{}

		query := "INSERT INTO votes (user_id, candidate_id, keyword, political_party)"

		for i := 0; i < count; i++ {
			valueStrings = append(valueStrings, "(?, ?, ?, ?)")

			valueArgs = append(valueArgs, userID)
			valueArgs = append(valueArgs, candidateID)
			valueArgs = append(valueArgs, keyword)
			valueArgs = append(valueArgs, politicalParty)
		}

		query = fmt.Sprintf(query, strings.Join(valueStrings, ","))

		db.Exec(query, valueArgs)
	}
}

func getVoiceOfSupporter(candidateIDs []int) (voices []string) {
	rows, err := db.Query(`
    SELECT keyword
    FROM votes
    WHERE candidate_id IN (?)
    GROUP BY keyword
    ORDER BY COUNT(*) DESC
    LIMIT 10`)
	if err != nil {
		return nil
	}

	defer rows.Close()
	for rows.Next() {
		var keyword string
		err = rows.Scan(&keyword)
		if err != nil {
			panic(err.Error())
		}
		voices = append(voices, keyword)
	}
	return
}
