package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

const apiKey = "AIzaSyDXCfbgY6DIqU72BZa8bpnOL4n8WyTX_AY"
const apiURL = "https://www.googleapis.com/youtube/v3/search"

type User struct {
	ID        int
	FirstName string
	LastName  string
}

type Session struct {
	ID   string
	User User
}

// Données à insérer dans le modèle HTML
type PageData struct {
	Title   string
	User    User
	Domaine []Domaine
	Content Content
	Erreur  string
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
	// IdPostant       int `json:"
}

// Response structure for YouTube API
type YouTubeResponse struct {
	Items []struct {
		ID struct {
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			Title        string `json:"title"`
			ChannelTitle string `json:"channelTitle"`
			PublishedAt  string `json:"publishedAt"`
		} `json:"snippet"`
	} `json:"items"`
}

// var idUtilisateur = 0

// Store to hold sessions in memory
var sessionStore = struct {
	sync.RWMutex
	sessions map[string]Session
}{
	sessions: make(map[string]Session),
}

// Générer une chaîne aléatoire pour les cookies de session
func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func authenticateUser(db *sql.DB, email string, password string) (User, bool) {
	idUtilisateur := 0
	name := ""
	firstname := ""
	hashedPassword := ""

	queryBDD := "SELECT id,prenom,nom,mot_de_passe FROM utilisateur WHERE mail = ?"
	err := db.QueryRow(queryBDD, email).Scan(&idUtilisateur, &firstname, &name, &hashedPassword)

	match := CheckPasswordHash(password, hashedPassword)

	if err != nil || !match {
		log.Println("Erreur de connexion", err)
		return User{ID: 1, FirstName: "John", LastName: "Doe"}, false
	}
	return User{ID: idUtilisateur, FirstName: firstname, LastName: name}, true
}

func loginHandler(w http.ResponseWriter, r *http.Request, user User) {
	sessionID, err := generateSessionID()
	if err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}

	// Créer une nouvelle session et la stocker
	session := Session{
		ID:   sessionID,
		User: user,
	}
	sessionStore.Lock()
	sessionStore.sessions[sessionID] = session
	sessionStore.Unlock()

	// Créer des cookies sécurisés pour la session
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HttpOnly: true,
		Secure:   false, // A mettre a true en passant en HTTPS
		Path:     "/",
		MaxAge:   3600, // 1 heure
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			// Si le cookie n'existe pas, rediriger vers la page de connexion
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}

	// Supprimer la session du magasin de sessions
	sessionStore.Lock()
	delete(sessionStore.sessions, sessionCookie.Value)
	sessionStore.Unlock()

	// Supprimer le cookie de session
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		HttpOnly: true,
		Secure:   false, // A mettre a true en passant en HTTPS
		Path:     "/",
		MaxAge:   -1, // Supprimer le cookie
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

