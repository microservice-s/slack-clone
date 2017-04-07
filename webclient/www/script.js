function loadOGP() {
    var title = document.getElementById("title")
    var description = document.getElementById("description")
    var image = document.getElementById("image")

    fetch('http://localhost:4000/v1/summary?url=http://ogp.me/', {
        method: 'get'
    }).then(function(resp) {
        
    })
}

loadOGP()