<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">

	<title>{{template "pagetitle" .}}</title>

    <script defer src="https://use.fontawesome.com/releases/v5.0.8/js/solid.js" integrity="sha384-+Ga2s7YBbhOD6nie0DzrZpJes+b2K1xkpKxTFFcx59QmVPaSA8c7pycsNaFwUK6l" crossorigin="anonymous"></script>
    <script defer src="https://use.fontawesome.com/releases/v5.0.8/js/fontawesome.js" integrity="sha384-7ox8Q2yzO/uWircfojVuCQOZl+ZZBg2D2J5nkpLqzH1HY0C1dHlTKIbpRz/LG23c" crossorigin="anonymous"></script>

    <link rel="stylesheet" href="/css/sceditor.default.min.css">
    <link rel="stylesheet" href="/css/content/default.min.css">
    <link rel="stylesheet" href="/css/style.css">

    <script type="text/javascript" src="/js/sceditor.min.js"></script>
    <script type="text/javascript" src="/js/bbcode.js"></script>
</head>
<body>
	<nav class="navbar">
        <a class="navbar-logo" href="/"><img src="/images/FM_vänster_neg.png" alt="Försvarsmakten"></a>

        <ul class="navbar-middle">
            <li>
                <a href="/">Hem</a>
            </li>
            <li>
                <a href="/">Forum</a>
            </li>
            <li>
                <a href="#">Kontakta oss</a>
            </li>
        </ul>

        <ul class="navbar-right">
            {{if not .loggedin}}
                <li><a href="/auth/register">Register</a></li>
                <li><a href="/auth/login"><i class="fas fa-sign-in-alt"></i> Login</a></li>
            {{else}}
                <li>
                    <a href="/profile">Välkommen {{.current_user_name}}! <span class="caret"></span></a>
                </li>
                <li>
                    <a href="/auth/logout">
                        <i class="fas fa-sign-out-alt"></i> Logout
                    </a>
                </li>
            {{end}}
        </ul>
	</nav>

    </br>
	{{with .flash_success}}<div class="alert alert-success">{{translate .}}</div>{{end}}
	{{with .flash_error}}<div class="alert alert-danger">{{translate .}}</div>{{end}}
	{{template "yield" .}}
	{{template "authboss" .}}
</body>
</html>
{{define "pagetitle"}}{{end}}
{{define "yield"}}{{end}}
{{define "authboss"}}{{end}}