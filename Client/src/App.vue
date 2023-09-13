<script setup>
  import { RouterLink, RouterView } from 'vue-router'
  import WebsocketProvider from './providers/WebsocketProvider.vue';
  import UserProvider from './providers/UserProvider.vue';
</script>

<template>
  <UserProvider>
    <WebsocketProvider>
      <RouterView />
    </WebsocketProvider>
  </UserProvider>
</template>

