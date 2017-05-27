var operations = {
    getParameterByName(name, url) {
        if (!url) url = window.location.href;
        name = name.replace(/[\[\]]/g, "\\$&");
        var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
            results = regex.exec(url);
        if (!results) return null;
        if (!results[2]) return '';
        return decodeURIComponent(results[2].replace(/\+/g, " "));
    },
    init: function() {
        // Init
        asticode.loader.init();
        asticode.notifier.init();

        // Refresh list operations
        operations.refreshListOperations();
    },
    refreshListOperations: function() {
        // Retrieve account id
        var accountID = operations.getParameterByName("account_id", window.location);

        // Send request to /api/accounts/:account_id/operations
        const req = new XMLHttpRequest();
        req.onreadystatechange = function() {
            if (this.readyState === XMLHttpRequest.DONE) {
                // Process errors
                if (req.status != 200) {
                    asticode.notifier.error(req.responseText);
                    return
                }

                // Loop through operations
                var operations = JSON.parse(req.responseText);
                if (typeof operations !== "undefined" && operations != null) {
                    var node = document.getElementById("operations");
                    var html = `<div class="operations-container"><table class="operations-table"><tbody>`;
                    for (var i = operations.length - 1; i >= 0; i--) {
                        var className = "amount-negative";
                        if (operations[i].amount > 0) {
                            className = "amount-positive";
                        }
                        html += `
                        <tr>
                            <td class="operations-cell" style="text-align: center; width: 100px">` + operations[i].date.split("T")[0] + `</td>
                            <td class="operations-cell" style="text-align: center; width: 200px">` + operations[i].subject + `</td>
                            <td class="operations-cell" style="text-align: center; width: 100px">` + operations[i].category + `</td>
                            <td class="operations-cell">` + operations[i].label + `</td>
                            <td class="operations-cell ` + className + `" style="text-align: right; width: 100px">` + operations[i].amount.toFixed(2) + `€</td>
                        </tr>
                    `;
                    }
                    html += "</tbody></table></div>";
                    node.innerHTML = html;
                }
            }
        };
        req.open('GET', '/api/accounts/' + accountID + "/operations", true);
        req.send(null);
    }
};