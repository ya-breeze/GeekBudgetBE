{{ template "header.tpl" . }}

<main>
    <h2>Bank Importers</h2>
    {{ range .BankImporters }}
    <div class="card" style="width: 18rem;">
        <div class="card-body">
            <h5 class="card-title">{{ .Name }}</h5>
            <p class="card-text">{{ .Description }}</p>
            <p><strong>Last successful import:</strong> {{ formatTime .LastSuccessfulImport "2006-01-02 15:04" }}</p>
        </div>
    </div>
    {{ end }}

</main>

{{ template "footer.tpl" . }}