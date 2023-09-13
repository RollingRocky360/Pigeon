<script>

    const userId = localStorage.getItem("userId")
    const websocket = new WebSocket(`ws://localhost:3000/ws/${userId}`)

    export default {
        provide: { websocket }
    }

</script>

<template>
    <slot></slot>
</template>
