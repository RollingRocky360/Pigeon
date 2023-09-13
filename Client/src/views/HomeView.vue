<script>
    import { RouterView } from 'vue-router';
    import ChatList from '../components/ChatList.vue';
    export default {
        data() {
            return {
                chats: []
            }
        },
        components: {
            ChatList,
        }
    }
</script>

<template>
    <div id="home">
        <ChatList :chats="chats" @initchats="(c) => chats = c"/>
        <router-view @newchat="(newchat) => chats.push(newchat)" :key="$route.fullPath"/>
    </div>
</template>

<style scoped>
    #home {
        width: 98%;
        height: 100%;
        margin: auto;
        box-sizing: border-box;
        display: grid;
        place-items: center;
        grid-template-columns: .3fr .7fr;
        grid-template-rows: 1fr;
    }
</style>