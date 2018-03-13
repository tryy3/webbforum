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
		{{range .categories}}
		<div class="panel panel-info">
			<div class="panel-heading">
				<div class="row">
                    <a href="/forums/{{.Name}}">{{.Name}}</a>
				</div>
			</div>
			<div class="panel-body">
                {{.Description}}
            </div>
		</div>
		{{end}}
	</div>
</div>