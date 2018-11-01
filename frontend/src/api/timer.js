import axios from "axios"
import { URL } from "./const"

export const create = () => {
    return axios.post(`${URL}/session/new`)
}

export const setTimer = (key, seconds) => {
    return axios.post(`${URL}/session/${key}/time`, { seconds })
}

export const get = (key) => {
    return  axios.get(`${URL}/session/${key}`)
}
