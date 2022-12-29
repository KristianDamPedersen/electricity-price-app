<script lang = "ts">
    import  {createRecords, FetchWithOptions } from "./stores";
    import { onMount } from 'svelte';
	import OptionPicker from "./+OptionPicker.svelte";
    import BarChart from "./+BarChart.svelte";
    let recordStore = createRecords();
    async function handleSubmit(options) {
        // Execute API call
        console.log(options)
        let json = await FetchWithOptions(options.detail)

        // Update the store
        recordStore.populate(json)
    }

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
        recordStore.populate(res)
    })


    </script>


<h2>This is the heading from Homepage</h2>



<BarChart/>
<OptionPicker on:optionsubmit={handleSubmit}/>
