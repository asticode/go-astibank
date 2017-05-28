var index = {
    init: function() {
        // Init
        asticode.loader.init();
        asticode.notifier.init();
        asticode.modaler.init();
        asticode.modaler.onclose = function() {
            index.sendAccountsList();
        };

        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {
            // Listen
            index.listen();

            // Get references
            index.sendReferencesList();

            // Refresh list accounts
            index.sendAccountsList();

            // Handle import
            document.getElementById("btn-import").onclick = index.onClickImport;
        })
    },
    listen: function() {
        astilectron.listen(function(message) {
            asticode.loader.hide();
            switch (message.name) {
                case "accounts.list":
                    index.listenAccountsList(message);
                    break;
                case "error":
                    index.listenError(message);
                    break;
                case "import":
                    index.listenImport(message);
                    break;
                case "operations.add":
                    index.listenOperationsAdd(message);
                    break;
                case "references.list":
                    index.listenReferencesList(message);
                    break;
            }
        });
    },
    listenAccountsList: function(message) {
        var node = document.getElementById("accounts");
        node.innerHTML = "";
        for (var i = 0; i < message.payload.length; i++) {
            var className = "amount-negative";
            if (message.payload[i].balance > 0) {
                className = "amount-positive";
            }
            node.innerHTML = node.innerHTML + `
                <div class="account-container">
                    <div class="account-wrapper">
                       <div class="account-table">
                            <div class="account-cell">` + message.payload[i].id + `</div>
                            <div class="account-cell ` + className + `">` + message.payload[i].balance.toFixed(0) + `€</div>
                            <div class="account-cell">
                                <a href="operations.html?account_id=` + message.payload[i].id + `"><i class="fa fa-list"></i></a>
                            </div>
                        </div>
                        <div class="account-footer">
                            Last updated on ` + message.payload[i].updated_at.split("T")[0] + `
                        </div>
                    </div>
                </div>
            `;
        }
    },
    listenError: function(message) {
        asticode.notifier.error(message.payload);
    },
    listenImport: function(message) {
        // No new operations detected
        if (message.payload.length == 0) {
            asticode.notifier.info("No new operations detected");
            return
        }

        // Set operations
        index.import = {
            operations: message.payload,
        };

        // Set modal content
        index.setModalContent();
    },
    listenOperationsAdd: function() {
        asticode.notifier.success("Operation successfully added!");
        index.nextOperation();
    },
    listenReferencesList: function(message) {
        index.references = message.payload;
    },
    nextOperation: function() {
        // Remove first operation
        index.import.operations.shift();

        // No operations left
        if (index.import.operations.length == 0) {
            index.sendAccountsList();
            asticode.modaler.hide();
            return
        }

        // Build modal content
        index.setModalContent();
    },
    onClickAdd: function() {
        // Get values
        var category = document.getElementById("content-category").value;
        var label = document.getElementById("content-label").value;
        var subject = document.getElementById("content-subject").value;

        // Send
        index.import.operations[0].operation.category = category;
        index.import.operations[0].operation.label = label;
        index.import.operations[0].operation.subject = subject;
        index.sendOperationsAdd(index.import.operations[0].account, index.import.operations[0].operation);
    },
    onClickImport: function() {
        astilectron.showOpenDialog({properties: ['openFile', 'multiSelections']}, function(paths) {
            index.sendImport(paths);
        })
    },
    onClickSkip: function() {
        index.nextOperation();
    },
    sendAccountsList: function() {
        asticode.loader.show();
        astilectron.send({name: "accounts.list"});
    },
    sendImport: function(paths) {
        asticode.loader.show();
        astilectron.send({name: "import", payload: paths});
    },
    sendOperationsAdd: function(account, operation) {
        asticode.loader.show();
        astilectron.send({name: "operations.add", payload: {account: account, operation: operation}});
    },
    sendReferencesList: function() {
        asticode.loader.show();
        astilectron.send({name: "references.list"});
    },
    setModalContent: function() {
        // Build content
        var html = `
        <div style="margin-bottom: 15px">
            <h3>Metadata</h3>
            <table style="width: 100%"><tbody>
                <tr>
                    <td>Account:</td>
                    <td>` + index.import.operations[0].account.id + `</td>
                </tr>
                <tr>
                    <td>Date:</td>
                    <td>` + index.import.operations[0].operation.date.split("T")[0] + `</td>
                </tr>
                <tr>
                    <td>Raw label:</td>
                    <td>` + index.import.operations[0].operation.raw_label + `</td>
                </tr>
                <tr>
                    <td>Amount:</td>
                    <td>` + index.import.operations[0].operation.amount + `€</td>
                </tr>
            </tbody></table>
        </div>
        <div style="margin-bottom: 15px">
            <h3>Custom data</h3>
            <label>Subject:</label>
            <input type="text" id="content-subject" value="` + index.import.operations[0].operation.subject + `"/>
            <label>Category:</label>
            <select id="content-category">`;
        for (var i = 0; i < index.references.categories.length; i++) {
            var selected = "";
            if (index.references.categories[i] == index.import.operations[0].operation.category) {
                selected = " selected"
            }
            html += `<option value="` + index.references.categories[i] + `"` + selected + `>` + index.references.categories[i] + `</option>`;
        }
        html += `</select>
            <label>Label:</label>
            <input type="text" id="content-label" value="` + index.import.operations[0].operation.label + `"/>
        </div>
        <div style="text-align: center">
            <div style="display: inline-block">
                <button class="btn-success" id="btn-add" onclick="index.onClickAdd()">Add</button>
            </div>
            <div style="display: inline-block">
                <button class="btn-danger" onclick="index.onClickSkip()">Skip</button>
            </div>
        </div>
        `;
        var content = document.createElement("div");
        content.innerHTML = html;
        content.style.textAlign = "left";

        // Update modal
        asticode.modaler.setContent(content);
        document.getElementById('content-label').focus();
        document.getElementById('content-label').onkeypress = function(e) {
            if (e.keyCode == 13) {
                document.getElementById("btn-add").click();
            }
        };
        asticode.modaler.show();
    }
};