{{define "pagetitle"}}Profile - Index{{end}}

<div class="row">
	<div class="col-md-offset-1 col-md-10">
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
                <td><img src="{{.user.ProfileImageURL}}" /></td>
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