{{define "pagetitle"}}Försvarsmakten forum{{end}}

<div class="page">
    <div class="page-toolbar">
        <form action="/forums/thread/{{.DisplayName}}/search" method="POST">
            <input type="text" class="search" name="search" placeholder="Sök">
            <input type="submit" value="Sök">
        </form>
        <a href="/forums/{{.category.DisplayName}}/create_thread">Ny tråd</a>
    </div>
    <div class="page-body">
        <h1>{{.category.Name}}</h1>
        <h3>{{.category.Description}}</h3>

        {{range .threads}}
            <div class="category-panel">
                <div class="icon">
                    <i class="fas fa-circle"></i>
                </div>

                <div class="body">
                    <div class="title">
                        <a href="/forums/thread/{{.DisplayName}}">{{.Name}}</a>
                    </div>
                    <div class="description">
                        <a href="/user/{{.CreatedBy.Username}}">{{.CreatedBy.Username}}</a>
                    </div>
                </div>

                <div class="information">
                    {{formatNumber .CountPost}} Posts<br>
                    123.315 Views
                </div>

                <div class="status">
                    <div class="title">
                        {{humanDate .LatestPost}}
                    </div>
                    <div class="description">
                        by <a href="/user/tryy3">tryy3</a>
                    </div>
                </div>
            </div>
        {{end}}
    </div>
</div>