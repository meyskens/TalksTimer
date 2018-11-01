import axios from "axios"
import { URL } from "./const"

export const send = (key, message) => {
    return axios.post(`${URL}/session/${key}/message`, { message })
}
