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
      <div v-if="server.status.locked && server.status.endTime">
        {{ server.status.name }} until
        {{ new Date(server.status.endTime).toLocaleString() }}
      </div>
      <div v-else>Free</div>
      <div class="row mx-0">
        <img class="col-auto px-0" :src="serverUrl" width="50" height="50" />
        <div class="col-auto d-flex align-items-center">
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
          <div class="material-symbols-outlined">
            {{ nic.status.locked ? 'lock' : 'lock_open_right' }}
          </div>
        </div>
      </div>
    </div>
    <div>{{ servers }}</div>
    <button @click="refreshStatus()">Refresh state</button>
  </div>
</template>

<style scoped></style>
