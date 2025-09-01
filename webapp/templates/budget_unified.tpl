{{template "header.tpl" .}}

<div class="container-fluid">
    <div class="row">
        <div class="col-12">
            <h2>Budget</h2>

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
                        <a href="/web/budget?month={{ timestamp (addMonths .MonthStart -1) }}" class="btn btn-outline-primary">
                            &larr; Previous Month
                        </a>
                        <span class="btn btn-primary disabled">
                            {{.MonthStart.Format "January 2006"}}
                        </span>
                        <a href="/web/budget?month={{ timestamp (addMonths .MonthStart 1) }}" class="btn btn-outline-primary">
                            Next Month &rarr;
                        </a>
                    </div>
                </div>
            </div>

            <div class="row">
                <!-- Planning Column -->
                <div class="col-lg-6">
                    <h4>Plan</h4>
                    <form method="POST" action="/web/budget">
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
                                                       value="{{ index $.PlannedByAccount $acc.Id }}"
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
                                <button type="submit" class="btn btn-outline-primary" name="action" value="copy">
                                    Copy from Previous Month
                                </button>
                            </div>
                        </div>
                    </form>
                </div>

                <!-- Comparison Column -->
                <div class="col-lg-6">
                    <h4>Compare</h4>
                    <div class="table-responsive">
                        <table class="table table-striped">
                            <thead>
                                <tr>
                                    <th>Account</th>
                                    <th class="text-end">Planned</th>
                                    <th class="text-end">Actual</th>
                                    <th class="text-end">Variance</th>
                                    <th class="text-end">Status</th>
                                </tr>
                            </thead>
                            <tbody>
                                {{range .Comparison.Rows}}
                                <tr>
                                    <td>{{.AccountName}}</td>
                                    <td class="text-end">{{money .Planned}}</td>
                                    <td class="text-end">{{money .Actual}}</td>
                                    <td class="text-end {{if lt .Delta 0.0}}text-success{{else if gt .Delta 0.0}}text-danger{{else}}text-secondary{{end}}">
                                        {{money .Delta}}
                                    </td>
                                    <td class="text-end">
                                        {{if lt .Delta 0.0}}
                                            <span class="badge bg-success">Under Budget</span>
                                        {{else if gt .Delta 0.0}}
                                            <span class="badge bg-danger">Over Budget</span>
                                        {{else}}
                                            <span class="badge bg-secondary">On Budget</span>
                                        {{end}}
                                    </td>
                                </tr>
                                {{else}}
                                <tr>
                                    <td colspan="5" class="text-center text-muted">
                                        No budget or expense data for this month
                                    </td>
                                </tr>
                                {{end}}
                            </tbody>
                            {{if .Comparison.Rows}}
                            <tfoot>
                                <tr class="table-dark">
                                    <th>Total</th>
                                    <th class="text-end">{{money .Comparison.TotalPlanned}}</th>
                                    <th class="text-end">{{money .Comparison.TotalActual}}</th>
                                    <th class="text-end {{if lt .Comparison.TotalDelta 0.0}}text-success{{else if gt .Comparison.TotalDelta 0.0}}text-danger{{else}}text-secondary{{end}}">
                                        {{money .Comparison.TotalDelta}}
                                    </th>
                                    <th class="text-end">
                                        {{if lt .Comparison.TotalDelta 0.0}}
                                            <span class="badge bg-success">Under Budget</span>
                                        {{else if gt .Comparison.TotalDelta 0.0}}
                                            <span class="badge bg-danger">Over Budget</span>
                                        {{else}}
                                            <span class="badge bg-secondary">On Budget</span>
                                        {{end}}
                                    </th>
                                </tr>
                            </tfoot>
                            {{end}}
                        </table>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>

{{template "footer.tpl" .}}

