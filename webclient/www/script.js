function loadOGP() {
    console.log("test");
    event.preventDefault();

    var title = document.getElementById("title")
    var description = document.getElementById("description")
    var image = document.getElementById("image")
    var urlForm = document.getElementById("urlForm")
    console.log("Why isn't this getting called?");

    fetch('localhost:4000/v1/summary?url=' + urlForm.value, {
        method: 'get'
    }).then(function(resp) {
        //console.log(resp)
        return resp.json()
    }).then(function(data) {
        console.log(data)
        title.innerHTML = data.title
        description.innerHTML = data.description
        image.src = data.image
        
    }).catch(function(err) {
        console.log(err)
    });
    // You must return false to prevent the default form behavior
    return false;
}

//loadOGP()