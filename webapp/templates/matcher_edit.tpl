{{ template "header.tpl" . }}

<main>
    <h2>Create matcher</h2>
    <div class="row">
        <!-- Left side: Matcher form -->
        <div class="col-md-6">
            <form method="POST" id="matcherForm">
                <input type="hidden" name="id" value="{{ .Matcher.Id }}">
                <input type="hidden" name="transaction_id" value="{{ .Transaction.ID }}">

                <div class="mb-3">
                    <label for="name" class="form-label">Name</label>
                    <input type="text" class="form-control" name="name" value="{{ .Matcher.Name }}">
                </div>

                <div class="mb-3">
                    <label for="outputDescription" class="form-label">Output description</label>
                    <input type="text" class="form-control" name="outputDescription" value="{{ .Matcher.OutputDescription }}">
                </div>
                <div class="mb-3">
                    <label for="outputTags" class="form-label">Output tags</label>
                    <input type="text" class="form-control" name="outputTags" value="{{ range .Matcher.OutputTags }}{{.}}, {{ end }}">
                </div>

                <div class="mb-3">
                    <label for="descriptionRegExp" class="form-label">Description regexp</label>
                    <input type="text" class="form-control" name="descriptionRegExp" value="{{ .Matcher.DescriptionRegExp }}">
                </div>
                <div class="mb-3">
                    <label for="partnerAccountNumberRegExp" class="form-label">Partner account regexp</label>
                    <input type="text" class="form-control" name="partnerAccountNumberRegExp" value="{{ .Matcher.PartnerAccountNumberRegExp }}">
                </div>
                <div class="mb-3">
                    <label for="account" class="form-label">Account</label>
                    <select class="form-select" name="account">
                        <option value="">Select account</option>
                        {{ range $.Accounts }}
                        <option value="{{ .Id }}" {{ if eq .Id $.Matcher.OutputAccountId }}selected{{ end }}>
                            {{.Type}}: {{ .Name }}
                        </option>
                        {{ end }}
                    </select>
                </div>

                {{ if ne .Transaction.ID "" }}
                <button class="btn btn-secondary" type="button" id="checkBtn">Check</button>
                {{ end }}
                <button class="btn btn-primary" type="submit">Save</button>
            </form>

            <!-- Check result display area -->
            {{ if ne .Transaction.ID "" }}
            <div id="checkResult" class="mt-3" style="display: none;">
                <div id="checkResultContent"></div>
            </div>
            {{ end }}
        </div>

        <!-- Right side: Transaction display (only visible when transaction_id is present) -->
        {{ if ne .Transaction.ID "" }}
        <div class="col-md-6">
            <h5>Transaction Details</h5>
            <div class="card">
                <div class="card-body">
                    <h6 class="card-subtitle mb-2 text-muted">{{ formatTime .Transaction.Date "2006-01-02" }}</h6>

                    <div class="mb-3">
                        <label class="form-label"><strong>Description:</strong></label>
                        <p>{{ .Transaction.Description }}</p>
                    </div>

                    {{ if ne .Transaction.PartnerName "" }}
                    <div class="mb-3">
                        <label class="form-label"><strong>Partner name:</strong></label>
                        <p>{{ .Transaction.PartnerName }}</p>
                    </div>
                    {{ end }}

                    {{ if ne .Transaction.PartnerAccount "" }}
                    <div class="mb-3">
                        <label class="form-label"><strong>Partner account:</strong></label>
                        <p>{{ .Transaction.PartnerAccount }}</p>
                    </div>
                    {{ end }}

                    {{ if ne .Transaction.Place "" }}
                    <div class="mb-3">
                        <label class="form-label"><strong>Place:</strong></label>
                        <p>{{ .Transaction.Place }}</p>
                    </div>
                    {{ end }}

                    {{ if .Transaction.Tags }}
                    <div class="mb-3">
                        <label class="form-label"><strong>Tags:</strong></label>
                        <p>{{ range .Transaction.Tags }}{{ . }}, {{ end }}</p>
                    </div>
                    {{ end }}

                    {{ if .Transaction.Movements }}
                    <div class="mb-3">
                        <label class="form-label"><strong>Movements:</strong></label>
                        {{ range $i, $m := .Transaction.Movements }}
                        <div class="card mt-2" style="width: 100%;">
                            <div class="card-body">
                                <p class="card-text">
                                    <strong>Amount:</strong> {{ $m.Amount }} {{ $m.CurrencyName }}<br>
                                    {{ if $m.AccountName }}<strong>Account:</strong> {{ $m.AccountName }}{{ end }}
                                </p>
                            </div>
                        </div>
                        {{ end }}
                    </div>
                    {{ end }}
                </div>
            </div>
        </div>
        {{ end }}
    </div>
</main>

<script>
document.addEventListener('DOMContentLoaded', function() {
    const checkBtn = document.getElementById('checkBtn');
    if (!checkBtn) return;

    checkBtn.addEventListener('click', async function() {
        // Collect matcher form data
        const form = document.getElementById('matcherForm');
        const formData = new FormData(form);

        const matcher = {
            name: formData.get('name'),
            outputDescription: formData.get('outputDescription'),
            outputTags: formData.get('outputTags').split(',').map(t => t.trim()).filter(t => t),
            descriptionRegExp: formData.get('descriptionRegExp'),
            partnerAccountNumberRegExp: formData.get('partnerAccountNumberRegExp'),
            outputAccountId: formData.get('account'),
        };

        const transactionId = formData.get('transaction_id');

        if (!transactionId) {
            alert('Transaction ID is missing');
            return;
        }

        // Show loading state
        const resultDiv = document.getElementById('checkResult');
        const resultContent = document.getElementById('checkResultContent');
        resultContent.innerHTML = '<div class="alert alert-info">Checking matcher...</div>';
        resultDiv.style.display = 'block';

        try {
            // Make web endpoint request (handles authentication via session)
            const response = await fetch('/web/matchers/check', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    matcher: matcher,
                    transaction: {
                        id: transactionId,
                    }
                })
            });

            if (!response.ok) {
                const errorText = await response.text();
                resultContent.innerHTML = '<div class="alert alert-danger">Error: ' + response.status + ' - ' + errorText + '</div>';
                return;
            }

            const result = await response.json();

            // Display result
            if (result.result) {
                resultContent.innerHTML = '<div class="alert alert-success" role="alert"><strong>Match!</strong> The matcher would match this transaction.</div>';
            } else {
                resultContent.innerHTML = '<div class="alert alert-danger" role="alert"><strong>No Match.</strong> The matcher would not match this transaction.</div>';
            }
        } catch (error) {
            resultContent.innerHTML = '<div class="alert alert-danger">Error: ' + error.message + '</div>';
        }
    });
});
</script>

{{ template "footer.tpl" . }}