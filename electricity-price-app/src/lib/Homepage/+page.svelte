<script lang = "ts">
    import  {createRecords, FetchWithOptions } from "./stores";
    import { onMount } from 'svelte';
	import OptionPicker from "./+OptionPicker.svelte";
    import BarChart from "./+BarChart.svelte";
    import { records } from './stores'

    function calcHours(recordList) {
        let hours = []
        recordList.map(record => hours.push(record.HourDK))
        return hours
    }
    let hours = calcHours($records)
    function calcPrices(recordList) {
        let prices = []
        recordList.map(record => prices.push(record.SpotPriceDKK))
        return prices
    }
    let prices = calcHours($records)
    async function handleSubmit(options) {
        // Execute API call
        let json = await FetchWithOptions(options.detail)

        // Update the store
        records.populate(json)
    }
    records.subscribe(() => {
        hours = calcHours($records)
        prices = calcPrices($records)
    })

    // Define onMount function
    onMount(async () => {

        // Get current date and format it correctly
        const fullDate = new Date();
        let day = fullDate.getDate();
        let month = fullDate.getMonth() + 1;
        let year = fullDate.getFullYear();

        const startDate = `${year}-${month}-${day}T00:00`
        const endDate = `${year}-${month}-${day}T23:59`

        // Execute API call with current date and default zone (DK1)
        let options = {
            start: startDate,
            end: endDate,
            area: "DK1"
        }
        let res = await FetchWithOptions(options)

        console.log(hours)
        // Populate our recordStore
        records.populate(res)


    })

    </script>


<h2>This is the heading from Homepage</h2>



<BarChart chartLabels = {hours} chartValues = {prices} context = {records} />
<OptionPicker on:optionsubmit={handleSubmit}/>
