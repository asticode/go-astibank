var operations = {
    init: function() {
        // Init
        asticode.loader.init();
        asticode.notifier.init();

        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {
            // Listen
            operations.listen();

            // Refresh list operations
            operations.sendOperationsList();
        });
    },
    listen: function() {
        astilectron.listen(function(message) {
            switch (message.name) {
                case "error":
                    operations.listenError(message);
                    break;
                case "operations.list":
                    operations.listenOperationsList(message);
                    break;
            }
        });
    },
    listenError: function(message) {
        asticode.notifier.error(message.payload);
    },
    listenOperationsList: function(message) {
        var node = document.getElementById("operations");
        var html = `<div class="operations-container"><table class="operations-table"><tbody>`;
        for (var i = message.payload.length - 1; i >= 0; i--) {
            var className = "amount-negative";
            if (message.payload[i].amount > 0) {
                className = "amount-positive";
            }
            html += `
                        <tr>
                            <td class="operations-cell" style="text-align: center; width: 100px">` + message.payload[i].date.split("T")[0] + `</td>
                            <td class="operations-cell" style="text-align: center; width: 200px">` + message.payload[i].subject + `</td>
                            <td class="operations-cell" style="text-align: center; width: 100px">` + message.payload[i].category + `</td>
                            <td class="operations-cell">` + message.payload[i].label + `</td>
                            <td class="operations-cell ` + className + `" style="text-align: right; width: 100px">` + message.payload[i].amount.toFixed(2) + `€</td>
                        </tr>
                    `;
        }
        html += "</tbody></table></div>";
        node.innerHTML = html;
    },
    sendOperationsList: function() {
        astilectron.send({name: "operations.list", payload: asticode.tools.getParameterByName("account_id", window.location)});
    }
};