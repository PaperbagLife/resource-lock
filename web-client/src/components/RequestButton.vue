<script setup lang="ts">
import {ref} from 'vue'
import { servers, refreshStatus } from '@/util/resource-status';
</script>

<template>
  <div>{{ servers }}</div>
  <button @click="refreshStatus()">Click me</button>
</template>
