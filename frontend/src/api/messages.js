import axios from "axios"
import { URL } from "./const"

export const get = (key) => {
    return axios.get(`${URL}/session/${key}/messages`)
}

export const set = (key, messages) => {
    return axios.post(`${URL}/session/${key}/messages`, { messages })
}

export const send = (key, message) => {
    return axios.post(`${URL}/session/${key}/message`, message)
}
