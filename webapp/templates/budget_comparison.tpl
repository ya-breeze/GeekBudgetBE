{{template "header.tpl" .}}

<div class="container-fluid">
    <div class="row">
        <div class="col-12">
            <h2>Budget vs Actual Comparison</h2>
            
            <!-- Month Navigation -->
            <div class="row mb-3">
                <div class="col-12">
                    <div class="btn-group" role="group">
                        <a href="{{addQueryParam .Request.URL "month" (timestamp .PrevMonth)}}" class="btn btn-outline-primary">
                            &larr; Previous Month
                        </a>
                        <span class="btn btn-primary disabled">
                            {{.MonthStart.Format "January 2006"}}
                        </span>
                        <a href="{{addQueryParam .Request.URL "month" (timestamp .NextMonth)}}" class="btn btn-outline-primary">
                            Next Month &rarr;
                        </a>
                    </div>
                </div>
            </div>

            <!-- Summary Cards -->
            <div class="row mb-4">
                <div class="col-md-4">
                    <div class="card">
                        <div class="card-body">
                            <h5 class="card-title">Total Planned</h5>
                            <h3 class="text-primary">{{money .Comparison.TotalPlanned}}</h3>
                        </div>
                    </div>
                </div>
                <div class="col-md-4">
                    <div class="card">
                        <div class="card-body">
                            <h5 class="card-title">Total Actual</h5>
                            <h3 class="text-info">{{money .Comparison.TotalActual}}</h3>
                        </div>
                    </div>
                </div>
                <div class="col-md-4">
                    <div class="card">
                        <div class="card-body">
                            <h5 class="card-title">Difference</h5>
                            <h3 class="{{if lt .Comparison.TotalDelta 0}}text-success{{else}}text-danger{{end}}">
                                {{money .Comparison.TotalDelta}}
                            </h3>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Comparison Table -->
            <div class="table-responsive">
                <table class="table table-striped">
                    <thead>
                        <tr>
                            <th>Account</th>
                            <th class="text-end">Planned</th>
                            <th class="text-end">Actual</th>
                            <th class="text-end">Difference</th>
                            <th class="text-end">Status</th>
                        </tr>
                    </thead>
                    <tbody>
                        {{range .Comparison.Rows}}
                        <tr>
                            <td>{{.AccountName}}</td>
                            <td class="text-end">{{money .Planned}}</td>
                            <td class="text-end">{{money .Actual}}</td>
                            <td class="text-end {{if lt .Delta 0}}text-success{{else}}text-danger{{end}}">
                                {{money .Delta}}
                            </td>
                            <td class="text-end">
                                {{if lt .Delta 0}}
                                    <span class="badge bg-success">Under Budget</span>
                                {{else if gt .Delta 0}}
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
                            <th class="text-end {{if lt .Comparison.TotalDelta 0}}text-success{{else}}text-danger{{end}}">
                                {{money .Comparison.TotalDelta}}
                            </th>
                            <th class="text-end">
                                {{if lt .Comparison.TotalDelta 0}}
                                    <span class="badge bg-success">Under Budget</span>
                                {{else if gt .Comparison.TotalDelta 0}}
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

            <div class="row mt-3">
                <div class="col-12">
                    <a href="/web/budget/plan?month={{timestamp .MonthStart}}" class="btn btn-primary">
                        Edit Budget for This Month
                    </a>
                    <a href="/web/budget/plan" class="btn btn-outline-secondary">
                        Plan New Budget
                    </a>
                </div>
            </div>
        </div>
    </div>
</div>

{{template "footer.tpl" .}}
