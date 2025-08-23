{{ template "header.tpl" . }}

<main>
    {{ with .Expenses }}
        <h2>Expenses</h2>
        <a href="{{ addQueryParam $.CurrentURL "from" (timestamp (addMonths $.To -2)) }}" class="btn btn-primary" tabindex="-1" role="button">
            <i class="bi-arrow-left-circle-fill"></i>
        </a>
        {{ formatTime $.From "2006-01-02" }} - {{ formatTime $.To "2006-01-02" }}
        <a href="{{ addQueryParam $.CurrentURL "from" $.Next }}" class="btn btn-primary" tabindex="-1" role="button">
            <i class="bi-arrow-right-circle-fill"></i>
        </a>

        {{ range .Currencies }}
            <table class="table table-sm table-hover">
                <thead>
                    <tr>
                        <th>
                            <a href="{{ addQueryParam $.CurrentURL "currency" .CurrencyName }}"
                                    class="btn btn-light {{if eq (index $.Query "currency") .CurrencyName}}active{{end}}" tabindex="-1" role="button">
                                <i class="bi-arrow-through-heart-fill"></i>
                            </a>
                            Currency: {{ .CurrencyName }}
                        </th>
                        {{ range .Intervals }}
                        <th>{{ formatTime . "2006-01-02" }}</th>
                        {{ end }}
                        <th>Total</th>
                    </tr>
                </thead>
                <tbody class="table-group-divider">
                    {{ $intervals := .Intervals }}
                    {{ range .Accounts }}
                    <tr>
                        {{ $accountID := .AccountID }}
                        <td>{{ .AccountName }}</td>
                        {{ range $i, $a := .Amounts }}
                        <td {{ if and (eq $accountID "") (ne $a 0.0) }}class="table-danger"{{end}}>
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
                            <strong>Total</strong>
                        </td>
                        {{ range .Total }}
                            <td>{{ money . }}</td>
                        {{ end }}
                </tfoot>
            </table>
        {{ end }}
    {{ end }}

</main>

<script>
document.addEventListener('DOMContentLoaded', function() {
    // Function to extract numeric value from text (handles money formatting)
    function extractNumericValue(text) {
        if (!text || text.trim() === '' || text.trim() === '0') {
            return 0;
        }
        // Remove currency symbols, commas, and extract number
        const cleaned = text.replace(/[^\d.-]/g, '');
        const value = parseFloat(cleaned);
        return isNaN(value) ? 0 : value; // Keep original sign
    }

    // Function to calculate 90th percentile
    function calculatePercentile(values, percentile) {
        if (values.length === 0) return 0;
        const sorted = values.slice().sort((a, b) => a - b);
        const index = (percentile / 100) * (sorted.length - 1);
        if (index === Math.floor(index)) {
            return sorted[index];
        } else {
            const lower = sorted[Math.floor(index)];
            const upper = sorted[Math.ceil(index)];
            return lower + (upper - lower) * (index - Math.floor(index));
        }
    }

    // Function to convert value to color using 90th percentile
    function valueToColor(value, percentile90) {
        if (percentile90 === 0 || value <= 0) return 'rgba(255, 255, 255, 0)'; // transparent for no data or negative values
        
        if (value > percentile90) {
            // For values above 90th percentile, use only red color with 30% opacity
            return 'rgba(255, 0, 0, 0.3)';
        }
        
        // For values below 90th percentile, use gradient from green to red
        const ratio = value / percentile90;
        const intensity = ratio;
        
        // Interpolate between green (small values) and red (large values)
        const red = Math.round(255 * intensity);
        const green = Math.round(155 * (1 - intensity));
        const blue = 50; // Keep some blue for better visibility
        
        return `rgba(${red}, ${green}, ${blue}, 0.1)`; // Reduced opacity to 10%
    }

    // Find all tables with expense data
    const tables = document.querySelectorAll('table.table');
    
    tables.forEach(function(table) {
        const cells = [];
        const values = [];
        
        // Collect all numeric cells and their values
        const rows = table.querySelectorAll('tbody tr');
        rows.forEach(function(row) {
            const dataCells = row.querySelectorAll('td');
            // Skip the first cell (account name) and last cell (total)
            for (let i = 1; i < dataCells.length - 1; i++) {
                const cell = dataCells[i];
                const value = extractNumericValue(cell.textContent);
                if (value > 0) {
                    cells.push({ element: cell, value: value });
                    values.push(value);
                }
            }
        });
        
        // Calculate 90th percentile
        const percentile90 = calculatePercentile(values, 90);
        
        // Apply color coding to cells
        cells.forEach(function(cellData) {
            const color = valueToColor(cellData.value, percentile90);
            cellData.element.style.backgroundColor = color;
        });
    });
});
</script>

{{ template "footer.tpl" . }}