// La fonction de hachage pour les mdp
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// La fonction de vérifications d'un mdp
func CheckPasswordHash(password, hash string) bool {
	// CompareHashAndPassword permet de comparer un mot de passe et un hash bcrypt.
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// SearchYouTube function to search videos based on a query
func SearchYouTube(query string) (*YouTubeResponse, error) {
	searchURL := fmt.Sprintf("%s?part=snippet&q=%s&type=video&key=%s", apiURL, url.QueryEscape(query), apiKey)
	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ytResponse YouTubeResponse
	if err := json.NewDecoder(resp.Body).Decode(&ytResponse); err != nil {
		return nil, err
	}

	return &ytResponse, nil
}

func getDomaine() ([]Domaine, error) {
	log.Println("Connexion à la base de données...")
	db, err := sql.Open("mysql", "avnadmin:AVNS_x1AB4PkPIRzS-yIr_bP@tcp(learnhub-learnhub.b.aivencloud.com:15055)/learnhub")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// On vérifie que la connexion à la base de données est réussie
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// On exécute notre requête SQL pour obtenir les domaines
	rows, err := db.Query("SELECT nom from domaine order by id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Créer un tableau pour stocker les résultats
	var domaines []Domaine

	// Parcourir les lignes de résultats
	for rows.Next() {
		var nomDomaine string
		// Scanner la valeur de la colonne dans une variable
		err := rows.Scan(&nomDomaine)
		if err != nil {
			return nil, err
		}

		var dom = Domaine{
			Nom:    nomDomaine,
			Themes: []string{},
		}

		// Ajouter le domaine à notre tableau de domaines
		domaines = append(domaines, dom)
	}

	// WORK
	// Afficher les résultats
	// log.Println("Résultats:")
	// for _, domaine := range domaines {
	// 	log.Println(domaine.Nom)
	// }

	// Vérifier s'il y a des erreurs lors de l'itération des résultats
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// On exécute notre requête SQL pour obtenir tous les thèmes
	rows, err = db.Query("SELECT nom, id_domaine from theme order by id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parcourir les lignes de résultats
	for rows.Next() {
		var nomTheme string
		var idDomaine int
		// Scanner la valeur de la colonne dans une variable
		err := rows.Scan(&nomTheme, &idDomaine)
		if err != nil {
			return nil, err
		}

		// Ajouter le theme à notre tableau de themes du tableau de domaines
		domaines[idDomaine-1].Themes = append(domaines[idDomaine-1].Themes, nomTheme)
	}

	// WORK
	// Afficher les résultats
	// log.Println("Résultats:")
	// for _, domaine := range domaines {
	// 	log.Println(domaine.Nom)
	// 	for _, theme := range domaine.Themes {
	// 		log.Println(theme)
	// 	}
	// }

	return domaines, nil
}

func generateJsonDomaines() {
	// Récupération des domaines de la BDD
	domaines, err := getDomaine()
	if err != nil {
		panic(err.Error())
	}

	// Encodage des données en JSON
	jsonData, err2 := json.Marshal(domaines)
	if err2 != nil {
		log.Fatal(err2)
	}

	// Écriture des données JSON dans un fichier
	err = os.WriteFile("static/domaines.json", jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func extractDomainesJSON() []Domaine {
	// Lecture du fichier JSON
	jsonData, err := os.ReadFile("static/domaines.json")
	if err != nil {
		log.Fatal(err)
	}

	// Décodage des données JSON dans une liste de structures Domaine
	var domaines []Domaine
	err = json.Unmarshal(jsonData, &domaines)
	if err != nil {
		log.Fatal(err)
	}
	return domaines
}

func getContent(db *sql.DB, query string) Content {
	var myContent Content

	search := "'%" + query + "%'"
	// On exécute nos requêtes SQL pour obtenir les documents
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
		video.Lien = "https://www.youtube.com/watch?v=" + video.Lien
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

func getContentByTheme(db *sql.DB, theme string) Content {
	var myContent Content

	idTheme := 0

	// on cherche notre idTheme
	queryIdTheme := "SELECT id FROM theme WHERE nom = ?"
	_ = db.QueryRow(queryIdTheme, theme).Scan(&idTheme)

	// On exécute nos requêtes SQL pour obtenir les documents
	rowsVideos, err1 := db.Query("select lien,titre,auteur,date_document from learnhub.document where (id_theme = ?) AND id_type_document = 1", idTheme)
	if err1 != nil {
		return myContent
	}
	defer rowsVideos.Close()

	rowsArticles, err2 := db.Query("select lien,titre,auteur,date_document from learnhub.document where (id_theme = ?) AND id_type_document = 2", idTheme)
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
		video.Lien = "https://www.youtube.com/watch?v=" + video.Lien
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

// Permet de charger les différentes pages html néccessaires
func handler(w http.ResponseWriter, r *http.Request, htmlName string) {
	// Charger le formulaire d'ajout de documents
	formulaire := template.Must(template.ParseFiles(htmlName))

	err := formulaire.Execute(w, nil)
	if err != nil {
		http.Error(w, "Impossible de charger "+htmlName, http.StatusInternalServerError)
		return
	}
}

func main() {
	// On indique l'emplacement de nos données static (Images, css, javascript, etc...)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Connexion à la base de données mySQL
	db, err := sql.Open("mysql", "avnadmin:AVNS_x1AB4PkPIRzS-yIr_bP@tcp(learnhub-learnhub.b.aivencloud.com:15055)/learnhub")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Vérifie que la connexion à la base de données est réussie
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Charger le modèle HTML à partir d'un fichier
	tmpl := template.Must(template.ParseFiles("index.html"))

	// Variables de la page qu'on ne veut pas recharger à chaque requête
	var dom []Domaine
	var content Content

	// Récupération des domaines pour le menu
	dom = extractDomainesJSON()

	http.HandleFunc("/deconnexion", func(w http.ResponseWriter, r *http.Request) {
		var user User
		user.ID = 0
		user.FirstName = ""
		user.LastName = ""

		logoutHandler(w, r)
	})

	// Définir la route pour la page de formulaire
	http.HandleFunc("/formulaire", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, "formulaire.html")
	})
	// Définir la route pour la page de connexion
	http.HandleFunc("/connexion", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, "connexion.html")
	})
	// Définir la route pour la page d'inscription
	http.HandleFunc("/inscription", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, "inscription.html")
	})
	// Définir la route pour la page du profil
	http.HandleFunc("/profil", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, "profil.html")
	})

	// 127.0.0.1:8080/miniature?id=JX1gUaRydFo
	http.HandleFunc("/miniature", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("miniature.html"))
		data := getYoutubeThumbnail(r.URL.Query().Get("id"))

		log.Println(r.URL.Query().Get("id"))
		log.Println(data)

		// Exécuter le modèle avec les données fournies
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Définir la route pour la page par défaut
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// On réinitialise la recherche
		var c Content
		content = c

		// On reset la dernière erreur
		erreur := ""

		// On récupére l'ID de l'utilisateur s'il est connecté
		// session, _ := store.Get(r, "session-name")
		// userID, _ := session.Values["userID"].(int)

		if r.Method == "GET" {
			// Extraire les paramètres de la requête
			query := r.URL.Query()

			// Récupération du type du formulaire
			formType := query.Get("formType")
			typeRequet := query.Get("typeRequete")

			if formType == "Connexion" { // Connexion du client
				log.Println("GET détecté, connexion en cours")

				email := query.Get("email")
				// motDePasse, err := HashPassword(query.Get("motDePasse"))
				motDePasse := query.Get("motDePasse")

				// On vérifie si l'utilisateur existe avec ce mdp
				user, userExist := authenticateUser(db, email, motDePasse)
				if userExist {
					loginHandler(w, r, user)
				} else {
					erreur = "Mail ou mot de passe incorrect"
				}

				// log.Println("Session :", session.Values["userID"], session.Values["name"], session.Values["firstname"])
			} else if typeRequet == "rechercheTheme" { // Recherche de documents par theme sur le site
				// theme := r.FormValue("query")
				theme := r.URL.Query().Get("query")
				domaine := r.URL.Query().Get("query2")

				if theme != "" { // On empeche de faire une recherche vide qui renvoie tous les resultats
					log.Println("Recherche :", theme)

					content = getContentByTheme(db, theme)

					log.Println("Résultats:")
					for _, res := range content.Videos {
						log.Println(res.Titre)
					}

					ytResponse, err := SearchYouTube(domaine + theme) // pour plus de pertinance on cherche le domaine suivi du theme
					if err != nil {
						log.Fatalf("Error searching YouTube: %v", err)
					}

					log.Println("Youtube search")

					for _, item := range ytResponse.Items {
						var doc DocumentVisiable
						doc.Lien = "https://www.youtube.com/watch?v=" + item.ID.VideoID
						doc.Titre = item.Snippet.Title
						doc.Auteur = item.Snippet.ChannelTitle
						doc.Date = item.Snippet.PublishedAt
						doc.Miniature = getYoutubeThumbnail(item.ID.VideoID)
						content.Videos = append(content.Videos, doc)
					}

					// http.Redirect(w, r, "/", http.StatusFound)
				}

				// http.Redirect(w, r, "/", http.StatusFound)
			} else { // Recherche de documents par titre et auteurs sur le site
				query := r.FormValue("query")

				if query != "" { // On empeche de faire une recherche vide qui renvoie tous les resultats
					log.Println("Recherche :", query)

					content = getContent(db, query)

					log.Println("Résultats:")
					for _, res := range content.Videos {
						log.Println(res.Titre)
					}
				}
			}
		} else if r.Method == "POST" {
			formType := r.FormValue("formType")

			if formType == "Inscription" { // Inscription de l'utilisateur
				nom := r.FormValue("nom")
				prenom := r.FormValue("prenom")
				dateNaissance := r.FormValue("date-de-naissance")
				niveauEducation := r.FormValue("niveauEducation")
				linkedin := r.FormValue("linkedin")
				// diplome := r.FormValue("diplome")
				email := r.FormValue("email")
				motDePasse, err := HashPassword(r.FormValue("motDePasse"))
				if err != nil {
					log.Println("Erreur de hashage du mdp")
				}

				// On vérifie si l'utilisateur existe
				queryCheckBDD := "SELECT COUNT(ID) FROM utilisateur WHERE mail=?"
				numberMail := -1
				_ = db.QueryRow(queryCheckBDD, email).Scan(&numberMail)

				if numberMail == 0 {
					// On l'ajoute le cas échéant
					queryIdEtude := "SELECT id FROM niveau_etude WHERE intitule = ?"
					idEtude := 0
					_ = db.QueryRow(queryIdEtude, niveauEducation).Scan(&idEtude)

					_, err = db.Exec("INSERT INTO utilisateur (mail,nom,prenom,mot_de_passe,date_naissance,id_niveau_etude,lien_linkedin) VALUES (?, ?, ?, ?, ?, ?, ?)", email, nom, prenom, motDePasse, dateNaissance, idEtude, linkedin)

					if err != nil {
						log.Println("Erreur inscription : impossible d'ajouter l'utilisateur", err)
					}

					idUtilisateur := 0
					nom := ""
					prenom := ""

					queryBDD := "SELECT id,prenom,nom FROM utilisateur WHERE mail = ?"
					err := db.QueryRow(queryBDD, email).Scan(&idUtilisateur, &prenom, &nom)

					var user User
					user.ID = idUtilisateur
					user.FirstName = prenom
					user.LastName = nom

					loginHandler(w, r, user)

					if err != nil {
						log.Println("Erreur connexion apres inscription", err)
					}
				} else {
					log.Printf("Erreur inscription : utilisateur deja existant")
					erreur = "Utilisateur deja existant"
				}
			}
		} else if r.Method == "PUT" { //Ajout d'un document par l'utilisateur
			log.Println("PUT détecté")

			var document Document
			if err := json.NewDecoder(r.Body).Decode(&document); err != nil {
				http.Error(w, "Données JSON invalides", http.StatusBadRequest)
				log.Println("Erreur de décodage :", err)
				return
			}
			defer r.Body.Close()

			// Répondre avec les données reçues
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(document)
			log.Printf("Reçu : %+v\n", document)

			// Récupération de l'id du theme
			var idTheme int
			query := "SELECT id FROM theme WHERE nom = ?"
			err := db.QueryRow(query, document.Theme).Scan(&idTheme)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Println("Erreur récupération id du thème :", err)
				return
			}

			log.Println("id du theme : ", idTheme)

			var idPostant = 0

			// On récupère les infos de l'utilisateur
			sessionCookie, err := r.Cookie("session_id")
			if err != nil {
				// http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
				log.Println("Erreur, un utilisateur ajoute un document sans etre connecté")
			} else {
				// Récupérer la session à partir du magasin de sessions
				sessionStore.RLock()
				session, _ := sessionStore.sessions[sessionCookie.Value]
				sessionStore.RUnlock()
				idPostant = session.User.ID
			}

			// Insertion dans la base de données
			_, err = db.Exec("INSERT INTO document (lien, titre, auteur, id_postant, id_theme, id_type_document, date) values (?, ?, ?, ?, ?, ?, ?)", document.Lien, document.Titre, document.Auteur, idPostant, idTheme, document.IdTypeDocument, document.Date)
			if err != nil {
				// http.Error(w, "Erreur lors de l'insertion en base de données", http.StatusInternalServerError)
				log.Println("Erreur d'insertion du document :", err)
				// return
				erreur = "Erreur d'insertion du document"
			}

			log.Println(document)
		}

		var user User
		user.ID = 0
		user.FirstName = ""
		user.LastName = ""

		// On récupère les infos de l'utilisateur
		sessionCookie, err := r.Cookie("session_id")
		if err == nil {
			// Récupérer la session à partir du magasin de sessions
			sessionStore.RLock()
			session, _ := sessionStore.sessions[sessionCookie.Value]
			sessionStore.RUnlock()
			user = session.User
		}

		// Données à insérer dans le modèle HTML
		data := PageData{
			Title:   "Accueil",
			User:    user,
			Domaine: dom,
			Content: content,
			Erreur:  erreur,
		}

		// Exécuter le modèle avec les données fournies
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// http.Redirect(w, r, "/", http.StatusFound)
	})

	// Démarrer le serveur sur le port 8080
	log.Println("Serveur démarré sur le port :8080")
	generateJsonDomaines()
	http.ListenAndServe(":8080", nil)

}
