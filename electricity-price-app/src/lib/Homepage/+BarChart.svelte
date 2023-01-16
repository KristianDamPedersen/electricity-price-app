<script>
    import Chart from 'chart.js/auto';
    import { onMount } from 'svelte';
    import '../../global.css';
    let chart;
    let loaded = false
    let fillColor = '';
    // Converts hex values to rgb. Credit: https://stackoverflow.com/questions/21646738/convert-hex-to-rgba
    const hex2rgba = (hex, alpha = 1) => {
        const [r, g, b] = hex.match(/\w\w/g).map(x => parseInt(x, 16));
        return `rgba(${r},${g},${b},${alpha})`;
    };

    //Takes a YYYY-MM-DD-T-HH:MM:SS and retrieves only the hour
    const GetHourFromUTC = (timeToday) => {
        return timeToday.split("T")[1].split(":")[0]
    }
    let numOfDaysToDisplay = 1;
    let chartCanvas;
    export let context;
    export let chartLabels;
    export let chartValues;
    chartValues.subscribe(() => {
        if(loaded){
            chart.data.labels = []
            chart.data.datasets[0] = []
            for(let label of $chartLabels){
                    /*
                    chart.data.labels.push(GetHourFromUTC(label)) 
                    */
                   chart.data.labels.push(label)           
            }
            let rawData = []
            for(let dataEntry of $chartValues){
                rawData.push(dataEntry)
            }
            chart.data.datasets[0] = {
                label: 'Price EUR',
                data: rawData,
                fill: true,
                pointStyle: 'circle',
                pointRadius: 10,
                backgroundColor: hex2rgba(fillColor, 0.5),
                pointHoverRadius: 15,
                borderWidth: 1
            }
        //update
        chart.update()

        }
       
    })
    onMount(() => {

        const ctx = document.getElementById('myChart');
        fillColor = getComputedStyle(document.getElementById("myChart")).getPropertyValue('--fill-color')

        chart = new Chart(ctx, {
            type: 'line',
            data: {
            labels: $chartLabels,
            datasets: [{
                label: 'Price EUR',
                fill: true,
                lineTension: 0.4,
                data: $chartValues,
                pointStyle: 'circle',
                pointRadius: 10,
                pointHoverRadius: 15,
                backgroundColor: hex2rgba(fillColor, 0.5),
                borderWidth: 1,
            }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    y: {
                        beginAtZero: false
                    }
                },
                plugins: {
                }
            }
        });
        loaded = true
    })
</script>
<canvas class="graph" id="myChart" ></canvas>

<style>
    .graph {
        width: 100%;
        height: 100%;
        --fill-color: var(--primary-green);
    }
</style>