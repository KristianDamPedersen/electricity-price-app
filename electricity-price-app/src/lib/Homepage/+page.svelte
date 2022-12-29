<script lang = "ts">
    import  {createRecords, FetchWithOptions } from "./stores";
	import OptionPicker from "./+OptionPicker.svelte";

    let recordStore = createRecords();
    async function handleSubmit(options) {
        // Execute API call
        console.log(options)
        let json = await FetchWithOptions(options.detail)

        // Update the store
        recordStore.populate(json)
    }
    </script>

<h2>This is the heading from Homepage</h2>

{#each $recordStore as rec}
    <h2>#### NEW RECORD ####</h2>
    <p>HourUTC: {rec.HourUTC}</p>
    <p>Hour DK: {rec.HourDK}</p>
    <p>Price area: {rec.PriceArea}</p>
    <p>Spot price DKK: {rec.SpotPriceDKK}</p>
    <p>Spot price EUR: {rec.SpotPriceEUR}</p>
{/each}

<OptionPicker on:optionsubmit={handleSubmit}/>