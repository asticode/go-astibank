var index = {
    init: function() {
        // Init
        asticode.loader.init();
        asticode.notifier.init();
        asticode.modaler.init();

        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {
            // Handle import
            document.getElementById("btn-import").onclick = index.onClickImport;

            // Get references
            index.getReferences();

            // Refresh list accounts
            index.refreshListAccounts();
        })
    },
    getReferences: function() {
        // Send request to /api/references
        const req = new XMLHttpRequest();
        req.onreadystatechange = function(event) {
            if (this.readyState === XMLHttpRequest.DONE) {
                // Process errors
                if (req.status != 200) {
                    asticode.notifier.error(req.responseText);
                    return
                }

                // Set references
                index.references = JSON.parse(req.responseText);
            }
        };
        req.open('GET', '/api/references', true);
        req.send(null);
    },
    onClickAdd: function(originalLength, operations) {
        return function() {
            // Get values
            var label = document.getElementById("content-label").value;
            var category = document.getElementById("content-category").value;

            // Show loader
            asticode.loader.show();

            // Send request to /api/operations
            const req = new XMLHttpRequest();
            req.onreadystatechange = function(event) {
                if (this.readyState === XMLHttpRequest.DONE) {
                    // Hide loader
                    asticode.loader.hide();

                    // Process errors
                    if (req.status != 204) {
                        asticode.notifier.error(req.responseText);
                        return
                    }

                    // Remove first operation
                    operations.shift();

                    // No operations left
                    if (operations.length == 0) {
                        index.refreshListAccounts();
                        asticode.modaler.hide();
                        asticode.notifier.success(originalLength + " operation(s) imported");
                        return
                    }

                    // Build modal content
                    index.setModalContent(originalLength, operations);
                }
            };
            req.open('POST', '/api/accounts/' + operations[0].account.id + '/operations', true);
            operations[0].operation.label = label;
            operations[0].operation.category = category;
            req.send(JSON.stringify(operations[0].operation));
        };
    },
    onClickImport: function() {
        // Show open dialog
        astilectron.showOpenDialog({properties: ['openFile', 'multiSelections']}, function(paths) {
            // Show loader
            asticode.loader.show();

            // Send request to /api/import
            const req = new XMLHttpRequest();
            req.onreadystatechange = function() {
                if (this.readyState === XMLHttpRequest.DONE) {
                    // Hide loader
                    asticode.loader.hide();

                    // Process errors
                    if (req.status != 200) {
                        asticode.notifier.error(req.responseText);
                        return
                    }

                    // No new operations detected
                    var operations = JSON.parse(req.responseText);
                    if (operations.length == 0) {
                        asticode.notifier.info("No new operations detected");
                        return
                    }

                    // Set modal content
                    index.setModalContent(operations.length, operations);
                }
            };
            req.open('POST', '/api/import', true);
            req.send(JSON.stringify({paths: paths}));
        })
    },
    refreshListAccounts: function() {
        // Send request to /api/accounts
        const req = new XMLHttpRequest();
        req.onreadystatechange = function() {
            if (this.readyState === XMLHttpRequest.DONE) {
                // Process errors
                if (req.status != 200) {
                    asticode.notifier.error(req.responseText);
                    return
                }

                // Loop through accounts
                var accounts = JSON.parse(req.responseText);
                if (typeof accounts !== "undefined" && accounts != null) {
                    var node = document.getElementById("accounts");
                    node.innerHTML = "";
                    for (var i = 0; i < accounts.length; i++) {
                        var className = "amount-negative";
                        if (accounts[i].balance > 0) {
                            className = "amount-positive";
                        }
                        node.innerHTML = node.innerHTML + `
                        <div class="account-container">
                            <div class="account-wrapper">
                               <div class="account-table">
                                    <div class="account-cell">` + accounts[i].id + `</div>
                                    <div class="account-cell ` + className + `">` + accounts[i].balance.toFixed(0) + `€</div>
                                    <div class="account-cell">
                                        <a href="/templates/operations?account_id=` + accounts[i].id + `"><i class="fa fa-list"></i></a>
                                    </div>
                                </div>
                                <div class="account-footer">
                                    Last updated on ` + accounts[i].updated_at.split("T")[0] + `
                                </div>
                            </div>
                        </div>
                    `;
                    }
                }
            }
        };
        req.open('GET', '/api/accounts', true);
        req.send(null);
    },
    setModalContent: function(originalLength, operations) {
        // Build button
        var btn = document.createElement("button");
        btn.innerText = "Add";
        btn.className = "btn-lg btn-success";
        btn.onclick = index.onClickAdd(originalLength, operations);

        // Build content
        var html = `
        <div style="margin-bottom: 15px">
            <h3>Metadata</h3>
            <table style="width: 100%"><tbody>
                <tr>
                    <td>Account:</td>
                    <td>` + operations[0].account.id + `</td>
                </tr>
                <tr>
                    <td>Date:</td>
                    <td>` + operations[0].operation.date.split("T")[0] + `</td>
                </tr>
                <tr>
                    <td>Raw label:</td>
                    <td>` + operations[0].operation.raw_label + `</td>
                </tr>
                <tr>
                    <td>Amount:</td>
                    <td>` + operations[0].operation.amount + `€</td>
                </tr>
            </tbody></table>
        </div>
        <div style="margin-bottom: 15px">
            <h3>Custom data</h3>
            <label>Subject:</label>
            <select id="content-subject">`;
        for (var i = 0; i < index.references.subjects.length; i++) {
            var selected = "";
            if (index.references.subjects[i] == operations[0].operation.subject) {
                selected = " selected"
            }
            html += `<option value="` + index.references.subjects[i] + `"` + selected + `>` + index.references.subjects[i] + `</option>`;
        }
        html += `</select>
        <label>Category:</label>
        <select id="content-category">`;
        for (var i = 0; i < index.references.categories.length; i++) {
            var selected = "";
            if (index.references.categories[i] == operations[0].operation.category) {
                selected = " selected"
            }
            html += `<option value="` + index.references.categories[i] + `"` + selected + `>` + index.references.categories[i] + `</option>`;
        }
        var label = "";
        if (operations[0].operation.label != "") {
            label = operations[0].operation.label;
        }
        html += `</select>
            <label>Label:</label>
            <input type="text" id="content-label" autofocus value="` + label + `"/>
        </div>
        `;
        var content = document.createElement("div");
        content.innerHTML = html;
        content.style.textAlign = "left";
        content.appendChild(btn);

        // Update modal
        asticode.modaler.setContent(content);
        document.getElementById('content-label').focus();
        document.getElementById('content-label').onkeypress = function(e) {
            if (e.keyCode == 13) {
                btn.click();
            }
        };
        asticode.modaler.show();
    }
};