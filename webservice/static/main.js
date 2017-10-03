function drawChart(options) {
    Highcharts.setOptions({
        global: {
            useUTC: false
        }
    });
    var chart = new Highcharts.Chart(options);
}

function fetchData(sensor, options, callback) {
    $.getJSON({
        url: '/sensor/' + sensor + '/lastvalue?count=500',
        beforeSend: function(xhr) {
            xhr.setRequestHeader("Authorization", "Basic " + btoa("foo" + ":" + "bar"));
        },
        success: function(data) {
            var mydata = [];
            for (i = 0; i < data.length; i += 1) {
                mydata.push([Date.parse(data[i].Date), data[i].Value]);
            }
            options.series.push({ data: mydata, name: "test" });
            callback(options);
        }
    });
}

function createChart(options, sensors, chart) {
    if (sensors.length > 0) {
        sensor = sensors.pop();
        fetchData(sensor, options,
            function(newOptions) {
                createChart(newOptions, sensors, chart)
            });
    } else {
        drawChart(options)
    }
}

$(document).ready(function() {
    var options = {
        chart: {
            renderTo: 'container',
            type: 'spline',
            zoomType: 'x',
            panning: true,
            panKey: 'shift'
        },
        series: [],
        title: {
            text: "Heizung"
        },
        xAxis: {
            type: 'datetime'
        },
        yAxis: {
            title: {
                text: "Temperatur"
            }
        }
    };

    createChart(options, ["10-000800355f27", "10-0008019453e2", "10-0008019462b6", "10-00080194662f", "10-0008019481df"]);

});