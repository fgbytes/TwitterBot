/*
 *   Copyright 2015 Benoit LETONDOR
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package db

import (
	"database/sql"
	"github.com/lib/pq"
	"time"
)

type Follow struct {
	id           int
	UserId       int64
	UserName     string
	TweetId      int64
	Status       string
	FollowDate   time.Time
	UnfollowDate time.Time
	LastAction   time.Time
}

const (
	_TABLE_FOLLOW = "follow"
)

func AlreadyFollow(userId int64) (bool, error) {
	stmtOut, err := database.Prepare("SELECT count(*) FROM " + _TABLE_FOLLOW + " WHERE userId = $1 LIMIT 1")
	if err != nil {
		return true, err
	}

	defer stmtOut.Close()

	var size int

	err = stmtOut.QueryRow(userId).Scan(&size)
	if err != nil {
		return true, err
	}

	return size > 0, nil
}

func (follow Follow) Persist() error {
	var stmtIns *sql.Stmt
	var err error

	if follow.id == 0 {
		stmtIns, err = database.Prepare("INSERT INTO " + _TABLE_FOLLOW + "(userId, userName, status, followDate, unfollowDate, lastAction) VALUES($1, $2, $3, $4, $5, $6)")
	} else {
		stmtIns, err = database.Prepare("UPDATE " + _TABLE_FOLLOW + " SET userId = $1, userName = $2, status = $3, followDate = $4, unfollowDate = $5, lastAction = $6 WHERE id = $7")
	}

	if err != nil {
		return err
	}

	defer stmtIns.Close()

	unfollowDate := pq.NullTime{Time: follow.UnfollowDate, Valid: !follow.UnfollowDate.IsZero()}

	if follow.id == 0 {
		_, err = stmtIns.Exec(follow.UserId, follow.UserName, follow.Status, follow.FollowDate, unfollowDate, time.Now())
	} else {
		_, err = stmtIns.Exec(follow.UserId, follow.UserName, follow.Status, follow.FollowDate, unfollowDate, follow.LastAction, follow.id)
	}

	return err
}

func GetNotUnfollowed(maxFollowDate time.Time, limit int) ([]Follow, error) {
	follows := make([]Follow, 0)

	stmtOut, err := database.Prepare("SELECT * FROM " + _TABLE_FOLLOW + " WHERE unfollowDate IS NULL AND followDate <= $1 ORDER BY lastAction LIMIT $2")
	if err != nil {
		return follows, err
	}

	defer stmtOut.Close()

	rows, err := stmtOut.Query(maxFollowDate, limit)
	if err != nil {
		return follows, err
	}

	defer rows.Close()

	for rows.Next() {
		follow, err := mapFollow(rows)
		if err != nil {
			return follows, err
		}

		follows = append(follows, follow)
	}

	return follows, nil
}

func mapFollow(rows *sql.Rows) (Follow, error) {
	var id int
	var userId int64
	var userName string
	var status string
	var followDate time.Time
	var unfollowDate pq.NullTime
	var lastAction time.Time

	err := rows.Scan(&id, &userId, &userName, &status, &followDate, &unfollowDate, &lastAction)
	if err != nil {
		return Follow{}, err
	}

	var unfollowTime time.Time
	if unfollowDate.Valid {
		unfollowTime = unfollowDate.Time
	}

	return Follow{id: id, UserId: userId, UserName: userName, Status: status, FollowDate: followDate, UnfollowDate: unfollowTime, LastAction: lastAction}, nil
}
