import { json } from '@sveltejs/kit'
import { writable, derived } from 'svelte/store'
import recordStore from './+page.svelte'


export interface Record{
    HourUTC: Date,
    HourDK: Date,
    PriceArea: string,
    SpotpriceDKK: number,
    SpotPriceEUR: number,
}
export interface option{
    start : Date,
    end: Date,
    PriceArea: string
}
//Makes an APICall to the danish public database, and populates the a store
export async function FetchWithOptions(option){
    let url = "https://api.energidataservice.dk/dataset/Elspotprices"
    + `?start=${option.start}`
    + `&end=${option.end}`
    +`&filter={"PriceArea":["${option.area}"]}`
    const res = await fetch(url)
    const jsonres = await res.json()
    console.log(jsonres)
    return jsonres["records"]
}

// Creates a store consisting of records (type Record).
export function createRecords() {
    const { subscribe, set, update } = writable<Record[]>([]);

    return {
        subscribe,
        add: (record) => update( oldRecords => [...oldRecords, record]),
        populate: (recordList) => set(recordList),
        reset: () => set([]),
        getHours: () => (function(oldRecords) {
            let hours = []
            oldRecords.map( record => hours.push(record.HourDK));
            return hours
        }),
        getPrices: () => (function(oldRecords) {
            let prices = []
            oldRecords.map( record => prices.push(record.SpotPriceDKK));
            return prices
        })
    }

}
export let records = createRecords()
// export const hours = derived(
//     records,
//     $records => (function(oldRecords) {
//         oldRecords.map( record => hours.push(record.HourDK));
//     })
// )
//
// export const prices = derived(
//     records,
//     $records => (function(oldRecords) {
//         let prices = []
//         oldRecords.map( record => hours.push(record.SpotPriceDKK));
//         prices  = prices
//         return prices
//     })
// )
export const count = writable(0)