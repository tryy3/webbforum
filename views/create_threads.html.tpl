{{define "pagetitle"}}Blogs - Index{{end}}

<div class="page">
    <div class="page-body">
        {{if .error}}
            <div class="alert alert-danger">{{translate .error}}</div>
        {{end}}

        <h1>Skapa en ny tråd</h1>

        <div class="new-post">
            <form method="post">
                <input type="hidden" name="csrf_token" value="{{.csrf_token}}" />

                {{with .errs}}{{with $errlist := index . "thread_name"}}<span class="has-error">{{end}}{{end}}
                <span class="input-title">Tråd namn</span>
                <input type="text" class="form-control" name="thread_name" placeholder="Namn" value="{{.thread_name}}" />
                {{with .errs}}{{with $errlist := index . "thread_name"}}{{range $errlist}}<div class="error-block">{{translate .}}</div>{{end}}</span>{{end}}{{end}}

                {{with .errs}}{{with $errlist := index . "thread_message"}}<span class="has-error">{{end}}{{end}}
                <span class="input-title">Kommentar</span>
                <textarea type="text" class="form-control" name="thread_message" id="commentField" style="width:100%;height:300px;">{{.thread_message}}</textarea>
                {{with .errs}}{{with $errlist := index . "thread_message"}}{{range $errlist}}<div class="error-block">{{translate .}}</div>{{end}}</span>{{end}}{{end}}

                <div class="buttonwrap"><button class="btn" type="submit">Skapa en ny tråd</button></div>
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