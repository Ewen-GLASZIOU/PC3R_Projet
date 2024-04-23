package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// // Définit le répertoire View ou se trouve nos fichiers html
	// fs := http.FileServer(http.Dir("View"))
	// http.Handle("/", fs)

	// Ouvre une connexion à la base de données MySQL le format est : "mysql", "user:password@tcp(127.0.0.1:3306)/dbname"
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

	// Route pour la page d'accueil
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Exemple de requête pour récupérer le nom de l'utilisateur depuis la base de données
		var name string
		query := "SELECT prenom FROM utilisateur WHERE id = ?"
		err := db.QueryRow(query, 1).Scan(&name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Affiche le nom de l'utilisateur dans la page HTML à l'emplacement spécifié
		// fmt.Fprintf(w, "<script>document.getElementById('user-name').innerText = '%s';</script>", name)

	})

	// Lance le serveur web sur le port 8080
	log.Println("Serveur démarré sur le port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
