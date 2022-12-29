<script>
    import Chart from 'chart.js/auto';
    import { onMount } from 'svelte';

    let chartValues = [20, 10, 5, 2, 20, 30, 45];
    let chartLabels = ['January', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul'];
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