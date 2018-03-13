{{define "pagetitle"}}Blogs - Index{{end}}

{{$loggedin := .loggedin}}
{{if $loggedin}}
<div class="row" style="margin-bottom: 20px;">
    <div class="col-md-offset-9 col-md-2 text-right">
        <a class="btn btn-primary" href="/blogs/new"><i class="fa fa-plus"></i> New Post</a>
    </div>
</div>
{{end}}

<div class="row">
    <div class="col-md-offset-1 col-md-10">
    {{range .threads}}
        <div class="panel panel-info">
            <div class="panel-heading">
                <div class="row">
                {{.Name}}
                    <div class="col-md-6"></div>
                    <div class="col-md-6 text-right">
                    {{if $loggedin}}
                        <a class="btn btn-xs btn-link" href="/blogs/edit">Test</a>
                    {{end}}
                    </div>
                </div>
            </div>
            <div class="panel-body">
                <form action="/thread/create" method="post">

                </form>
            </div>
            <div class="panel-footer">
                <div class="row">
                    <div class="col-md-6">Beep</div>
                    <div class="col-md-6 text-right">Boop</div>
                </div>
            </div>
        </div>
    {{end}}
    </div>
</div>

{{define "thread_create"}}
<form action="/thread/create" method="post">
    <input type="hidden" name="csrf_token" value="{{.csrf_token}}" />
    <input type="hidden" name="category_id" value="{{.category.ID}}" />
    <div class="form-group {{with .errs}}{{with $errlist := index . "category_name"}}has-error{{end}}{{end}}">
        <input type="text" class="form-control" name="category_name" placeholder="Namn" value="{{.category.Name}}" />
    {{with .errs}}{{with $errlist := index . "category_name"}}{{range $errlist}}<span class="help-block">{{.}}</span>{{end}}{{end}}{{end}}
    </div>
    <div class="form-group {{with .errs}}{{with $errlist := index . "category_description"}}has-error{{end}}{{end}}">
        <input type="text" class="form-control" name="category_description" placeholder="Beskrivning" value="{{.category.Description}}" />
    {{with .errs}}{{with $errlist := index . "category_description"}}{{range $errlist}}<span class="help-block">{{.}}</span>{{end}}{{end}}{{end}}
    </div>
    <div class="row">
        <div class="col-md-offset-1 col-md-10">
            <button class="btn btn-primary btn-block" type="submit">Uppdatera kategorin</button>
        </div>
    </div>
</form>
{{end}}

{{define "thread_edit"}}
<form action="/admin/kategori/uppdatera" method="post">
    <input type="hidden" name="csrf_token" value="{{.csrf_token}}" />
    <input type="hidden" name="category_id" value="{{.category.ID}}" />
    <div class="form-group {{with .errs}}{{with $errlist := index . "category_name"}}has-error{{end}}{{end}}">
        <input type="text" class="form-control" name="category_name" placeholder="Namn" value="{{.category.Name}}" />
    {{with .errs}}{{with $errlist := index . "category_name"}}{{range $errlist}}<span class="help-block">{{.}}</span>{{end}}{{end}}{{end}}
    </div>
    <div class="form-group {{with .errs}}{{with $errlist := index . "category_description"}}has-error{{end}}{{end}}">
        <input type="text" class="form-control" name="category_description" placeholder="Beskrivning" value="{{.category.Description}}" />
    {{with .errs}}{{with $errlist := index . "category_description"}}{{range $errlist}}<span class="help-block">{{.}}</span>{{end}}{{end}}{{end}}
    </div>
    <div class="row">
        <div class="col-md-offset-1 col-md-10">
            <button class="btn btn-primary btn-block" type="submit">Uppdatera kategorin</button>
        </div>
    </div>
</form>
{{end}}

{{define "thread_delete"}}
<form action="/admin/kategori/ta_bort" method="post">
    <input type="hidden" name="csrf_token" value="{{.csrf_token}}" />
    <input type="hidden" name="category_id" value="{{.category.ID}}" />
    <h3>{{.category.Name}}</h3>
    <div class="row">
        <div class="col-md-offset-1 col-md-10">
            <button class="btn btn-primary btn-block" type="submit" name="delete">Ta bort kategori</button>
        </div>
    </div>
</form>
{{end}}