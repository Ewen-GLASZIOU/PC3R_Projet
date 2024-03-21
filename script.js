document.getElementById('formInscription').addEventListener('submit', function(event) {
    event.preventDefault(); // Empêche la soumission classique du formulaire

    console.log("Formulaire soumis !");

    // Récupération des données du formulaire
    const nom = document.getElementById('nom').value;
    const prenom = document.getElementById('prenom').value;
    const age = document.getElementById('age').value;
    const niveauEducation = document.getElementById('niveauEducation').value;
    
    // Affichage des données dans la console
    console.log(`Nom : ${nom}, Prénom : ${prenom}, Âge : ${age}, Niveau d'éducation : ${niveauEducation}`);
    window.location.href = 'page-daccueil.html';
});

document.addEventListener('DOMContentLoaded', function() {
    var menuItems = document.querySelectorAll('.menu-item');

    menuItems.forEach(function(item) {
        item.addEventListener('click', function() {
            // Toggle la classe 'active' pour afficher/cacher le sous-menu
            this.classList.toggle('active');
        });
    });
});
