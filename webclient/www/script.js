var form = document.getElementById("form")
var title = document.getElementById("title");
var description = document.getElementById("description");
var image = document.getElementById("image");
var urlSubmit = document.getElementById("urlSubmit");


if (form.attachEvent) {
    form.attachEvent("submit", loadOGP);
} else {
    form.addEventListener("submit", loadOGP);
}
function status(response) {
    if (response.status >= 200 && response.status < 300) {  
        return Promise.resolve(response)  
    } else { 
        return response.text().then(function(errorMessage) {     
            console.log(errorMessage)   
            return Promise.reject(new Error(errorMessage))  
        })
    }  
}

function json(response) {  
    return response.json()  
}

function loadOGP(e) {
    if (e.preventDefault) e.preventDefault();
    fetch('http://138.68.55.2/v1/summary?url=' + urlSubmit.value)
        .then(status)
        .then(json)
        .then(function(data) {
            if(data.title === "" || data.title === undefined) {
                title.innerHTML = "No ogp title"
            } else {
                title.innerHTML = data.title
            }
            if(data.description === "" || data.description === undefined) {
                description.innerHTML = "No ogp description"
            } else {   
                description.innerHTML = data.description
            }   
            if(data.image === "" || data.image == undefined){
                image.src = "images/no-image.png"
            }else {
                image.src = data.image
            }        
        }).catch(function(err) {
            console.log(err)
            title.innerHTML = "Error"
            description.innerHTML = err.message 
            image.src = "images/shrug.jpg"
        });

    return false;
}