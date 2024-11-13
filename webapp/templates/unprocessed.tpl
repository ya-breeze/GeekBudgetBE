{{ template "header.tpl" . }}

<main>
    <h1>{{ .Title }}</h1>

    <h2>Unprocessed</h2>
    {{ with .Unprocessed }}
    <form method="GET">
        <input type="hidden" name="id" value="{{ .Transaction.ID }}">
        <button type="submit">Skip</button>
    </form>

    <form action="/web/matchers/edit" method="GET">
        <input type="hidden" name="transaction_id" value="{{ .Transaction.ID }}">
        <button type="submit">Create matcher</button>
    </form>

    <form action="/web/unprocessed/convert" method="POST">
        <input type="hidden" name="transaction_id" value="{{ .Transaction.ID }}">
        <h5>{{ formatTime .Transaction.Date "2006-01-02" }} {{ .Transaction.Description }}</h5>
        <p>
            {{ range $i, $m := .Transaction.Movements }}
                <div class="mb-3">
                    <label for="account" class="form-label">{{ $m.Amount }} {{ $m.CurrencyName }}</label>
                    <select class="form-select" name="account_{{ $i }}">
                        <option value="">Select account</option>
                        {{ range $.Accounts }}
                        <option value="{{ .Id }}" {{ if eq .Id $m.AccountID }}selected{{ end }}>
                            {{.Type}}: {{ .Name }}
                        </option>
                        {{ end }}
                    </select>
                </div>
            {{ end }}
        </p>
        <button type="submit">Convert</button>
    </form>
    {{ with .Matched }}
    <h3>Matched</h3>
    {{ range . }}
    <form action="/web/unprocessed/convert" method="POST">
        <input type="hidden" name="transaction_id" value="{{ .Transaction.ID }}">
        {{ range $i, $m := .Transaction.Movements }}
        <input type="hidden" name="account_{{ $i }}" value="{{ $m.AccountID }}">
        {{ end }}
        <button type="submit">Convert</button>
    </form>
    <h5>{{ formatTime .Transaction.Date "2006-01-02" }} {{ .Transaction.Description }}</h5>
    <p>
        {{ range .Transaction.Movements }}
        {{ .Amount }} {{ .CurrencyName }} [{{ .AccountName }}] <br>
        {{ end }}
    </p>
    {{ end }}
    {{ end }}


    {{ end }}

</main>

{{ template "footer.tpl" . }}