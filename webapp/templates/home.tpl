{{ template "header.tpl" . }}

<main>
    {{ with .UserID }}
        <h2>This is the Home Page for {{ . }}</h2>
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
        <h3>From: {{ formatTime .From "2006-01-02" }}</h3>
        <h3>To: {{ formatTime .To "2006-01-02" }}</h3>

        {{ range .Currencies }}
            <h4>Currency: {{ .CurrencyName }}</h4>

            <table class="table table-sm table-hover">
                <thead>
                    <tr>
                        <th></th>
                        {{ range .Intervals }}
                        <th>{{ formatTime . "2006-01-02" }}</th>
                        {{ end }}
                        <th>Total for year</th>
                    </tr>
                </thead>
                <tbody class="table-group-divider">
                    {{ $intervals := .Intervals }}
                    {{ range .Accounts }}
                    <tr>
                        {{ $accountID := .AccountID }}
                        <td>{{ .AccountName }}</td>
                        {{ range $i, $a := .Amounts }}
                        <td>
                            {{ if eq $a 0.0 }}
                            0
                            {{ else }}
                                <a class="link-dark link-offset-2 link-offset-3-hover link-underline link-underline-opacity-0 link-underline-opacity-75-hover"
                                    href="/web/transactions?from={{timestamp (index $intervals $i)}}&accountID={{$accountID}}">
                                        {{ money $a }}
                                </a>                    
                            {{ end }}
                        </td>
                        {{ end }}
                        <td>
                            <strong>{{ money .TotalForYear }}</strong>
                        </td>
                    </tr>
                    {{ end }}
                </tbody>
                <tfoot class="table-secondary">
                        <td>
                            <strong>Total for month</strong>
                        </td>
                        {{ range .Total }}
                            <td>{{ money . }}</td>
                        {{ end }}
                </tfoot>
            </table>
        {{ end }}
    {{ end }}

</main>

{{ template "footer.tpl" . }}