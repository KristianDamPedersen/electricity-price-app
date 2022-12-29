import { json } from '@sveltejs/kit'
import { writable } from 'svelte/store'


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

export function createRecords() {
    const { subscribe, set, update } = writable<Record[]>([]);

    return {
        subscribe,
        add: (record) => update( oldRecords => [...oldRecords, record]),
        populate: (recordList) => set(recordList),
        reset: () => set([]),
    }

}

export const count = writable(0)