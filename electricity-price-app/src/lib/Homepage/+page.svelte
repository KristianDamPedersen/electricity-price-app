<script lang = "ts">
    import  {createRecords, FetchWithOptions } from "./stores";
    import { onMount } from 'svelte';
	import OptionPicker from "./+OptionPicker.svelte";
    import { writable } from 'svelte/store';
    import { createEventDispatcher } from 'svelte';
    import BarChart from "./+BarChart.svelte";
    import { records } from './stores'
    import '../../global.css'
    let hours = writable([]);
    let prices = writable([]);
    // Creates an array of dates + times in the current record store, timezone UTC+1 (Copenhagen)
    function calcHours(recordList) {
        hours.set([])
        recordList.map(record => hours.update( prevHours => [...prevHours, record.HourDK]));
        return hours
    }
    // Creates an array of all the prices from the record store.
    function calcPrices(recordList) {
        prices.set([]);
        recordList.map(record => prices.update( prices => [...prices, record.SpotPriceDKK]))
        return prices
    }
    async function handleSubmit(options) {
        // Execute API call
        let json = await FetchWithOptions(options.detail)

        // Update the store
        records.populate(json)
    }

    // Subscribe to our record store, so that hours and prices are updated when new data arrives into it.
    records.subscribe(() => {
        hours = calcHours($records)
        prices = calcPrices($records)
    });

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

        // Populate our recordStore
        records.populate(res)
    })
    </script>

<div class="navbar">
    <h1 class="heading">SmartPower (name)?</h1>
</div>
<BarChart chartLabels = {hours} chartValues = {prices}/>
<OptionPicker on:optionsubmit={handleSubmit}/>

<style>
    .navbar {
        width: auto;
        display: flex;
        padding-left: 3em;
        justify-content: start;
        align-items: center;
        height: 3em;
        box-shadow: rgba(0, 0, 0, 0.24) 0px 3px 8px;;
        background-color: var(--main-bg-color);
    }
    .heading {
        margin: 0;
        color: var(--text-black);
        font: var(--main-font-family);
    }
</style>