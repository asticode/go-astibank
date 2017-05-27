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
        })
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
                    var json = JSON.parse(req.responseText);
                    if (json.operations.length == 0) {
                        asticode.notifier.info("No new operations detected");
                        return
                    }

                    // Set modal content
                    index.setModalContent(json.operations.length, json.operations);
                }
            };
            req.open('POST', '/api/import', true);
            req.send(JSON.stringify({paths: paths}));
        })
    },
    setModalContent: function(originalLength, operations) {
        // Build button
        var btn = document.createElement("button");
        btn.innerText = "Add";
        btn.className = "btn-lg btn-success";
        btn.onclick = index.onClickAdd(originalLength, operations);

        // Build content
        var content = document.createElement("div");
        content.innerHTML = `
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
            <label>Label:</label>
            <input type="text" id="content-label"/>
            <label>Category:</label>
            <select id="content-category" style="margin-bottom: 10px">
                <option value="Food">Food</option>
                <option value="Pleasure">Pleasure</option>
                <option value="Rent">Rent</option>
            </select>
        </div>
        `;
        content.style.textAlign = "left";
        content.appendChild(btn);

        // Update modal
        asticode.modaler.setContent(content);
        asticode.modaler.show();
    }
};