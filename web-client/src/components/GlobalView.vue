<script setup lang="ts">
import { ref } from 'vue'
import { type Server, useResourceStatus } from '@/util/resource-status'
const { servers, refreshStatus } = useResourceStatus()

const focusedServerIdx = ref<null | number>(null)
</script>

<template>
  <div>
    <template v-if="focusedServerIdx != null">
      <div>{{ servers[focusedServerIdx].name }}</div>
    </template>
    <div v-for="server in servers" :key="server.name">
      <div>
        {{ server.name }}
      </div>
      <div>{{ server.status }}</div>
    </div>
    <div>{{ servers }}</div>
    <button @click="refreshStatus()">Refresh state</button>
  </div>
</template>

<style scoped></style>
