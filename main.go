package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Données à insérer dans le modèle HTML
type PageData struct {
	Title     string
	Firstname string
	Name      string
	Domaine   []Domaine
	Content   Content
}

type Domaine struct {
	Nom    string
	Themes []string
}

type Content struct {
	Videos   []string
	Articles []string
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

/* // Save
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
*/

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

func getContent(query string) Content {
	var myContent Content

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

	search := "'%" + query + "%'"

	// On exécute nos requêtes SQL pour obtenir les documents
	rowsVideos, err1 := db.Query("select lien from learnhub.document where (titre LIKE " + search + " or auteur LIKE " + search + ") AND id_type_document = 1")
	if err1 != nil {
		return myContent
	}
	defer rowsVideos.Close()

	rowsArticles, err2 := db.Query("select lien from learnhub.document where (titre LIKE " + search + " or auteur LIKE " + search + ") AND id_type_document = 2")
	if err2 != nil {
		return myContent
	}
	defer rowsArticles.Close()

	for rowsVideos.Next() {
		var lien string

		err1 := rowsVideos.Scan(&lien)
		if err1 != nil {
			return myContent
		}

		myContent.Videos = append(myContent.Videos, lien)
	}

	for rowsArticles.Next() {
		var lien string

		err2 := rowsArticles.Scan(&lien)
		if err2 != nil {
			return myContent
		}

		myContent.Articles = append(myContent.Articles, lien)
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
	var name string
	var firstname string
	var dom []Domaine
	var content Content

	// Définir la route pour la page par défaut
	http.HandleFunc("/aled", func(w http.ResponseWriter, r *http.Request) {
		// On récupère le nom et prenom de l'utilisateur
		query := "SELECT prenom, nom FROM utilisateur WHERE id = ?"
		err := db.QueryRow(query, 1).Scan(&firstname, &name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Récupération des domaines pour le menu
		dom := extractDomainesJSON()

		// Données à insérer dans le modèle HTML
		data := PageData{
			Title:     "Accueil",
			Firstname: firstname,
			Name:      name,
			Domaine:   dom,
			Content:   content,
		}

		// TODO : make an add domaines and themes form
		// generateJsonDomaines()

		// Exécuter le modèle avec les données fournies
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Définir la route pour la page de formulaire
	http.HandleFunc("/formulaire", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, "formulaire.html")
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {

			// Données à insérer dans le modèle HTML
			data := PageData{
				Title:     "Accueil",
				Firstname: firstname,
				Name:      name,
				Domaine:   dom,
				Content:   content,
			}

			// Exécuter le modèle avec les données fournies
			err = tmpl.Execute(w, data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		} else if r.Method == "PUT" {
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

			log.Println("id du theme : ?", idTheme)

			// Insertion dans la base de données
			_, err = db.Exec("INSERT INTO document (lien, titre, auteur, id_postant, id_theme, id_type_document) values (?, ?, ?, 1, ?, ?)", document.Lien, document.Titre, document.Auteur, idTheme, document.IdTypeDocument)
			if err != nil {
				http.Error(w, "Erreur lors de l'insertion en base de données", http.StatusInternalServerError)
				log.Println("Erreur d'insertion du document :", err)
				return
			}

			log.Println(document)
		} else {
			// r.ParseForm()
			// firstname = r.FormValue("firstname")
			// name = r.FormValue("name")
			// log.Println(firstname)
			// log.Println(name)
			query := r.FormValue("query")
			log.Println("Recherche :", query)

			content = getContent(query)

			log.Println(err)

			log.Println("Résultats:")
			for _, res := range content.Videos {
				log.Println(res)
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	})

	// Démarrer le serveur sur le port 8080
	log.Println("Serveur démarré sur le port :8080")
	http.ListenAndServe(":8080", nil)

}
