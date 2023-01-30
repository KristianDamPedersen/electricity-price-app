import {writable} from 'svelte/store'


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

//Makes an API Call to the danish public database, and populates the a store
export async function FetchWithOptions(option){
    const url = "https://api.energidataservice.dk/dataset/Elspotprices"
    + `?start=${option.start}`
    + `&end=${option.end}`
    +`&filter={"PriceArea":["${option.area}"]}`
    const res = await fetch(url)
    const jsonres = await res.json()
    return jsonres["records"].reverse()
}

// Creates a store consisting of records (type Record).
export function createRecords() {
    const { subscribe, set } = writable<Record[]>([]);

    return {
        subscribe,
        populate: (recordList) => set(recordList),
        reset: () => set([]),
    }

}
export const records = createRecords()

const convertToDateStringIncludingHours = (date) => {
    let monthstr =""
    let datestr = (date.getDate() < 10) ? "0" + date.getDate() : date.getDate()
    let monthString = ( (date.getMonth() +1) < 10) ? "0" + (date.getMonth() + 1) : date.getMonth() + 1
    let hourString = date.getHours()
    console.log(`${date.getFullYear()}-${monthString}-${datestr}T${hourString}:00`)
    return `${date.getFullYear()}-${monthString}-${datestr}T${hourString}:00`
}

//since the API  uses a weird time format, we need to first convert it
const convertDatetoStringExcludingHours = (date) => {
    let monthstr =""

    let datestr = (date.getDate() < 10) ? "0" + date.getDate() : date.getDate()
    let monthString = ( (date.getMonth() +1) < 10) ? "0" + (date.getMonth() + 1) : date.getMonth() + 1
    console.log(`${date.getFullYear()}-${monthString}-${datestr}T00:00`)
    return `${date.getFullYear()}-${monthString}-${datestr}T00:00`
}

// Fetches data from EnergiDataService's endpoint.
// numberOfDays: Is the number of days back in time we want to fetch. (we also count tomorrow as part of this)
// priceArea: Is the desired price area we want to retrieve from (currently supports DK1 and DK2)
// Returns: the raw json response from the endpoint
const fetchData = async (numberOfDays, priceArea) => {
    const today = new Date(Date.now())
    //calculating days behind is from  karim79 on stackoverflow https://stackoverflow.com/questions/1296358/how-to-subtract-days-from-a-plain-date
    const dateOffset = (24 * 60 * 60 * 1000) * (numberOfDays - 1); // We subtract one from number of days to account for including tomorrow as well
    const previousDate = new Date(today.getTime() - dateOffset)

    //add one with the same as above
    const offset = 24*60*60*1000;
    const tomorrow = new Date(today.getTime() + offset)

    //the function "fetch with options" does not give enough control so therefor we do a manual fetch
    const res = await fetch("https://api.energidataservice.dk/dataset/Elspotprices?" +
        `start=${convertDatetoStringExcludingHours(previousDate)}`
        + `&end=${convertDatetoStringExcludingHours(tomorrow)}&timezone=dk&filter={"PriceArea" : ["${priceArea}"]}`)
    return await res.json()
}

// Splits an array up into X chunks.
const splitArray = (array, chunks) => {
   const newArray = array.reduce((array, item, index) => {
       const chunkIndex = Math.floor(index/chunks)
       if (!array[chunkIndex]) {
           array[chunkIndex] = [] // Create a new chunk here
       }
       array[chunkIndex].push(item)
       return array
       }, [])
    return newArray
}

// Given an array of chunks containing floats ( [[1.2,2.3,3.1],[3.8,4.9,5.3]] etc.) it will return a new array
// containing the averages ([avg1, avg2, avg3) etc.
const averageChunkArrayFloats = (array) => {
    const averaged = []
    array.map(element => {
        const average = element.reduce( (a, b) => a + b, 0) / element.length
        averaged.push(average)
    })
    return averaged
}

// Given an array of chunks containing datetimes, it will return a new array containing the average time of each chunk.
const averageChunkArrayDatetime = (array) => {
    const averagedTimes = []
    array.map(element => {
        const milliseconds = []
        element.map(item => {
            milliseconds.push((new Date(item).getTime()))
        })

        const avgTimeInMilliseconds = (milliseconds.reduce( (a, b) => a+b, 0))/milliseconds.length
        const avgTime = new Date(avgTimeInMilliseconds)
        averagedTimes.push(avgTime)
    })
    return averagedTimes
}

// Given an array of floats (prices) it splits the array up into appropriate chunks and averages those chunks.
const averagePrices = (priceArray, averageOverHours) => {
    let priceArr = splitArray(priceArray, averageOverHours)
    let resultsArray = averageChunkArrayFloats(priceArr)
    return resultsArray
}

// Given an array of datetimes, it chops the array into appropriate chunks and averages the chunks
const averageTime = (timeArray, averageOverHours) => {
    let timeArr = splitArray(timeArray, averageOverHours)
    let resultsArray = averageChunkArrayDatetime(timeArr)
    return resultsArray
}
export const retrieveData = async (numberOfDays, averagedOverHours) => {

    // Fetch JSON from api endpoint
    const json = await fetchData(numberOfDays, "DK1")

    // Data treatment:
    const timesDKK: any[] = []
    const timesUTC: any[] = []
    const pricesEUR: any[] = []
    const pricesDKK: any[] = []

    const recordsJson = json["records"].reverse()
    recordsJson.map(record => {
        timesDKK.push(record["HourDK"])
        timesUTC.push(record["HourUTC"])
        pricesEUR.push(record["SpotPriceEUR"])
        pricesDKK.push(record["SpotPriceDKK"])
    })

    // Split up prices into appropriate chunks and average the chunks
    // Average prices in DKK
    const averagedPrices = averagePrices(pricesDKK, averagedOverHours)

    // Average prices in EUR
    const averagedPricesEUR = averagePrices(pricesEUR, averagedOverHours)

    // Split up timeslots into appropriate chunks and average them
    // Average times UTC + 1
    const averagedTimes = averageTime(timesDKK, averagedOverHours)

    // Average times UTC
    const averagedTimesUTC = averageTime(timesUTC, averagedOverHours)

    // Construct the record list
    const recordList = []
    for (let i=0; i < 24/averagedOverHours*numberOfDays; i++) {
        console.log(i)
        const record = {
            HourUTC: convertToDateStringIncludingHours(averagedTimesUTC[i]),
            HourDK: convertToDateStringIncludingHours(averagedTimes[i]),
            PriceArea: "DK1",
            SpotPriceDKK: averagedPrices[i],
            SpotPriceEUR: averagedPricesEUR[i],
        }
        recordList.push(record)
    }
    return recordList
}
