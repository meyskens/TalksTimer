import axios from "axios"
import { URL } from "./const"

export const get = (key) => {
    return axios.get(`${URL}/session/${key}/colors`)
}

export const set = (key, options) => {
    return axios.post(`${URL}/session/${key}/colors`, { options })
}
