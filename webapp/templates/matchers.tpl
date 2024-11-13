{{ template "header.tpl" . }}

<main>
    <h1>{{ .Title }}</h1>

    <h2>Matchers</h2>
    {{ range .Matchers }}
    <div class="card" style="width: 18rem;">
        <div class="card-body">
            <h5 class="card-title">{{ .Name }}</h5>
            <a class="link-offset-2 link-offset-3-hover link-underline link-underline-opacity-0 link-underline-opacity-75-hover"
                href="/web/matchers/edit?id={{.Id}}">
                Edit matcher
            </a>
            <form action="/web/matchers" method="DELETE">
                <input type="hidden" name="id" value="{{ .Id }}">
                <button type="submit" class="btn btn-danger">Delete</button>
        </div>
    </div>
    {{ end }}

</main>

{{ template "footer.tpl" . }}