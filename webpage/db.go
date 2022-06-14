package webpage

import (
	"Naming-Service/mnameutil"
	"Naming-Service/search"
	"bytes"
	"encoding/gob"
	"errors"
	"strconv"

	bolt "go.etcd.io/bbolt"
)

const (
	dbName      = "webpage/dashBoard.db"
	linebucket  = "line"
	boardbucket = "board"
	lastLine    = "lastLine"
)

var (
	db             *bolt.DB
	lineNo         int
	wholeDashboard []OneLine
)

func DB() {
	if db == nil {
		dbPointer, err := bolt.Open(dbName, 0600, nil)
		db = dbPointer
		mnameutil.Errchk(err)
		err = db.Update(func(t *bolt.Tx) error {
			_, err := t.CreateBucketIfNotExists([]byte(linebucket))
			mnameutil.Errchk(err)
			_, err = t.CreateBucketIfNotExists([]byte(boardbucket))
			return err
		})
		mnameutil.Errchk(err)
		lineNo, err = strconv.Atoi(getNewestLine().No)
		mnameutil.Errchk(err)
		if lineNo > 0 {
			getDashboard(lineNo)
		}
	}
}

func Close() {
	db.Close()
}

func SaveUserSays(usersay OneLine) {
	lineNo++
	usersay.No = strconv.Itoa(lineNo)

	err := db.Update(func(t *bolt.Tx) error {
		lineB := t.Bucket([]byte(linebucket))
		err := lineB.Put([]byte(lastLine), mnameutil.Encode(usersay))
		mnameutil.Errchk(err)
		boardB := t.Bucket([]byte(boardbucket))
		return boardB.Put([]byte(usersay.No), mnameutil.Encode(usersay))
	})
	mnameutil.Errchk(err)
	appendNewLine(lineNo)
}

func (d *DataStruct) getAsideUsers() {
	leftUser := getNewestLine()
	d.LeftUser = search.GetUserInfo(leftUser.User)
	d.LeftUser.Says = leftUser.Says

	leftUserNum, err := strconv.Atoi(leftUser.No)
	mnameutil.Errchk(err)
	rightUserKey := strconv.Itoa(leftUserNum - 1)

	rightUser, err := getSpecificLine(rightUserKey)
	if err != nil {
		rightUser = leftUser
	}

	d.RightUser = search.GetUserInfo(rightUser.User)
	d.RightUser.Says = rightUser.Says
}

func getNewestLine() OneLine {
	var lastlineAsByte []byte
	err := db.View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(linebucket))
		lastlineAsByte = bucket.Get([]byte(lastLine))
		if lastlineAsByte == nil {
			return errors.New("lastline dosen't exist")
		}
		return nil
	})

	if err != nil {
		return OneLine{No: "0"}
	}
	return decode(lastlineAsByte)
}

func getSpecificLine(keyString string) (OneLine, error) {
	var dataAsByte []byte

	err := db.View(func(t *bolt.Tx) error {
		lineB := t.Bucket([]byte(boardbucket))
		dataAsByte = lineB.Get([]byte(keyString))
		if dataAsByte == nil {
			return errors.New(keyString + "User dosen't exist")
		}
		return nil
	})
	if err != nil {
		return OneLine{}, err
	}
	return decode(dataAsByte), nil
}

func getDashboard(lastnum int) {
	err := db.View(func(t *bolt.Tx) error {
		boardB := t.Bucket([]byte(boardbucket))
		for i := lastnum; i > 0; i-- {
			onelineAsByte := boardB.Get([]byte(strconv.Itoa(i)))
			foundOneLine := decode(onelineAsByte)
			wholeDashboard = append(wholeDashboard, foundOneLine)
		}
		return nil
	})
	mnameutil.Errchk(err)
}

func appendNewLine(lastnum int) {
	err := db.View(func(t *bolt.Tx) error {
		boardB := t.Bucket([]byte(boardbucket))
		onelineAsByte := boardB.Get([]byte(strconv.Itoa(lastnum)))
		foundOneLine := decode(onelineAsByte)
		wholeDashboard = append([]OneLine{foundOneLine}, wholeDashboard...)
		return nil
	})
	mnameutil.Errchk(err)
}

func decode(databytes []byte) OneLine {
	var line OneLine
	var buffer bytes.Buffer
	_, err := buffer.Write(databytes)
	mnameutil.Errchk(err)
	decoder := gob.NewDecoder(&buffer)
	err = decoder.Decode(&line)
	mnameutil.Errchk(err)
	return line
}
