<script>
    import Chart from 'chart.js/auto';
    let chartData;
    import { onMount } from 'svelte';

    let chartValues = [20, 10, 5, 2, 20, 30, 45];
    let chartLabels = ['January', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul'];
    let ctx;
    let chartCanvas;

    onMount(() => {
        ctx = chartCanvas.getContext('2d');
        var chart = new Chart(ctx, {
             type: "bar",
            data: {
                labels: chartLabels,
                datasets: [{
                    label: 'Revenue',
                    backgroundColor: 'rgb(255,99,132)',
                    borderColor: 'rgb(255,99,132)',
                    data: chartValues
                }]
            }
        });
    });
</script>
<canvas bind:this={chartCanvas} id="myChart"></canvas>