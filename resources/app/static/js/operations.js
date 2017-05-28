var operations = {
    init: function() {
        // Init
        asticode.loader.init();
        asticode.notifier.init();
        asticode.modaler.init();

        // Get account id
        operations.account_id = asticode.tools.getParameterByName("account_id", window.location);

        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {
            // Listen
            operations.listen();

            // Get references
            operations.sendReferencesList();

            // Refresh list operations
            operations.sendOperationsList();
        });
    },
    listen: function() {
        astilectron.listen(function(message) {
            asticode.loader.hide();
            switch (message.name) {
                case "error":
                    operations.listenError(message);
                    break;
                case "operations.list":
                    operations.listenOperationsList(message);
                    break;
                case "operations.one":
                    operations.listenOperationsOne(message);
                    break;
                case "operations.update":
                    operations.listenOperationsUpdate(message);
                    break;
                case "references.list":
                    operations.listenReferencesList(message);
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
                <tr style="cursor: pointer" onclick="operations.sendOperationsOne(` + message.payload[i].id + `)">
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
    listenOperationsOne: function(message) {
        // Build button
        var btn = document.createElement("button");
        btn.innerText = "Update";
        btn.className = "btn-lg btn-success";
        btn.onclick = operations.onClickUpdate(message.payload);

        // Build content
        var html = `
        <label>Subject:</label>
        <input type="text" id="content-subject" value="` + message.payload.subject + `"/>
        <label>Category:</label>
        <select id="content-category">`;
        for (var i = 0; i < operations.references.categories.length; i++) {
            var selected = "";
            if (operations.references.categories[i] == message.payload.category) {
                selected = " selected"
            }
            html += `<option value="` + operations.references.categories[i] + `"` + selected + `>` + operations.references.categories[i] + `</option>`;
        }
        html += `</select>
            <label>Label:</label>
            <input type="text" id="content-label" value="` + message.payload.label + `"/>
        </div>
        `;
        var content = document.createElement("div");
        content.innerHTML = html;
        content.style.textAlign = "left";
        content.appendChild(btn);

        // Update modal
        asticode.modaler.setContent(content);
        document.getElementById('content-label').onkeypress = function(e) {
            if (e.keyCode == 13) {
                btn.click();
            }
        };
        asticode.modaler.show();
    },
    listenOperationsUpdate: function() {
        asticode.modaler.hide();
        operations.sendOperationsList();
    },
    listenReferencesList: function(message) {
        operations.references = message.payload;
    },
    onClickUpdate: function(operation) {
        return function() {
            operation.category = document.getElementById("content-category").value;
            operation.label = document.getElementById("content-label").value;
            operation.subject = document.getElementById("content-subject").value;
            operations.sendOperationsUpdate(operation);
        };
    },
    sendOperationsList: function() {
        asticode.loader.show();
        astilectron.send({name: "operations.list", payload: operations.account_id});
    },
    sendOperationsOne: function(id) {
        asticode.loader.show();
        astilectron.send({name: "operations.one", payload: {account: {id: operations.account_id}, operation: {id: id}}});
    },
    sendOperationsUpdate: function(operation) {
        asticode.loader.show();
        astilectron.send({name: "operations.update", payload: {account: {id: operations.account_id}, operation: operation}});
    },
    sendReferencesList: function() {
        asticode.loader.show();
        astilectron.send({name: "references.list"});
    }
};