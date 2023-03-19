package router

import (
	"net/http"
	"text/template"
)

const TEMPLATE = `<!DOCTYPE html>
<html>
<head>
<title>Monolith</title>
<meta charset="utf-8">
</head>
<body>
<h1>Fruits service</h1>
{{if .User}}
<form action="/fruit" method="POST">
	<label for="fruit">Fruit:</label>
	<select name="fruit" id="fruit">
		<option value="apple">Apple 🍎</option>
		<option value="banana">Banana 🍌</option>
		<option value="orange">Orange 🍊</option>
		<option value="pear">Pear 🍐</option>
		<option value="pineapple">*Pineapple 🍍</option>
		<option value="kiwi">*Kiwi 🥝</option>
	</select>
	<input type="submit" value="Set fruit">
</form>
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
<ul>
{{range $key, $value := .Fruits}}
	<li>{{$key}}: {{$value}}</li>
{{end}}
</ul>
</body>
</html>`

var t = template.Must(template.New("index").Parse(TEMPLATE))

func printIndexPage(fruits map[string]string, user int, w http.ResponseWriter) error {
	return t.Execute(w, struct {
		Fruits map[string]string
		User   bool
	}{
		fruits,
		user > 0,
	})
}
