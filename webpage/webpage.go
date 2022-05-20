package webpage

import (
	"Naming-Service/mnameutil"
	"Naming-Service/search"
	"Naming-Service/search/nyehing"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
)

const (
	port1       string = ":80"
	port2       string = ":443"
	onelineList int    = 10
	robots      string = "robots.txt"
	sitemap     string = "sitemap.xml"
)

var (
	robotTemplates   *template.Template
	sitemapTemplates *template.Template
	Templates        *template.Template
	data             DataStruct
)

type DataStruct struct {
	Result    []nyehing.Nonames
	LeftUser  search.SearchUserdata
	RightUser search.SearchUserdata
	Dashboard []OneLine
	Pagenum   []int
}

type OneLine struct {
	No   string
	User string
	Says string
}

func home(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := Templates.ExecuteTemplate(rw, "home.html", nil)
		if err != nil {
			fmt.Println(err)
		}

	case "POST":
		r.ParseForm()
		namelenstr := r.Form.Get("nameLen")
		namelen, err := strconv.Atoi(namelenstr)
		if err != nil {
			log.Panic(err)
		}
		data.Result = search.MapleGG(namelen)
		http.Redirect(rw, r, "/mname", http.StatusPermanentRedirect)
	}

}

func mName(rw http.ResponseWriter, r *http.Request) {
	data.getAsideUsers()

	switch r.Method {
	case "GET":
		err := Templates.ExecuteTemplate(rw, "mname.html", data)
		if err != nil {
			fmt.Println(err)
		}
	case "POST":
		err := Templates.ExecuteTemplate(rw, "mname.html", data)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func promotePage(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		var totalpages int
		vars := mux.Vars(r)
		pagenum, err := strconv.Atoi(vars["page"])
		mnameutil.Errchk(err)
		totalpages, data.Dashboard = dashboardSlicing(wholeDashboard, pagenum)
		data.Pagenum = make([]int, totalpages)
		for i := totalpages; i > 0; i-- {
			data.Pagenum[i-1] = i
		}

		err = Templates.ExecuteTemplate(rw, "promote.html", data)
		if err != nil {
			fmt.Println(err)
		}
	case "POST":
		r.ParseForm()
		user := r.Form.Get("user")
		says := r.Form.Get("says")

		newestLine := OneLine{
			User: user,
			Says: says,
		}
		SaveUserSays(newestLine)
		r.Method = http.MethodGet
		http.Redirect(rw, r, "/mname", http.StatusPermanentRedirect)
	}
}

func dashboardSlicing(wholedashBoard []OneLine, reqPage int) (totalpage int, onePageDashboad []OneLine) {
	wholePages := int(len(wholedashBoard)/onelineList) + 1
	if reqPage > wholePages {
		reqPage = wholePages
	}
	if len(wholedashBoard[(reqPage-1)*onelineList:]) < onelineList {
		startidx := (reqPage - 1) * onelineList
		return wholePages, wholedashBoard[startidx:]

	} else {
		startidx, endidx := (reqPage-1)*onelineList, (reqPage * onelineList)
		return wholePages, wholedashBoard[startidx:endidx]
	}
}

func robottxt(rw http.ResponseWriter, r *http.Request) {
	err := robotTemplates.ExecuteTemplate(rw, robots, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func sitemapxml(rw http.ResponseWriter, r *http.Request) {
	err := sitemapTemplates.ExecuteTemplate(rw, sitemap, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func Start() {
	router := mux.NewRouter()
	Templates = template.Must(template.ParseGlob("*.html"))
	robotTemplates = template.Must(template.ParseGlob(robots))
	sitemapTemplates = template.Must(template.ParseGlob(sitemap))
	fmt.Printf("Listening on http://localhost%s\n", port1)
	fmt.Printf("Listening on http://localhost%s\n", port2)
	router.HandleFunc("/", home).Methods("GET", "POST")
	router.HandleFunc("/mname", mName).Methods("GET", "POST")
	router.HandleFunc("/promote/{page:[0-9]+}", promotePage).Methods("GET", "POST")
	router.HandleFunc("/robots.txt", robottxt)
	router.HandleFunc("/sitemap.xml", sitemapxml)
	go http.ListenAndServe(port1, router)
	http.ListenAndServe(port2, router)
}
