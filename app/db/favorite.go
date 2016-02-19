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
	//"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"time"
)

type Favorite struct {
	id         int
	UserId     int64
	UserName   string
	TweetId    int64
	Status     string
	FavDate    time.Time
	UnfavDate  time.Time
	LastAction time.Time
}

const (
	_TABLE_FAVORITE = "favorite"
)

func (fav Favorite) Persist() error {
	var stmtIns *sql.Stmt
	var err error

	if fav.id == 0 {
		stmtIns, err = database.Prepare("INSERT INTO " + _TABLE_FAVORITE + "(userId, userName, tweetId, status, favDate, unfavDate, lastAction) VALUES( $1, $2, $3, $4, $5, $6, $7 )")
	} else {
		stmtIns, err = database.Prepare("UPDATE " + _TABLE_FAVORITE + " SET userId = $1, userName = $2, tweetId = $3, status = $4, favDate = $5, unfavDate = $6, lastAction = $7 WHERE id = $8")
	}

	if err != nil {
		return err
	}

	defer stmtIns.Close()

	unfavDate := pq.NullTime{Time: fav.UnfavDate, Valid: !fav.UnfavDate.IsZero()}

	if fav.id == 0 {
		_, err = stmtIns.Exec(fav.UserId, fav.UserName, fav.TweetId, fav.Status, fav.FavDate, unfavDate, time.Now())
	} else {
		_, err = stmtIns.Exec(fav.UserId, fav.UserName, fav.TweetId, fav.Status, fav.FavDate, unfavDate, fav.LastAction, fav.id)
	}

	return err
}

func HasAlreadyFav(tweetId int64) (bool, error) {
	stmtOut, err := database.Prepare("SELECT count(*) FROM " + _TABLE_FAVORITE + " WHERE tweetId = $1 LIMIT 1")
	if err != nil {
		return true, err
	}

	defer stmtOut.Close()

	var size int

	err = stmtOut.QueryRow(tweetId).Scan(&size)
	if err != nil {
		return true, err
	}

	return size > 0, nil
}

func GetNotUnfavorite(maxFavDate time.Time, limit int) ([]Favorite, error) {
	favs := make([]Favorite, 0)

	stmtOut, err := database.Prepare("SELECT * FROM " + _TABLE_FAVORITE + " WHERE unfavDate IS NULL AND favDate <= $1 ORDER BY lastAction LIMIT $2")
	if err != nil {
		return favs, err
	}

	defer stmtOut.Close()

	rows, err := stmtOut.Query(maxFavDate, limit)
	if err != nil {
		return favs, err
	}

	defer rows.Close()

	for rows.Next() {
		fav, err := mapFav(rows)
		if err != nil {
			return favs, err
		}

		favs = append(favs, fav)
	}

	return favs, nil
}

func mapFav(rows *sql.Rows) (Favorite, error) {
	var id int
	var userId int64
	var userName string
	var tweetId int64
	var status string
	var favDate time.Time
	var unfavDate pq.NullTime
	var lastAction time.Time

	err := rows.Scan(&id, &userId, &userName, &tweetId, &status, &favDate, &unfavDate, &lastAction)
	if err != nil {
		return Favorite{}, err
	}

	var unfavTime time.Time
	if unfavDate.Valid {
		unfavTime = unfavDate.Time
	}

	return Favorite{id: id, UserId: userId, UserName: userName, TweetId: tweetId, Status: status, FavDate: favDate, UnfavDate: unfavTime, LastAction: lastAction}, nil
}
