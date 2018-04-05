<form method="POST" class="box">
    <input type="hidden" name="token" value="{{.token}}" />
    <input type="hidden" name="{{.xsrfName}}" value="{{.xsrfToken}}" />

    <div class="box-title">Återställ lösenord</div>

    {{if .error}}
        <div class="alert alert-danger">
            {{translate .error}}
        </div>
    {{end}}


    {{with .errs}}{{with $errlist := index . "password"}}<span class="has-error">{{end}}{{end}}
        <span class="input-title">Lösenord</span>
        <input type="text" name="password" placeholder="Lösenord" value="{{.password}}"/>
    {{with .errs}}{{with $errlist := index . "password"}}{{range $errlist}}<div class="error-block">{{translate .}}</div>{{end}}</span>{{end}}{{end}}

    {{with .errs}}{{with $errlist := index . "confirm_password"}}<span class="has-error">{{end}}{{end}}
        <span class="input-title">Bekräfta lösenord</span>
        <input type="text" name="confirm_password" placeholder="Bekräfta lösenord" value="{{.confirmPassword}}"/>
    {{with .errs}}{{with $errlist := index . "confirm_password"}}{{range $errlist}}<div class="error-block">{{translate .}}</div>{{end}}</span>{{end}}{{end}}

    <div class="recover-buttons">
        <button class="btn" type="submit">Återställ</button>
        <a class="btn" href="{{mountpathed "login"}}">Avbryt</a>
    </div>
</form>
