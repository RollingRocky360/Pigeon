<script>
    export default {
        inject: ['websocket', 'user', 'currentChat'],
        props: ['chats'],
        data() {
            return {
                searchterm: '',
            }
        },
        async mounted() {
            console.log('mounted Chat list')
            const resp = await fetch('http://localhost:3000/chat', {
                headers: {
                    'Authorization': 'BEARER ' + localStorage.getItem("userId")
                }
            })
            if (resp.ok)
            this.$emit('initchats', await resp.json())
        },
        computed: {
            filteredChats() {
                if (this.searchterm.length === 0) return this.chats
                return this.chats.filter(chat => {
                    for (let user of chat.users) {
                        if (user.username.toLowerCase().startsWith(this.searchterm.toLowerCase())) 
                            return true
                    } 
                    return false
                })
            }
        }
    }
</script>

<template>
    <div v-if="chats.length > 0" id="chatlist-container">
        <div id="title-and-create">
            <p><b>Pigeon Chat</b></p>
            <router-link :to="'/new'" id="create-chat" title="New Chat">
                <span class="material-symbols-outlined">
                    add
                </span>
            </router-link>
        </div>
        <input type="text" placeholder="search" id="searchterm" v-model="searchterm">
        <ul id="chatlist">
            <li v-for="chat in filteredChats" :key="chat.id">
                <router-link :to="`/chat/${chat.id}`">
                    <div class="chat" @click="() => $emit('changeChat', chat)">
                        <div class="pfp" :style="{ backgroundColor: chat.pfp_color }">
                            {{ (chat.users[0].id !== user.id ? chat.users[0].username : chat.users[1].username)[0] }}
                        </div>
                        <p class="chatname">
                            {{ chat.users[0].id !== user.id ? chat.users[0].username : chat.users[1].username }}
                        </p>
                    </div>
                </router-link>
            </li>
        </ul>
    </div>
</template>

<style scoped>
    #chatlist-container {
        width: 99%;
        height: 95%;
        background-color: #ffffff;
        box-sizing: border-box;
        padding: 1rem;
        border-radius: .5rem;
    }

    #searchterm {
        all: unset;
        box-sizing: border-box;
        display: block;
        width: 95%;
        height: 2.2rem;
        padding: .2rem 1rem;
        margin: 0 auto;
        border: 1px solid #727375;
        border-radius: 5rem;
    }
    #searchterm::placeholder {
        color: #72737596
    }
    #searchterm:focus::placeholder {
        color: transparent;
    }

    #title-and-create {
        width: 95%;
        box-sizing: border-box;
        margin: auto;
        margin-bottom: 1rem;
        font-size: 1.2rem;
        font-weight: 700;
        display: flex;
        align-items: center;
        justify-content: space-between;
    }
    
    #create-chat {
        all: unset;
        background-color: #9090ff;
        color: white;
        width: 2rem;
        height: 2rem;
        display: flex;
        align-items: center;
        justify-content: center;
        border-radius: 50%;
        transition-duration: .2s;
        cursor: pointer;
    }
    #create-chat:hover {
        box-shadow: 0 0 15px 1px #9090ff80;
        transition-duration: .2s;
    }

    #chatlist {
        padding-left: 0;
        width: 100%;
        list-style-type: none;
        margin: 1rem auto;
        overflow-y: auto;
    }

    li {
        width: 100%;
        height: 3.5rem;
    }

    router-link, .chat {
        width: 100%;
        height: 100%;
        text-decoration: none;
    }

    a {
        all: unset;
    }

    .chat {
        display: flex;
        align-items: center;
        justify-content: flex-start;
        box-sizing: border-box;
        padding: 0 .5rem;
        gap: .7rem;
        border-radius: 5rem;
        cursor: default;
    }
    .chat:hover {
        background-color: #00000010
    }

    .pfp {
        font-size: 1.5rem;
        color: white;
        display: flex;
        align-items: center;
        justify-content: center;
        width: 2.5rem;
        height: 2.5rem;
        border-radius: 50%;
    }

    .chatname {
        font-size: 1.15rem;
    }
</style>