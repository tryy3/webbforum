{{define "pagetitle"}}Profile - Index{{end}}

<div class="page">
	<div class="page-body">
        {{if .error}}
            <div class="alert alert-danger">{{.error}}</div>
        {{end}}

        <h1>Inställningar</h1>

        <div class="profileWrapper"><a class="profile" href="/user/{{.user.Username}}">Publik profil</a></div>

        <div class="settings">
            <form method="POST">
                <input type="hidden" name="csrf_token" value="{{.csrf_token}}" />

                {{with .errs}}{{with $errlist := index . "email"}}<span class="has-error">{{end}}{{end}}
                <span class="input-title">E-postadress</span>
                <input type="text" class="form-control" name="email" placeholder="E-postadress" value="{{.user.Email}}" />
                {{with .errs}}{{with $errlist := index . "email"}}{{range $errlist}}<div class="error-block">{{translate .}}</div>{{end}}</span>{{end}}{{end}}


                {{with .errs}}{{with $errlist := index . "first_name"}}<span class="has-error">{{end}}{{end}}
                <span class="input-title">Namn</span>
                <input type="text" class="form-control" name="first_name" placeholder="Namn" value="{{.user.FirstName}}" />
                {{with .errs}}{{with $errlist := index . "first_name"}}{{range $errlist}}<div class="error-block">{{translate .}}</div>{{end}}</span>{{end}}{{end}}

                {{with .errs}}{{with $errlist := index . "last_name"}}<span class="has-error">{{end}}{{end}}
                <span class="input-title">Efter namn</span>
                <input type="text" class="form-control" name="last_name" placeholder="Efternamn" value="{{.user.LastName}}" />
                {{with .errs}}{{with $errlist := index . "last_name"}}{{range $errlist}}<div class="error-block">{{translate .}}</div>{{end}}</span>{{end}}{{end}}

                {{with .errs}}{{with $errlist := index . "description"}}<span class="has-error">{{end}}{{end}}
                <span class="input-title">Beskrivning</span>
                <textarea type="text" class="form-control" name="description" id="commentField" style="width:100%;height:300px;">{{.user.Description}}</textarea>
                {{with .errs}}{{with $errlist := index . "description"}}{{range $errlist}}<div class="error-block">{{translate .}}</div>{{end}}</span>{{end}}{{end}}

                <div class="buttonwrap"><button class="btn" type="submit">Updatera information</button></div>
            </form>
        </div>

        <div class="picture">
            <h1>Profil bild</h1>

            <form action="/profil/upload" enctype="multipart/form-data" method="POST">
                <input type="hidden" name="csrf_token" value="{{.csrf_token}}" />

                <div class="inputwrap">
                    <input type="file" name="uploadfile" id="uploadfile" class="inputfile" />
                    <label for="uploadfile"><i class="fas fa-upload"></i> <span>Välj en bild</span></label>
                </div>

                <div class="buttonwrap"><button class="btn" type="submit">Ladda upp en profil bild</button></div>
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

    var inputs = document.querySelectorAll( '.inputfile' );
    Array.prototype.forEach.call( inputs, function( input )
    {
        var label	 = input.nextElementSibling,
            labelVal = label.innerHTML;

        input.addEventListener( 'change', function( e )
        {
            var fileName = e.target.value.split( '\\' ).pop();

            if( fileName )
                label.querySelector( 'span' ).innerHTML = fileName;
            else
                label.innerHTML = labelVal;
        });

        input.addEventListener( 'focus', function(){ input.classList.add( 'has-focus' ); });
        input.addEventListener( 'blur', function(){ input.classList.remove( 'has-focus' ); });
    });
</script>