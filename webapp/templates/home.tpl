{{ template "header.tpl" . }}

<main>
    <h1>{{ .Title }}</h1>

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
    <p>Welcome to the home page of our web app.</p>

    <h2>Accounts</h2>
    {{ range .Accounts }}
    <div>
        {{ .Name }} - {{ .Description }}
    </div>
    {{ end }}

    {{ with .Expenses }}
    <h2>Expenses</h2>
    <h3>From: {{ formatTime .From "2006-01-02" }}</h3>
    <h3>To: {{ formatTime .To "2006-01-02" }}</h3>

    {{ range .Currencies }}
    <h4>Currency: {{ .CurrencyName }}</h4>

    <table>
        <thead>
            <tr>
                <th style="border: 1px solid black;"></th>
                {{ range .Intervals }}
                <th style="border: 1px solid black;">{{ formatTime . "2006-01-02" }}</th>
                {{ end }}
            </tr>
        </thead>
        <tbody>
            {{ range .Accounts }}
            <tr>
                <td style="border: 1px solid black;">{{ .AccountName }}</td>
                {{ range .Amounts }}
                <td style="border: 1px solid black;">{{ . }}</td>
                {{ end }}
            </tr>
            {{ end }}
        </tbody>
    </table>

    {{ end }}

    {{ end }}

</main>

{{ template "footer.tpl" . }}