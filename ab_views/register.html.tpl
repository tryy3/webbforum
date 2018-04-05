<form method="POST" class="box">
    <div class="box-title">Registration</div>

    <input type="hidden" name="{{.xsrfName}}" value="{{.xsrfToken}}" />

    {{with .errs}}{{with $errlist := index . "first_name"}}<span class="has-error">{{end}}{{end}}
    <span class="input-title">Namn</span>
    <input type="text" name="first_name" placeholder="Namn" value="{{.first_name}}"/>
    {{with .errs}}{{with $errlist := index . "first_name"}}{{range $errlist}}<div class="error-block">{{translate .}}</div>{{end}}</span>{{end}}{{end}}

    {{with .errs}}{{with $errlist := index . "last_name"}}<span class="has-error">{{end}}{{end}}
    <span class="input-title">Efternamn</span>
    <input type="text" name="last_name" placeholder="Efternamn" value="{{.last_name}}"/>
    {{with .errs}}{{with $errlist := index . "last_name"}}{{range $errlist}}<div class="error-block">{{translate .}}</div>{{end}}</span>{{end}}{{end}}

    {{with .errs}}{{with $errlist := index . "email"}}<span class="has-error">{{end}}{{end}}
    <span class="input-title">E-postadress</span>
    <input type="text" name="email" placeholder="E-postadress" value="{{.email}}"/>
    {{with .errs}}{{with $errlist := index . "email"}}{{range $errlist}}<div class="error-block">{{translate .}}</div>{{end}}</span>{{end}}{{end}}

    {{$pid := .primaryID}}
    {{with .errs}}{{with $errlist := index . $pid}}<span class="has-error">{{end}}{{end}}
    <span class="input-title">Användarnamn</span>
    <input type="text" name="{{.primaryID}}" placeholder="Användarnamn" value="{{.primaryIDValue}}"/>
    {{with .errs}}{{with $errlist := index . $pid}}{{range $errlist}}<div class="error-block">{{translate .}}</div>{{end}}</span>{{end}}{{end}}

    {{with .errs}}{{with $errlist := index . "password"}}<span class="has-error">{{end}}{{end}}
    <span class="input-title">Lösenord</span>
    <input type="password" name="password" placeholder="Lösenord" value="{{.password}}"/>
    {{with .errs}}{{with $errlist := index . "password"}}{{range $errlist}}<div class="error-block">{{translate .}}</div>{{end}}</span>{{end}}{{end}}

    {{with .errs}}{{with $errlist := index . "confirm_password"}}<span class="has-error">{{end}}{{end}}
    <span class="input-title">Bekräfta lösenord</span>
    <input type="password" name="confirm_password" placeholder="Bekräfta lösenord" value="{{.confirmPassword}}"/>
    {{with .errs}}{{with $errlist := index . "confirm_password"}}{{range $errlist}}<div class="error-block">{{translate .}}</div>{{end}}</span>{{end}}{{end}}

    <div class="btnWrapper"><button class="btn" type="submit">Registrera</button></div>

    <div class="forgotWrapper"><a class="forgot" href="{{mountpathed "login"}}">Redan medlem?</a></div>
</form>