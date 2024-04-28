package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

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

		dom, err2 := getDomaine()

		if err2 != nil {
			panic(err.Error())
		}

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
	})

	// Définir la route pour la page de formulaire
	http.HandleFunc("/formulaire", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, "formulaire.html")
	})
	http.HandleFunc("/null", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, "null.html")
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
