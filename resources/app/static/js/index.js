astibank.loadIndex = function() {
    document.body.innerHTML = `
<div class="table">
    <div class="row">
        <div class="cell">
            <table>
                <tr>
                    <td>
                        <div class="index-panel success" onclick="astibank.import()">
                            <i class="fa fa-file-import"></i>
                            Import
                        </div>
                    </td>
                    <td>
                        <div class="index-panel success">
                            <i class="fa fa-tags"></i>
                            Labels
                        </div>
                    </td>
                </tr>
                <tr>
                    <td>
                        <div class="index-panel success">
                            <i class="fa fa-chart-pie"></i>
                            Graphs
                        </div>
                    </td>
                    <td>
                        <div class="index-panel success" onclick="astibank.quit()">
                            <i class="fa fa-sign-out-alt"></i>
                            Quit
                        </div>
                    </td>
                </tr>
            </table>
        </div>
    </div>
</div>
    `;
};
astibank.import = function() {
    astilectron.showOpenDialog({properties: ['openFile', 'multiSelections']}, function(paths) {
        astilectron.sendMessage({"name": "js.import", "payload": paths}, function() {
            astibank.loadTags();
        });
    });
};
astibank.loadTags = function() {
    document.body.innerHTML = "<a onclick='astibank.loadIndex()'>Back</a>"
};
astibank.quit = function() {
    astilectron.sendMessage({"name": "js.quit"});
};