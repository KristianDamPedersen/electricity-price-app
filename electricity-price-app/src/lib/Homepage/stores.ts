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
export const records = writable<Record[]>()

//Makes an APICall to the danish public database, and populates the a store
export function populateWithAPICall(option : option){
    let url = "https://api.energidataservice.dk/dataset/Elspotpriceshttps://api.energidataservice.dk/dataset/Elspotprices"
    + `?start=${option.start}`
    + `&end=${option.end}`
    +`&filter={"PriceArea":["${option.PriceArea}"]}`

}

export function createRecords() {
    const { subscribe, set, update } = writable<Record[]>([]);

    return {
        subscribe,
        add: (record) => update( oldRecords => [...oldRecords, record]),
        reset: () => set([]),
    }

}

export const count = writable(0)