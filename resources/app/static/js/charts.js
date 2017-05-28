var charts = {
    countCharts: 0,
    init: function() {
        // Init
        asticode.loader.init();
        asticode.notifier.init();

        // Get account id
        charts.account_id = asticode.tools.getParameterByName("account_id", window.location);

        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {
            // Listen
            charts.listen();

            // Refresh charts
            charts.sendChartsAll();
        });
    },
    listen: function() {
        astilectron.listen(function(message) {
            asticode.loader.hide();
            switch (message.name) {
                case "error":
                    charts.listenError(message);
                    break;
                case "charts.all":
                    charts.listenChartsAll(message);
                    break;
            }
        });
    },
    listenError: function(message) {
        asticode.notifier.error(message.payload);
    },
    sendChartsAll: function() {
        asticode.loader.show();
        astilectron.send({name: "charts.all", payload: charts.account_id});
    },
    listenChartsAll: function(message) {
        var node = document.getElementById("charts");
        node.innerHTML = "";
        for (var i = 0; i < message.payload.length; i++) {
            charts.countCharts++;
            var canvas = document.createElement("canvas");
            canvas.id = "chart-" + charts.countCharts;
            document.getElementById("charts").append(canvas);
            new Chart(canvas, message.payload[i]);
        }
    }
};