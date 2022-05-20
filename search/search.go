package search

import (
	"Naming-Service/search/nyehing"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type SearchUserdata struct {
	Img    string
	Server string
	Name   string
	Says   string
}

const (
	opggurl      string = "https://maple.gg/search?q="
	mapleURL     string = "https://maplestory.nexon.com/Ranking/World/Total?c="
	mapleURLtail string = "&w=0"
)

func respheader(url string) (resp *http.Response) {
	Client := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := Client.Do(req)
	nyehing.Errchk(err)
	return resp
}

var wg sync.WaitGroup
var nameCh = make(chan nyehing.Nonames, 8)

func MapleGG(n int) []nyehing.Nonames {
	searchname := nyehing.Mkneyhings(n)
	wg.Add(8)
	for _, name := range searchname {
		go opggExist(name)
	}
	wg.Wait()
	for i := 0; i < 8; i++ {
		searchname[i] = <-nameCh
	}
	return searchname
}

func opggExist(searchname nyehing.Nonames) {
	SerchingOpgg := opggurl + searchname.Name
	resp := respheader(SerchingOpgg)
	time.Sleep(2 * time.Second)
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	nyehing.Errchk(err)
	_, exist := doc.Find(".user-profile").Attr("data-nick")
	failText := doc.Find(".search-result__item").Text()
	if strings.Contains(failText, searchname.Name) {
		exist = true
	}
	if exist {
		searchname.OPGG = "Not available"
	} else {
		if strings.Contains(failText, "캐릭터 이름을 다시 한 번 확인해주세요.") || strings.Contains(failText, "갱신에 실패했습니다. (캐릭터를 찾을 수 없습니다.)") {
			searchname.OPGG = "Available !"
		} else {
			searchname.OPGG = "Not available"
		}
	}
	nameCh <- searchname
	wg.Done()
}

// func mapleExist(searchname nyehing.Nonames) {
// 	SerchingMaple := mapleURL + searchname.Name + mapleURLtail
// 	resp := respheader(SerchingMaple)
// 	time.Sleep(1 * time.Second)
// 	doc, err := goquery.NewDocumentFromReader(resp.Body)
// 	nyehing.Errchk(err)
// 	names := doc.Find(".rank_table_wrap").Find("tbody").Text()
// 	var exist bool
// 	if strings.Contains(names, searchname.Name) {
// 		exist = true
// 	}
// 	if exist {
// 		searchname.OPGG = "Not available"
// 	} else {
// 		searchname.OPGG = "Available !"
// 	}
// 	nameCh <- searchname
// 	wg.Done()
// }

func GetUserInfo(userName string) SearchUserdata {
	var user SearchUserdata

	userURL := "https://maple.gg/u/" + userName
	resp := respheader(userURL)
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	nyehing.Errchk(err)
	img, _ := doc.Find(".character-avatar-row").Find("img").Attr("src")
	server, _ := doc.Find(".container").Find("h3").Find("img").Attr("src")
	name := doc.Find(".container").Find("h3").Find("b").Text()
	user = SearchUserdata{img, server, name, ""}
	return user
}
