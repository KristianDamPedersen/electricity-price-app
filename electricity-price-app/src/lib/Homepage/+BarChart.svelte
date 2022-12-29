<script>
    import Chart from 'chart.js/auto';
    import { onMount } from 'svelte';

    export let chartValues ;
    export let chartLabels ;
    let chartCanvas;

    onMount(() => {
        const ctx = document.getElementById('myChart');

        new Chart(ctx, {
            type: 'bar',
            data: {
            labels: chartLabels,
            datasets: [{
                label: 'Stonks',
                data: chartValues,
                borderWidth: 1
            }]
            },
            options: {
                scales: {
                    y: {
                        beginAtZero: true
                    }
            }
            }
        });
    })
</script>
<div>


<canvas id="myChart" ></canvas>

</div>