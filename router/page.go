package router

import (
	"monolith/config"
	"net/http"
	"text/template"
)

const TEMPLATE = `<!DOCTYPE html>
<html>
<head>
<title>Monolith</title>
<meta charset="utf-8">
<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/3.6.4/jquery.min.js" integrity="sha512-pumBsjNRGGqkPzKHndZMaAG+bir374sORyzM3uulLV14lN5LyykqNk8eEeUlUkB3U0M4FApyaHraT65ihJhDpQ==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
<script type="text/javascript">

function requestToken() {
    return new Promise((resolve, reject) => {
        $.get("/token")
            .done(function (data) {
                localStorage.setItem("token", data);
                resolve(data);
            })
            .fail(function (xhr, status, error) {
                alert("Got error code " + xhr.status + " while requesting token");
                reject(error);
            });
    });
}

/* Checks if the token in localStorage is still usable or requests a new one */
function getToken() {
    return new Promise((resolve, reject) => {
        let token = localStorage.getItem("token");

        if (token == null)
            requestToken()
                .then( (newToken) => resolve(newToken) )
                .catch( (error) => reject(error) );
        else {
            let body = JSON.parse(atob(token.split('.')[1]));
            let now = Math.floor(Date.now() / 1000) + 2; // at most 2s until expire
            
            if (now > body['exp'])
                requestToken()
                    .then( (newToken) => resolve(newToken) )
                    .catch( (error) => reject(error) );
            else resolve(token);
        }
    });
}

const FRUITS_MICROSERVICE_ENDPOINT="{{.FruitsEndpoint}}";

/* Get the list of fruits and fill the list */
function getFruits() {
	$("#fruits-error").css({"display": "none"});

	$.get(FRUITS_MICROSERVICE_ENDPOINT + "/")
	
		.done((data) => {
			$("ul#fruits").empty();
			data.forEach( (pair) => {
				let li = $("<li>").text(pair.username + ": "+ pair.fruit);
				$("ul#fruits").append(li);
			});
		})
		
		.fail((xhr, status, error) => {
			$("#fruits-error-content").text(error);
			$("#fruits-error").css({"display": "block"});
		});
}

async function setFruit() {
	let fruit = $("#fruit").val();
	$("#set-fruit").prop("disabled", true);
	try {
		let token = await getToken();
		$.ajax({
			url: FRUITS_MICROSERVICE_ENDPOINT + "/fruit",
			method: "PUT",
			headers: {
				"X-Auth-Token": token,
			},
			data: {
				fruit: fruit,
			},
		})
		.done( () => getFruits() )
		.fail( (x, s, e) => alert("Error setting fruit [" + x.status + "]: " + e) )
		.always( () => $("#set-fruit").prop("disabled", false) );

	} catch(error) {
		alert("Error setting fruit: " + error);
		$("#set-fruit").prop("disabled", false);
	}
}

$(()=>{
	$("#fruits-retry").on("click", ()=> getFruits() );
	$("#set-fruit").on("click", ()=> setFruit() );
	getFruits();
});

</script>
<body>
<h1>Fruits service</h1>
{{if .User}}
<label for="fruit">Fruit:</label>
<select name="fruit" id="fruit">
	<option value="apple">Apple 🍎</option>
	<option value="banana">Banana 🍌</option>
	<option value="orange">Orange 🍊</option>
	<option value="pear">Pear 🍐</option>
	<option value="pineapple">*Pineapple 🍍</option>
	<option value="kiwi">*Kiwi 🥝</option>
</select>
<button id="set-fruit">Set fruit</button>
<form action="/logout" method="GET">
	<input type="submit" value="🚪 Logout">
</form>
{{else}}
<form action="/login" method="POST">
	<label for="username">🐱 Username:</label>
	<input type="text" name="username" id="username">
	<label for="password">🔑 Password:</label>
	<input type="password" name="password" id="password">
	<input type="submit" value="Login ➡️">
</form>
{{end}}
<div id="fruits-error" style="display:none">
	Error loading fruits: <span id="fruits-error-content">unknown</span>
	<button id="fruits-retry">Retry</button>
</div>
<ul id="fruits">
</ul>
</body>
</html>`

var t = template.Must(template.New("index").Parse(TEMPLATE))

func printIndexPage(user int, w http.ResponseWriter) error {
	return t.Execute(w, struct {
		FruitsEndpoint string
		User           bool
	}{
		config.FruitsEndpointExternal,
		user > 0,
	})
}
