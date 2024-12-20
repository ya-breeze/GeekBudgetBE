<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ block "title" . }}My Web App{{ end }}</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet"
        integrity="sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH" crossorigin="anonymous">
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"
        integrity="sha384-YvpcrYf0tY3lHB60NNkmXc5s9fDVZLESaAA55NDzOxhy9GkcIdslK1eN7N6jIeHz"
        crossorigin="anonymous"></script>
</head>

<body>
    <header>
        <h1>Welcome to My Web App</h1>
        <nav>
            <a href="/">Home</a>
            |
            <a href="/web/transactions">Transactions</a>
            |
            <a href="/web/accounts">Accounts</a>
            |
            <a href="/web/bank-importers">Importers</a>
            |
            <a href="/web/matchers">Matchers</a>
            |
            <a href="/web/unprocessed">Unprocessed</a>
            |
            <a href="/web/about">About</a>
        </nav>
    </header>