{{define "pagetitle"}}Profile - Index{{end}}

<div class="page">
    <div class="page-body">
        <h1>{{.user.Username}}</h1>

        <div class="profile-body">
            <div class="left">
                <img src="{{.user.ProfileImageURL}}" alt="{{.user.Username}}"><br>
                32,421 Inlägg<br>
                3,324 Trådar
            </div>
            <div class="right">
                <b>Namn: </b>{{.user.FirstName}} {{.user.LastName}}<br>
                <b>E-postadress: </b>{{.user.Email}}<br>
                <b>Grupp: </b>
                {{if .user.Group}}
                    {{.user.Group.Name}}
                {{else}}
                    User
                {{end}}
                <br>
                <br>
                <b>Beskrivning: </b>{{.user.ParsedDescription}}
            </div>
        </div>
    </div>
</div>