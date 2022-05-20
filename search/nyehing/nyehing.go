package nyehing

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/LoperLee/golang-hangul-toolkit/hangul"
)

func Errchk(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func array() []string {
	var hanword []string
	var a = []string{"ㄱ", "ㄴ", "ㄷ", "ㄹ", "ㅁ", "ㅂ", "ㅅ", "ㅇ", "ㅈ", "ㅊ", "ㅋ", "ㅌ", "ㅍ", "ㅎ"}
	var b = []string{"ㅏ", "ㅑ", "ㅐ", "ㅔ", "ㅗ", "ㅛ", "ㅜ", "ㅠ", "ㅡ", "ㅣ", "ㅏ", "ㅓ", "ㅗ", "ㅜ", "ㅡ", "ㅣ"}
	var c = []string{"ㄱ", "ㄴ", "ㄹ", "ㅁ", "ㅂ", "ㅅ", "ㅇ", "ㅈ", "", "", "", "", ""}

	for i := 0; i < 3; i++ {
		rand.Seed(time.Now().UnixNano())
		r1 := rand.Intn(len(a))
		r2 := rand.Intn(len(b))
		r3 := rand.Intn(len(c))
		switch i {
		case 0:
			hanword = append(hanword, a[r1])
		case 1:
			hanword = append(hanword, b[r2])
		case 2:
			hanword = append(hanword, c[r3])
		}
	}
	return hanword
}

func combine(n int) []string {
	var gul hangul.Hangul
	var guls []string

	for i := 0; i < n; i++ {
		for {
			same := true
			hanword := array()
			if hanword[2] == "" {
				gul = hangul.Hangul{Chosung: hanword[0], Jungsung: hanword[1]}
				err := hangul.CombineHangul(&gul)
				if err != nil {
					panic(err)
				}
			} else {
				gul = hangul.Hangul{Chosung: hanword[0], Jungsung: hanword[1], Jongsung: hanword[2]}
				err := hangul.CombineHangul(&gul)
				if err != nil {
					panic(err)
				}
			}
			for i := 0; i < len(guls); i++ {
				if guls[i] == gul.Word {
					same = false
				}
			}
			if same {
				break
			}
		}
		guls = append(guls, gul.Word)
	}
	return guls
}

func wordlen(strslice []string) string {

	var name string
	for _, namesl := range strslice {
		name += namesl
	}
	return name

}

type Nonames struct {
	No   string
	Name string
	//Maple string
	OPGG string
}

func Mkneyhings(n int) []Nonames {
	var namechart []Nonames

	for i := 0; i < 8; i++ {
		noname := Nonames{
			No:   fmt.Sprintf("%02d", i+1),
			Name: wordlen(combine(n)),
			//Maple: "",
			OPGG: "",
		}
		namechart = append(namechart, noname)
	}
	return namechart
}

func Writefile(n []Nonames) {
	const Filename string = "저잉녜힝.csv"
	file, err := os.Create(Filename)
	Errchk(err)
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()
	header := []string{"No", "NyeHing" /*"MapleHome", */, "OPGG"}
	Errchk(w.Write(header))
	for _, idx := range n {
		var temp []string = []string{idx.No, idx.Name /* idx.Maple,*/, idx.OPGG}
		Errchk(w.Write(temp))
	}
	var time []string = []string{"작성 시간", fmt.Sprintf("%02d:%02d:%02d", time.Now().Hour(), time.Now().Minute(), time.Now().Second())}
	Errchk(w.Write(time))
}
