{{define "pagetitle"}}Forum Admin{{end}}

<div class="page">
    <div class="page-body">
        {{if .error}}
            <div class="alert alert-danger">{{translate .error}}</div>
        {{end}}

        <h1>Skapa en ny kategori</h1>

        <div class="category-form">
            <form action="/admin/category/create" method="post">
                <input type="hidden" name="csrf_token" value="{{.csrf_token}}" />

                {{with .errs}}{{with $errlist := index . "category_name"}}<span class="has-error">{{end}}{{end}}
                <span class="input-title">Kategori namn</span>
                <input type="text" class="form-control" name="category_name" placeholder="Kategori namn" value="{{.category_name}}" />
                {{with .errs}}{{with $errlist := index . "category_name"}}{{range $errlist}}<div class="error-block">{{translate .}}</div>{{end}}</span>{{end}}{{end}}

                {{with .errs}}{{with $errlist := index . "category_description"}}<span class="has-error">{{end}}{{end}}
                <span class="input-title">Beskrivning</span>
                <textarea type="text" class="form-control" name="category_description" id="commentField" style="width:100%;height:300px;">{{.category_description}}</textarea>
                {{with .errs}}{{with $errlist := index . "category_description"}}{{range $errlist}}<div class="error-block">{{translate .}}</div>{{end}}</span>{{end}}{{end}}

                <div class="buttonwrap"><button class="btn" type="submit">Skapa en ny kategori</button></div>
            </form>
        </div>

        <br>
        <br>
        <br>
        <br>

        {{$csrf_token := .csrf_token}}
        {{$catErrs := .category_errs}}
        {{$cat := .category_edit}}

        <h1>Uppdatera kategorier</h1>

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
        {{else}}
            <h3>Inga kategorier hittades</h3>
        {{end}}

        <br>
        <br>
        <br>
        <br>

        <h1>Skapa en ny grupp</h1>

        <div class="category-form">
            <form action="/admin/group/create" method="post">
                <input type="hidden" name="csrf_token" value="{{.csrf_token}}" />

            {{with .errs}}{{with $errlist := index . "group_name"}}<span class="has-error">{{end}}{{end}}
                <span class="input-title">Grupp namn</span>
                <input type="text" class="form-control" name="group_name" placeholder="Grupp namn" value="{{.group_name}}" />
            {{with .errs}}{{with $errlist := index . "group_name"}}{{range $errlist}}<div class="error-block">{{translate .}}</div>{{end}}</span>{{end}}{{end}}

            {{with .errs}}{{with $errlist := index . "group_description"}}<span class="has-error">{{end}}{{end}}
                <span class="input-title">Beskrivning</span>
                <textarea type="text" class="form-control" name="group_description" id="commentField" style="width:100%;height:300px;">{{.group_description}}</textarea>
            {{with .errs}}{{with $errlist := index . "group_description"}}{{range $errlist}}<div class="error-block">{{translate .}}</div>{{end}}</span>{{end}}{{end}}

                <div class="buttonwrap"><button class="btn" type="submit">Skapa en ny kategori</button></div>
            </form>
        </div>

        <br>
        <br>
        <br>
        <br>

        {{$groupsErrs := .group_errs}}
        {{$groups := .group_edit}}

        <h1>Uppdatera Grupper</h1>

        {{range $index, $group := .groups}}
            {{with $groups}}
                {{if eq $cat.ID $group.ID}}
                    {{template "group_edit" map "group" $cat "csrf_token" $csrf_token "errs" $groupsErrs}}
                {{else}}
                    {{template "group_edit" map "group" $group "csrf_token" $csrf_token}}
                {{end}}
            {{else}}
                {{template "group_edit" map "group" $group "csrf_token" $csrf_token}}
            {{end}}
        {{else}}
            <h3>Inga grupper hittades</h3>
        {{end}}
    </div>
</div>

{{define "category_edit"}}
    <div class="category-form">
        <form action="/admin/category/modify" method="post">
            <input type="hidden" name="csrf_token" value="{{.csrf_token}}" />
            <input type="hidden" name="category_id" value="{{.category.ID}}" />

        {{with .errs}}{{with $errlist := index . "category_name"}}<span class="has-error">{{end}}{{end}}
            <span class="input-title">Kategori namn</span>
            <input type="text" class="form-control" name="category_name" placeholder="Kategori namn" value="{{.category.Name}}" />
        {{with .errs}}{{with $errlist := index . "category_name"}}{{range $errlist}}<div class="error-block">{{translate .}}</div>{{end}}</span>{{end}}{{end}}

        {{with .errs}}{{with $errlist := index . "category_description"}}<span class="has-error">{{end}}{{end}}
            <span class="input-title">Beskrivning</span>
            <textarea type="text" class="form-control" name="category_description" id="commentField" style="width:100%;height:300px;">{{.category.Description}}</textarea>
        {{with .errs}}{{with $errlist := index . "category_description"}}{{range $errlist}}<div class="error-block">{{translate .}}</div>{{end}}</span>{{end}}{{end}}

            <div class="buttonwrap">
                <button class="btn" type="submit">Uppdatera kategorin</button>
                <a href="/admin/category/remove/{{.category.ID}}" class="btn">Ta bort kategorin</a>
            </div>
        </form>
    </div>
{{end}}

{{define "group_edit"}}
    <div class="group-form">
        <form action="/admin/group/modify" method="post">
            <input type="hidden" name="csrf_token" value="{{.csrf_token}}" />
            <input type="hidden" name="group_id" value="{{.group.ID}}" />

        {{with .errs}}{{with $errlist := index . "group_name"}}<span class="has-error">{{end}}{{end}}
            <span class="input-title">Kategori namn</span>
            <input type="text" class="form-control" name="group_name" placeholder="Kategori namn" value="{{.group.Name}}" />
        {{with .errs}}{{with $errlist := index . "group_name"}}{{range $errlist}}<div class="error-block">{{translate .}}</div>{{end}}</span>{{end}}{{end}}

        {{with .errs}}{{with $errlist := index . "group_description"}}<span class="has-error">{{end}}{{end}}
            <span class="input-title">Beskrivning</span>
            <textarea type="text" class="form-control" name="group_description" id="commentField" style="width:100%;height:300px;">{{.group.Description}}</textarea>
        {{with .errs}}{{with $errlist := index . "group_description"}}{{range $errlist}}<div class="error-block">{{translate .}}</div>{{end}}</span>{{end}}{{end}}

            {{$group := .group}}

            {{range .group.Permissions}}
                <div class="permissionWrapper">
                    <span class="rememberCheckbox">
                        <input type="checkbox" id="{{$group.ID}}_rememberCheckboxInput_{{.Permission}}" class="input-box" name="permission_{{.Permission}}" {{if .Has}}checked{{end}}>
                        <label for="{{$group.ID}}_rememberCheckboxInput_{{.Permission}}"></label>
                    </span>
                    <span class="remember">{{.Title}}</span>
                </div>
            {{end}}

            <div class="buttonwrap">
                <button class="btn" type="submit">Uppdatera kategorin</button>
                <a href="/admin/group/remove/{{.group.ID}}" class="btn">Ta bort kategorin</a>
            </div>
        </form>
    </div>
{{end}}


<script>
    var textarea = document.getElementById("commentField")
    sceditor.create(textarea, {
        format: "bbcode",
        style: "/css/content/default.min.css",
        emoticonsRoot: "/images/"
    })
</script>