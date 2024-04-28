document.addEventListener('DOMContentLoaded', function() {
    // Check if the form exists before adding an event listener
    let formInscription = document.getElementById('formInscription');
    if (formInscription) {
        formInscription.addEventListener('submit', function(event) {
            event.preventDefault(); // Prevents the default form submission

            console.log("Formulaire soumis !");

            // Fetch form data
            const nom = document.getElementById('nom').value;
            const prenom = document.getElementById('prenom').value;
            const date_naissance = document.getElementById('date-de-naissance').value;
            const niveauEducation = document.getElementById('niveauEducation').value;
            
            // Log data to console
            console.log(`Nom : ${nom}, Prénom : ${prenom}, Date de Naissance : ${date_naissance}, Niveau d'éducation : ${niveauEducation}`);
            window.location.href = 'page-daccueil.html';
        });
    }

    // Adding event listeners to menu items as before
    let menuItems = document.querySelectorAll('.menu-item');
    menuItems.forEach(function(item) {
        item.addEventListener('click', function() {
            this.classList.toggle('active');
        });
    });

    // Handle formConnexion submission
    let formConnexion = document.getElementById('formConnexion');
    if (formConnexion) {
        formConnexion.addEventListener('submit', function(event) {
            event.preventDefault(); // Prevents the default form submission

            // Dummy validation and redirection for demonstration
            const email = document.getElementById('email').value;
            const password = document.getElementById('motDePasse').value;
            console.log(`Login Attempt: Email - ${email}, Password - ${password}`);
            
            // Here, you would typically send a request to your server for validation
            // For demonstration, let's just redirect to the home page on any submission
            window.location.href = 'page-daccueil.html';
        });
    }
});
