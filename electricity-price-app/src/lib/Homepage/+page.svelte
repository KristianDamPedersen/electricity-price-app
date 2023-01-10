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

    let lazyLoad = async () => {
        // Get current date and format it correctly
        const today = new Date();
        const yesterday = new Date(today);
        yesterday.setDate(yesterday.getDate() -1)

        console.log(today.toISOString().split('T')[0])
        console.log(yesterday.toISOString().split('T')[0])

        const startDate = yesterday.toISOString().split('T')[0]
        const endDate = today.toISOString().split('T')[0]

        // Execute API call with current date and default zone (DK1)
        let options = {
            start: startDate,
            end: endDate,
            area: "dk1"
        }
        let promise = await FetchWithOptions(options)
        records.populate(promise)
        return promise
    }
    </script>

<div class="navbar">
    <h1 class="heading">SmartPower (name)?</h1>
</div>
<div class="content-container">
    {#await lazyLoad()}
        <div class="graph">
            <h1>Loading</h1>
        </div>
    {:then}
        <div class="graph">
            <BarChart chartLabels = {hours} chartValues = {prices}/>
        </div>
    {:catch error}
            <p>oops something went wrong: {error.message}</p>
    {/await}
    <div class="button-container">
        <button class="button">24 hours</button>
        <button class="button">3 days</button>
        <button class="button">1 week</button>
        <button class="button">1 month</button>
        <button class="button">6 months</button>
        <button class="button">12 months</button>
    </div>
    <OptionPicker on:optionsubmit={handleSubmit}/>
</div>

<style>
    .content-container {
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        grid-gap: 2em;
        max-height: 80%;
        margin: 2em;
        z-index: 2;
        grid-template-columns: 20% 80%;
    }

    .button-container {
        align-self: center;
        margin-top: 2em;
        justify-content: center;
        display: grid;
        grid-auto-flow: column;
        grid-template-columns: repeat(5, 1fr);
        height: 100%;
        max-width: 100em;
        border-style: inset;
        border-color: var(--text-black);
        border-radius: 1em;
        background-color: var(--main-accent-yellow);
        box-shadow: rgba(0, 0, 0, 0.24) 0px 3px 8px;
    }

    .button {
        margin: 2em;
        height: 6em;
        min-width: 10em;
        max-width: 40em;
        cursor: pointer;
        background-color: var(--secondary-green);
        border-radius: 1em;
        border-color: var(--text-black);
        font: var(--main-font-family);
        font-size: medium;
        color: var(--text-black);
        box-shadow: rgba(0, 0, 0, 0.24) 0px 3px 8px;
    }

    .button:hover {
        transition-duration: 0.1s;
        scale: 105%;
        background-color: var(--secondary-green-highlighted);
    }

    .button:active {
        transition-duration: 0.1s;
        scale: 95%;
    }
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

    .graph {
        height: 50em;
    }
</style>