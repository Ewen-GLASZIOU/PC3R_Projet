$(document).ready(function(){
    
    // Appel de la fonction pour charger le formulaire au chargement de la page
    // chargerFormulaire();
    attacherEvenement();
});

// Fonction pour charger le formulaire via AJAX
function chargerFormulaire() {
    $.ajax({
        url: "/formulaire", // URL du fichier contenant le formulaire
        type: "GET",
        success: function(response) {
            $("#formulaireContainer").html(response); // Insérer le contenu récupéré dans la div
            attacherEvenement(); // Appeler la fonction pour attacher l'événement après avoir chargé le formulaire
            document.getElementById("buttonFormulaire").style.display = "none";
        },
        error: function(xhr, status, error) {
            console.error("Erreur lors du chargement du formulaire: " + error);
        }
    });
}

function fermerFormulaire() {
    $("#formulaireContainer").html("");
    document.getElementById("buttonFormulaire").style.display = "block";
}

// Fonction pour attacher l'événement au formulaire après qu'il a été chargé via AJAX
function attacherEvenement() {
    $("#documentLink").on("input", function() {
        // Appeler la fonction lorsque du texte est collé dans l'input
        checkLink();
    });
}

function getYoutubeID(videoIdPiece) {
    // On récupere l'ID : ce qui se trouve avant le premier "&" ou "?" dans l'URL s'il y en a
    return videoIdPiece.split("&")[0].split("?")[0];
}

function getYoutubeInformations(videoId) {
    // Récupérer les informations de la vidéo via l'API YouTube
    $.ajax({
        url: `https://www.googleapis.com/youtube/v3/videos?id=${videoId}&key=${API_KEY}&part=snippet`,
        type: "GET",
        success: function(response) {
            // Extraire le titre, les auteurs et la date de la vidéo
            let videoTitle = response.items[0].snippet.title;
            let videoAuthors = response.items[0].snippet.channelTitle;
            let videoDate = response.items[0].snippet.publishedAt;

            // Mettre à jour les champs de titre, auteurs et date avec les informations récupérées
            document.getElementById("documentTitle").value = videoTitle;
            document.getElementById("documentAuthors").value = videoAuthors;
            document.getElementById("documentDate").value = videoDate;
        },
        error: function(xhr, status, error) {
            console.error("Erreur lors de la récupération des informations de la vidéo YouTube: " + error);
        }
    });

}

function checkLink() {
    // Récupérer la valeur du lien de la vidéo
    let documentLink = document.getElementById("documentLink").value;
    let lastVideoID = null;

    
    document.getElementById("documentTitle").disabled = true;
    document.getElementById("documentAuthors").disabled = true;
    document.getElementById("documentDate").disabled = true;

    document.getElementById("documentTitle").value = "";
    document.getElementById("documentAuthors").value = "";
    document.getElementById("documentDate").value = "";

    // Expression régulière pour rechercher "http://" ou "https://"
    let regex = /https?:\/\//;

    // On cache la miniature youtube
    $("#documentThumbnail").html("");

    // Vérifier si la chaîne contient "http://" ou "https://"
    if (regex.test(documentLink)) {
        // Vérifier si le lien est un lien youtube
        if (documentLink.includes("youtube.com") || documentLink.includes("youtu.be")) {
            console.log("youtube site");
            //https://www.youtube.com/watch?v=JX1gUaRydFo
            //https://youtu.be/JX1gUaRydFo?si=2ps9vIu7AiNgGm1X

            // Récupérer les informations de la vidéo
            //getYoutubeInformations(videoId);

            let videoID = getYoutubeVideoId(url);

            if(videoID != lastVideoID) {
                lastVideoID = videoID;
                afficherMiniatureYoutube(documentLink);
            }            
        } else if (documentLink.includes(".pdf")) {            
            if (documentLink.includes("dblp.org")) {
                console.log("dblp site");
            } else {
                // alert("Veuillez entrer un lien valide vers une vidéo YouTube ou un document DBLP.");
                console.log("other site");
                document.getElementById("documentTitle").disabled = false;
                document.getElementById("documentAuthors").disabled = false;
                document.getElementById("documentDate").disabled = false;
            }
        }
    }
}


function validerFormulaire() {
    // Récupérer la valeur du lien de la vidéo
    let documentLink = document.getElementById("documentLink").value;
    
    let documentTitle = "Titre de la vidéo";
    let documentAuthors = "Auteurs de la vidéo";
    let documentDate = "Date de la vidéo";

    // // document.getElementById("documentTitle").disabled = false;
    
    // Mettre à jour les champs de titre, auteurs et date avec les informations récupérées
    // document.getElementById("documentLink").value = "videoTitle";
    document.getElementById("documentTitle").value = documentLink;

    document.getElementById("documentAuthors").value = documentAuthors;
    document.getElementById("documentDate").value = documentDate;
}

function resetFormulaire() {
    // Réinitialiser les champs du formulaire
    document.getElementById("documentForm").reset();
}

// Fonction pour extraire l'identifiant de la vidéo YouTube à partir de l'URL
function getYoutubeVideoId(url) {
    var videoId = null;
    var match = url.match(/(?:https?:\/\/)?(?:www\.)?(?:youtube\.com\/(?:[^\/\n\s]+\/\S+\/|(?:v|e(?:mbed)?)\/|\S*?[?&]v=)|youtu\.be\/)([a-zA-Z0-9_-]{11})/);
    if (match) {
        videoId = match[1];
    }
    return videoId;
}

// Fonction pour afficher l'image miniature de la vidéo YouTube
function afficherMiniatureYoutube(url) {    
    // Afficher l'image miniature dans un élément HTML
    $.ajax({
        url: "/miniature?id=" + videoId,
        type: "GET",
        success: function(response) {
            $("#documentThumbnail").html(response); // Insérer le contenu récupéré dans la div
        },
        error: function(xhr, status, error) {
            console.error("Erreur lors du chargement de l'image miniature: " + error);
        }
    });
}