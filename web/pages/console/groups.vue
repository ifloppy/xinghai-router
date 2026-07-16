<script setup lang="ts">
import RouterApp from '~/components/RouterApp.vue'

definePageMeta({ middleware: 'console-auth' })
</script>

<template>
  <RouterApp active-view="groups" />
</template>
