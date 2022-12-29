<script>
    import Chart from 'chart.js/auto';
    import { onMount } from 'svelte';
    let chart;
    let loaded = false
    let chartCanvas;
    export let context;
    export let chartLabels;
    export let chartValues;
    chartValues.subscribe(() => {
        if(loaded){
            chart.data.labels = []
            chart.data.datasets[0] = []
            for(let label of $chartLabels){
                chart.data.labels.push(label)
            }
            let rawData = []
            for(let dataEntry of $chartValues){
                rawData.push(dataEntry)
            }
            chart.data.datasets[0] = {
                label: 'Stonks',
                data: rawData,
                borderWidth: 1
            }
        //update
        chart.update()

        }
       
    })
    onMount(() => {

        const ctx = document.getElementById('myChart');

        chart = new Chart(ctx, {
            type: 'bar',
            data: {
            labels: $chartLabels,
            datasets: [{
                label: 'Stonks',
                data: $chartValues,
                borderWidth: 1
            }]
            },
            options: {
                scales: {
                    y: {
                        beginAtZero: false
                    }
            }
            }
        });
        loaded = true
    })
</script>
<div>


<canvas id="myChart" ></canvas>

</div>