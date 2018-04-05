{{define "pagetitle"}}Blogs - Index{{end}}

{{$thread := .thread}}
{{$editPost := .edit_post}}
{{$csrfToken := .csrf_token}}
{{$errs := .errs}}


<div class="page">
    <div class="page-body">
        {{if .error}}
            <div class="alert alert-danger">{{translate .error}}</div>
        {{end}}

        <h1>{{.thread.Name}}</h1>

        {{range $i, $v := .posts}}
            <div class="post-panel">
                <div class="user">
                    <div class="profile">
                        <a href="/user/{{.User.Username}}">{{.User.Username}}</a>
                    </div>
                    <div class="group">
                        {{if .User.Group}}
                            {{.User.Group.Name}}
                        {{else}}
                            User
                        {{end}}
                    </div>
                    <div class="name">
                        {{.User.FirstName}} {{.User.LastName}}
                    </div>
                    <div class="picture">
                        <img src="{{.User.ProfileImageURL}}" alt="{{.User.Username}}">
                    </div>
                </div>
                <div class="comment">
                    <div class="status">
                        <div class="left">
                            {{formatDate .CreatedAt}}
                        </div>
                        <div class="right">
                            #{{$i}}
                        </div>
                    </div>
                    <div class="message">
                        {{if not $editPost}}
                            {{.DisplayComment}}
                        {{else if eq $editPost.ID .ID}}
                            <form method="post" action="/forums/thread/{{$thread.DisplayName}}/edit/{{.ID}}">
                                <input type="hidden" name="csrf_token" value="{{$csrfToken}}" />
                                <div class="form-group {{with $errs}}{{with $errlist := index . "comment"}}has-error{{end}}{{end}}">
                                    <textarea class="form-control" name="comment" id="commentField" style="width:80%;height:300px;">{{$editPost.Comment}}</textarea>
                                {{with $errs}}{{with $errlist := index . "comment"}}{{range $errlist}}<span class="help-block">{{.}}</span>{{end}}{{end}}{{end}}
                                </div>
                                <div class="row">
                                    <div class="col-md-offset-1 col-md-10">
                                        <button class="btn btn-primary btn-block" type="submit">Modifiera kommentar</button>
                                    </div>
                                </div>
                            </form>
                        {{else}}
                            {{.DisplayComment}}
                        {{end}}
                    </div>
                </div>
            </div>
        {{end}}

        <div class="new-post">
            <form method="post" action="/forums/thread/{{.thread.DisplayName}}/new_post">
                <input type="hidden" name="csrf_token" value="{{.csrf_token}}" />

                {{with .errs}}{{with $errlist := index . "comment"}}<span class="has-error">{{end}}{{end}}
                <textarea type="text" class="form-control" name="comment" id="commentField" style="width:100%;height:300px;">{{.comment}}</textarea>
                {{with .errs}}{{with $errlist := index . "comment"}}{{range $errlist}}<div class="error-block">{{translate .}}</div>{{end}}</span>{{end}}{{end}}

                <div class="buttonwrap"><button class="btn" type="submit">Skapa en ny kommentar</button></div>
            </form>
        </div>
    </div>
</div>


<script>
    var textarea = document.getElementById("commentField")
    sceditor.create(textarea, {
        format: "bbcode",
        style: "/css/content/default.min.css",
        emoticonsRoot: "/images/"
    })
</script>