{{define "pagetitle"}}Försvarsmakten forum{{end}}

<div class="page">
    <div class="page-body">
        <h1>Kategorier</h1>
        {{range .categories}}
            <div class="category-panel">
                <div class="icon">
                    <i class="fas fa-circle"></i>
                </div>

                <div class="body">
                    <div class="title">
                        <a href="/forums/{{.DisplayName}}">{{.Name}}</a>
                    </div>
                    <div class="description">
                        {{.Description}}
                    </div>
                </div>

                <div class="information">
                    {{formatNumber .CountThreads}} Trådar<br>
                    {{formatNumber .CountPost}} Kommentarer
                </div>

                <div class="status">
                    {{if not .LatestUpdate}}
                        <div class="title"><a href="#">Ingen har skapat en tråd</a></div>
                    {{else}}
                        <div class="title">
                            <a href="/forums/threads/{{.LatestUpdate.Thread.DisplayName}}">{{.LatestUpdate.Thread.Name}}</a>
                        </div>
                        <div class="description">
                            Av <a href="/user/{{.LatestUpdate.User.Username}}">{{.LatestUpdate.User.Username}}</a> {{humanDate .LatestUpdate.UpdatedAt}}
                        </div>
                    {{end}}
                </div>
            </div>
        {{end}}
    </div>
</div>