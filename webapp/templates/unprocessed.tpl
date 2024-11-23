{{ template "header.tpl" . }}

<main>
    <h2>Unprocessed</h2>
    {{ with .Unprocessed }}
        <form method="GET">
            <input type="hidden" name="id" value="{{ .Transaction.ID }}">
            <button type="submit">Skip ({{ decrease $.UnprocessedCount }} left)</button>
        </form>

        <form action="/web/matchers/edit" method="GET">
            <input type="hidden" name="transaction_id" value="{{ .Transaction.ID }}">
            <button type="submit">Create matcher</button>
        </form>

        <form action="/web/unprocessed/convert" method="POST">
            <input type="hidden" name="transaction_id" value="{{ .Transaction.ID }}">
            <h5>{{ formatTime .Transaction.Date "2006-01-02" }}</h5>
            <div class="mb-3">
                <label for="description" class="form-label">Description: </label>
                <input type="text" class="form-control" name="description" value="{{ .Transaction.Description }}">
            </div>

            {{ if ne .Transaction.PartnerName "" }}
            <h6>Partner name: {{ .Transaction.PartnerName }}</h6>
            {{ end }}
            {{ if ne .Transaction.PartnerAccount "" }}
            <h6>Partner account: {{ .Transaction.PartnerAccount }}</h6>
            {{ end }}
            {{ if ne .Transaction.Place "" }}
            <h6>Place: {{ .Transaction.Place }}</h6>
            {{ end }}
            <h6>Tags: 
            {{ range .Transaction.Tags }}
            {{ . }}
            {{ end }}
            </h6>
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
                {{ if ne .Transaction.PartnerName "" }}
                    <h6>Partner name: {{ .Transaction.PartnerName }}</h6>
                {{ end }}
                {{ if ne .Transaction.PartnerAccount "" }}
                    <h6>Partner account: {{ .Transaction.PartnerAccount }}</h6>
                {{ end }}
                <h6>Tags: 
                {{ range .Transaction.Tags }}
                    {{ . }},
                {{ end }}
                </h6>
                <p>
                    {{ range .Transaction.Movements }}
                        {{ .Amount }} {{ .CurrencyName }} [{{ .AccountName }}] <br>
                    {{ end }}
                </p>
            {{ end }}
        {{ end }}


        {{ if .Duplicates }}
            <h3>Duplicates</h3>
            {{ range .Duplicates }}
                <h5>{{ formatTime .Date "2006-01-02" }} {{ .Description }}</h5>
                <p>
                    {{ range .Movements }}
                        {{ .Amount }} {{ .CurrencyName }} [{{ .AccountName }}] <br>
                    {{ end }}
                </p>
                <a href="/web/unprocessed/delete?id={{$.Unprocessed.Transaction.ID}}&duplicateOf={{.ID}}" class="btn btn-primary" tabindex="-1" role="button">
                    Mark as duplicate
                </a>
            {{ end }}
        {{ end }}

    {{ end }}

</main>

{{ template "footer.tpl" . }}