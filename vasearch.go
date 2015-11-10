package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"gopkg.in/olivere/elastic.v2"
	"github.com/facebookgo/inject"
	"vasearch/model"
	"github.com/gorilla/mux"
	"github.com/BurntSushi/toml"
	"runtime"
	"log"
)

type Sample struct {
	Name string `json:"name"`
	Variants int `json:"variants"`
	Genome string `json:"genome"`
}

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func ViewHandler(w http.ResponseWriter, r *http.Request, ctx model.AppLoader) {
	//title := r.URL.Path[len("/view/"):]
	vars := mux.Vars(r)
	title := vars["title"]

	log.Println(title)
	p, _ := loadPage(title)
	fmt.Fprintf(w, "<h1>%s</h1>bla2<div>%s</div>", p.Title, p.Body)
}

func sampleHandler(w http.ResponseWriter, r *http.Request, ctx model.AppLoader) {
	m := Sample{"Liver Test", 1024, "hg18"}

	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	fmt.Printf("client=%v", ctx.Client)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, string(b))
}

func fakeResponseHandler(w http.ResponseWriter, r *http.Request, ctx model.AppLoader) {
	fmt.Printf("fakeResponse\n")

	data, err := ioutil.ReadFile("search_response.json")
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, string(data))
}

func searchHandler(w http.ResponseWriter, r *http.Request, ctx model.AppLoader) {
	path := r.URL.Path[1:]
	log.Println(path)

	// curl -X GET "http://localhost:9200/vaindex/sample/292175?pretty=true"
	//fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "Search", "TODO")

	//searchTerm := r.URL.Path[len("/search/"):]
	vars := mux.Vars(r)
	searchTerm := vars["searchTerm"]
	fmt.Printf("Search Term: [%s]", searchTerm)
	termQuery := elastic.NewQueryStringQuery(searchTerm).DefaultOperator("AND")

	// Facets
	genderFacet := elastic.NewTermsFacet().Field("gender")
	genomeFacet := elastic.NewTermsFacet().Field("genome")
	ethnicityFacet := elastic.NewTermsFacet().Field("annotations.ethnicity")

	// Specify highlighter
	hl := elastic.NewHighlight()
	hl = hl.Fields(elastic.NewHighlighterField("*"))
	hl = hl.PreTags("<mark>").PostTags("</mark>")

	searchResult, err := ctx.Client.Search().
	    Index(ctx.Config.Index).   // search in index "vaindex"
	    Type("sample").
	    Query(&termQuery).  // specify the query
		Facet("gender", genderFacet).
		Facet("genome", genomeFacet).
		Facet("ethnicity", ethnicityFacet).
	    Highlight(hl).
	    Sort("name", true). // sort by "name" field, ascending
	    From(0).Size(10).   // take documents 0-9
	    Pretty(true).       // pretty print request and response JSON
	    Do()                // execute
	if err != nil {
	    // Handle error
	    panic(err)
	}
	
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)
	fmt.Printf("Found a total of %d samples\n", searchResult.TotalHits())
	
	if searchResult.Hits != nil {
	    fmt.Printf("Found a total of %d samples\n", searchResult.Hits.TotalHits)

	    // Iterate through results
	    for _, hit := range searchResult.Hits.Hits {
	        // hit.Index contains the name of the index
	
	        // Deserialize hit.Source into a Sample
	        var t Sample
	        err := json.Unmarshal(*hit.Source, &t)
			//fmt.Fprintf(w, "[json=%s]", *hit.Source) // json output
	        if err != nil {
	            // Deserialization failed
	        }
	
	        fmt.Printf("id=%s: ", hit.Id)
	        fmt.Printf("Sample [name %s], [variants %d], [genome %s] \n", t.Name, t.Variants, t.Genome)
	    }

		jsonResult, err := json.Marshal(searchResult)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		results := "{\"results\":" + string(jsonResult) + "}"
		fmt.Fprintf(w, results)
	} else {
	    // No hits
	    fmt.Print("Found no samples\n")
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, model.AppLoader), app model.AppLoader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, app)
	}
}

func TestHandler(w http.ResponseWriter, r *http.Request, ctx model.AppLoader) {
	log.Printf("I'm running on %s with an %s CPU ", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Test=%d\n", ctx.Test)
	fmt.Printf("ctx=%v\n", ctx.Config)
	fmt.Printf("ctx.host=%s\n", ctx.Config.Host)
	fmt.Fprintf(w, "I'm running on %s with an %s CPU ", runtime.GOOS, runtime.GOARCH)
}

func main() {
	// Read config file properties
	var conf model.Config
	if _, err := toml.DecodeFile("application.properties", &conf); err != nil {
		panic(err)
	}

	// Create Elastic Search service
	client, err := elastic.NewClient(elastic.SetURL(conf.Url))
	if err != nil {
		fmt.Printf("error")
    	panic(err)
	}

	// Inject services & config
	var app model.AppLoader
	err = inject.Populate(client, &app);
	if err != nil {
		panic(err);
	}
	app.Test = 5
	app.Config = conf

	// Set up routes
	router := mux.NewRouter()
	router.HandleFunc("/view/{title}", makeHandler(ViewHandler, app))
	router.HandleFunc("/test", makeHandler(TestHandler, app))
	router.HandleFunc("/fake/{fakeTerm}", makeHandler(fakeResponseHandler, app))
	router.HandleFunc("/sample/{sampleId}", makeHandler(sampleHandler, app))
	router.HandleFunc("/search/{searchTerm}", makeHandler(searchHandler, app))

	// http://stackoverflow.com/questions/15834278/serving-static-content-with-a-root-url-with-the-gorilla-toolkit
	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	router.PathPrefix("/static/").Handler(s)

	http.Handle("/", router)

	//var port string = ":" + strconv.Itoa(app.Config.Port)
	var port string = ":" + fmt.Sprintf("%v", app.Config.Port)
	fmt.Println(port)
	log.Printf("I'm running on %s with an %s CPU ", runtime.GOOS, runtime.GOARCH)
	http.ListenAndServe(port, router)

	//m := new(model.Context);
}
