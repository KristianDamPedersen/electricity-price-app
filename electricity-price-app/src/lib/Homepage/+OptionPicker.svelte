<script>


    let options  ={
        start : "",
        end : "",
        area : ""
    };
    function handleSubmit(e){
        e.preventDefault()
        alert("submiiiit")
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
<h1>{options.start} and {options.end} and {options.area}</h1>