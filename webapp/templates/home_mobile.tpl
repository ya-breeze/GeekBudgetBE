{{ template "header.tpl" . }}

<main>
    {{ with .UserID }}
    {{ else }}
        <h2>This is the Home Page. Please login</h2>
        <form action="/" method="POST">
            <label for="username">Username:</label>
            <input type="text" id="username" name="username" required>
            <br>
            <label for="password">Password:</label>
            <input type="password" id="password" name="password" required>
            <br>
            <button type="submit">Login</button>
        </form>
    {{ end }}

    {{ with .Expenses }}
        <h2>Expenses</h2>
        <a href="/?from={{$.Last}}" class="btn btn-primary" tabindex="-1" role="button">
            <i class="bi-arrow-left-circle-fill"></i>
        </a>
        {{ formatTime $.From "2006-01-02" }} - {{ formatTime $.To "2006-01-02" }}
        <a href="/?from={{$.Next}}" class="btn btn-primary" tabindex="-1" role="button">
            <i class="bi-arrow-right-circle-fill"></i>
        </a>

        {{ range .Currencies }}
            <h4>Currency: {{ .CurrencyName }}</h4>

            {{ $last_index := decrease (len .Intervals) }}
            <table class="table table-sm table-hover">
                <thead>
                    <tr>
                        <th></th>
                        <th>{{ formatTime (index .Intervals $last_index) "2006-01-02" }}</th>
                    </tr>
                </thead>
                <tbody class="table-group-divider">
                    {{ $intervals := .Intervals }}
                    {{ range .Accounts }}
                    <tr>
                        {{ $accountID := .AccountID }}
                        <td>{{ .AccountName }}</td>
                        {{ $a := index .Amounts $last_index }}
                        <td {{ if and (eq $accountID "") (ne $a 0.0) }}class="table-danger"{{end}}>
                            {{ if eq $a 0.0 }}
                                0
                            {{ else }}
                                <a class="link-dark link-offset-2 link-offset-3-hover link-underline link-underline-opacity-0 link-underline-opacity-75-hover"
                                    href="/web/transactions?from={{timestamp (index $intervals $last_index)}}&accountID={{$accountID}}">
                                        {{ money $a }}
                                </a>                    
                            {{ end }}
                        </td>
                        {{ end }}
                    </tr>
                </tbody>
            </table>
        {{ end }}
    {{ end }}

</main>

{{ template "footer.tpl" . }}