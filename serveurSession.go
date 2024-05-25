package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

// Données à insérer dans le modèle HTML
type PageData struct {
	Title     string
	Firstname string
	Name      string
	Id        int
	Domaine   []Domaine
	Content   Content
}

type Domaine struct {
	Nom    string
	Themes []string
}

type DocumentVisiable struct {
	Lien      string
	Titre     string
	Auteur    string
	Date      string
	Miniature string
}

type Content struct {
	Videos   []DocumentVisiable
	Articles []DocumentVisiable
}

type Document struct {
	Lien           string `json:"documentLink"`
	Titre          string `json:"documentTitle"`
	Auteur         string `json:"documentAuthors"`
	Date           string `json:"documentDate"`
	Theme          string `json:"documentTheme"`
	IdTypeDocument int    `json:"documentType"`
}

var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func getDomaine() ([]Domaine, error) {
	log.Println("Connexion à la base de données...")
	db, err := sql.Open("mysql", "avnadmin:AVNS_x1AB4PkPIRzS-yIr_bP@tcp(learnhub-learnhub.b.aivencloud.com:15055)/learnhub")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT nom from domaine order by id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domaines []Domaine

	for rows.Next() {
		var nomDomaine string
		err := rows.Scan(&nomDomaine)
		if err != nil {
			return nil, err
		}

		var dom = Domaine{
			Nom:    nomDomaine,
			Themes: []string{},
		}

		domaines = append(domaines, dom)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	rows, err = db.Query("SELECT nom, id_domaine from theme order by id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var nomTheme string
		var idDomaine int
		err := rows.Scan(&nomTheme, &idDomaine)
		if err != nil {
			return nil, err
		}

		domaines[idDomaine-1].Themes = append(domaines[idDomaine-1].Themes, nomTheme)
	}

	return domaines, nil
}

func generateJsonDomaines() {
	domaines, err := getDomaine()
	if err != nil {
		panic(err.Error())
	}

	jsonData, err2 := json.Marshal(domaines)
	if err2 != nil {
		log.Fatal(err2)
	}

	err = os.WriteFile("static/domaines.json", jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func extractDomainesJSON() []Domaine {
	jsonData, err := os.ReadFile("static/domaines.json")
	if err != nil {
		log.Fatal(err)
	}

	var domaines []Domaine
	err = json.Unmarshal(jsonData, &domaines)
	if err != nil {
		log.Fatal(err)
	}
	return domaines
}

func getContent(query string) Content {
	var myContent Content

	log.Println("Connexion à la base de données...")
	db, err := sql.Open("mysql", "avnadmin:AVNS_x1AB4PkPIRzS-yIr_bP@tcp(learnhub-learnhub.b.aivencloud.com:15055)/learnhub")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	search := "'%" + query + "%'"

	rowsVideos, err1 := db.Query("select lien,titre,auteur,date_document from learnhub.document where (titre LIKE " + search + " or auteur LIKE " + search + ") AND id_type_document = 1")
	if err1 != nil {
		return myContent
	}
	defer rowsVideos.Close()

	rowsArticles, err2 := db.Query("select lien,titre,auteur,date_document from learnhub.document where (titre LIKE " + search + " or auteur LIKE " + search + ") AND id_type_document = 2")
	if err2 != nil {
		return myContent
	}
	defer rowsArticles.Close()

	for rowsVideos.Next() {
		var video DocumentVisiable

		err1 := rowsVideos.Scan(&video.Lien, &video.Titre, &video.Auteur, &video.Date)
		if err1 != nil {
			return myContent
		}

		video.Miniature = getYoutubeThumbnail(video.Lien)
		myContent.Videos = append(myContent.Videos, video)
	}

	for rowsArticles.Next() {
		var article DocumentVisiable

		err2 := rowsArticles.Scan(&article.Lien, &article.Titre, &article.Auteur, &article.Date)
		if err2 != nil {
			return myContent
		}

		myContent.Articles = append(myContent.Articles, article)
	}

	return myContent
}

func getYoutubeThumbnail(videoID string) string {
	return "https://img.youtube.com/vi/" + videoID + "/0.jpg"
}

func handler(w http.ResponseWriter, r *http.Request, htmlName string) {
	formulaire := template.Must(template.ParseFiles(htmlName))

	err := formulaire.Execute(w, nil)
	if err != nil {
		http.Error(w, "Impossible de charger "+htmlName, http.StatusInternalServerError)
		return
	}
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	db, err := sql.Open("mysql", "avnadmin:AVNS_x1AB4PkPIRzS-yIr_bP@tcp(learnhub-learnhub.b.aivencloud.com:15055)/learnhub")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(template.ParseFiles("index.html"))

	var dom []Domaine
	var content Content

	dom = extractDomainesJSON()

	http.HandleFunc("/aled", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")

		userID, ok := session.Values["userID"].(int)
		if !ok {
			http.Redirect(w, r, "/connexion", http.StatusFound)
			return
		}

		var firstname, name string
		query := "SELECT prenom, nom FROM utilisateur WHERE id = ?"
		err := db.QueryRow(query, userID).Scan(&firstname, &name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dom := extractDomainesJSON()

		data := PageData{
			Title:     "Accueil",
			Firstname: firstname,
			Name:      name,
			Id:        userID,
			Domaine:   dom,
			Content:   content,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/deconnexion", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")

		session.Values["userID"] = 0
		session.Save(r, w)

		data := PageData{
			Title:     "Accueil",
			Firstname: "",
			Name:      "",
			Id:        0,
			Domaine:   dom,
			Content:   content,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/profil", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")

		userID, ok := session.Values["userID"].(int)
		if !ok {
			http.Redirect(w, r, "/connexion", http.StatusFound)
			return
		}

		var firstname, name string
		query := "SELECT prenom, nom FROM utilisateur WHERE id = ?"
		err := db.QueryRow(query, userID).Scan(&firstname, &name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dom := extractDomainesJSON()

		data := PageData{
			Title:     "Profil",
			Firstname: firstname,
			Name:      name,
			Id:        userID,
			Domaine:   dom,
			Content:   content,
		}

		tmpl.Execute(w, data)
	})

	http.HandleFunc("/recherche", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")

		userID, ok := session.Values["userID"].(int)
		if !ok {
			http.Redirect(w, r, "/connexion", http.StatusFound)
			return
		}

		var firstname, name string
		query := "SELECT prenom, nom FROM utilisateur WHERE id = ?"
		err := db.QueryRow(query, userID).Scan(&firstname, &name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dom := extractDomainesJSON()

		content := getContent(r.FormValue("search"))

		data := PageData{
			Title:     "Accueil",
			Firstname: firstname,
			Name:      name,
			Id:        userID,
			Domaine:   dom,
			Content:   content,
		}

		tmpl.Execute(w, data)
	})

	http.HandleFunc("/inscription", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, "inscription.html")
	})

	http.HandleFunc("/connexion", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, "connexion.html")
	})

	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, "contact.html")
	})

	http.HandleFunc("/apropos", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, "apropos.html")
	})

	http.HandleFunc("/form_inscription", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		prenom := r.FormValue("prenom")
		nom := r.FormValue("nom")
		motDePasse := r.FormValue("motdepasse")
		email := r.FormValue("email")

		_, err = db.Exec("INSERT INTO utilisateur (prenom, nom, email, mot_de_passe) VALUES (?, ?, ?, ?)", prenom, nom, email, motDePasse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/connexion", http.StatusFound)
	})

	http.HandleFunc("/form_connexion", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		motDePasse := r.FormValue("motdepasse")
		email := r.FormValue("email")

		var id int
		query := "SELECT id FROM utilisateur WHERE email = ? AND mot_de_passe = ?"
		err = db.QueryRow(query, email, motDePasse).Scan(&id)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		session, _ := store.Get(r, "session-name")
		session.Values["userID"] = id
		session.Save(r, w)

		http.Redirect(w, r, "/aled", http.StatusFound)
	})

	log.Println("Server started on: http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
