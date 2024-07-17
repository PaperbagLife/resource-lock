<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import nicUrl from '../assets/nic.jpg'
import serverUrl from '../assets/server.png'
import { type Server, useResourceStatus } from '@/util/resource-status'
const { servers, refreshStatus } = useResourceStatus()
const expanded = ref<boolean[]>([])
watch(servers, (servers) => {
  console.log('server updated')
  expanded.value.length = 0
  servers.forEach((_) => expanded.value.push(false))
})

onMounted(() => {
  refreshStatus()
})
</script>

<template>
  <div class="row">
    <div
      v-for="(server, i) in servers"
      :key="server.name"
      class="mx-1 py-1 my-1 col-auto rounded server-outline"
    >
      <div class="row">
        <img class="col-auto" :src="serverUrl" width="50" height="50" />
        <div class="col-auto">
          <div class="lock-status-icon row material-symbols-outlined">
            {{ server.status.locked ? 'lock' : 'lock_open_right' }}
          </div>
          <button class="row" @click="expanded[i] = !expanded[i]">
            {{ server.name + (expanded[i] ? '△' : '▼') }}
          </button>
        </div>
      </div>

      <div v-if="expanded[i]" class="row">
        <div v-for="nic in server.nics" :key="nic.name" class="col-auto">
          <img :src="nicUrl" width="30" height="30" />
          <div>{{ nic.name }}</div>
          <div class="material-symbols-outlined">
            {{ nic.status.locked ? 'lock' : 'lock_open_right' }}
          </div>
        </div>
      </div>
    </div>
    <button @click="refreshStatus()">Refresh state</button>
  </div>
</template>

<style scoped></style>
