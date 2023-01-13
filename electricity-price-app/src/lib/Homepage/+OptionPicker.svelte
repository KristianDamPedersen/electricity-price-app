<script>
    import {createEventDispatcher} from 'svelte'

    const dispatch = createEventDispatcher()
    let options  ={
        start : "",
        end : "",
        area : ""
    };
    function handleSubmit(e){
        e.preventDefault()
        console.log(options)
        dispatch('optionsubmit',options)
    }
</script>

<form>
    <label> start</label>
    <input bind:value={options.start} type ="date">
    <label> end </label>
    <input bind:value={options.end} type = "date">
    <label>Area</label>
    <select bind:value = {options.area}>
        <option value = {"dk1"}>DK1</option>
        <option value = {"dk2"}>DK2</option>
    </select>
    <button on:click={handleSubmit}>Submit</button>
</form>