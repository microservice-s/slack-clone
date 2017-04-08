var form = document.getElementById("form")
if (form.attachEvent) {
    form.attachEvent("submit", loadOGP);
} else {
    form.addEventListener("submit", loadOGP);
}
function status(response) {
    if (response.status >= 200 && response.status < 300) {  
        return Promise.resolve(response)  
    } else {  
        return Promise.reject(new Error(response.status))  
    }  
}

function json(response) {  
    return response.json()  
}

function loadOGP(e) {
    if (e.preventDefault) e.preventDefault();

    var title = document.getElementById("title")
    var description = document.getElementById("description")
    var image = document.getElementById("image")
    var urlSubmit = document.getElementById("urlSubmit")
    console.log(urlSubmit.value);

    fetch('http://localhost:4000/v1/summary?url=' + urlSubmit.value)
        .then(status)
        .then(json)
        .then(function(data) {
            title.innerHTML = data.title
            description.innerHTML = data.description
            image.src = data.image
        }).catch(function(err) {
            // TODO add response to user that the server wasn't found
            console.log(err)
        });

    
    // .then(function(resp) {
    //     if(resp.status !== 200){
    //         console.log("error")
    //         return
    //     } else {
    //         return resp.json();
    //     }
        
    // }).then(function(data) {
    //     title.innerHTML = data.title
    //     description.innerHTML = data.description
    //     image.src = data.image
    // }).catch(function(err) {
    //     // TODO add response to user that the server wasn't found
    //     console.log(err)
    // });
    // You must return false to prevent the default form behavior
    return false;
}

//loadOGP()