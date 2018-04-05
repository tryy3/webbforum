<form method="POST" class="box">
    <div class="box-title">Kontoåterställning</div>

    {{if .error}}
        <div class="alert alert-danger">
            {{translate .error}}
        </div>
    {{end}}

    <input type="hidden" name="{{.xsrfName}}" value="{{.xsrfToken}}" />


    {{$pid := .primaryID}}
    {{with .errs}}{{with $errlist := index . $pid}}<span class="has-error">{{end}}{{end}}
        <span class="input-title">Användarnamn</span>
        <input type="text" name="{{$pid}}" placeholder="Användarnamn" value="{{.primaryIDValue}}"/>
    {{with .errs}}{{with $errlist := index . $pid}}{{range $errlist}}<div class="error-block">{{translate .}}</div>{{end}}</span>{{end}}{{end}}

    {{$cpid := .primaryID | printf "confirm_%s"}}
    {{with .errs}}{{with $errlist := index . $cpid}}<span class="has-error">{{end}}{{end}}
        <span class="input-title">Bekräfta användarnamn</span>
        <input type="text" name="{{$cpid}}" placeholder="Bekräfta användarnamn" value="{{.confirmPrimaryIDValue}}"/>
    {{with .errs}}{{with $errlist := index . $cpid}}{{range $errlist}}<div class="error-block">{{translate .}}</div>{{end}}</span>{{end}}{{end}}

    <div class="recover-buttons">
        <button class="btn" type="submit">Återställ</button>
        <a class="btn" href="{{mountpathed "login"}}">Avbryt</a>
    </div>
</form>
