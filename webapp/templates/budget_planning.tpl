{{template "header.tpl" .}}

<div class="container-fluid">
    <div class="row">
        <div class="col-12">
            <h2>Budget Planning</h2>
            
            {{ if (index .Query "error") }}
            <div class="alert alert-danger" role="alert">
                {{ index .Query "error" }}
            </div>
            {{ end }}

            {{ if (index .Query "success") }}
            <div class="alert alert-success" role="alert">
                {{ index .Query "success" }}
            </div>
            {{ end }}

            <!-- Month Navigation -->
            <div class="row mb-3">
                <div class="col-12">
                    <div class="btn-group" role="group">
                        <a href="{{ addQueryParam $.CurrentURL "month" (timestamp (addMonths .MonthStart -1)) }}" class="btn btn-outline-primary">
                            &larr; Previous Month
                        </a>
                        <span class="btn btn-primary disabled">
                            {{.MonthStart.Format "January 2006"}}
                        </span>
                        <a href="{{ addQueryParam $.CurrentURL "month" (timestamp (addMonths .MonthStart 1)) }}" class="btn btn-outline-primary">
                            Next Month &rarr;
                        </a>
                    </div>
                </div>
            </div>

            <!-- Budget Form -->
            <form method="POST" action="/web/budget/plan">
                <input type="hidden" name="month" value="{{timestamp .MonthStart}}">
                
                <div class="table-responsive">
                    <table class="table table-striped">
                        <thead>
                            <tr>
                                <th>Account</th>
                                <th>Budgeted Amount</th>
                            </tr>
                        </thead>
                        <tbody>
                            {{ range .ExpenseAccounts }}
                            {{ $acc := . }}
                            <tr>
                                <td>{{ $acc.Name }}</td>
                                <td>
                                    <div class="input-group">
                                        <input type="number"
                                               class="form-control"
                                               name="{{ $acc.Id }}"
                                               step="0.01"
                                               min="0"
                                               value="{{ range $.BudgetItems }}{{ if eq .AccountId $acc.Id }}{{ .Amount }}{{ end }}{{ end }}"
                                               placeholder="0.00">
                                    </div>
                                </td>
                            </tr>
                            {{ end }}
                        </tbody>
                    </table>
                </div>

                <div class="row">
                    <div class="col-12">
                        <button type="submit" class="btn btn-primary">Save Budget</button>
                        <a href="/web/budget/compare?month={{timestamp .MonthStart}}" class="btn btn-outline-secondary">
                            View Budget vs Actual
                        </a>
                    </div>
                </div>

            <!-- Copy Budget Form -->
            <div class="row mt-3">
                <div class="col-12">
                    <form method="POST" action="/web/budget/plan" class="d-inline">
                        <input type="hidden" name="month" value="{{timestamp .MonthStart}}">
                        <input type="hidden" name="action" value="copy">
                        <button type="submit" class="btn btn-outline-primary">
                            Copy from Previous Month
                        </button>
                    </form>
                </div>
            </div>
            </form>
        </div>
    </div>
</div>

{{template "footer.tpl" .}}
