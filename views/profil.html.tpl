{{define "pagetitle"}}Profile - Index{{end}}

<div class="row">
	<div class="col-md-offset-1 col-md-10">
        {{if .error}}
            <div class="alert alert-danger">{{.error}}</div>
        {{end}}
        <form method="POST">
            <input type="hidden" name="csrf_token" value="{{.csrf_token}}" />
            <div class="form-group {{with .errs}}{{with $errlist := index . "email"}}has-error{{end}}{{end}}">
                <input type="text" class="form-control" name="email" placeholder="Email" value="{{.user.Email}}" />
            {{with .errs}}{{with $errlist := index . "email"}}{{range $errlist}}<span class="help-block">{{.}}</span>{{end}}{{end}}{{end}}
            </div>
            <div class="form-group {{with .errs}}{{with $errlist := index . "first_name"}}has-error{{end}}{{end}}">
                <input type="text" class="form-control" name="first_name" placeholder="Namn" value="{{.user.FirstName}}" />
            {{with .errs}}{{with $errlist := index . "first_name"}}{{range $errlist}}<span class="help-block">{{.}}</span>{{end}}{{end}}{{end}}
            </div>
            <div class="form-group {{with .errs}}{{with $errlist := index . "last_name"}}has-error{{end}}{{end}}">
                <input type="text" class="form-control" name="last_name" placeholder="Efter namn" value="{{.user.LastName}}" />
            {{with .errs}}{{with $errlist := index . "last_name"}}{{range $errlist}}<span class="help-block">{{.}}</span>{{end}}{{end}}{{end}}
            </div>
            <div class="form-group {{with .errs}}{{with $errlist := index . "description"}}has-error{{end}}{{end}}">
                <input type="text" class="form-control" name="description" placeholder="Beskrivning" value="{{.user.Description}}" />
            {{with .errs}}{{with $errlist := index . "description"}}{{range $errlist}}<span class="help-block">{{.}}</span>{{end}}{{end}}{{end}}
            </div>
            <div class="row">
                <div class="col-md-offset-1 col-md-10">
                    <button class="btn btn-primary btn-block" type="submit">Updatera information</button>
                </div>
            </div>
        </form>

        <form action="/profil/upload" enctype="multipart/form-data" method="POST">
            <input type="file" name="uploadfile" />
            <input type="hidden" name="csrf_token" value="{{.csrf_token}}" />
            <div class="row">
                <div class="col-md-offset-1 col-md-10">
                    <button class="btn btn-primary btn-block" type="submit">Ladda upp bilden</button>
                </div>
            </div>
        </form>
        <table style="width:100%">
            <tr>
                <td>Username</td>
                <td>{{.user.Username}}</td>
            </tr>
            <tr>
                <td>Email</td>
                <td>{{.user.Email}}</td>
            </tr>
            <tr>
                <td>Name</td>
                <td>{{.user.FirstName}} {{.user.LastName}}</td>
            </tr>
            <tr>
                <td>Profile Image</td>
                <td>{{.user.ProfileImage}}</td>
            </tr>
            <tr>
                <td>Description</td>
                <td>{{.user.Description}}</td>
            </tr>
            <tr>
                <td>Group</td>
                <td>{{.user.Group.Name}}</td>
            </tr>
        </table>
    </div>
</div>