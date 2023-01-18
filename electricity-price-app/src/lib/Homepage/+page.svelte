<script lang = "ts">
    import  {createRecords, FetchWithOptions } from "./stores";
    import { onMount } from 'svelte';
    import Button from "./Button/+Button.svelte"
    import { writable } from 'svelte/store';
    import { createEventDispatcher } from 'svelte';
    import BarChart from "./+BarChart.svelte";
    import { records } from './stores'
    import '../../global.css'
    let hours = writable([]);
    let prices = writable([]);
    //since the API  uses a weird time format, we need to first convert it
    const convertDatetoString = (date) => {
        let monthstr =""
        
        let datestr = (date.getDate() < 10) ? "0" + date.getDate() : date.getDate()
        let monthString = ( (date.getMonth() +1) < 10) ? "0" + (date.getMonth() + 1) : date.getMonth() + 1
       return `${date.getFullYear()}-${monthString}-${datestr}T00:00`
    }
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
    //fetches the specified days behind today.
    //https://api.energidataservice.dk/dataset/Elspotprices?offset=0&start=2023-01-11T00:00&end=2023-01-14T00:00&sort=HourUTC%20DESC&timezone=dk
    async function fetchPreviousRecords(numberOfDays){
            let today = new Date(Date.now())

            //calculating days behind is from  karim79 on stackoverflow https://stackoverflow.com/questions/1296358/how-to-subtract-days-from-a-plain-date
            var dateOffset = (24*60*60*1000) * numberOfDays; 
            let previousDate = new Date(today.getTime() - dateOffset) 

            //add one with the same as above
            var offset = 24*60*60*1000;
            let tommorow = new Date(today.getTime() + offset)

            //now we can fetch the prices, here we need to choose in the future which area we want to see, but for now we will use DK1 for testing

            //the function "fetch with options" does not give enough control so therefor we do a manual fetch
            var res = await fetch("https://api.energidataservice.dk/dataset/Elspotprices?" + 
            `start=${convertDatetoString(previousDate)}` 
            + `&end=${convertDatetoString(tommorow)}&timezone=dk&PriceArea="DK1"`)
            var json = await  res.json()
            records.populate(json["records"])
    }

    // Subscribe to our record store, so that hours and prices are updated when new data arrives into it.
    records.subscribe(() => {
        hours = calcHours($records)
        prices = calcPrices($records)
    });
    // Lazy load the data for the last 24 hours
    let lazyLoad = async () => {
        // Get current date and format it correctly
        const today = new Date();
        const yesterday = new Date(today);
        yesterday.setDate(yesterday.getDate() -1)

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
    <div class="button-container-no-color">
        <Button text="24 hours" on:click = {() => fetchPreviousRecords(1)}/>
        <Button text="3 days" on:click = {() => fetchPreviousRecords(3)}/>
        <Button text="1 week" on:click = {() => fetchPreviousRecords(7)}/>
        <Button text="1 month" on:click = {() => fetchPreviousRecords(31)}/>
        <Button text="6 months" on:click={() => fetchPreviousRecords(180)}/>
        <Button text="12 months" on:click = {() =>fetchPreviousRecords(365)}/>
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