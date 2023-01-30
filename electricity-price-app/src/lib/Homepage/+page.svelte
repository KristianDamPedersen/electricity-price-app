<script lang = "ts">
    import  {createRecords, FetchWithOptions } from "./stores";
    // import { onMount } from 'svelte';
    import Button from "./Button/+Button.svelte"
    import { writable } from 'svelte/store';
    import BarChart from "./+BarChart.svelte";
    import { records } from './stores'
    import { retrieveData } from "./stores"
    import '../../global.css'
    let hours = writable([]);
    let prices = writable([]);

    const handleFetch = async (numberOfDays, averagedOverHours) => {
        let res = await retrieveData(numberOfDays, averagedOverHours)
        records.reset()
        records.populate(res)
    }

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

    // Subscribe to our record store, so that hours and prices are updated when new data arrives into it.
    records.subscribe(() => {
        hours = calcHours($records)
        prices = calcPrices($records)
    });

    // Lazy load the data for the last 24 hours
    let lazyLoad = async () => {
        // Get current date and format it correctly
        let promise = await retrieveData(1, 1)
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
    <div class="button-container-no-color">
        <Button text="24 hours" on:click = {() => handleFetch(1, 1)}/>
        <Button text="3 days" on:click = {() => handleFetch(3, 3)}/>
        <Button text="1 week" on:click = {() => handleFetch(7, 6)}/>
        <Button text="1 month" on:click = {() => handleFetch(31, 24)}/>
        <Button text="6 months" on:click={() => handleFetch(180, (24*3))}/>
        <Button text="12 months" on:click = {() => handleFetch(365, (24*7))}/>
    </div>
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
    .button-container-no-color{
        align-self: center;
        margin-top: 2em;
        justify-content: center;
        display: grid;
        grid-auto-flow: column;
        grid-template-columns: repeat(5, 1fr);
        height: 100%;
        max-width: 100em;
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