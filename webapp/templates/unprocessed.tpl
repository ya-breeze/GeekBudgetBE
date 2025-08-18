{{ template "header.tpl" . }}

<main>
    {{ with .Unprocessed }}
        <div class="container-fluid">
            <div class="row">
                <!-- First Column: Main Form -->
                <div class="col-md-4">
                    <form class="g-3" action="/web/unprocessed/convert" method="POST">
                        <input type="hidden" name="transaction_id" value="{{ .Transaction.ID }}">

                        <strong><i class="bi-calendar-event"></i> Date: {{ formatTime .Transaction.Date "2006-01-02" }}</strong>

                        <div class="mb-3">
                            <label for="description" class="form-label">Description: </label>
                            <input type="text" class="form-control" name="description" value="{{ .Transaction.Description }}">
                        </div>

                        <div class="mb-3">
                            <i class="bi-bank"></i> Movements:
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
                        </div>

                        <div class="mb-3">
                            {{ if ne .Transaction.PartnerName "" }}
                                Partner name: {{ .Transaction.PartnerName }}
                            {{ end }}
                        </div>
                        <div class="mb-3">
                            {{ if ne .Transaction.PartnerAccount "" }}
                                <i class="bi-bank"></i> Partner account: {{ .Transaction.PartnerAccount }}
                            {{ end }}
                        </div>
                        <div class="mb-3">
                            {{ if ne .Transaction.Place "" }}
                                <i class="bi-pin-map-fill"></i> Place: {{ .Transaction.Place }}
                            {{ end }}
                        </div>
                        <div class="mb-3">
                            <i class="bi-tags"></i> Tags: 
                            {{ range .Transaction.Tags }}
                                {{ . }}
                            {{ end }}
                        </div>
                        

                        <div class="mb-3 d-grid gap-2">
                            <button class="btn btn-success btn-lg" type="submit">Convert</button>
                            <a href="/web/matchers/edit?transaction_id={{ .Transaction.ID }}" class="btn btn-outline-primary">Create matcher</a>
                            <a href="?id={{ .Transaction.ID }}" class="btn btn-outline-secondary">Skip ({{ decrease $.UnprocessedCount }} left)</a>
                        </div>
                    </form>
                </div>

                <!-- Second Column: Matched Transactions -->
                {{ with .Matched }}
                    <div class="col-md-4">
                        <h3>Matched</h3>
                        <div class="border rounded p-3" style="max-height: 70vh; overflow-y: auto;">
                            {{ range . }}
                                <div class="card mb-3">
                                    <div class="card-body">
                                        <form action="/web/unprocessed/convert" method="POST">
                                            <input type="hidden" name="transaction_id" value="{{ .Transaction.ID }}">
                                            <input type="hidden" name="matcher_id" value="{{ .MatcherID }}">
                                            <input type="hidden" name="other_matchers" value="{{ join .OtherMatcherIDs "," }}">
                                            {{ range $i, $m := .Transaction.Movements }}
                                                <input type="hidden" name="account_{{ $i }}" value="{{ $m.AccountID }}">
                                            {{ end }}
                                            <div class="d-grid gap-2 mb-2">
                                                <button type="submit" class="btn btn-success">Convert</button>
                                            </div>
                                        </form>
                                        <h6 class="card-title">{{ formatTime .Transaction.Date "2006-01-02" }} {{ .Transaction.Description }}</h6>
                                        {{ if ne .Transaction.PartnerName "" }}
                                            <small class="text-muted">Partner: {{ .Transaction.PartnerName }}</small><br>
                                        {{ end }}
                                        {{ if ne .Transaction.PartnerAccount "" }}
                                            <small class="text-muted"><i class="bi-bank"></i> Account: {{ .Transaction.PartnerAccount }}</small><br>
                                        {{ end }}
                                        <small class="text-muted">Tags: 
                                        {{ range .Transaction.Tags }}
                                            {{ . }},
                                        {{ end }}
                                        </small>
                                        <div class="mt-2">
                                            {{ range .Transaction.Movements }}
                                                <small class="text-primary">{{ .Amount }} {{ .CurrencyName }} [{{ .AccountName }}]</small><br>
                                            {{ end }}
                                        </div>
                                    </div>
                                </div>
                            {{ end }}
                        </div>
                    </div>
                {{ end }}

                <!-- Third Column: Duplicates -->
                {{ if .Duplicates }}
                    <div class="col-md-4">
                        <h3>Duplicates</h3>
                        <div class="border rounded p-3" style="max-height: 70vh; overflow-y: auto;">
                            {{ range .Duplicates }}
                                <div class="card mb-3">
                                    <div class="card-body">
                                        <h6 class="card-title">{{ formatTime .Date "2006-01-02" }} {{ .Description }}</h6>
                                        <div class="mb-2">
                                            {{ range .Movements }}
                                                <small class="text-primary">{{ .Amount }} {{ .CurrencyName }} [{{ .AccountName }}]</small><br>
                                            {{ end }}
                                        </div>
                                        <div class="d-grid gap-2 mb-2">
                                            <a href="/web/unprocessed/delete?id={{$.Unprocessed.Transaction.ID}}&duplicateOf={{.ID}}" class="btn btn-primary btn-sm" tabindex="-1" role="button">
                                                Mark as duplicate
                                            </a>
                                        </div>                                        
                                    </div>
                                </div>
                            {{ end }}
                        </div>
                    </div>
                {{ end }}
            </div>
        </div>
    {{ end }}

</main>

{{ template "footer.tpl" . }}