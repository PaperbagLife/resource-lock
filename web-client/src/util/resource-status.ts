import { ref } from 'vue'
import axios from 'axios'
import REQUEST_URL from '../../requestURL.txt?raw'
export enum ResourceType {
  NIC = 0,
  SERVER = 1
}


type ResourceStatus = {
  locked: boolean
  name?: string // who locked
  endTime?: number // Unix timestamp in milliseconds, or -1 for infite length
}

export type Nic = {
  type: ResourceType.NIC
  name: string // mlx name of the nic
  status: ResourceStatus
}

export type Server = {
  type: ResourceType.SERVER
  name: string // aa, smokey, oracle etc
  status: ResourceStatus
  nics: Nic[]
}

const servers = ref<Server[]>([])

async function refreshStatus() {
  console.log(REQUEST_URL)
  // async function sendPostRequest() {
  try {
    const { data } = await axios.post<Server[]>(`${REQUEST_URL}/status`, {})
    servers.value = data
  } catch (error) {
    console.log(error)
  }
}

export async function setStatus(newServers: Server[]) {
  try {
    console.log(JSON.stringify(newServers))
    await axios.post(`${REQUEST_URL}/update`, JSON.stringify(newServers))
    console.log('update complete')
    refreshStatus()
  } catch (error) {
    console.log(error)
  }
}

export function useResourceStatus() {
  return {
    servers: servers,
    refreshStatus: refreshStatus,
    setStatus: setStatus
  }
}
