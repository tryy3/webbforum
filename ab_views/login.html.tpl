<form method="POST" class="box">
    <div class="box-title">Logga in</div>

    {{if .error}}
        <div class="alert alert-danger">
            {{translate .error}}
        </div>
    {{end}}
    <input type="hidden" name="{{.xsrfName}}" value="{{.xsrfToken}}" />

    <span class="input-title">Användarnamn</span>
    <input type="text" name="{{.primaryID}}" placeholder="Användarnamn" value="{{.primaryIDValue}}">

    <span class="input-title">Lösenord</span>
    <input  type="password" class="input-box" name="password" placeholder="Lösenord">

    {{if .showRemember}}
        <div class="rememberWrapper">
            <span class="rememberCheckbox">
                <input type="checkbox" id="rememberCheckboxInput" class="input-box" name="rm" value="true">
                <label for="rememberCheckboxInput"></label>
            </span>
            <span class="remember">Kom ihåg mig</span>
        </div>
    {{end}}

    <div class="btnWrapper"><button class="btn" type="submit">Logga in</button></div>

    {{if .showRecover}}
        <div class="forgotWrapper"><a class="forgot" href="{{mountpathed "recover"}}">Glömt lösenord?</a></div>
    {{end}}
</form>
