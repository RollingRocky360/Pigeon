<script>
    // const resp = await fetch('http://localhost:3000/user', {
    //     headers: {
    //         'Authorization': 'BEARER ' + localStorage.getItem("userId")
    //     }
    // })
    // if (!resp.ok) console.log("User Fetch failed")
    const user = {
        id: localStorage.getItem("userId")
    }

    export default {
        provide: { user }
    }
</script>

<template>
    <slot></slot>
</template>