<script>
	import { createEventDispatcher } from "svelte";


    export let text = "";
    const dispatch = createEventDispatcher()

    function handleClick(event){
        dispatch('click', event.detail)
    }
</script>

<button class="button" on:click={handleClick}>{text}</button>

<style>
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
</style>