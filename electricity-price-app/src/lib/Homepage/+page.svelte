<script lang = "ts">
    import  { records, createRecords } from "./stores";
    import Optionpicker from './+OptionPicker.svelte'
	import OptionPicker from "./+OptionPicker.svelte";
    // import  { Record } from "./stores"


    let recordStore = createRecords();
    let testRecord = {
        HourUTC: new Date(Date.parse("2022-12-18T00:00")),
        HourDK: new Date(Date.parse("2022-12-18T00:00")),
        PriceArea: "DK1",
        SpotPriceDKK: 150,
        SpotPriceEUR: 20,
    };

     recordStore.add(testRecord);
     recordStore.add(testRecord);
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

<OptionPicker />