{{define "pagetitle"}}admin - Index{{end}}

<div class="row">
	<div class="col-md-offset-1 col-md-10">
        {{if .error}}
            <div class="alert alert-danger">{{.error}}</div>
        {{end}}
        <form action="/admin/category/create" method="post">
            <input type="hidden" name="csrf_token" value="{{.csrf_token}}" />
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
                    <button class="btn btn-primary btn-block" type="submit">Skapa en ny kategori</button>
                </div>
            </div>
        </form>

        <br>
        <br>
        <br>
        <br>

        {{$csrf_token := .csrf_token}}
        {{$catErrs := .category_errs}}
        {{$cat := .category_edit}}

        {{range $index, $category := .categories}}
            {{with $cat}}
                {{if eq $cat.ID $category.ID}}
                    {{template "category_edit" map "category" $cat "csrf_token" $csrf_token "errs" $catErrs}}
                {{else}}
                    {{template "category_edit" map "category" $category "csrf_token" $csrf_token}}
                {{end}}
            {{else}}
                {{template "category_edit" map "category" $category "csrf_token" $csrf_token}}
            {{end}}

            {{template "category_delete" map "category" $category "csrf_token" $csrf_token}}
        {{else}}
            <h3>No categories found</h3>
        {{end}}

        {{define "category_edit"}}
            <form action="/admin/category/modify" method="post">
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

        {{define "category_delete"}}
            <form action="/admin/category/remove" method="post">
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
	</div>
</div>