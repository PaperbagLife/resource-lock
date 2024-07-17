<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import nicUrl from '../assets/nic.jpg'
import serverUrl from '../assets/server.png'
import { type Server, type Nic, useResourceStatus } from '@/util/resource-status'
const { servers, refreshStatus, setStatus } = useResourceStatus()
const expanded = ref<boolean[]>([])
watch(servers, (servers) => {
  expanded.value.length = 0
  if (servers) {
    servers.forEach((_) => expanded.value.push(false))
  }
})

onMounted(() => {
  refreshStatus()
})

const selectedResources = ref(new Set<string>())

function getResourceID(server: Server, nic?: Nic) {
  if (nic) {
    return `${server.name}~${nic.name}`
  }
  return server.name
}

function getNameFromID(ID: string) {
  if (ID.split('~').length === 1) {
    return { server: ID }
  } else return { server: ID.split('~')[0], nic: ID.split('~')[1] }
}

function addOrDelete(event: InputEvent, server: Server, nic?: Nic) {
  console.log(event)
  if (event.target?.checked) {
    selectedResources.value.add(getResourceID(server, nic))
  } else {
    selectedResources.value.delete(getResourceID(server, nic))
  }
}

const username = ref('')

const lockTime = ref(0)

const newServers = ref<Server[]>([])

const showWarning = ref(false)

const lockOwners = ref<string[]>([])

function updateResources() {
  const endTime = new Date().getTime() + lockTime.value * 60000
  console.log(username.value, endTime)
  let occupied = false
  newServers.value.length = 0
  lockOwners.value.length = 0
  servers.value.forEach((server) => {
    const nics: Nic[] = []
    server.nics.forEach((nic) => {
      const resourceID = getResourceID(server, nic)
      if (selectedResources.value.has(resourceID)) {
        if (nic.status.locked) {
          occupied = true
          lockOwners.value.push(nic.status.name ?? '')
        }
        nics.push({ ...nic, status: { locked: true, name: username.value, endTime } })
      } else {
        nics.push({ ...nic })
      }
    })
    const serverID = getResourceID(server)
    if (selectedResources.value.has(serverID)) {
      if (server.status.locked) {
        occupied = true
        lockOwners.value.push(server.status.name ?? '')
      }
      newServers.value.push({
        ...server,
        nics: nics,
        status: { locked: true, name: username.value, endTime }
      })
    } else {
      newServers.value.push({ ...server, nics: nics })
    }
  })
  if (!occupied) {
    setStatus(newServers.value)
  } else {
    showWarning.value = true
  }
}
</script>

<template>
  <div class="row">
    <button @click="refreshStatus()">Refresh state</button>
    <div>Selected resources: {{ selectedResources }}</div>
    <div>Click on name of server to see nics</div>
    <div
      v-for="(server, i) in servers"
      :key="server.name"
      class="mx-1 py-1 my-1 col-auto rounded server-outline"
    >
      <div v-if="server.status.locked && server.status.endTime">
        {{ server.status.name }} until
        {{ new Date(server.status.endTime).toLocaleString() }}
      </div>
      <div v-else>Free</div>
      <div class="row mx-0">
        <img class="col-auto px-0" :src="serverUrl" width="50" height="50" />
        <div class="col-auto d-flex align-items-center">
          <input type="checkbox" @change="(event) => addOrDelete(event, server)" />
          <div
            class="lock-status-icon material-symbols-outlined"
            :class="server.status.locked ? 'locked' : 'unlocked'"
          >
            {{ server.status.locked ? 'lock' : 'lock_open_right' }}
          </div>

          <button class="h-50" @click="expanded[i] = !expanded[i]">
            {{ server.name + (expanded[i] ? '△' : '▼') }}
          </button>
        </div>
      </div>

      <div v-if="expanded[i]" class="row">
        <div v-for="nic in server.nics" :key="nic.name" class="col-auto">
          <img :src="nicUrl" width="30" height="30" />
          <div>{{ nic.name }}</div>
          <div
            class="lock-status-icon material-symbols-outlined"
            :class="nic.status.locked ? 'locked' : 'unlocked'"
            :title="
              nic.status.locked && nic.status.endTime
                ? `${nic.status.name} until ${new Date(nic.status.endTime).toLocaleString()}`
                : 'free'
            "
          >
            {{ nic.status.locked ? 'lock' : 'lock_open_right' }}
          </div>
          <input type="checkbox" @change="(event) => addOrDelete(event, server, nic)" />
        </div>
      </div>
    </div>
    <!-- Button trigger modal -->
    <button
      type="button"
      class="btn btn-primary"
      data-bs-toggle="modal"
      data-bs-target="#exampleModal"
    >
      Lock Selected Resource
    </button>

    <!-- Modal -->
    <div
      class="modal fade"
      id="exampleModal"
      tabindex="-1"
      aria-labelledby="exampleModalLabel"
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="exampleModalLabel">Take lock on resources</h5>
            <button
              type="button"
              class="btn-close"
              data-bs-dismiss="modal"
              aria-label="Close"
            ></button>
          </div>
          <div class="modal-body">
            <div>Name:</div>
            <input v-model="username" />
            <div>Take the lock for</div>
            <button @click="lockTime = 10">10 min</button>
            <button @click="lockTime = 60">1hour</button>
            <button @click="lockTime = 240">4hours</button>
            <div></div>
            <input type="number" v-model="lockTime" />
            <span class="mx-2">Minutes</span>
            <div v-if="showWarning">
              WARNING: You are attempting to lock a resource currently in use. please make sure you
              coordinate with the current lock owners {{ lockOwners }}
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
            <button type="button" class="btn btn-primary" @click="updateResources">
              Lock the resources
            </button>
            <button
              v-if="showWarning"
              type="button"
              class="btn btn-primary"
              data-bs-dismiss="modal"
              @click="setStatus(newServers)"
            >
              I understand and will proceed
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
